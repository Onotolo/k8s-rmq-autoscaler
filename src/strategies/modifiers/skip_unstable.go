package modifiers

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"k8s.io/klog"
)

var SkipUnstable = strategy.ResultModifier{
	Name:               "skip-unstable",
	RequiredParameters: strategy.RequiredParameters{},
	Execute: func(app scalable.App, params parameter.Values, prev strategy.Result) (strategy.Result, error) {
		if app.ReadyReplicas != app.Replicas {
			if klog.V(2) {
				klog.Infof(
					"%s is unstable, skipping scaling",
					app.Name,
				)
			}
			return strategy.Result{Skip: true}, nil
		}
		return prev, nil
	},
}
