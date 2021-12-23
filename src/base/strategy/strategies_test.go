package strategy

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

var correctParams = []RequiredParameters{
	{},
	{
		"int":    {Type: parameter.Int},
		"float":  {Type: parameter.Float},
		"string": {Type: parameter.String},
	},
	{
		"int":    {Type: parameter.Int, DefaultValue: 42},
		"float":  {Type: parameter.Float, DefaultValue: 4.2},
		"string": {Type: parameter.String, DefaultValue: "42"},
	},
}

var incorrectParams = []RequiredParameters{
	{"int": {Type: parameter.Int, DefaultValue: "42"}},
	{"float": {Type: parameter.Float, DefaultValue: "42"}},
	{"string": {Type: parameter.String, DefaultValue: 4.2}},
}

func Test_Ready(t *testing.T) {
	type testCases struct {
		parameters parameter.Values
		ready      bool
	}

	config := Config{
		Name: "test strategy",
		RequiredParameters: RequiredParameters{
			"int":    {Type: parameter.Int},
			"float":  {Type: parameter.Float},
			"string": {Type: parameter.String},
		},
		ResultModifiers: []ResultModifier{
			{
				RequiredParameters: RequiredParameters{
					"modifier_int": {Type: parameter.Int},
				},
			},
		},
	}
	cases := []testCases{
		{
			parameter.Values{
				Ints: map[parameter.Name]int{
					"int":          1,
					"modifier_int": 1,
				},
				Floats: map[parameter.Name]float64{
					"float": 1.,
				},
				Strings: map[parameter.Name]string{
					"string": "some value",
				},
			},
			true,
		},
		{
			parameter.Values{
				Ints: map[parameter.Name]int{
					"int":          1,
					"modifier_int": 1,
				},
				Floats: map[parameter.Name]float64{
					"float": 1.,
				},
				Strings: map[parameter.Name]string{
					"other-string": "some value",
				},
			},
			false,
		},
		{
			parameter.Values{
				Ints: map[parameter.Name]int{
					"int": 1,
				},
				Floats: map[parameter.Name]float64{
					"float": 1.,
				},
				Strings: map[parameter.Name]string{
					"string": "some value",
				},
			},
			false,
		},
	}
	for _, testCase := range cases {
		if Ready(config, testCase.parameters) != testCase.ready {
			t.Errorf("Expected Ready(config, params) to be '%t'\nParameters: %+v", testCase.ready, testCase.parameters)
		}
	}
}

func TestRequiredParameters_Validate(t *testing.T) {
	for _, params := range correctParams {
		err := params.Validate()
		require.NoError(t, err)
	}
	for _, params := range incorrectParams {
		err := params.Validate()
		require.Errorf(t, err, "%+v", params)
	}
}

func TestConfig_Validate(t *testing.T) {
	type testCases struct {
		params         RequiredParameters
		modifierParams RequiredParameters
		fail           bool
	}
	makeConfig := func(params RequiredParameters, modifierParams RequiredParameters) Config {
		return Config{
			Name:               "test strategy",
			RequiredParameters: params,
			ResultModifiers: []ResultModifier{
				{
					RequiredParameters: modifierParams,
					Execute: func(app scalable.App, values parameter.Values, result Result) (Result, error) {
						panic("test implementation")
					},
				},
			},
			Execute: func(app scalable.App, values parameter.Values) (Result, error) {
				panic("test implementation")
			},
		}
	}
	cases := []testCases{
		{
			params:         RequiredParameters{},
			modifierParams: RequiredParameters{},
			fail:           false,
		},
		{
			params:         correctParams[0],
			modifierParams: RequiredParameters{},
			fail:           false,
		},
		{
			params:         RequiredParameters{},
			modifierParams: correctParams[1],
			fail:           false,
		},
		{
			params:         correctParams[0],
			modifierParams: correctParams[1],
			fail:           false,
		},
		{
			params:         incorrectParams[0],
			modifierParams: RequiredParameters{},
			fail:           true,
		},
		{
			params:         RequiredParameters{},
			modifierParams: incorrectParams[0],
			fail:           true,
		},
		{
			params:         correctParams[0],
			modifierParams: incorrectParams[0],
			fail:           true,
		},
		{
			params:         incorrectParams[0],
			modifierParams: correctParams[0],
			fail:           true,
		},
		{
			params:         incorrectParams[0],
			modifierParams: incorrectParams[1],
			fail:           true,
		},
	}
	for _, testCase := range cases {
		cfg := makeConfig(testCase.params, testCase.modifierParams)
		err := cfg.Validate()
		if testCase.fail {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func Test_Execute(t *testing.T) {
	cfg := Config{
		Name:     "test-config",
		YAMLName: "test-config",
		RequiredParameters: RequiredParameters{
			"int": {Type: parameter.Int},
		},
		ResultModifiers: []ResultModifier{
			{
				RequiredParameters: RequiredParameters{
					"float": {Type: parameter.Float},
				},
				Execute: func(app scalable.App, params parameter.Values, prev Result) (Result, error) {
					return Result{RequiredReplicas: int(math.Ceil(float64(prev.RequiredReplicas) / params.Floats["float"]))}, nil
				},
			},
		},
		Execute: func(app scalable.App, params parameter.Values) (Result, error) {
			return Result{RequiredReplicas: params.Ints["int"]}, nil
		},
	}
	result, err := Execute(cfg, scalable.App{}, parameter.Values{
		Ints: map[parameter.Name]int{
			"int": 42,
		},
		Floats: map[parameter.Name]float64{
			"float": 4.2,
		},
	})
	require.NoError(t, err)
	require.Equal(t, 10, result.RequiredReplicas)
}
