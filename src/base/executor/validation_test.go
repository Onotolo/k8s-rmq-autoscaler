package executor

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/stretchr/testify/require"
	"testing"
)

func makeConfig(strategyParameters map[parameter.Name]strategy.ParameterSpec) Config {
	return Config{
		EnabledStrategies: []strategy.Config{
			makeStrategyConfig(strategyParameters),
		},
	}
}

func TestConfig_Validate(t *testing.T) {
	config := makeConfig(
		map[parameter.Name]strategy.ParameterSpec{
			"int":    {Type: parameter.Int},
			"float":  {Type: parameter.Float},
			"string": {Type: parameter.String},
		},
	)
	require.Empty(t, config.Validate())

	// Wrong parameter's default value
	config = makeConfig(
		map[parameter.Name]strategy.ParameterSpec{
			"int": {Type: parameter.Int, DefaultValue: "213"},
		},
	)
	require.NotEmpty(t, config.Validate())
}

func TestConfig_Validate_withDefaultStrategy(t *testing.T) {
	// Default strategy that isn't present in enabled strategies
	config := makeConfig(
		map[parameter.Name]strategy.ParameterSpec{
			"int":    {Type: parameter.Int},
			"string": {Type: parameter.String, DefaultValue: "str value"},
		},
	)
	config.DefaultStrategy = "default_strategy"
	errs := config.Validate()
	require.NotEmpty(t, errs)
}

func TestConfig_Validate_withDefaultParamProviders(t *testing.T) {
	// Default provider for 'int' param present in enabled providers
	config := makeConfig(
		map[parameter.Name]strategy.ParameterSpec{
			"int": {Type: parameter.Int},
		},
	)
	config.EnabledProviders = []provider.Config{
		{
			Name: "int_provider",
			AvailableParameters: map[parameter.Name]parameter.Type{
				"int": parameter.Int,
			},
		},
	}
	config.DefaultParametersProviders = map[parameter.Name]provider.Name{
		"int": "int_provider",
	}
	require.Empty(t, config.Validate())

	// Default provider for 'int' param isn't present in enabled providers
	config.DefaultParametersProviders = map[parameter.Name]provider.Name{
		"int": "missing_provider",
	}
	errs := config.Validate()
	require.NotEmpty(t, errs)
}
