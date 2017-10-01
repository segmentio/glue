package generator

import (
	"go/types"
	"path/filepath"
	"strconv"
	"strings"
)

type resolver struct {
	// imports is a map of package name to package path to an incrementing integer
	// for a given package name, there may be multiple paths--
	// e.g. "mypkg" is the name of paths "github.com/x/mypkg" and "github.com/y/mypkg"
	//
	// The incrementing integer is used to handle collisions. If you need to import
	// - "github.com/x/mypackage"
	// - "github.com/y/mypackage"
	// - "github.com/z/mypackage"
	// they will be imported as
	// mypackage "github.com/x/mypackage"
	// mypackage1 "github.com/y/mypackage"
	// mypackage2 "github.com/z/mypackage"
	imports map[string]map[string]importMapping
}

type importMapping struct {
	Name    string
	counter int
}

func newResolver() *resolver {
	return &resolver{
		imports: map[string]map[string]importMapping{},
	}
}

func (r *resolver) GetTypeString(t types.Type) string {
	return types.TypeString(t, func(pkg *types.Package) string {
		name := pkg.Name()
		path := stripVendor(pkg.Path())

		if existingPaths, ok := r.imports[name]; ok {
			mapping, ok := existingPaths[path]
			if ok {
				return mapping.Name
			}

			n := len(existingPaths)
			rename := name + strconv.Itoa(n)
			r.imports[name][path] = importMapping{
				Name:    rename,
				counter: n,
			}

			return rename
		}

		r.imports[name] = map[string]importMapping{}
		r.imports[name][path] = importMapping{Name: name}
		return name
	})
}

func (r *resolver) GetImports() []Import {
	ret := make([]Import, 0, len(r.imports))
	for originalName, pathsMap := range r.imports {
		for path, info := range pathsMap {
			var name string
			if originalName != info.Name {
				name = info.Name
			}

			ret = append(ret, Import{
				Name: name,
				Path: path,
			})
		}
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
