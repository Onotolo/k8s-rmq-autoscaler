package executor

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
)

type Config struct {
	EnabledStrategies          []strategy.Config
	EnabledProviders           []provider.Config
	AnnotationsPrefix          string
	DefaultStrategy            strategy.YAMLName
	DefaultParametersProviders map[parameter.Name]provider.Name
	FallbackToDefaultStrategy  bool
}

func Launch(config Config, apps []scalable.App) (<-chan strategy.Result, <-chan Error) {
	if config.DefaultParametersProviders == nil {
		config.DefaultParametersProviders = map[parameter.Name]provider.Name{}
	}
	ex, results, errs := makeExecutor(config, apps)
	go ex.start()
	return results, errs
}
