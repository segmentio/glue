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
	case *types.Interface:
		if !specific.Empty() {
			// TODO(tejasmanohar): This shouldn't really be a "panic", as it's not a
			// developer error but a user one. Figure that out later.
			panic("non-empty interfaces are not cannot be JSON-encoded")
		}

		ret = provider.TypeInfo{Identifier: "interface{}"}
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
