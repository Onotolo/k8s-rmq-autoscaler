package parameter

import (
	"reflect"
	"strconv"
	"time"
)

type Name string

type Type struct {
	Name        string
	ReflectType reflect.Type
	StrConv     func(string) (interface{}, error)
}

func (t Type) EqualTo(other Type) bool {
	return t.Name == other.Name && t.ReflectType == other.ReflectType
}

var (
	Int = Type{
		Name:        "int",
		ReflectType: reflect.TypeOf(1),
		StrConv: func(s string) (interface{}, error) {
			v, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			return v, nil
		},
	}
	Float = Type{
		Name:        "float",
		ReflectType: reflect.TypeOf(1.),
		StrConv: func(s string) (interface{}, error) {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, err
			}
			return v, nil
		},
	}
	String = Type{
		Name:        "string",
		ReflectType: reflect.TypeOf(""),
		StrConv: func(s string) (interface{}, error) {
			return s, nil
		},
	}
	Bool = Type{
		Name:        "bool",
		ReflectType: reflect.TypeOf(false),
		StrConv: func(s string) (interface{}, error) {
			v, err := strconv.ParseBool(s)
			if err != nil {
				return nil, err
			}
			return v, nil
		},
	}
	Duration = Type{
		Name:        "duration",
		ReflectType: reflect.TypeOf(time.Duration(1)),
		StrConv: func(s string) (interface{}, error) {
			v, err := time.ParseDuration(s)
			if err != nil {
				return nil, err
			}
			return v, nil
		},
	}
)
