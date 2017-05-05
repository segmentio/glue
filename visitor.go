package glue

import (
	"go/ast"
	"go/types"

	"github.com/tejasmanohar/glue/provider"

	"golang.org/x/tools/go/loader"
)

// Visitor traverses a Go package's AST, visits declarations that satisfy the
// structure of an RPC service, and extracts their RPC methods.
type Visitor struct {
	pkg      *loader.PackageInfo
	methods  map[string][]*types.Func
	provider provider.Provider

	target string
}

// VisitorConfig is used to create a Visitor.
type VisitorConfig struct {
	// Pkg contains metadata for the target package.
	Pkg *loader.PackageInfo
	// Provider determines which RPC methods are suitable.
	Provider provider.Provider
	// Declaration is the name of the target RPC declaration (method receiver).
	Declaration string
}

// NewVisitor creates a Visitor.
func NewVisitor(cfg VisitorConfig) *Visitor {
	return &Visitor{
		pkg:      cfg.Pkg,
		provider: cfg.Provider,
		methods:  map[string][]*types.Func{},
		target:   cfg.Declaration,
	}
}

// Go starts Visitor's trip around the supplied package. Upon return, it
// sends a mapping of receiver identifiers to RPC methods.
func (p *Visitor) Go() map[string][]*types.Func {
	for _, file := range p.pkg.Files {
		ast.Walk(p, file)
	}

	return p.methods
}

// Visit extracts functions from RPC declarations. It satisfies go/ast.Visitor.
func (p *Visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case nil:
		return nil
	case *ast.TypeSpec:
		p.visitType(n)
	}

	return p
}

func (p *Visitor) visitType(ts *ast.TypeSpec) {
	obj := p.pkg.Info.ObjectOf(ts.Name)
	if obj == nil {
		return
	}

	namedType, ok := obj.Type().(*types.Named)
	if !ok {
		return
	}

	if obj.Name() != p.target {
		return
	}

	for i := 0; i < namedType.NumMethods(); i++ {
		method := namedType.Method(i)
		if p.provider.IsSuitableMethod(method) {

			recv := namedType.Obj().Name()
			p.methods[recv] = append(p.methods[recv], method)
		}
	}
}
