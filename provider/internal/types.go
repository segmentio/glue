package internal

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/tejasmanohar/glue/provider"
)

// IsExportedOrBuiltin returns true if a type is either exported or primitive.
func IsExportedOrBuiltin(t types.Type) bool {
	t = Unpack(t)

	if _, isPrimitive := t.(*types.Basic); isPrimitive {
		return true
	}

	namedType, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := namedType.Obj()
	if obj == nil {
		return false
	}

	return obj.Exported()
}

// Unpack unpacks the underlying type from a pointer if wrapped.
func Unpack(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}

// GetTypeInfo derives provider.TypeInfo, a Glue construct for metadata about a type,
// from a types.Type.
func GetTypeInfo(t types.Type) provider.TypeInfo {
	t = Unpack(t)

	switch specific := t.(type) {
	case *types.Named:
		return provider.TypeInfo{
			Name:    specific.Obj().Name(),
			Package: specific.Obj().Pkg(),
		}
	case *types.Basic:
		return provider.TypeInfo{
			Name: specific.Name(),
		}
	}

	panic(fmt.Errorf("unexpected type: %s", reflect.TypeOf(t)))
}
