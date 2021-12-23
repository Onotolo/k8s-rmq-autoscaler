package modifiers

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"k8s.io/klog"
)

var MinMax = strategy.ResultModifier{
	Name: "min-max",
	RequiredParameters: strategy.RequiredParameters{
		parameters.Max: {Type: parameter.Int},
		parameters.Min: {Type: parameter.Int},
	},
	Execute: func(app scalable.App, params parameter.Values, prev strategy.Result) (strategy.Result, error) {
		max, min := params.Ints[parameters.Max], params.Ints[parameters.Min]
		replicas := prev.RequiredReplicas

		switch {
		case prev.Skip:
			return prev, nil
		case replicas > max:
			if klog.V(2) {
				klog.Infof(
					"%s's number of required replicas (%d) exceeds configured maximum (%d)",
					app.Name, replicas, max,
				)
			}
			return strategy.Result{RequiredReplicas: max}, nil
		case replicas < min:
			if klog.V(2) {
				klog.Infof(
					"%s's number of required replicas (%d) is less than configured minimum (%d)",
					app.Name, replicas, min,
				)
			}
			return strategy.Result{RequiredReplicas: min}, nil
		}
		return prev, nil
	},
}
