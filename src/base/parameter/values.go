package parameter

import (
	"fmt"
	"github.com/medal-labs/k8s-rmq-autoscaler/common"
	"reflect"
	"time"
)

type Values struct {
	Ints      map[Name]int
	Floats    map[Name]float64
	Strings   map[Name]string
	Booleans  map[Name]bool
	Durations map[Name]time.Duration
}

func EmptyValues() Values {
	return Values{
		Ints:      map[Name]int{},
		Floats:    map[Name]float64{},
		Strings:   map[Name]string{},
		Booleans:  map[Name]bool{},
		Durations: map[Name]time.Duration{},
	}
}

func (p Values) Contains(name Name, t Type) bool {
	m, err := p.MapValueOfType(t)
	if err != nil {
		panic(err)
	}
	return m.MapIndex(reflect.ValueOf(name)).IsValid()
}

func (p Values) Insert(name Name, v interface{}, t Type) (err error) {
	defer common.HandlePanics(&err, func(rec interface{}) error {
		return fmt.Errorf("could not insert value %v", v)
	})
	m, err := p.MapValueOfType(t)
	if err != nil {
		return err
	}
	m.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(v).Convert(t.ReflectType))
	return nil
}

func (p Values) Merge(other Values) Values {
	newParams := EmptyValues()

	newParamsRef := reflect.ValueOf(newParams)
	oldParamsRef := reflect.ValueOf(p)
	otherParamsRef := reflect.ValueOf(other)

	var pn Name
	paramNameType := reflect.ValueOf(pn).Type()
	for i := 0; i < newParamsRef.NumField(); i++ {
		field := newParamsRef.Field(i)
		fieldType := field.Type()
		if field.Kind() != reflect.Map || fieldType.Key() != paramNameType {
			continue
		}
		oldField := oldParamsRef.Field(i)
		oldRange := oldField.MapRange()
		if oldField.Len() > 0 {
			for oldRange.Next() {
				field.SetMapIndex(oldRange.Key(), oldRange.Value())
			}
		}
		otherField := otherParamsRef.Field(i)
		otherRange := otherField.MapRange()
		if otherField.Len() > 0 {
			for otherRange.Next() {
				field.SetMapIndex(otherRange.Key(), otherRange.Value())
			}
		}
	}
	return newParams
}

func (p Values) MapValueOfType(t Type) (reflect.Value, error) {
	v := reflect.ValueOf(p)
	var pn Name
	paramNameType := reflect.ValueOf(pn).Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := field.Type()
		if field.Kind() != reflect.Map || fieldType.Key() != paramNameType {
			continue
		}
		if fieldType.Elem() == t.ReflectType {
			return field, nil
		}
	}
	return reflect.Value{}, fmt.Errorf(
		"field with map with a proper type not specified for '%s' parameter type", t.Name,
	)
}

func (p Values) Len() int {
	total := 0
	v := reflect.ValueOf(p)
	var pn Name
	paramNameType := reflect.ValueOf(pn).Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := field.Type()
		if field.Kind() != reflect.Map || fieldType.Key() != paramNameType {
			continue
		}
		total += field.Len()
	}
	return total
}
