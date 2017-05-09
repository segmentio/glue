package gorilla

import (
	"go/types"

	"github.com/tejasmanohar/glue/provider"
)

// Provider is a Glue provider for gorilla/rpc.
// The main difference between stl and gorilla's rpc method format is the optional,
// request first argument in Gorilla so we shift that and proxy to stl provider in
// most of the methods.
type Provider struct {
	BaseProvider provider.Provider
}

// New creates a new gorilla/rpc Provider.
func New(base provider.Provider) *Provider {
	return &Provider{BaseProvider: base}
}

// IsSuitableMethod determines if a receiver method is structured as a gorilla/rpc method.
func (p *Provider) IsSuitableMethod(method *types.Func) bool {
	newMethod := p.shiftReqParam(method)
	return p.BaseProvider.IsSuitableMethod(newMethod)
}

// GetArgType proxies stl.GetArgType with a shifted function.
func (p *Provider) GetArgType(f *types.Func) provider.TypeInfo {
	newMethod := p.shiftReqParam(f)
	return p.BaseProvider.GetArgType(newMethod)
}

// GetReplyType proxies stl.GetReplyType with a shifted function.
func (p *Provider) GetReplyType(f *types.Func) provider.TypeInfo {
	newMethod := p.shiftReqParam(f)
	return p.BaseProvider.GetReplyType(newMethod)
}

// shiftReqParam returns a new *types.Func without the *http.Request param.
func (p *Provider) shiftReqParam(method *types.Func) *types.Func {
	originalSignature := method.Type().(*types.Signature)
	originalParams := originalSignature.Params()

	if originalParams.Len() == 3 {
		// rebuild *types.Func _without_ *http.Req param
		var vars []*types.Var
		for i := 1; i < originalParams.Len(); i++ {
			param := originalParams.At(i)
			vars = append(vars, param)
		}

		params := types.NewTuple(vars...)
		signature := types.NewSignature(originalSignature.Recv(), params,
			originalSignature.Results(), originalSignature.Variadic())
		return types.NewFunc(method.Pos(), method.Pkg(), method.Name(), signature)
	}

	return method
}
