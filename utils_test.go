package subcommands

import (
	"errors"
	"io"
	"maps"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type arrayFlag [1]string

func (a arrayFlag) String() string { return a[0] }

func (a *arrayFlag) Set(v string) error {
	if v == "BAD" {
		return io.EOF
	}

	a[0] = v

	return nil
}

type sliceFlag []string

func (s sliceFlag) String() string {
	if len(s) > 0 {
		return s[0]
	}

	return ""
}

func (s *sliceFlag) Set(v string) error {
	*s = append(*s, v)

	return nil
}

type mapFlag map[string]struct{}

func (m mapFlag) String() string {
	return strings.Join(slices.Collect(maps.Keys(m)), ", ")
}

func (m mapFlag) Set(v string) error {
	m[v] = struct{}{}

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
			reflect.TypeFor[*arrayFlag](),
			[]string{"next"},
			reflect.ValueOf(&arrayFlag{"next"}),
			true,
			nil,
		},
		{ // 6
			reflect.TypeFor[*arrayFlag](),
			[]string{"BAD"},
			reflect.Value{},
			false,
			io.EOF,
		},
		{ // 5
			reflect.TypeFor[sliceFlag](),
			[]string{"next"},
			reflect.ValueOf(sliceFlag{"next"}),
			true,
			nil,
		},
		{ // 7
			reflect.TypeFor[mapFlag](),
			[]string{"next"},
			reflect.ValueOf(mapFlag{"next": struct{}{}}),
			true,
			nil,
		},
		{ // 8
			reflect.TypeFor[*arrayFlag](),
			[]string{"BAD"},
			reflect.Value{},
			false,
			io.EOF,
		},
		{ // 9
			reflect.TypeFor[arrayFlag](),
			[]string{"next"},
			reflect.ValueOf(arrayFlag{"next"}),
			true,
			nil,
		},
		{ // 10
			reflect.TypeFor[arrayFlag](),
			[]string{"BAD"},
			reflect.Value{},
			false,
			io.EOF,
		},
		{ // 11
			reflect.TypeFor[string](),
			nil,
			reflect.Value{},
			false,
			ErrMissingValue,
		},
		{ // 12
			reflect.TypeFor[string](),
			[]string{"next"},
			reflect.ValueOf("next"),
			true,
			nil,
		},
		{ // 13
			reflect.TypeFor[chan string](),
			[]string{"next"},
			reflect.Value{},
			false,
			ErrInvalidType,
		},
		{ // 14
			reflect.TypeFor[*string](),
			[]string{"next"},
			reflect.ValueOf(new("next")),
			true,
			nil,
		},
		{ // 15
			reflect.TypeFor[*chan string](),
			[]string{"next"},
			reflect.Value{},
			false,
			ErrInvalidType,
		},
		{ // 16
			reflect.TypeFor[[]string](),
			[]string{"next"},
			reflect.ValueOf([]string{"next"}),
			true,
			nil,
		},
		{ // 16
			reflect.TypeFor[[]chan string](),
			[]string{"next"},
			reflect.Value{},
			false,
			ErrInvalidType,
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
