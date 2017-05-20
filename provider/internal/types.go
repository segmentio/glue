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

func GetTypeInfo(t types.Type) provider.TypeInfo {
	switch specific := t.(type) {
	case *types.Pointer:
		return GetTypeInfo(specific.Elem())
	case *types.Named:
		obj := specific.Obj()
		return provider.TypeInfo{
			Identifier: fmt.Sprintf("%s.%s", obj.Pkg().Name(), obj.Name()),
			Imports:    []string{obj.Pkg().Path()},
		}
	case *types.Basic:
		return provider.TypeInfo{
			Identifier: specific.Name(),
		}
	case *types.Slice:
		underlying := GetTypeInfo(specific.Elem())
		return provider.TypeInfo{
			Identifier: fmt.Sprintf("[]%s", underlying.Identifier),
			Imports:    underlying.Imports,
		}
	case *types.Array:
		underlying := GetTypeInfo(specific.Elem())
		return provider.TypeInfo{
			Identifier: fmt.Sprintf("[]%s", underlying.Identifier),
			Imports:    underlying.Imports,
		}
	case *types.Map:
		key := GetTypeInfo(specific.Key())
		val := GetTypeInfo(specific.Elem())

		return provider.TypeInfo{
			Identifier: fmt.Sprintf("map[%s]%s", key.Identifier, val.Identifier),
			Imports:    append(key.Imports, val.Imports...),
		}
	default:
		panic(fmt.Errorf("unexpected type: %s", reflect.TypeOf(t)))
	}
}
