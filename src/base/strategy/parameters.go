package strategy

import (
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"reflect"
)

type RequiredParameters map[parameter.Name]ParameterSpec

type ParameterSpec struct {
	DefaultValue interface{}
	Type         parameter.Type
}

func (p RequiredParameters) Validate() error {
	for param, spec := range p {
		if err := spec.Validate(); err != nil {
			return fmt.Errorf("validation failed for '%s' parameter with '%+v' spec: %w", param, spec, err)
		}
	}
	return nil
}

func (cfg ParameterSpec) Validate() error {
	if cfg.DefaultValue == nil {
		return nil
	}
	defaultValue := cfg.DefaultValue
	if reflect.ValueOf(defaultValue).Type().ConvertibleTo(cfg.Type.ReflectType) {
		return nil
	}
	return fmt.Errorf(
		"could not convert default value '%v' to specified type %s",
		defaultValue, cfg.Type.Name,
	)
}
