package subcommands

import (
	"errors"
	"reflect"
)

func (c *Commands) run(args []string) error {
	cmd := &c.command

	var subcommand []string

	for n, arg := range args {
		if sub, ok := cmd.subcommands[arg]; ok {
			cmd = sub
		} else {
			subcommand = args[:n]
			args = args[n:]

			break
		}
	}

	if !c.fn.IsValid() {
		return cmd.wrapHelp(subcommand, ErrInvalidSubcommand)
	}

	at := c.fn.Type().In(0)

	var fv [2]reflect.Value

	if at.Kind() == reflect.Pointer {
		fv[0] = reflect.New(at.Elem())
		fv[1] = fv[0]
	} else {
		fv[0] = reflect.New(at).Elem()
		fv[1] = fv[0].Addr()
	}

	if err := c.flags.set(fv[1], args); err != nil {
		return cmd.wrapHelp(subcommand, err)
	}

	return cmd.wrapHelp(subcommand, c.fn.Call(fv[:])[0].Interface().(error))
}

var ErrInvalidSubcommand = errors.New("invalid subcommand")
