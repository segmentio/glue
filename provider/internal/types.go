package internal

import (
	"fmt"
	"go/types"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/segmentio/glue/provider"
)

// IsExportedOrBuiltin returns true if a type is either exported or primitive.
//
// Note(tejasmanohar): If the type is a map, both the key and value must be exported
// since they're both necessary to represent the type.
func IsExportedOrBuiltin(_t types.Type) bool {
	for _, t := range Unpack(_t) {
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

// Unpack unpacks the most basic, underlying types from pointers, slices, maps, etc.
// There can be multiple types in the case of maps (key and value).
// Note(tejasmanohar): Are there any other structures that have multiple types
// that we need to validate are exported? Structs don't have to have all fields exported
// as we can just not encode the ones that aren't.
func Unpack(t types.Type) []types.Type {
	ret := make([]types.Type, 0, 1)

	if ptr, ok := t.(*types.Pointer); ok {
		ret = append(ret, Unpack(ptr.Elem())...)
	}

	if slice, ok := t.(*types.Slice); ok {
		ret = append(ret, Unpack(slice.Elem())...)
	}

	if array, ok := t.(*types.Array); ok {
		ret = append(ret, Unpack(array.Elem())...)
	}

	return ret
}

func GetTypeInfo(t types.Type) provider.TypeInfo {
	var ret provider.TypeInfo

	switch specific := t.(type) {
	case *types.Pointer:
		ret = GetTypeInfo(specific.Elem())
	case *types.Named:
		obj := specific.Obj()
		ret = provider.TypeInfo{
			Identifier: fmt.Sprintf("%s.%s", obj.Pkg().Name(), obj.Name()),
			Imports:    []string{obj.Pkg().Path()},
		}
	case *types.Basic:
		ret = provider.TypeInfo{
			Identifier: specific.Name(),
		}
	case *types.Slice:
		underlying := GetTypeInfo(specific.Elem())
		ret = provider.TypeInfo{
			Identifier: fmt.Sprintf("[]%s", underlying.Identifier),
			Imports:    underlying.Imports,
		}
	case *types.Array:
		underlying := GetTypeInfo(specific.Elem())
		ret = provider.TypeInfo{
			Identifier: fmt.Sprintf("[]%s", underlying.Identifier),
			Imports:    underlying.Imports,
		}
	case *types.Map:
		key := GetTypeInfo(specific.Key())
		val := GetTypeInfo(specific.Elem())

		ret = provider.TypeInfo{
			Identifier: fmt.Sprintf("map[%s]%s", key.Identifier, val.Identifier),
			Imports:    append(key.Imports, val.Imports...),
		}
	default:
		panic(fmt.Errorf("unexpected type: %s", reflect.TypeOf(t)))
	}

	for i, pkg := range ret.Imports {
		ret.Imports[i] = stripVendor(pkg)
	}

	return ret
}

// github.com/x/y/vendor/github.com/a/b -> github.com/a/b
func stripVendor(path string) string {
	dirs := strings.Split(path, string(filepath.Separator))
	num := len(dirs)
	vendorIndex := -1
	for i := 1; i <= num; i++ {
		dir := dirs[num-i]
		if dir == "vendor" {
			vendorIndex = num - i
			break
		}
	}

	if vendorIndex == -1 {
		return path
	}

	return filepath.Join(dirs[vendorIndex+1:]...)
}
