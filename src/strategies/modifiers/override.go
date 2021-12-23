package modifiers

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"k8s.io/klog"
)

var OverrideLimits = strategy.ResultModifier{
	Name: "override-limits",
	RequiredParameters: map[parameter.Name]strategy.ParameterSpec{
		parameters.Override: {Type: parameter.Bool, DefaultValue: false},
		parameters.Min:      {Type: parameter.Int},
		parameters.Max:      {Type: parameter.Int},
	},
	Execute: func(app scalable.App, params parameter.Values, prev strategy.Result) (strategy.Result, error) {
		repl := app.Replicas
		override := params.Booleans[parameters.Override]
		max, min := params.Ints[parameters.Max], params.Ints[parameters.Min]

		switch {
		case prev.Skip:
			return prev, nil
		case override && (repl > max || repl < min):
			if klog.V(2) {
				klog.Infof("%s limits are override, do nothing", app.Key)
			}
			return strategy.Result{Skip: true}, nil
		}
		return prev, nil
	},
}
