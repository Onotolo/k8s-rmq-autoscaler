package executor

import (
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
)

func (cfg Config) Validate() []error {
	var errs []error
	foundDefault := false
	defaultStrategyName := cfg.DefaultStrategy
	for _, strategyCfg := range cfg.EnabledStrategies {
		if err := strategyCfg.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("validation failed for '%s' strtategy: %w", strategyCfg.Name, err))
		}
		if !foundDefault && strategyCfg.YAMLName == defaultStrategyName {
			foundDefault = true
		}
	}
	if !foundDefault && len(defaultStrategyName) != 0 {
		errs = append(errs, fmt.Errorf("default strategy '%s' not found among enabled strategies", defaultStrategyName))
	}

	enabledProviders := map[provider.Name]provider.Config{}
	for _, enabledProvider := range cfg.EnabledProviders {
		enabledProviders[enabledProvider.Name] = enabledProvider
	}
	for paramName, provName := range cfg.DefaultParametersProviders {
		prov, ok := enabledProviders[provName]
		if !ok {
			errs = append(errs, fmt.Errorf("default '%s' parameter provider '%s' not found among enabled providers", paramName, provName))
			continue
		}
		if _, ok := prov.AvailableParameters[paramName]; !ok {
			errs = append(errs, fmt.Errorf("provider '%s' set as default for '%s' parameter doesn't have available parameter with this name", provName, paramName))
		}
	}
	return errs
}
