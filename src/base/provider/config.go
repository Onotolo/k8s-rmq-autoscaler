package provider

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
)

type Name string

type Config struct {
	Name                Name
	AvailableParameters map[parameter.Name]parameter.Type
	Provide             func(appsCtx map[scalable.App]AppContext)
}

func Launch(config Config, params map[scalable.App][]parameter.Name) map[scalable.App]ResultAppContext {

	appsCtx := map[scalable.App]AppContext{}
	resultAppsCtx := map[scalable.App]ResultAppContext{}

	for app, paramsNames := range params {
		appsCtx[app], resultAppsCtx[app] = newConnectedContexts(app, config, paramsNames)
	}
	go config.Provide(appsCtx)
	return resultAppsCtx
}
