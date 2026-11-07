// Package subcommands allows for simple subcommand and flag parsing.
package subcommands // import "vimagination.zapto.org/subcommands"

import (
	"os"
	"reflect"
)

type (
	Required  struct{}
	Mutual    struct{}
	Exclusive struct{}
	Requires  struct{}
	Forbids   struct{}
	OneOf     struct{}
	AnyOf     struct{}
)

type Commands struct {
	command
	preperations []prep
}

var globalCommand Commands

func (c *Commands) Register[T any](fn func(T) error, opts ...Option) error {
	var f flags

	if err := f.process(reflect.TypeFor[T](), nil); err != nil {
		return err
	}

	var o options

	for _, opt := range opts {
		opt(&o)
	}

	return c.command.register(o.description, o.long, fn, f, o.subcommand...)
}

func (c *Commands) Prepare[T, U any](fn func(T) (U, error)) error {
	var f flags

	if err := f.process(reflect.TypeFor[T](), nil); err != nil {
		return err
	}

	p := prep{
		flags:  f,
		result: reflect.TypeFor[U](),
		fn:     reflect.ValueOf(fn),
	}

	c.preperations = append(c.preperations, p)

	return nil
}

func (c *Commands) Run() error {
	return c.RunWithArgs(os.Args)
}

func (c *Commands) RunWithArgs(args []string) error {
	return c.run(args)
}

func Register[T any](fn func(T) error, opts ...Option) error {
	return globalCommand.Register(fn, opts...)
}

func Prepare[T, U any](fn func(T) (U, error)) error {
	return globalCommand.Prepare(fn)
}

func Run() error {
	return globalCommand.RunWithArgs(os.Args)
}

func RunWithArgs(args []string) error {
	return globalCommand.RunWithArgs(args)
}
