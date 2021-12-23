package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type TestStruct struct {
	BoolField   bool    `nameTag:"bool_field"`
	IntField    int     `nameTag:"int_field"`
	StringField string  `nameTag:"string_field" default:"stringField default value"`
	FloatField  float64 `nameTag:"float_field" default:"4.25"`
}

func TestParseK8sAnnotations_WithPrefix(t *testing.T) {
	var testStruct TestStruct
	annotations := map[string]string{
		"prefix/bool_field":   "true",
		"prefix/string_field": "stringValue",
		"prefix/int_field":    "20",
		"prefix/float_field":  "0.5",
		"prefix/extra_field":  "false",
		"extra_field":         "some extra field",
	}
	parser := NewFlatMapParser("nameTag", "default")
	err := parser(annotations, &testStruct, "prefix/")
	require.NoError(t, err, "ParseK8sAnnotations should not fail given correct input")
	require.True(t, testStruct.BoolField, "boolField must be true")
	require.Equal(t, "stringValue", testStruct.StringField, "stringField value must be 'stringValue'")
	require.Equal(t, 20, testStruct.IntField, "intField value must be 20")
	require.InDelta(t, 0.5, testStruct.FloatField, 10e-7, "field value must be around 0.5")
}

func TestParseK8sAnnotations_WithoutPrefix(t *testing.T) {
	var testStruct TestStruct
	annotations := map[string]string{
		"bool_field":   "true",
		"string_field": "stringValue",
		"int_field":    "20",
		"float_field":  "0.5",
		"extra_field":  "some extra field",
	}
	parser := NewFlatMapParser("nameTag", "default")
	err := parser(annotations, &testStruct)
	require.NoError(t, err, "ParseK8sAnnotations should not fail given correct input")
	require.True(t, testStruct.BoolField, "boolField must be true")
	require.Equal(t, "stringValue", testStruct.StringField, "stringField value must be 'stringValue'")
	require.Equal(t, 20, testStruct.IntField, "intField value must be 20")
	require.InDelta(t, 0.5, testStruct.FloatField, 10e-7, "field value must be around 0.5")
}

func TestParseK8sAnnotations_WithoutDefaults(t *testing.T) {
	var testStruct TestStruct
	annotations := map[string]string{
		"bool_field":  "true",
		"int_field":   "20",
		"extra_field": "some extra field",
	}
	parser := NewFlatMapParser("nameTag", "default")
	err := parser(annotations, &testStruct)
	require.NoError(t, err, "ParseK8sAnnotations should not fail given correct input")
	require.True(t, testStruct.BoolField, "boolField must be true")
	require.Equal(t, "stringField default value", testStruct.StringField,
		"stringField value must be equal to specified default")
	require.Equal(t, 20, testStruct.IntField, "intField value must be 20")
	require.InDelta(t, 4.25, testStruct.FloatField, 10e-7,
		"field value must be around specified default value")
}

func TestParseK8sAnnotations_WithMissingAnnotation(t *testing.T) {
	var testStruct TestStruct
	annotations := map[string]string{
		"bool_field":   "true",
		"string_field": "stringValue",
		"float_field":  "0.5",
		"extra_field":  "some extra field",
	}
	parser := NewFlatMapParser("nameTag", "default")
	err := parser(annotations, &testStruct)
	require.Error(t, err,
		"ParseK8sAnnotations must fail when annotation with no specified default value is missing")
}

func TestParseK8sAnnotations_WithMalformedBool(t *testing.T) {
	var testStruct TestStruct
	annotations := map[string]string{
		"bool_field":  "tue",
		"int_field":   "20",
		"extra_field": "some extra field",
	}
	parser := NewFlatMapParser("nameTag", "default")
	err := parser(annotations, &testStruct)
	require.Error(t, err,
		"ParseK8sAnnotations must fail when bool annotation value is malformed")
}

func TestParseK8sAnnotations_WithMalformedInt(t *testing.T) {
	var testStruct TestStruct
	annotations := map[string]string{
		"bool_field":  "true",
		"int_field":   "20a",
		"extra_field": "some extra field",
	}
	parser := NewFlatMapParser("nameTag", "default")
	err := parser(annotations, &testStruct)
	require.Error(t, err,
		"ParseK8sAnnotations must fail when int annotation value is malformed")
}

func TestParseK8sAnnotations_WithMalformedFloat(t *testing.T) {
	var testStruct TestStruct
	annotations := map[string]string{
		"bool_field":  "true",
		"int_field":   "20",
		"float_field": "20.a",
		"extra_field": "some extra field",
	}
	parser := NewFlatMapParser("nameTag", "default")
	err := parser(annotations, &testStruct)
	require.Error(t, err,
		"ParseK8sAnnotations must fail when float annotation value is malformed")
}
