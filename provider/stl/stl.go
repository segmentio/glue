package stl

import (
	"go/types"

	"github.com/tejasmanohar/glue/log"
	"github.com/tejasmanohar/glue/provider"
	"github.com/tejasmanohar/glue/provider/internal"
)

// Provider is a Glue provider for net/rpc.
type Provider struct{}

// IsSuitableMethod determines if a receiver method is structured as a net/rpc method.
// The criteria is net/rpc.suitableMethods ported from reflect to types.Type.
// https://github.com/golang/go/blob/release-branch.go1.8/src/net/rpc/server.go#L292
func (p *Provider) IsSuitableMethod(method *types.Func) bool {
	if !method.Exported() {
		log.Debugf("skipping %s: unexported", method.Name())
		return false
	}

	signature := method.Type().(*types.Signature)
	params := signature.Params()

	if params.Len() != 2 {
		log.Debugf("skipping %s: expected 2 params, found %d", method.Name(), params.Len())
		return false
	}

	arg := params.At(0)
	if !internal.IsExportedOrBuiltin(arg.Type()) {
		log.Debugf("skipping %s: argument parameter's type is not exported", method.Name())
		return false
	}

	reply := params.At(1)
	replyType := reply.Type()
	if !internal.IsExportedOrBuiltin(replyType) {
		log.Debugf("skipping %s: reply parameter's type is not exported", method.Name())
		return false
	}

	if _, ok := replyType.(*types.Pointer); !ok {
		log.Debugf("skipping %s: reply type is not a pointer")
		return false
	}

	returns := signature.Results()
	if returns.Len() != 1 {
		log.Debugf("skipping %s: expected 1 return value, found %d", method.Name(), returns.Len())
		return false
	}

	err := returns.At(0)
	if err.Type().String() != "error" {
		log.Debugf("skipping %s: expected func to return `error`, found %s", method.Name(), err.Type().String())
	}

	return true
}

// GetArgType extracts metadata about the response type from an RPC method.
func (p *Provider) GetArgType(f *types.Func) provider.TypeInfo {
	signature := f.Type().(*types.Signature)
	params := signature.Params()
	arg := params.At(0)
	return internal.GetTypeInfo(arg.Type())
}

// GetReplyType extracts metadata about the response type from an RPC method.
func (p *Provider) GetReplyType(f *types.Func) provider.TypeInfo {
	signature := f.Type().(*types.Signature)
	params := signature.Params()
	reply := params.At(1)
	return internal.GetTypeInfo(reply.Type())
}
