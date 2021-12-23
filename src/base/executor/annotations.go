package executor

import (
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
)

const StrategyAnnotationName = "strategy"

type strategySelectionConfig struct {
	annotationPrefix string
	strategies       map[strategy.YAMLName]strategy.Config
	defaultStrategy  *strategy.Config
}

type providerSelectionConfig struct {
	annotationPrefix string
	providers        map[provider.Name]provider.Config
	defaultProviders map[parameter.Name]provider.Name
}

func (cfg strategySelectionConfig) selectAppStrategy(appAnnotations map[string]string) (strategy.Config, error) {
	name, ok := appAnnotations[cfg.annotationPrefix+StrategyAnnotationName]
	if !ok {
		if cfg.defaultStrategy != nil {
			return *cfg.defaultStrategy, nil
		}
		return strategy.Config{}, fmt.Errorf("strategy not specified in app's annotations and default strategy isn't configured")
	}
	selected, ok := cfg.strategies[strategy.YAMLName(name)]
	if !ok {
		return strategy.Config{}, fmt.Errorf("'%s' strategy specified in annotations doesn't exist", name)
	}
	return selected, nil
}

func (cfg providerSelectionConfig) selectFor(
	strategyCfg strategy.Config, annotations map[string]string) (map[provider.Name][]parameter.Name, parameter.Values, error) {

	requiredParams := map[provider.Name][]parameter.Name{}
	yamlProvided := parameter.EmptyValues()

	for paramName, spec := range strategyCfg.GetRequiredParameters() {
		annValue, ok := annotations[cfg.annotationPrefix+string(paramName)]
		if !ok {
			provName, ok := cfg.defaultProviders[paramName]
			if !ok {
				if err := cfg.trySpecDefault(paramName, spec, yamlProvided); err != nil {
					return nil, parameter.Values{}, err
				}
				continue
			}
			providerCfg, ok := cfg.providers[provName]
			if !ok {
				if err := cfg.trySpecDefault(paramName, spec, yamlProvided); err != nil {
					return nil, parameter.Values{}, err
				}
				continue
			}
			paramType, ok := providerCfg.AvailableParameters[paramName]
			if !ok || !paramType.EqualTo(spec.Type) {
				return nil, parameter.Values{}, fmt.Errorf(
					"'%s' provider doesn't have available parameter '%s' with type %s",
					providerCfg.Name, paramName, paramType.Name,
				)
			}
			requiredParams[providerCfg.Name] = append(requiredParams[providerCfg.Name], paramName)
			continue

		}
		providerCfg, ok := cfg.providers[provider.Name(annValue)]
		if ok {
			paramType, ok := providerCfg.AvailableParameters[paramName]
			if !ok || !paramType.EqualTo(spec.Type) {
				return nil, parameter.Values{}, fmt.Errorf(
					"'%s' provider doesn't have available parameter '%s' with type %s",
					providerCfg.Name, paramName, paramType.Name,
				)
			}
			requiredParams[providerCfg.Name] = append(requiredParams[providerCfg.Name], paramName)
			continue
		}
		strConv := spec.Type.StrConv
		converted, err := strConv(annValue)
		if err != nil {
			return nil, parameter.Values{}, fmt.Errorf("failed to set value for required '%s' parameter: %w", paramName, err)
		}
		if err := yamlProvided.Insert(paramName, converted, spec.Type); err != nil {
			return nil, parameter.Values{}, fmt.Errorf("failed to set value for required '%s' parameter: %w", paramName, err)
		}
		continue
	}
	return requiredParams, yamlProvided, nil
}

func (cfg providerSelectionConfig) trySpecDefault(paramName parameter.Name, spec strategy.ParameterSpec, yamlProvided parameter.Values) error {
	if spec.DefaultValue == nil {
		return fmt.Errorf("required %s parameter with no default value is not specified in annotations", paramName)
	}
	if err := yamlProvided.Insert(paramName, spec.DefaultValue, spec.Type); err != nil {
		return err
	}
	return nil
}
