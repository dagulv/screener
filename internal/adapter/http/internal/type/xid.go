package types

import (
	"reflect"
	"unsafe"

	"github.com/rs/xid"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

type idType struct{}

func XID() registry.TypeRegistrar {
	return idType{}
}

// Type implements registry.TypeRegistrar.
func (i idType) Type() reflect.Type {
	return reflect.TypeFor[xid.ID]()
}

// TypeDescription implements registry.TypeRegistrar.
func (i idType) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(tags reflect.StructTag) (openapi.Schema, error) {
			return &openapi.String{}, nil
		},
		Decoder: func(tags reflect.StructTag) (registry.Decoder, error) {
			return func(p unsafe.Pointer, s string) (err error) {
				id := (*xid.ID)(p)
				*id, err = xid.FromString(s)
				return
			}, nil
		},
	}
}
