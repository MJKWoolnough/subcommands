package subcommands

import (
	"errors"
	"flag"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"unicode/utf8"
)

var flagI = reflect.TypeFor[flag.Value]()

type flagData struct {
	aliases      []string
	typ          reflect.Type
	defaultValue string
	index        []int
	usage        string
}

func (f *flagData) aliasLength() int {
	c := len(f.aliases) - 1

	for _, f := range f.aliases {
		c += utf8.RuneCountInString(f)
	}

	return c
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
		case reflect.Interface, reflect.Func, reflect.Chan:
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

	fd.defaultValue = t.Tag.Get("default")
	fd.usage = t.Tag.Get("usage")

	return fd
}

func (f *flags) set(v reflect.Value, args []string) error {
	for i := 0; i < len(args); i++ {
		flag, ok := f.flags[args[i]]
		if !ok {
			return fmt.Errorf("%w: %s", ErrInvalidArgument, args[i])
		}

		w, next, err := setValue(flag.typ, args[i+1:])
		if err != nil {
			return err
		}

		if next {
			i++
		}

		v.FieldByIndex(flag.index).Set(w)
	}

	return nil
}

func (f *flags) list() []*flagData {
	flags := slices.Collect(maps.Keys(maps.Collect(func(yield func(*flagData, struct{}) bool) {
		for _, v := range f.flags {
			yield(v, struct{}{})
		}
	})))

	slices.SortFunc(flags, func(a, b *flagData) int {
		return strings.Compare(a.aliases[0], b.aliases[0])
	})

	return flags
}

var (
	ErrFlagStructRequired = errors.New("flag type must be a struct")
	ErrAmbiguousFlags     = errors.New("ambiguous flags")
	ErrInvalidArgument    = errors.New("invalid argument")
)
