package strategies

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"github.com/medal-labs/k8s-rmq-autoscaler/strategies/modifiers"
	"k8s.io/klog"
	"math"
)

var SimpleQueueBased = strategy.Config{
	Name:     "simple-queue-based",
	YAMLName: "simple-queue-based",
	RequiredParameters: strategy.RequiredParameters{
		parameters.MessagesPerWorker: {Type: parameter.Int, DefaultValue: 1},
		parameters.Offset:            {Type: parameter.Int, DefaultValue: 2},
		parameters.QueueLength:       {Type: parameter.Int},
	},
	ResultModifiers: []strategy.ResultModifier{
		modifiers.WithSteps,
		modifiers.MinMax,
		modifiers.QuickUnscale,
		modifiers.SafeUnscale,
		modifiers.OverrideLimits,
		modifiers.Cooldown,
		modifiers.SkipUnstable,
	},
	Execute: func(app scalable.App, params parameter.Values) (strategy.Result, error) {

		queueLen, messagesPerWorker := float64(params.Ints[parameters.QueueLength]), float64(params.Ints[parameters.MessagesPerWorker])
		offset := params.Ints[parameters.Offset]

		reqRepl := int(math.Ceil(queueLen/messagesPerWorker)) + offset

		if reqRepl == app.Replicas {
			if klog.V(2) {
				klog.Infof(
					"%s's required replicas number is equal to its current replicas (%d), skipping scaling",
					app.Name, app.Replicas,
				)
			}
			return strategy.Result{Skip: true}, nil
		}
		if klog.V(2) {
			klog.Infof(
				"%s's required replicas number will be changed to %d. "+
					"Parameters: current replicas - %d, queue length - %d, "+
					"messages per worker - %d, offset - %d",
				app.Name, reqRepl, app.Replicas, int(queueLen), int(messagesPerWorker), offset,
			)
		}
		return strategy.Result{RequiredReplicas: reqRepl}, nil
	},
}
