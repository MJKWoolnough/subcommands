package subcommands

import (
	"errors"
	"flag"
	"reflect"
	"strconv"
)

func setValue(typ reflect.Type, args []string) (reflect.Value, bool, error) {
	var (
		next string
		ok   bool
	)

	if len(args) > 0 {
		next = args[0]
		ok = true
	}

	if typ.Kind() == reflect.Bool {
		if !ok {
			return reflect.ValueOf(true), false, nil
		} else if v, err := strconv.ParseBool(next); err == nil {
			return reflect.ValueOf(v), true, nil
		}

		return reflect.ValueOf(true), false, nil
	} else if !ok {
		return reflect.Value{}, false, ErrMissingValue
	} else if typ.Implements(flagI) {
		pv, v := makeType(typ)

		if err := pv.Interface().(flag.Value).Set(next); err != nil {
			return reflect.Value{}, false, err
		}

		return v, true, nil
	} else if reflect.PointerTo(typ).Implements(flagI) {
		pv, v := makeType(typ)

		if err := pv.Interface().(flag.Value).Set(next); err != nil {
			return reflect.Value{}, false, err
		}

		return v, true, nil
	}

	switch typ.Kind() {
	case reflect.Int:
		return parseInt[int](next, 0)
	case reflect.Int8:
		return parseInt[int8](next, 8)
	case reflect.Int16:
		return parseInt[int16](next, 16)
	case reflect.Int32:
		return parseInt[int32](next, 32)
	case reflect.Int64:
		return parseInt[int64](next, 64)
	case reflect.Uint:
		return parseUint[uint](next, 0)
	case reflect.Uint8:
		return parseUint[uint8](next, 6)
	case reflect.Uint16:
		return parseUint[uint16](next, 16)
	case reflect.Uint32:
		return parseUint[uint32](next, 32)
	case reflect.Uint64:
		return parseUint[uint64](next, 64)
	case reflect.Float32:
		return parseFloat[float32](next, 32)
	case reflect.Float64:
		return parseFloat[float64](next, 64)
	case reflect.Pointer:
		v, n, err := setValue(typ.Elem(), args)
		if err != nil {
			return reflect.Value{}, false, err
		}

		p := reflect.New(typ.Elem())

		p.Elem().Set(v)

		return p, n, nil
	case reflect.Slice:
		v, n, err := setValue(typ.Elem(), args)
		if err != nil {
			return reflect.Value{}, false, err
		}

		sv := reflect.MakeSlice(typ, 1, 1)
		sv.Index(0).Set(v)

		return sv, n, nil
	case reflect.String:
		return reflect.ValueOf(next), true, nil
	}

	return reflect.Value{}, false, ErrInvalidType
}

func makeType(typ reflect.Type) (reflect.Value, reflect.Value) {
	switch typ.Kind() {
	case reflect.Map:
		v := reflect.MakeMap(typ)

		return v, v
	case reflect.Slice:
		pv := reflect.New(typ)

		pv.Elem().Set(reflect.MakeSlice(typ, 0, 0))

		return pv, pv.Elem()
	case reflect.Pointer:
		v := reflect.New(typ.Elem())

		return v, v
	default:
		v := reflect.New(typ)

		return v, v.Elem()
	}
}

func parseInt[T int8 | int16 | int32 | int64 | int](next string, bits int) (reflect.Value, bool, error) {
	v, err := strconv.ParseInt(next, 10, bits)

	return reflect.ValueOf(T(v)), true, err
}

func parseUint[T uint8 | uint16 | uint32 | uint64 | uint](next string, bits int) (reflect.Value, bool, error) {
	v, err := strconv.ParseUint(next, 10, bits)

	return reflect.ValueOf(T(v)), true, err
}

func parseFloat[T float32 | float64](next string, bits int) (reflect.Value, bool, error) {
	v, err := strconv.ParseFloat(next, bits)

	return reflect.ValueOf(T(v)), true, err
}

var (
	ErrInvalidType  = errors.New("invalid type")
	ErrMissingValue = errors.New("missing value")
)
