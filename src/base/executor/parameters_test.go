package executor

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_ConvertParameters(t *testing.T) {
	provided := provider.ProvidedParameters{
		"int":        1,
		"float":      1.,
		"string":     "120",
		"bool_true":  true,
		"bool_false": false,
		"duration":   time.Duration(2520000000000),
	}
	config := strategy.RequiredParameters{
		"int":        {Type: parameter.Int},
		"float":      {Type: parameter.Float},
		"string":     {Type: parameter.String},
		"bool_true":  {Type: parameter.Bool},
		"bool_false": {Type: parameter.Bool},
		"duration":   {Type: parameter.Duration},
	}
	converted, err := convertParameters(provided, config)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t,
		parameter.Values{
			Ints:      map[parameter.Name]int{"int": 1},
			Floats:    map[parameter.Name]float64{"float": 1.},
			Strings:   map[parameter.Name]string{"string": "120"},
			Booleans:  map[parameter.Name]bool{"bool_true": true, "bool_false": false},
			Durations: map[parameter.Name]time.Duration{"duration": time.Duration(2520000000000)},
		},
		converted,
		"Expected provided parameters to be properly converted",
	)
}
