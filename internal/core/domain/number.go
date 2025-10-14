package domain

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

type Number interface {
	int | float32
}

func valueFrom[T Number](s string) (t Nullable[T], err error) {
	if s == "" {
		t = Optional(T(0))
		return
	}

	switch any(t.Content).(type) {
	case int:
		var i int
		i, err = strconv.Atoi(s)
		if err != nil {
			return
		}
		t = Optional(T(i), true)
	case float32:
		var f float64
		f, err = strconv.ParseFloat(s, 32)
		if err != nil {
			return
		}
		t = Optional(T(f), true)
	default:
		return t, errors.New("invalid type")
	}

	return
}

type MinMax[T Number] struct {
	Min Nullable[T] `json:"min"`
	Max Nullable[T] `json:"max"`
}

func (m MinMax[T]) IsZero() bool {
	return !m.Min.Valid && !m.Max.Valid
}

// TypeDescription implements registry.TypeDescriber.
func (MinMax[T]) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(tags reflect.StructTag) (schema openapi.Schema, err error) {
			return reg.Schema(reflect.TypeFor[T](), tags)
		},
		Parser: func(tags reflect.StructTag) (registry.Parser, error) {
			return func(p unsafe.Pointer, s string) error {
				ptr := (*MinMax[T])(p)

				raw := strings.Split(s, ",")
				min, err := valueFrom[T](raw[0])
				if err != nil {
					return err
				}
				max, err := valueFrom[T](raw[1])
				if err != nil {
					return err
				}
				if min.Valid && max.Valid && min.Content > max.Content {
					min.Valid = false
				}
				if min.Valid && max.Valid && max.Content < min.Content {
					max.Valid = false
				}
				*ptr = MinMax[T]{
					Min: min,
					Max: max,
				}

				return err
			}, nil
		},
	}
}
