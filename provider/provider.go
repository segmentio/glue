package provider

import (
	"fmt"
	"go/types"
)

// A Provider surfaces RPC-implementation-specific details.
type Provider interface {
	IsSuitableMethod(*types.Func) bool
	GetArgType(*types.Func) TypeInfo
	GetReplyType(*types.Func) TypeInfo
}

// TypeInfo contains metadata about input/output types.
type TypeInfo struct {
	Name    string
	Package *types.Package
}

// Identifier returns the identifier path for a type (e.g. pkg.SomeStruct, int64).
func (t TypeInfo) Identifier() string {
	if t.Package == nil {
		return t.Name
	}

	return fmt.Sprintf("%s.%s", t.Package.Name(), t.Name)
}
