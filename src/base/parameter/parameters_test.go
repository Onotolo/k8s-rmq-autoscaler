package parameter

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestValues_Merge(t *testing.T) {
	pv0, pv1 :=
		Values{
			Ints:      map[Name]int{"int_0": 0},
			Floats:    map[Name]float64{"float_0": 0.},
			Strings:   map[Name]string{"string_0": "0"},
			Booleans:  map[Name]bool{"bool_true": true},
			Durations: map[Name]time.Duration{"duration_0": time.Duration(42)},
		},
		Values{
			Ints:      map[Name]int{"int_1": 1},
			Floats:    map[Name]float64{"float_1": 1.},
			Strings:   map[Name]string{"string_1": "1"},
			Booleans:  map[Name]bool{"bool_false": false},
			Durations: map[Name]time.Duration{"duration_1": time.Duration(420)},
		}
	merged := pv0.Merge(pv1)
	require.Equal(t,
		Values{
			Ints:      map[Name]int{"int_0": 0, "int_1": 1},
			Floats:    map[Name]float64{"float_0": 0., "float_1": 1.},
			Strings:   map[Name]string{"string_0": "0", "string_1": "1"},
			Booleans:  map[Name]bool{"bool_true": true, "bool_false": false},
			Durations: map[Name]time.Duration{"duration_0": time.Duration(42), "duration_1": time.Duration(420)},
		},
		merged,
		"Expected values to be properly merged",
	)
}

func TestValues_Insert(t *testing.T) {
	v := Values{
		Ints:      map[Name]int{"int_0": 0},
		Floats:    map[Name]float64{"float_0": 0},
		Strings:   map[Name]string{"string_0": "0"},
		Booleans:  map[Name]bool{"bool_true": true},
		Durations: map[Name]time.Duration{"duration_0": time.Duration(42)},
	}
	cases := []struct {
		name Name
		v    interface{}
		t    Type
	}{
		{"int_1", 1, Int},
		{"float_1", 1., Float},
		{"string_1", "1", String},
		{"bool_false", false, Bool},
		{"duration_1", time.Duration(420), Duration},
	}
	for _, c := range cases {
		err := v.Insert(c.name, c.v, c.t)
		require.NoError(t, err)
	}
	require.Equal(t,
		Values{
			Ints:      map[Name]int{"int_0": 0, "int_1": 1},
			Floats:    map[Name]float64{"float_0": 0., "float_1": 1.},
			Strings:   map[Name]string{"string_0": "0", "string_1": "1"},
			Booleans:  map[Name]bool{"bool_true": true, "bool_false": false},
			Durations: map[Name]time.Duration{"duration_0": time.Duration(42), "duration_1": time.Duration(420)},
		},
		v,
	)
}

func TestStrConversions(t *testing.T) {

	type testCase struct {
		t        Type
		s        string
		expected interface{}
		fail     bool
	}
	testCases := []testCase{
		{t: Int, s: "43242", expected: 43242},
		{t: Int, s: "-123321", expected: -123321},
		{t: Int, s: "0", expected: 0},
		{t: Int, s: "str", fail: true},
		{t: Int, s: "_", fail: true},
		{t: Int, s: "123.321", fail: true},
		{t: Float, s: "0", expected: 0.},
		{t: Float, s: "4.2", expected: 4.2},
		{t: Float, s: "-4.2", expected: -4.2},
		{t: Float, s: "1234", expected: 1234.},
		{t: Float, s: "-1234", expected: -1234.},
		{t: Float, s: "str", fail: true},
		{t: String, s: "str", expected: "str"},
	}

	trueStrings := []string{"1", "t", "T", "TRUE", "true", "True"}
	falseStrings := []string{"0", "f", "F", "FALSE", "false", "False"}

	for _, str := range trueStrings {
		testCases = append(testCases, testCase{
			t:        Bool,
			s:        str,
			expected: true,
		})
	}
	for _, str := range falseStrings {
		testCases = append(testCases, testCase{
			t:        Bool,
			s:        str,
			expected: false,
		})
	}
	for _, tc := range testCases {
		v, err := tc.t.StrConv(tc.s)
		msgPrefix := fmt.Sprintf("Test case: type=%s, str='%s'", tc.t.Name, tc.s)
		if tc.fail {
			require.Error(t, err, msgPrefix)
			continue
		}
		require.NoError(t, err)
		require.Equalf(t, tc.t.ReflectType, reflect.ValueOf(v).Type(), "%s, %v, %v", msgPrefix, tc.t.Name, v)
		require.Equal(t, tc.expected, v, msgPrefix)
	}
}
