package modifiers

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"k8s.io/klog"
)

var ScalingThreshold = strategy.ResultModifier{
	Name: "scaling-threshold",
	RequiredParameters: map[parameter.Name]strategy.ParameterSpec{
		parameters.ScalingThreshold: {Type: parameter.Int, DefaultValue: -1},
		parameters.QueueLength:      {Type: parameter.Int},
	},
	Execute: func(app scalable.App, params parameter.Values, prev strategy.Result) (strategy.Result, error) {
		threshold := params.Ints[parameters.ScalingThreshold]
		queueLen := params.Ints[parameters.QueueLength]

		switch {
		case threshold < 0:
			return prev, nil
		case prev.Skip:
			return prev, nil
		case prev.RequiredReplicas < app.Replicas:
			return prev, nil
		case queueLen > threshold:
			return prev, nil
		}
		if klog.V(2) {
			klog.Infof(
				"Skipping %s upscaling as its queue contains less messages than configured threshold. Messages: %d, scaling-threshold: %d",
				app.Name, queueLen, threshold,
			)
		}
		return strategy.Result{Skip: true}, nil
	},
}