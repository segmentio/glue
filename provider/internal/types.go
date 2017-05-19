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

// Unpack unpacks the most basic, underlying type from pointers, slices, etc.
func Unpack(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return Unpack(ptr.Elem())
	}

	if slice, ok := t.(*types.Slice); ok {
		return Unpack(slice.Elem())
	}

	if array, ok := t.(*types.Array); ok {
		return Unpack(array.Elem())
	}

	return t
}

// GetTypeInfo derives provider.TypeInfo, a Glue construct for metadata about a type,
// from a types.Type.
func GetTypeInfo(t types.Type) provider.TypeInfo {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

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
	case *types.Slice:
		return provider.TypeInfo{
			// e.g. String() => []int
			Name:    specific.String(),
			Package: GetTypeInfo(specific.Elem()).Package,
		}
	case *types.Array:
		return provider.TypeInfo{
			Name:    specific.String(),
			Package: GetTypeInfo(specific.Elem()).Package,
		}
	}

	panic(fmt.Errorf("unexpected type: %s", reflect.TypeOf(t)))
}
