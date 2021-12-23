package providers

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/providers/rmqhttp"
)

type Config struct {
	RMQHTTP rmqhttp.Config
}

func Configure(config Config) []provider.Config {
	return []provider.Config{
		rmqhttp.ProviderConfig(config.RMQHTTP),
	}
}
