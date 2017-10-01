package internal

import (
	"go/types"
)

// IsExportedOrBuiltin returns true if a type is either exported or primitive.
//
// Note(tejasmanohar): If the type is a map, both the key and value must be exported
// since they're both necessary to represent the type.
func IsExportedOrBuiltin(_t types.Type) bool {
	for _, t := range unpack(_t) {
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

		if !obj.Exported() {
			return false
		}
	}

	return true
}

// unpack unpacks the most basic, underlying types from pointers, slices, maps, etc.
// There can be multiple types in the case of maps (key and value).
func unpack(t types.Type) []types.Type {
	ret := make([]types.Type, 0, 1)

	if ptr, ok := t.(*types.Pointer); ok {
		ret = append(ret, unpack(ptr.Elem())...)
	}

	if slice, ok := t.(*types.Slice); ok {
		ret = append(ret, unpack(slice.Elem())...)
	}

	if array, ok := t.(*types.Array); ok {
		ret = append(ret, unpack(array.Elem())...)
	}

	return ret
}

// Dereference dereferences pointers as needed.
func Dereference(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return Dereference(ptr.Elem())
	}

	return t
}
