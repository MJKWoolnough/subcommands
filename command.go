package subcommands

import (
	"errors"
	"reflect"
)

type command struct {
	subcommands map[string]*command
	flags       flags
	fn          reflect.Value
}

func (c *command) register[T any](short, long string, fn func(T) error, f flags, subs ...string) error {
	if c.subcommands == nil {
		c.subcommands = make(map[string]*command)
	}

	if len(subs) > 0 {
		d := new(command)
		c.subcommands[subs[0]] = d

		return d.register(short, long, fn, f, subs[1:]...)
	}

	if !c.fn.IsValid() {
		return ErrSubcommandRegistered
	}

	c.flags = f
	c.fn = reflect.ValueOf(fn)

	return nil
}

var ErrSubcommandRegistered = errors.New("subcommand already registered")
