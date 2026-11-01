package subcommands

import "reflect"

type prep struct {
	flags, result reflect.Type
	fn            reflect.Value
}
