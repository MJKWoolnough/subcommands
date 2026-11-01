package subcommands

import "reflect"

func isStruct(typ reflect.Type) bool {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return typ.Kind() == reflect.Struct
}
