package executor

import (
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"sync"
)

type executor struct {
	apps                  []scalable.App
	config                Config
	appsStrategiesConfigs map[scalable.App]strategy.Config
	out                   output
	done                  chan struct{}
}

type output struct {
	results  chan<- strategy.Result
	errors   chan<- Error
	errorsWg *sync.WaitGroup
}

func makeExecutor(config Config, apps []scalable.App) (executor, <-chan strategy.Result, <-chan Error) {
	errs := make(chan Error)
	results := make(chan strategy.Result)

	ex := executor{
		apps:   apps,
		config: config,
		out: output{
			results:  results,
			errors:   errs,
			errorsWg: &sync.WaitGroup{},
		},
		done: make(chan struct{}),
	}
	appsStrategies := map[scalable.App]strategy.Config{}
	strategyConfigs := map[strategy.YAMLName]strategy.Config{}

	for _, strategyCfg := range ex.config.EnabledStrategies {
		strategyConfigs[strategyCfg.YAMLName] = strategyCfg
	}
	reportError := func(app scalable.App, err error) {
		ex.out.errors <- BaseError{
			App: app,
			Err: err,
		}
	}
	strategySelector := strategySelectionConfig{
		annotationPrefix: config.AnnotationsPrefix,
		strategies:       strategyConfigs,
	}
	defaultStrategy, ok := strategyConfigs[config.DefaultStrategy]
	if ok {
		strategySelector.defaultStrategy = &defaultStrategy
	}

	for _, app := range apps {
		selected, err := strategySelector.selectAppStrategy(*app.Annotations)
		if err != nil {
			go reportError(app, fmt.Errorf("could not select strategy: %w", err))
			continue
		}
		appsStrategies[app] = selected
	}
	ex.appsStrategiesConfigs = appsStrategies
	return ex, results, errs
}

func (ex executor) start() {
	provSchedulingResult := ex.scheduleProviders()
	done := ex.scheduleStrategies(provSchedulingResult)
	ex.scheduleCleanup(done)
}

func (ex executor) cleanup() {
	close(ex.out.errors)
	close(ex.out.results)
	close(ex.done)
}
