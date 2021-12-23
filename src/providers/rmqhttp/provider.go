package rmqhttp

import (
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
	"github.com/medal-labs/k8s-rmq-autoscaler/common"
	"github.com/medal-labs/k8s-rmq-autoscaler/parameters"
	"net/http"
)

func ProviderConfig(config Config) provider.Config {
	client := rmqHTTPClient{
		Client: &http.Client{},
		config: config,
	}
	return provider.Config{
		Name: config.Name,
		AvailableParameters: map[parameter.Name]parameter.Type{
			parameters.QueueLength: parameter.Int,
		},
		Provide: func(appsCtx map[scalable.App]provider.AppContext) {
			for app, ctx := range appsCtx {
				go func(app scalable.App, ctx provider.AppContext) {
					if ctx.IsCanceled() {
						return
					}
					var appConfig AppConfig
					if err := app.ParseAnnotations(&appConfig, common.AnnotationPrefix); err != nil {
						err = fmt.Errorf("failed to parse annotations: %w", err)
						ctx.Error(err)
						return
					}
					info, err := client.getQueueInfo(appConfig.QueueName, appConfig.Vhost)
					if err != nil {
						err = fmt.Errorf("failed to get queue info: %w", err)
						ctx.Error(err)
						return
					}
					params := provider.ProvidedParameters{}
					for _, param := range ctx.Parameters {
						switch param {
						case parameters.QueueLength:
							params.Set(parameters.QueueLength, info.Messages)
						}
					}
					ctx.PutResult(params)
					ctx.Finish()
				}(app, ctx)
			}
		},
	}
}
