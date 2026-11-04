package subcommands

import (
	"errors"
	"reflect"
	"testing"
)

func TestSetValue(t *testing.T) {
	for n, test := range [...]struct {
		Type     reflect.Type
		Args     []string
		Value    reflect.Value
		UsedNext bool
		Err      error
	}{
		{
			reflect.TypeFor[bool](),
			nil,
			reflect.ValueOf(true),
			false,
			nil,
		},
		{
			reflect.TypeFor[bool](),
			[]string{"unknown"},
			reflect.ValueOf(true),
			false,
			nil,
		},
		{
			reflect.TypeFor[bool](),
			[]string{"True"},
			reflect.ValueOf(true),
			true,
			nil,
		},
		{
			reflect.TypeFor[bool](),
			[]string{"0"},
			reflect.ValueOf(false),
			true,
			nil,
		},
	} {
		if v, used, err := setValue(test.Type, test.Args); !errors.Is(err, test.Err) {
			t.Errorf("test %d: expecting error %v, got %v", n+1, test.Err, err)
		} else if !v.Equal(test.Value) {
			t.Errorf("test %d: expecting value %v, got %v", n+1, test.Value, v)
		} else if used != test.UsedNext {
			t.Errorf("test %d: expecting next used to be %v, got %v", n+1, test.UsedNext, used)
		}
	}
}
