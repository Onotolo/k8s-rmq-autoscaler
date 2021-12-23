package main

import (
	"context"
	"flag"
	"github.com/kelseyhightower/envconfig"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/executor"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/loop"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"github.com/medal-labs/k8s-rmq-autoscaler/providers"
	"github.com/medal-labs/k8s-rmq-autoscaler/providers/rmqhttp"
	"github.com/medal-labs/k8s-rmq-autoscaler/strategies"
	"k8s.io/klog"
	"os"
	"regexp"
)

type EnvConfig struct {
	Namespaces      string `envconfig:"NAMESPACES" default:""`
	InCluster       bool   `envconfig:"IN_CLUSTER" default:"false"`
	RMQUrl          string `envconfig:"RMQ_URL" required:"true"`
	RMQUser         string `envconfig:"RMQ_USER" required:"true"`
	RMQPassword     string `envconfig:"RMQ_PASSWORD" required:"true"`
	Tick            int    `envconfig:"TICK" default:"10"`
	LogLevel        string `envconfig:"MDL_COMN_LOGLEVEL" default:"INFO"`
	DefaultStrategy string `envconfig:"K8S_AUTOSCALER_DEFAULT_STRATEGY" default:"simple-queue-based"`
}

func main() {
	ctx := context.Background()

	var cfg EnvConfig

	if err := envconfig.Process("", &cfg); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	configureLogLevel(cfg)
	flag.Parse()

	enabledProviders := providers.Configure(
		providers.Config{
			RMQHTTP: rmqhttp.Config{
				Name:     "rmq-http-provider",
				Url:      cfg.RMQUrl,
				User:     cfg.RMQUser,
				Password: cfg.RMQPassword,
			},
		},
	)
	executorCfg := executor.Config{
		EnabledStrategies: []strategy.Config{
			strategies.SimpleQueueBased,
		},
		EnabledProviders:  enabledProviders,
		AnnotationsPrefix: "k8s-rmq-autoscaler/",
		DefaultStrategy:   strategy.YAMLName(cfg.DefaultStrategy),
		DefaultParametersProviders: map[parameter.Name]provider.Name{
			parameters.QueueLength: "rmq-http-provider",
		},
	}
	errs := executorCfg.Validate()
	if len(errs) > 0 {
		for _, err := range errs {
			klog.Error(err)
			os.Exit(1)
		}
	}

	loopCfg := loop.Config{
		ExecutorCfg:     executorCfg,
		InCluster:       cfg.InCluster,
		Namespaces:      cfg.Namespaces,
		LoopTickSeconds: cfg.Tick,
	}
	err := loop.Launch(ctx, loopCfg)
	if err != nil {
		klog.Error(err)
		os.Exit(128)
	}
	<-ctx.Done()
}

func configureLogLevel(cfg EnvConfig) {
	klog.InitFlags(nil)

	for _, arg := range os.Args[1:] {

		if ok, _ := regexp.MatchString(`^-{1,2}v(=\d+)?$`, arg); ok {
			// Log level configured via command line arguments, no need to do anything
			return
		}
	}
	if cfg.LogLevel == "DEBUG" {
		// Simulate configuration via command line arguments
		_ = flag.Set("v", "2")
	}
}
