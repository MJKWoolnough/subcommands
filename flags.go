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
			g := parseFlag(field.Name, field, level)

			for _, alias := range g.aliases {
				if e, ok := f.flags[alias]; ok && e.depth == level {
					return fmt.Errorf("%w: %s", ErrAmbiguousFlags, alias)
				} else if ok && e.depth < level {
					continue
				}

				f.flags[alias] = g
			}
		}
	}

	return nil
}

func parseFlag(name string, t reflect.StructField, level int) *flagData {
	fd := &flagData{
		depth: level,
		typ:   t.Type,
		index: t.Index,
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
