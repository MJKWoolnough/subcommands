package subcommands

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
)

var flagI = reflect.TypeFor[flag.Value]()

type flagData struct {
	aliases      []string
	typ          reflect.Type
	defaultValue string
	index        []int
}

type flags struct {
	flags map[string]*flagData
}

func (f *flags) process(typ reflect.Type, index []int) error {
	if f.flags == nil {
		f.flags = make(map[string]*flagData)
	}

	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return ErrFlagStructRequired
	}

	for field := range typ.Fields() {
		if !field.IsExported() {
			continue
		}

		if field.Anonymous {
			if err := f.process(field.Type, append(index, field.Index...)); errors.Is(err, ErrFlagStructRequired) {
				continue
			} else if err != nil {
				return err
			}
		}

		typ := field.Type

		if typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}

		switch typ.Kind() {
		case reflect.Interface:
			continue
		case reflect.Struct:
			if !field.Type.Implements(flagI) && !reflect.PointerTo(field.Type).Implements(flagI) {
				continue
			}

			fallthrough
		default:
			g := parseFlag(field.Name, field, append(index, field.Index...))

			for _, alias := range g.aliases {
				if e, ok := f.flags[alias]; ok && len(e.index) == len(index)+1 {
					return fmt.Errorf("%w: %s", ErrAmbiguousFlags, alias)
				} else if ok && len(e.index) <= len(index) {
					continue
				}

				f.flags[alias] = g
			}
		}
	}

	return nil
}

func parseFlag(name string, t reflect.StructField, index []int) *flagData {
	fd := &flagData{
		typ:   t.Type,
		index: index,
	}

	if aliases, ok := t.Tag.Lookup("flags"); ok {
		fd.aliases = strings.Split(aliases, ",")
	} else {
		fd.aliases = []string{"--" + strings.ToLower(name)}
	}

	if defaultValue, ok := t.Tag.Lookup("default"); ok {
		fd.defaultValue = defaultValue
	}

	return fd
}

var (
	ErrFlagStructRequired = errors.New("flag type must be a struct")
	ErrAmbiguousFlags     = errors.New("ambiguous flags")
)
