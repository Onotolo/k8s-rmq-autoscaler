package common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	DefaultNameTag   = "k8s-annotation"
	DefaultValueTag  = "default"
	AnnotationPrefix = "k8s-rmq-autoscaler/"
)

type FlatMapParser func(m map[string]string, v interface{}, prefixes ...string) error

var ParseK8sAnnotations = NewFlatMapParser(DefaultNameTag, DefaultValueTag)

func NewFlatMapParser(tag string, defaultValueTag string) FlatMapParser {
	return func(annotations map[string]string, v interface{}, prefixes ...string) error {
		ptrValue := reflect.ValueOf(v)
		if ptrValue.Kind() != reflect.Ptr {
			return fmt.Errorf("expected filled value to be a pointer, got %t", v)
		}
		innerValue := ptrValue.Elem()
		vType := innerValue.Type()
		for i := 0; i < vType.NumField(); i++ {
			fieldInfo := vType.Field(i)
			field := innerValue.Field(i)

			name, ok := fieldInfo.Tag.Lookup(tag)
			if !ok {
				continue
			}
			annotationName := strings.Join(append(prefixes, name), "")
			fieldValue, ok := annotations[annotationName]

			if !ok {
				defaultValue, hasDefaultValue := fieldInfo.Tag.Lookup(defaultValueTag)
				if !hasDefaultValue {
					return fmt.Errorf("'%s' annotation is not specified and default value is not provided", annotationName)
				}
				fieldValue = defaultValue
			}
			if err := setFieldValue(field, fieldValue); err != nil {
				return fmt.Errorf("failed to convert value provided for '%s' annotation: %w", annotationName, err)
			}
		}
		return nil
	}
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		value, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(value))
	case reflect.Float64:
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(value)
	case reflect.Bool:
		value, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(value)
	}
	return nil
}
