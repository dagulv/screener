package domain

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

type Nullable[T any] struct {
	Content T
	Valid   bool
}

// Value implements the database/sql/driver Valuer interface.
func (src Nullable[T]) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}

	return src.Content, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Nullable[T]) Scan(src any) error {
	if src == nil {
		*dst = Nullable[T]{}
		return nil
	}

	var t T

	// TODO: Support other custom types
	switch v := any(&t).(type) {
	case sql.Scanner:
		if err := v.Scan(src); err != nil {
			return err
		}

		src = t
	}

	switch v := src.(type) {
	case T:
		*dst = Nullable[T]{Content: v, Valid: true}
		return nil
	case []byte:
		switch any(t).(type) {
		case string:
			*dst = Nullable[T]{Content: any(string(v)).(T), Valid: true}
			return nil
		}
		*dst = Nullable[T]{Content: any(v).(T), Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan type")
}

func (src Nullable[T]) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(src.Content)
}

func (dst *Nullable[T]) UnmarshalJSON(b []byte) error {
	var t *T

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	if t == nil {
		*dst = Nullable[T]{}
	} else {
		*dst = Nullable[T]{Content: *t, Valid: true}
	}

	return nil
}

// TypeDescription implements registry.TypeDescriber.
func (Nullable[T]) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(tags reflect.StructTag) (schema openapi.Schema, err error) {
			return reg.Schema(reflect.TypeFor[T](), tags)
		},
		Decoder: func(tags reflect.StructTag) (registry.Decoder, error) {
			v, err := reg.Decoder(reflect.TypeFor[T](), tags)

			if err != nil {
				return nil, err
			}

			return func(p unsafe.Pointer, s string) error {
				ptr := (*Nullable[T])(p)
				err := v(p, s)
				ptr.Valid = err == nil

				return err
			}, nil
		},
	}
}

func (v Nullable[T]) IsZero() bool {
	return !v.Valid
}

func (v Nullable[T]) IsNil() bool {
	return !v.Valid
}
