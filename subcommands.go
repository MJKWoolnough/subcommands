package subcommands

import (
	"errors"
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

func (c *Commands) Register[T any](short, long string, fn func(T) error, subs ...string) error {
	if !isStruct(reflect.TypeFor[T]()) {
		return ErrFlagStructRequired
	}

	return c.command.Register(short, long, fn, subs...)
}

func (c *Commands) Prepare[T, U any](fn func(T) (U, error)) error {
	p := prep{
		flags:  reflect.TypeFor[T](),
		result: reflect.TypeFor[U](),
		fn:     reflect.ValueOf(fn),
	}

	if !isStruct(p.flags) {
		return ErrFlagStructRequired
	}

	c.preperations = append(c.preperations, p)

	return nil
}

func (c *Commands) Run() error {
	return c.RunWithArgs(os.Args)
}

func (c *Commands) RunWithArgs(args []string) error {
	return nil
}

func Register[T any](short, long string, fn func(T) error, subs ...string) error {
	return globalCommand.Register(short, long, fn, subs...)
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

var ErrFlagStructRequired = errors.New("flag type must be a struct")
