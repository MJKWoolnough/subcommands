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
		{
			Input:  int(0),
			Output: map[string]*flagData{},
			Err:    ErrFlagStructRequired,
		},
		{
			Input:  struct{}{},
			Output: map[string]*flagData{},
		},
		{
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
		{
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
		{
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
		{
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
		{
			Input: struct {
				FlagA string `flags:"--aaa,--bbb"`
			}{},
			Output: map[string]*flagData{
				"--aaa": aliased,
				"--bbb": aliased,
			},
		},
		{
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
		{
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
		{
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
		{
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
		{
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
		{
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
		{
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
		{
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
		{
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
		{
			Input: struct {
				time.Duration
			}{},
			Output: map[string]*flagData{},
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
