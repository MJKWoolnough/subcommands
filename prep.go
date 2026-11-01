package subcommands

import "reflect"

type prep struct {
	flags  flags
	result reflect.Type
	fn     reflect.Value
}
