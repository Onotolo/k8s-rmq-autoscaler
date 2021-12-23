package modifiers

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"k8s.io/klog"
	"time"
)

var Cooldown = strategy.ResultModifier{
	Name: "cooldown-delay",
	RequiredParameters: map[parameter.Name]strategy.ParameterSpec{
		parameters.CooldownDelay: {Type: parameter.Duration, DefaultValue: time.Duration(0)},
	},
	Execute: func(app scalable.App, params parameter.Values, prev strategy.Result) (strategy.Result, error) {
		delay := params.Durations[parameters.CooldownDelay]

		if prev.Skip || delay <= 0 || time.Now().Sub(app.UpdatedDate) > delay {
			return prev, nil
		}
		if klog.V(2) {
			klog.Infof("%s is cooled down, waiting more (date %s, duration %s)", app.Name, app.UpdatedDate, delay)
		}
		return strategy.Result{Skip: true}, nil
	},
}
