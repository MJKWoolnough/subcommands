package subcommands

import (
	"errors"
	"flag"
	"reflect"
	"testing"
	"time"
)

type customFlag struct {
	v string
}

func (c *customFlag) String() string {
	return c.v
}

func (c *customFlag) Set(v string) error {
	c.v = v

	return nil
}

var _ flag.Value = &customFlag{}

func TestFlagsProcess(t *testing.T) {
	type Global struct {
		GlobalFlag string
	}

	type Sub struct {
		Another string
		Global
	}

	aliased := &flagData{
		aliases: []string{"--aaa", "--bbb"},
		typ:     reflect.TypeFor[string](),
		index:   []int{0},
	}

	for n, test := range [...]struct {
		Input  any
		Output map[string]*flagData
		Err    error
	}{
		{ // 1
			Input:  int(0),
			Output: map[string]*flagData{},
			Err:    ErrFlagStructRequired,
		},
		{ // 2
			Input:  struct{}{},
			Output: map[string]*flagData{},
		},
		{ // 3
			Input: struct {
				FlagA string
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases: []string{"--flaga"},
					typ:     reflect.TypeFor[string](),
					index:   []int{0},
				},
			},
		},
		{ // 4
			Input: &struct {
				FlagA string
				FlagB int
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases: []string{"--flaga"},
					typ:     reflect.TypeFor[string](),
					index:   []int{0},
				},
				"--flagb": {
					aliases: []string{"--flagb"},
					typ:     reflect.TypeFor[int](),
					index:   []int{1},
				},
			},
		},
		{ // 5
			Input: &struct {
				FlagA string
				flagB int
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases: []string{"--flaga"},
					typ:     reflect.TypeFor[string](),
					index:   []int{0},
				},
			},
		},
		{ // 6
			Input: struct {
				FlagA string `flags:"--aaa"`
			}{},
			Output: map[string]*flagData{
				"--aaa": {
					aliases: []string{"--aaa"},
					typ:     reflect.TypeFor[string](),
					index:   []int{0},
				},
			},
		},
		{ // 7
			Input: struct {
				FlagA string `flags:"--aaa,--bbb"`
			}{},
			Output: map[string]*flagData{
				"--aaa": aliased,
				"--bbb": aliased,
			},
		},
		{ // 8
			Input: struct {
				FlagA string
				Global
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases: []string{"--flaga"},
					typ:     reflect.TypeFor[string](),
					index:   []int{0},
				},
				"--globalflag": {
					aliases: []string{"--globalflag"},
					typ:     reflect.TypeFor[string](),
					index:   []int{1, 0},
				},
			},
		},
		{ // 9
			Input: struct {
				FlagA string
				Sub
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases: []string{"--flaga"},
					typ:     reflect.TypeFor[string](),
					index:   []int{0},
				},
				"--another": {
					aliases: []string{"--another"},
					typ:     reflect.TypeFor[string](),
					index:   []int{1, 0},
				},
				"--globalflag": {
					aliases: []string{"--globalflag"},
					typ:     reflect.TypeFor[string](),
					index:   []int{1, 1, 0},
				},
			},
		},
		{ // 10
			Input: struct {
				FlagA *customFlag
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases: []string{"--flaga"},
					typ:     reflect.TypeFor[*customFlag](),
					index:   []int{0},
				},
			},
		},
		{ // 11
			Input: struct {
				FlagA customFlag
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases: []string{"--flaga"},
					typ:     reflect.TypeFor[customFlag](),
					index:   []int{0},
				},
			},
		},
		{ // 12
			Input: struct {
				FlagA string `default:"123"`
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases:      []string{"--flaga"},
					typ:          reflect.TypeFor[string](),
					index:        []int{0},
					defaultValue: "123",
				},
			},
		},
		{ // 13
			Input: struct {
				FlagA string
				FlagB string `flags:"--flaga"`
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases: []string{"--flaga"},
					typ:     reflect.TypeFor[string](),
					index:   []int{0},
				},
			},
			Err: ErrAmbiguousFlags,
		},
		{ // 14
			Input: struct {
				FlagA error
				FlagB func()
				FlagC chan int
				FlagD string
			}{},
			Output: map[string]*flagData{
				"--flagd": {
					aliases: []string{"--flagd"},
					typ:     reflect.TypeFor[string](),
					index:   []int{3},
				},
			},
		},
		{ // 15
			Input: struct {
				FlagA int `flags:"--globalflag"`
				Global
			}{},
			Output: map[string]*flagData{
				"--globalflag": {
					aliases: []string{"--globalflag"},
					typ:     reflect.TypeFor[int](),
					index:   []int{0},
				},
			},
		},
		{ // 16
			Input: struct {
				Global
				FlagA int `flags:"--globalflag"`
			}{},
			Output: map[string]*flagData{
				"--globalflag": {
					aliases: []string{"--globalflag"},
					typ:     reflect.TypeFor[int](),
					index:   []int{1},
				},
			},
		},
		{ // 17
			Input: struct {
				time.Duration
			}{},
			Output: map[string]*flagData{},
		},
		{ // 18
			Input: struct {
				FlagA string `usage:"Flag Usage Text"`
			}{},
			Output: map[string]*flagData{
				"--flaga": {
					aliases: []string{"--flaga"},
					typ:     reflect.TypeFor[string](),
					index:   []int{0},
					usage:   "Flag Usage Text",
				},
			},
		},
	} {
		var flags flags

		err := flags.process(reflect.TypeOf(test.Input), nil)
		if !errors.Is(err, test.Err) {
			t.Errorf("test %d: expecting error %v, got %v", n+1, test.Err, err)
		} else if !reflect.DeepEqual(flags.flags, test.Output) {
			t.Errorf("test %d: expected flags %v, got %v", n+1, test.Output, flags.flags)
		}
	}
}

func TestFlagsSet(t *testing.T) {
	for n, test := range [...]struct {
		Args   []string
		Output any
		Err    error
	}{
		{
			[]string{"--a", "abc"},
			struct {
				A string
			}{
				A: "abc",
			},
			nil,
		},
		{
			[]string{"--a", "abc", "--b"},
			struct {
				A string
				B bool
			}{
				A: "abc",
				B: true,
			},
			nil,
		},
		{
			[]string{"--a", "xxx"},
			struct {
				A arrayFlag
			}{
				A: arrayFlag{"xxx"},
			},
			nil,
		},
		{
			[]string{"--a", "abc"},
			struct {
				A *string
			}{
				A: new("abc"),
			},
			nil,
		},
		{
			[]string{"--z", "abc"},
			struct {
				A *string
			}{},
			ErrInvalidArgument,
		},
		{
			[]string{"--a"},
			struct {
				A *string
			}{},
			ErrMissingValue,
		},
	} {
		var flags flags

		v := reflect.New(reflect.TypeOf(test.Output)).Elem()

		if err := flags.process(reflect.TypeOf(test.Output), nil); err != nil {
			t.Errorf("test %d: unexpected error: %v", n+1, err)
		} else if err = flags.set(v, test.Args); !errors.Is(err, test.Err) {
			t.Errorf("test %d: expecting error %v, got %v", n+1, test.Err, err)
		} else if !reflect.DeepEqual(v.Interface(), test.Output) {
			t.Errorf("test %d: expecting %+v, got %+v", n+1, test.Output, v.Interface())
		}
	}
}
