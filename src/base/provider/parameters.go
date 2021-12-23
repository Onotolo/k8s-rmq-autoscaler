package provider

import (
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
)

type RequiredAppsParameters map[scalable.App][]parameter.Name

type Result struct {
	Parameters ProvidedParameters
	Error      error
}

type ProvidedParameters map[parameter.Name]interface{}

func (p ProvidedParameters) Set(paramName parameter.Name, value interface{}) {
	p[paramName] = value
}
