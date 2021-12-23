package strategy

import (
	"errors"
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/scalable"
)

type Name string
type YAMLName string

type Config struct {
	Name               Name
	YAMLName           YAMLName
	RequiredParameters RequiredParameters
	ResultModifiers    []ResultModifier
	Execute            func(app scalable.App, params parameter.Values) (Result, error)
}

type Result struct {
	App              scalable.App
	RequiredReplicas int
	Skip             bool
}

type ResultModifier struct {
	Name               string
	RequiredParameters RequiredParameters
	Execute            func(app scalable.App, params parameter.Values, prev Result) (Result, error)
}

func Execute(config Config, app scalable.App, params parameter.Values) (Result, error) {
	result, err := config.Execute(app, params)
	if err != nil {
		return Result{}, err
	}
	for _, modifier := range config.ResultModifiers {
		result, err = modifier.Execute(app, params, result)
		if err != nil {
			return Result{}, err
		}
	}
	result.App = app
	return result, nil
}

func (sc Config) GetRequiredParameters() RequiredParameters {
	req := map[parameter.Name]ParameterSpec{}
	for name, spec := range sc.RequiredParameters {
		req[name] = spec
	}
	for _, modifier := range sc.ResultModifiers {
		for name, spec := range modifier.RequiredParameters {
			req[name] = spec
		}
	}
	return req
}

func Ready(sc Config, parameters parameter.Values) bool {
	for param, spec := range sc.GetRequiredParameters() {
		if !parameters.Contains(param, spec.Type) {
			return false
		}
	}
	return true
}

func (sc Config) Validate() error {
	if err := sc.RequiredParameters.Validate(); err != nil {
		return err
	}
	if sc.Execute == nil {
		return errors.New("required Execute method is not defined")
	}
	for _, modifier := range sc.ResultModifiers {
		if err := modifier.Validate(); err != nil {
			return fmt.Errorf("validation failed for '%s' result modified: %w", modifier.Name, err)
		}
	}
	return nil
}

func (m ResultModifier) Validate() error {
	if err := m.RequiredParameters.Validate(); err != nil {
		return err
	}
	if m.Execute == nil {
		return errors.New("Execute method is not defined")
	}
	return nil
}
