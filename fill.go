// Package feign automatically fills types with random data. It can be useful
// when testing, where you only need to verify persistence and retrieval and
// aren't concerned with valid data values.
package feign

import (
	"errors"
	"reflect"
	"time"
)

var typeTime = reflect.TypeOf((*time.Time)(nil)).Elem()

// ErrUnhandledType is an error returned when the type is not fillable.
// Unfillable types include funcs, chans, and interfaces. When filling structs,
// maps, or slices, these types will be ignored.
var ErrUnhandledType = errors.New("unhandled type")

// Filler is a func used to provide the value used to fill a struct field.
type Filler func(path string) (val interface{}, ok bool)

// Fill fills a type with random data.
func Fill(val interface{}, fillers ...Filler) error {
	t := reflect.TypeOf(val)
	if t.Kind() != reflect.Ptr || reflect.ValueOf(val).IsNil() {
		return errors.New("not a pointer value")
	}
	v := reflect.ValueOf(val)
	result, err := getValue("", val, fillers...)
	if err != nil {
		return err
	}
	v.Elem().Set(result.Elem().Convert(t.Elem()))
	return nil
}

// MustFill fills a type with random data and panics if there is an error.
func MustFill(val interface{}, fillers ...Filler) {
	err := Fill(val, fillers...)
	if err != nil {
		panic(err)
	}
}

func getValue(path string, a interface{}, fillers ...Filler) (reflect.Value, error) {
	if path != "" {
		for _, fn := range fillers {
			if v, ok := fn(path); ok {
				if v == nil {
					return reflect.Zero(reflect.TypeOf(a)), nil
				}
				return reflect.ValueOf(v), nil
			}
		}
	}
	t := reflect.TypeOf(a)
	if t == nil {
		return reflect.Value{}, ErrUnhandledType
	}
	switch t.Kind() {
	case reflect.Ptr:
		v := reflect.New(t.Elem())
		var val reflect.Value
		var err error
		if a != reflect.Zero(reflect.TypeOf(a)).Interface() {
			val, err = getValue(path, reflect.ValueOf(a).Elem().Interface(), fillers...)
			if err != nil {
				return reflect.Value{}, err
			}
		} else {
			val, err = getValue(path, v.Elem().Interface(), fillers...)
			if err != nil {
				return reflect.Value{}, err
			}
		}
		v.Elem().Set(val.Convert(t.Elem()))
		return v, nil

	case reflect.Struct:
		switch t {
		case typeTime:
			ft := time.Time{}.Add(time.Duration(random.Int63()))
			return reflect.ValueOf(ft), nil
		default:
			v := reflect.New(t).Elem()
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				if !field.CanSet() {
					continue // avoid panic to set on unexported field in struct
				}
				val, err := getValue(path+"."+t.Field(i).Name, field.Interface(), fillers...)
				if err == ErrUnhandledType {
					continue
				}
				if err != nil {
					return reflect.Value{}, err
				}
				val = val.Convert(field.Type())
				v.Field(i).Set(val)
			}
			return v, nil
		}

	case reflect.String:
		return reflect.ValueOf(randomString()), nil

	case reflect.Array:
		v := reflect.New(t).Elem()
		for i := 0; i < t.Len(); i++ {
			val, err := getValue(path, v.Index(i).Interface(), fillers...)
			if err == ErrUnhandledType {
				continue
			}
			if err != nil {
				return reflect.Value{}, err
			}
			v.Index(i).Set(val)
		}
		return v, nil

	case reflect.Slice:
		len := randomSliceAndMapSize()
		if len == 0 {
			return reflect.Zero(t), nil
		}
		v := reflect.MakeSlice(t, len, len)
		for i := 0; i < v.Len(); i++ {
			val, err := getValue(path, v.Index(i).Interface(), fillers...)
			if err == ErrUnhandledType {
				continue
			}
			if err != nil {
				return reflect.Value{}, err
			}
			v.Index(i).Set(val)
		}
		return v, nil

	case reflect.Map:
		len := randomSliceAndMapSize()
		if len == 0 {
			return reflect.Zero(t), nil
		}
		v := reflect.MakeMapWithSize(t, len)
		for i := 0; i < len; i++ {
			keyInstance := reflect.New(t.Key()).Elem().Interface()
			key, err := getValue(path, keyInstance, fillers...)
			if err == ErrUnhandledType {
				continue
			}
			if err != nil {
				return reflect.Value{}, err
			}

			valueInstance := reflect.New(t.Elem()).Elem().Interface()
			val, err := getValue(path, valueInstance, fillers...)
			if err == ErrUnhandledType {
				continue
			}
			if err != nil {
				return reflect.Value{}, err
			}
			v.SetMapIndex(key, val)
		}
		return v, nil

	case reflect.Bool:
		return reflect.ValueOf(random.Intn(2) > 0), nil

	case reflect.Int:
		return reflect.ValueOf(randomInteger()), nil
	case reflect.Int8:
		return reflect.ValueOf(int8(randomInteger())), nil
	case reflect.Int16:
		return reflect.ValueOf(int16(randomInteger())), nil
	case reflect.Int32:
		return reflect.ValueOf(int32(randomInteger())), nil
	case reflect.Int64:
		return reflect.ValueOf(int64(randomInteger())), nil

	case reflect.Float32:
		return reflect.ValueOf(random.Float32()), nil
	case reflect.Float64:
		return reflect.ValueOf(random.Float64()), nil

	case reflect.Uint:
		return reflect.ValueOf(uint(randomInteger())), nil
	case reflect.Uint8:
		return reflect.ValueOf(uint8(randomInteger())), nil
	case reflect.Uint16:
		return reflect.ValueOf(uint16(randomInteger())), nil
	case reflect.Uint32:
		return reflect.ValueOf(uint32(randomInteger())), nil
	case reflect.Uint64:
		return reflect.ValueOf(uint64(randomInteger())), nil

	default:
		return reflect.Value{}, ErrUnhandledType
	}
}
