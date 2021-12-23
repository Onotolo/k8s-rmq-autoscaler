package executor

import (
	"errors"
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/parameter"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/provider"
	"github.com/medal-labs/k8s-rmq-autoscaler/base/strategy"
	"github.com/medal-labs/k8s-rmq-autoscaler/common"
	"reflect"
)

func convertParameters(p provider.ProvidedParameters, config strategy.RequiredParameters) (parameters parameter.Values, err error) {
	defer common.HandlePanics(&err, func(_ interface{}) error {
		return errors.New("could not convert provided parameters to specified types")
	})
	parameters = parameter.EmptyValues()
	conversionError := func(param parameter.Name, expectedType parameter.Type, value interface{}) error {
		return fmt.Errorf(
			"could not convert value '%v' provided for '%s' with expected type %s",
			value, param, expectedType.Name,
		)
	}
	for param, value := range p {
		refValue := reflect.ValueOf(value)
		specType := config[param].Type
		paramType := specType.ReflectType
		mapOfType, err := parameters.MapValueOfType(specType)
		if err != nil {
			return parameter.Values{}, fmt.Errorf("%s:%w", conversionError(param, specType, value), err)
		}
		if !refValue.Type().ConvertibleTo(paramType) {
			return parameter.Values{}, conversionError(param, specType, value)
		}
		mapOfType.SetMapIndex(reflect.ValueOf(param), refValue.Convert(paramType))
	}
	return parameters, nil
}
