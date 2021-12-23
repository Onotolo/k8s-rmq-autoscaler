package modifiers

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"k8s.io/klog"
	"math"
)

var WithSteps = strategy.ResultModifier{
	Name: "with-steps",
	RequiredParameters: map[parameter.Name]strategy.ParameterSpec{
		parameters.Steps: {Type: parameter.Int, DefaultValue: 2},
	},
	Execute: func(app scalable.App, params parameter.Values, prev strategy.Result) (strategy.Result, error) {
		maxStep := params.Ints[parameters.Steps]
		scale := prev.RequiredReplicas - app.Replicas

		var reqRepl int

		switch {
		case prev.Skip || absLess(scale, maxStep):
			return prev, nil
		case scale < 0:
			reqRepl = app.Replicas - maxStep
		case scale > 0:
			reqRepl = app.Replicas + maxStep
		}
		if klog.V(2) {
			klog.Infof(
				"%s's required replicas change (%d) exceeds maximum step (%d), modifying required replicas to be (%d)",
				app.Name, scale, maxStep, reqRepl,
			)
		}
		return strategy.Result{RequiredReplicas: reqRepl}, nil
	},
}

func absLess(a, b int) bool {
	return math.Abs(float64(a)) < math.Abs(float64(b))
}
