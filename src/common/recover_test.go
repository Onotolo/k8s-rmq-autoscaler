package common

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandlePanics(t *testing.T) {
	expected := func() (err error) {
		defer HandlePanics(&err,
			func(rec interface{}) error {
				return fmt.Errorf("failed with panic: %v", rec)
			},
		)
		panic("something wrong, I can feel it")
	}()
	require.Error(t, expected)
}

func ExampleHandlePanics() {
	handled := func() (err error) {
		defer HandlePanics(&err,
			func(rec interface{}) error {
				return fmt.Errorf("failed with panic: %v", rec)
			},
		)
		panic("something wrong, I can feel it")
	}()
	fmt.Println(handled)
	// Output: failed with panic: something wrong, I can feel it
}
