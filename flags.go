package subcommands

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
)

var flagI = reflect.TypeFor[flag.Value]()

type flagData struct {
	depth        int
	aliases      []string
	typ          reflect.Type
	defaultValue string
	index        []int
}

type flags struct {
	flags map[string]*flagData
}

func (f *flags) process(typ reflect.Type, level int) error {
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
			if err := f.process(field.Type, level+1); errors.Is(err, ErrFlagStructRequired) {
				continue
			} else if err != nil {
				return err
			}
		}

		switch field.Type.Kind() {
		case reflect.Interface:
			continue
		case reflect.Struct:
			if !field.Type.Implements(flagI) {
				continue
			}

			fallthrough
		default:
			g, err := parseFlag(field.Name, field.Type)
			if err != nil {
				return err
			}

			for _, alias := range g.aliases {
				if e, ok := f.flags[alias]; ok && e.depth == level {
					return fmt.Errorf("%w: %s", ErrAmbiguousFlags, alias)
				} else if e.depth < level {
					continue
				}

				f.flags[alias] = g
			}
		}
	}

	return nil
}

func parseFlag(name string, t reflect.Type) (*flagData, error) {
	return nil, nil
}

var (
	ErrFlagStructRequired = errors.New("flag type must be a struct")
	ErrAmbiguousFlags     = errors.New("ambiguous flags")
)
