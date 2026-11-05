package subcommands

import (
	"errors"
	"io"
	"reflect"
	"testing"
)

type sliceFlag [1]string

func (s sliceFlag) String() string { return s[0] }

func (s *sliceFlag) Set(v string) error {
	if v == "BAD" {
		return io.EOF
	}

	s[0] = v

	return nil
}

func TestSetValue(t *testing.T) {
	for n, test := range [...]struct {
		Type     reflect.Type
		Args     []string
		Value    reflect.Value
		UsedNext bool
		Err      error
	}{
		{ // 1
			reflect.TypeFor[bool](),
			nil,
			reflect.ValueOf(true),
			false,
			nil,
		},
		{ // 2
			reflect.TypeFor[bool](),
			[]string{"unknown"},
			reflect.ValueOf(true),
			false,
			nil,
		},
		{ // 3
			reflect.TypeFor[bool](),
			[]string{"True"},
			reflect.ValueOf(true),
			true,
			nil,
		},
		{ // 4
			reflect.TypeFor[bool](),
			[]string{"0"},
			reflect.ValueOf(false),
			true,
			nil,
		},
		{ // 5
			reflect.TypeFor[*sliceFlag](),
			[]string{"next"},
			reflect.ValueOf(&sliceFlag{"next"}),
			true,
			nil,
		},
		{ // 6
			reflect.TypeFor[*sliceFlag](),
			[]string{"BAD"},
			reflect.Value{},
			false,
			io.EOF,
		},
		{ // 7
			reflect.TypeFor[sliceFlag](),
			[]string{"next"},
			reflect.ValueOf(sliceFlag{"next"}),
			true,
			nil,
		},
		{ // 8
			reflect.TypeFor[sliceFlag](),
			[]string{"BAD"},
			reflect.Value{},
			false,
			io.EOF,
		},
		{ // 9
			reflect.TypeFor[string](),
			nil,
			reflect.Value{},
			false,
			ErrMissingValue,
		},
		{ // 10
			reflect.TypeFor[string](),
			[]string{"next"},
			reflect.ValueOf("next"),
			true,
			nil,
		},
	} {
		if v, used, err := setValue(test.Type, test.Args); !errors.Is(err, test.Err) {
			t.Errorf("test %d: expecting error %v, got %v", n+1, test.Err, err)
		} else if v.IsValid() != test.Value.IsValid() || v.IsValid() && !reflect.DeepEqual(v.Interface(), test.Value.Interface()) {
			t.Errorf("test %d: expecting value %v, got %v", n+1, test.Value, v)
		} else if used != test.UsedNext {
			t.Errorf("test %d: expecting next used to be %v, got %v", n+1, test.UsedNext, used)
		}
	}
}
