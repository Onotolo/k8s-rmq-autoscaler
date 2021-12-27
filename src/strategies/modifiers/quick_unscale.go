package modifiers

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"k8s.io/klog"
	"time"
)

var QuickUnscale = strategy.ResultModifier{
	Name: "quick-unscale",
	RequiredParameters: map[parameter.Name]strategy.ParameterSpec{
		parameters.ScaleToZeroIn: {Type: parameter.Duration, DefaultValue: time.Duration(-1)},
		parameters.ScaleToMinIn:  {Type: parameter.Duration, DefaultValue: time.Duration(-1)},
		parameters.Min:           {Type: parameter.Int, DefaultValue: -1},
		parameters.QueueLength:   {Type: parameter.Int},
	},
	Execute: func() func(scalable.App, parameter.Values, strategy.Result) (strategy.Result, error) {

		withoutMessagesFrom := map[string]time.Time{}

		return func(app scalable.App, params parameter.Values, prev strategy.Result) (strategy.Result, error) {

			scaleToZeroIn := params.Durations[parameters.ScaleToZeroIn]
			scaleToMinIn := params.Durations[parameters.ScaleToMinIn]
			min := params.Ints[parameters.Min]

			if params.Ints[parameters.QueueLength] > 0 {
				delete(withoutMessagesFrom, app.Key)
				return prev, nil
			}
			appWithoutMessagesFrom, ok := withoutMessagesFrom[app.Key]

			switch {
			case !ok:
				withoutMessagesFrom[app.Key] = time.Now()
				return prev, nil
			case scaleToZeroIn > 0 && time.Now().Sub(appWithoutMessagesFrom) >= scaleToZeroIn:
				if klog.V(2) {
					klog.Infof(
						"%s's queue has been empty since %s, which is longer than configured 'scale-to-zero-in' duration: %s. Scaling to zero",
						app.Name, appWithoutMessagesFrom, scaleToZeroIn,
					)
				}
				return strategy.Result{RequiredReplicas: 0}, nil
			case scaleToMinIn > 0 && time.Now().Sub(appWithoutMessagesFrom) >= scaleToMinIn:
				if min < 0 {
					klog.Errorf(
						"%s's has its scale-to-min duration set without setting min workers, setting min to be 0",
						app.Name,
					)
					min = 0
				}
				if prev.RequiredReplicas <= min || app.Replicas <= min && prev.Skip {
					return prev, nil
				}
				if klog.V(2) {
					klog.Infof(
						"%s's queue has been empty since %s, which is longer than configured 'scale-to-min-in' duration: %s. Scaling to %d workers",
						app.Name, appWithoutMessagesFrom, scaleToMinIn, min,
					)
				}
				return strategy.Result{RequiredReplicas: min}, nil
			}
			return prev, nil
		}
	}(),
}
