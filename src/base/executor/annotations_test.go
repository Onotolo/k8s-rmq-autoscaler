package executor

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/stretchr/testify/require"
	"testing"
)

var providerCfgs = map[provider.Name]provider.Config{
	"int_provider": {
		Name: "int_provider",
		AvailableParameters: map[parameter.Name]parameter.Type{
			"int": parameter.Int,
		},
	},
	"string_provider": {
		Name: "string_provider",
		AvailableParameters: map[parameter.Name]parameter.Type{
			"string": parameter.String,
		},
	},
	"float_provider": {
		Name: "float_provider",
		AvailableParameters: map[parameter.Name]parameter.Type{
			"float": parameter.Float,
		},
	},
}

func makeProvSelectionConfig(defaultProviders map[parameter.Name]provider.Name) providerSelectionConfig {
	return providerSelectionConfig{
		annotationPrefix: "prefix/",
		providers:        providerCfgs,
		defaultProviders: defaultProviders,
	}
}

var providerSelectionCfg = makeProvSelectionConfig(map[parameter.Name]provider.Name{})

func makeStrategyConfig(parameters map[parameter.Name]strategy.ParameterSpec) strategy.Config {
	return strategy.Config{
		Name:               "test_strategy",
		YAMLName:           "test_strategy",
		RequiredParameters: parameters,
		Execute: func(app scalable.App, values parameter.Values) (strategy.Result, error) {
			panic("test implementation")
		},
	}
}

func TestStrategySelectorConfig(t *testing.T) {
	cfg := strategySelectionConfig{
		annotationPrefix: "prefix/",
		strategies: map[strategy.YAMLName]strategy.Config{
			"strategy_1": {
				Name:     "strategy_1",
				YAMLName: "strategy_1",
			},
			"strategy_2": {
				Name:     "strategy_2",
				YAMLName: "strategy_2",
			},
		},
		defaultStrategy: &strategy.Config{
			Name: "default_strategy",
		},
	}
	selected, err := cfg.selectAppStrategy(
		map[string]string{
			"prefix/strategy": "strategy_1",
		},
	)
	require.NoError(t, err)
	require.Equal(t, strategy.Name("strategy_1"), selected.Name)

	selected, err = cfg.selectAppStrategy(
		map[string]string{
			"prefix/strategy": "strategy_2",
		},
	)
	require.NoError(t, err)
	require.Equal(t, strategy.Name("strategy_2"), selected.Name)

	selected, err = cfg.selectAppStrategy(
		map[string]string{},
	)
	require.NoError(t, err)
	require.Equal(t, strategy.Name("default_strategy"), selected.Name)

	_, err = cfg.selectAppStrategy(
		map[string]string{
			"prefix/strategy": "nonexistent strategy",
		},
	)
	require.Error(t, err)

	// Without default strategy
	cfg.defaultStrategy = nil
	_, err = cfg.selectAppStrategy(
		map[string]string{},
	)
	require.Error(t, err)

	_, err = cfg.selectAppStrategy(
		map[string]string{
			"prefix/strategy": "nonexistent strategy",
		},
	)
	require.Error(t, err)
}

func TestProviderSelectorConfig(t *testing.T) {
	params, yamlProvided, err := providerSelectionCfg.selectFor(
		makeStrategyConfig(
			map[parameter.Name]strategy.ParameterSpec{
				"int":    {Type: parameter.Int},
				"float":  {Type: parameter.Float},
				"string": {Type: parameter.String},
			},
		),
		map[string]string{
			"prefix/int":    "int_provider",
			"prefix/float":  "float_provider",
			"prefix/string": "string_provider",
		},
	)
	require.NoError(t, err)
	require.Len(t, params, 3)
	require.Equal(t, 0, yamlProvided.Len())
}

func TestProviderSelectorConfig_withYAMLProvidedValues(t *testing.T) {
	params, yamlProvided, err := providerSelectionCfg.selectFor(
		makeStrategyConfig(
			map[parameter.Name]strategy.ParameterSpec{
				"int":    {Type: parameter.Int},
				"float":  {Type: parameter.Float},
				"string": {Type: parameter.String},
			},
		),
		map[string]string{
			"prefix/int":    "5",
			"prefix/float":  "5.1",
			"prefix/string": "some_string",
		},
	)
	require.NoError(t, err)
	require.Len(t, params, 0)
	require.Equal(t, 3, yamlProvided.Len())
}

func TestProviderSelectorConfig_withDefaults(t *testing.T) {
	params, yamlProvided, err := providerSelectionCfg.selectFor(
		makeStrategyConfig(
			map[parameter.Name]strategy.ParameterSpec{
				"int":    {Type: parameter.Int, DefaultValue: 1},
				"float":  {Type: parameter.Float, DefaultValue: 1.},
				"string": {Type: parameter.String, DefaultValue: "default"},
			},
		),
		map[string]string{},
	)
	require.NoError(t, err)
	require.Len(t, params, 0)
	require.Equal(t, 3, yamlProvided.Len())

	providerSelectionCfg := makeProvSelectionConfig(
		map[parameter.Name]provider.Name{
			"int":    "int_provider",
			"float":  "float_provider",
			"string": "string_provider",
		},
	)
	params, yamlProvided, err = providerSelectionCfg.selectFor(
		makeStrategyConfig(
			map[parameter.Name]strategy.ParameterSpec{
				"int":    {Type: parameter.Int},
				"float":  {Type: parameter.Float},
				"string": {Type: parameter.String},
			},
		),
		map[string]string{},
	)
	require.NoError(t, err)
	require.Len(t, params, 3)
	require.Equal(t, 0, yamlProvided.Len())
}

func TestProviderSelectorConfig_withError(t *testing.T) {
	// Parameter with no default value not specified in annotations
	_, _, err := providerSelectionCfg.selectFor(
		makeStrategyConfig(
			map[parameter.Name]strategy.ParameterSpec{
				"unprovided_int": {Type: parameter.Int},
			},
		),
		map[string]string{},
	)
	require.Error(t, err)

	// Wrong parameter type specified in annotations
	_, _, err = providerSelectionCfg.selectFor(
		makeStrategyConfig(
			map[parameter.Name]strategy.ParameterSpec{
				"int": {Type: parameter.Int},
			},
		),
		map[string]string{
			"int": "1.2",
		},
	)
	require.Error(t, err)

	// Different types for parameter in strategy and specified provider
	_, _, err = providerSelectionCfg.selectFor(
		makeStrategyConfig(
			map[parameter.Name]strategy.ParameterSpec{
				"int": {Type: parameter.Float},
			},
		),
		map[string]string{
			"int": "int_provider",
		},
	)
	require.Error(t, err)
}
