package provider

import (
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
	Identifier string
	Imports    []string
}
