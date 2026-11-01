package subcommands

import (
	"errors"
	"reflect"
)

type command struct {
	subcommands map[string]*command
	flags       reflect.Type
	fn          reflect.Value
}

func (c *command) Register[T any](short, long string, fn func(T) error, subs ...string) error {
	if c.subcommands == nil {
		c.subcommands = make(map[string]*command)
	}

	if len(subs) > 0 {
		d := new(command)
		c.subcommands[subs[0]] = d

		return d.Register(short, long, fn, subs[1:]...)
	}

	if !c.fn.IsValid() {
		return ErrSubcommandRegistered
	}

	// validate flags now.

	c.flags = reflect.TypeFor[T]()
	c.fn = reflect.ValueOf(fn)

	return nil
}

var ErrSubcommandRegistered = errors.New("subcommand already registered")
