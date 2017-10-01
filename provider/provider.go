package provider

import (
	"go/types"
)

// A Provider surfaces RPC-implementation-specific details.
type Provider interface {
	IsSuitableMethod(*types.Func) bool
	GetArgType(*types.Func) types.Type
	GetReplyType(*types.Func) types.Type
}
