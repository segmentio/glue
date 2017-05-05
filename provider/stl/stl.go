package stl

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/apex/log"
	"github.com/tejasmanohar/glue/provider"
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
	if !isExportedOrBuiltin(arg.Type()) {
		log.Debugf("skipping %s: argument parameter's type is not exported", method.Name())
		return false
	}

	reply := params.At(1)
	replyType := reply.Type()
	if !isExportedOrBuiltin(replyType) {
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
	return getTypeInfo(arg.Type())
}

// GetReplyType extracts metadata about the response type from an RPC method.
func (p *Provider) GetReplyType(f *types.Func) provider.TypeInfo {
	signature := f.Type().(*types.Signature)
	params := signature.Params()
	reply := params.At(1)
	return getTypeInfo(reply.Type())
}

func getTypeInfo(t types.Type) provider.TypeInfo {
	t = unpack(t)

	switch specific := t.(type) {
	case *types.Named:
		return provider.TypeInfo{
			Name:    specific.Obj().Name(),
			Package: specific.Obj().Pkg(),
		}
	case *types.Basic:
		return provider.TypeInfo{
			Name: specific.Name(),
		}
	}

	panic(fmt.Errorf("unexpected type: %s", reflect.TypeOf(t)))
}

func isExportedOrBuiltin(t types.Type) bool {
	t = unpack(t)

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

func unpack(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}

var _ provider.Provider = (*Provider)(nil)
