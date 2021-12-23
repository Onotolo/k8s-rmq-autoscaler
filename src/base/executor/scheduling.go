package executor

import (
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"sync"
)

type providerSchedulingResult struct {
	appsProvidersResults map[scalable.App][]provider.ResultAppContext
	staticAppParameters  map[scalable.App]parameter.Values
}

func (ex executor) scheduleProviders() providerSchedulingResult {
	providersCfg := map[provider.Name]provider.Config{}

	for _, providerCfg := range ex.config.EnabledProviders {
		providersCfg[providerCfg.Name] = providerCfg
	}

	requiredParams := map[provider.Name]provider.RequiredAppsParameters{}
	staticAppParameters := map[scalable.App]parameter.Values{}

	providerSelection := providerSelectionConfig{
		annotationPrefix: ex.config.AnnotationsPrefix,
		providers:        providersCfg,
		defaultProviders: ex.config.DefaultParametersProviders,
	}
	for _, app := range ex.apps {
		strategyCfg := ex.appsStrategiesConfigs[app]
		appProvidersParameters, providedValues, err := providerSelection.selectFor(strategyCfg, *app.Annotations)
		if err != nil {
			ex.out.errors <- BaseError{
				App:      app,
				Strategy: strategyCfg,
				Err:      err,
			}
			continue
		}
		for provName, paramNames := range appProvidersParameters {
			providerParams, ok := requiredParams[provName]
			if !ok {
				providerParams = map[scalable.App][]parameter.Name{}
				requiredParams[provName] = providerParams
			}
			providerParams[app] = paramNames
		}
		staticAppParameters[app] = providedValues
	}

	appsProvidersResults := map[scalable.App][]provider.ResultAppContext{}

	for name, appsParameters := range requiredParams {
		provResults := provider.Launch(providersCfg[name], appsParameters)
		for app, provResult := range provResults {
			appsProvidersResults[app] = append(appsProvidersResults[app], provResult)
		}
	}
	return providerSchedulingResult{
		appsProvidersResults: appsProvidersResults,
		staticAppParameters:  staticAppParameters,
	}
}

func (ex executor) scheduleStrategies(schedulingResult providerSchedulingResult) <-chan struct{} {

	done := make(chan struct{}, 1)

	appsParameters := schedulingResult.staticAppParameters
	appsWg := sync.WaitGroup{}

	for _, app := range ex.apps {
		strategyCfg := ex.appsStrategiesConfigs[app]
		providersAppContexts := schedulingResult.appsProvidersResults[app]

		appParams := make(chan parameter.Values)

		cancelAppProviders := func() {
			for _, appContext := range providersAppContexts {
				appContext.Cancel()
			}
		}
		wg := sync.WaitGroup{}

		for _, providerContext := range providersAppContexts {
			wg.Add(1)
			go ex.collectProvidersResults(app, providerContext, strategyCfg, appParams, cancelAppProviders, wg.Done)
		}
		go func() {
			wg.Wait()
			close(appParams)
		}()
		appsWg.Add(1)
		go ex.collectAppParams(app, strategyCfg, appsParameters[app], appParams, cancelAppProviders, appsWg.Done)
	}
	go func() {
		appsWg.Wait()
		close(done)
	}()
	return done
}

func (ex executor) collectProvidersResults(
	app scalable.App,
	providerContext provider.ResultAppContext,
	strategyCfg strategy.Config,
	appParams chan<- parameter.Values,
	cancelAppProviders, done func()) {

	defer done()
	for {
		result, ok := providerContext.GetNextResult()
		if !ok {
			return
		}
		if result.Error != nil {
			cancelAppProviders()
			ex.out.errors <- ProviderError{
				BaseError: BaseError{
					App:      app,
					Strategy: strategyCfg,
					Err:      result.Error,
				},
				ProviderName: providerContext.ProviderName,
			}
			return
		}
		converted, err := convertParameters(result.Parameters, strategyCfg.GetRequiredParameters())
		if err != nil {
			cancelAppProviders()
			ex.out.errors <- BaseError{
				App:      app,
				Strategy: strategyCfg,
				Err:      err,
			}
			return
		}
		appParams <- converted
	}
}

func (ex executor) collectAppParams(
	app scalable.App,
	strategyCfg strategy.Config,
	yamlProvided parameter.Values,
	appParams <-chan parameter.Values,
	cancelAppProviders, done func()) {

	defer done()
	collectedParams := yamlProvided

	for appParam := range appParams {
		merged := collectedParams.Merge(appParam)
		collectedParams = merged

		if strategy.Ready(strategyCfg, merged) {
			cancelAppProviders()
			result, err := strategy.Execute(strategyCfg, app, merged)
			if err != nil {
				ex.out.errors <- BaseError{App: app, Strategy: strategyCfg, Err: err}
				return
			}
			ex.out.results <- result
			return
		}
	}
	requiredParamsNames := make([]parameter.Name, 0, len(strategyCfg.RequiredParameters))
	for name := range strategyCfg.GetRequiredParameters() {
		requiredParamsNames = append(requiredParamsNames, name)
	}
	ex.out.errors <- BaseError{
		App:      app,
		Strategy: strategyCfg,
		Err:      fmt.Errorf("could not collect all required parameters: required %v, provided %+v", requiredParamsNames, collectedParams),
	}
}

func (ex executor) scheduleCleanup(done <-chan struct{}) {
	go func() {
		<-done
		ex.cleanup()
	}()
}
