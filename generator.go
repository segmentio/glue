package glue

import (
	"bytes"
	"fmt"
	"go/types"
	"html/template"

	"golang.org/x/tools/imports"

	"github.com/apex/log"
	"github.com/tejasmanohar/glue/internal/gen"
	"github.com/tejasmanohar/glue/provider"
)

//go:generate go-bindata -nomemcopy -pkg glue templates/...
var tmpl = mustParseTemplate("templates/client.gohtml")

// Generator creates the output Golang net/rpc client code.
type Generator struct {
	Provider  provider.Provider
	SourcePkg string
}

// TemplateData structures input to the template/client.gohtml template.
type TemplateData struct {
	// Package is the name of the output package.
	Package string
	// Service is the name of the service.
	Service string
	// Imports is a list of package paths to import.
	Imports []string
	// Identifier is the name of the RPC client struct.
	Identifier string
	// Methods is a list of method metadata.
	Methods []MethodTemplate
}

// MethodTemplate describes the structure of an RPC method.
type MethodTemplate struct {
	// Name is the name of the RPC method.
	Name string
	// ArgType is the name of the RPC argument type (e.g. `string`).
	ArgType string
	// ReplyType is the name of the RPC response type (e.g. `string`).
	ReplyType string
}

// GenerateInput is used to generate code.
type GenerateInput struct {
	PackageName string
	Service     string

	Funcs []*types.Func
}

// Generate generates code in a streaming fashion.
func (g *Generator) Generate(in GenerateInput) ([]byte, error) {
	data := TemplateData{
		Package: in.PackageName,
		Service: in.Service,
	}

	pkgs := gen.NewStringSet()
	for _, f := range in.Funcs {
		data.Methods = append(data.Methods, MethodTemplate{
			Name:      f.Name(),
			ArgType:   g.Provider.GetArgType(f).Identifier(),
			ReplyType: g.Provider.GetReplyType(f).Identifier(),
		})

		fImports := g.getFuncImports(f)
		pkgs = pkgs.Union(fImports)
	}

	pkgs.Discard(g.SourcePkg)
	data.Imports = pkgs.AsList()

	var src bytes.Buffer
	err := tmpl.Execute(&src, data)
	if err != nil {
		log.WithError(err).Error("failed to render template")
		return nil, err
	}

	formatted, err := imports.Process(fmt.Sprintf("client.go", in.Service), src.Bytes(), nil)
	return formatted, err
}

func (g *Generator) getFuncImports(f *types.Func) *gen.StringSet {
	packages := gen.NewStringSet()
	signature := f.Type().(*types.Signature)

	params := signature.Params()
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		path := param.Pkg().Path()
		packages.Add(path)
	}

	results := signature.Results()
	for i := 0; i < results.Len(); i++ {
		res := results.At(i)
		path := res.Pkg().Path()
		packages.Add(path)
	}

	return packages
}

func mustParseTemplate(path string) *template.Template {
	data, err := Asset(path)
	if err != nil {
		panic(err)
	}

	return template.Must(
		template.New(path).Parse(string(data)),
	)
}
