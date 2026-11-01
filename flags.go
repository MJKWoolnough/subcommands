package subcommands

import (
	"errors"
	"reflect"
)

type flags struct{}

func (f *flags) process(typ reflect.Type) error {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Struct {
		return ErrFlagStructRequired
	}

	for field := range typ.Fields() {
		if !field.IsExported() {
			continue
		}

		if field.Anonymous {
			if err := f.process(field.Type); errors.Is(err, ErrFlagStructRequired) {
				continue
			} else if err != nil {
				return err
			}
		}
	}

	return nil
}

var ErrFlagStructRequired = errors.New("flag type must be a struct")
