package modifiers

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"k8s.io/klog"
)

var SafeUnscale = strategy.ResultModifier{
	Name: "safe-unscale",
	RequiredParameters: map[parameter.Name]strategy.ParameterSpec{
		parameters.QueueLength: {Type: parameter.Int},
		parameters.SafeUnscale: {Type: parameter.Bool, DefaultValue: true},
	},
	Execute: func(app scalable.App, params parameter.Values, prev strategy.Result) (strategy.Result, error) {
		safeUnscale := params.Booleans[parameters.SafeUnscale]
		queueLen := params.Ints[parameters.QueueLength]

		switch {
		case prev.Skip || !safeUnscale:
			return prev, nil
		case prev.RequiredReplicas < app.Replicas && queueLen < 0:
			if klog.V(2) {
				klog.Infof(
					"Skipping %s downscaling as its queue contains messages and safe unscale is enabled",
					app.Name,
				)
			}
			return strategy.Result{Skip: true}, nil
		}
		return prev, nil
	},
}
