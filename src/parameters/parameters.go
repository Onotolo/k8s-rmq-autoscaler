package parameters

import "github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"

const (
	QueueLength       parameter.Name = "queue-length"
	MessagesPerWorker                = "messages-per-worker"
	Offset                           = "offset"
	Steps                            = "steps"
	CooldownDelay                    = "cooldown-delay"
	Min                              = "min-workers"
	Max                              = "max-workers"
	Override                         = "override"
	SafeUnscale                      = "safe-unscale"
	ScalingThreshold                 = "scaling-threshold"
)
