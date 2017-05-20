package glue

import (
	"bytes"
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
	Provider provider.Provider
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
		argInfo := g.Provider.GetArgType(f)
		replyInfo := g.Provider.GetReplyType(f)
		data.Methods = append(data.Methods, MethodTemplate{
			Name:      f.Name(),
			ArgType:   argInfo.Identifier,
			ReplyType: replyInfo.Identifier,
		})

		pkgs.AddList(argInfo.Imports)
		pkgs.AddList(replyInfo.Imports)
	}

	data.Imports = pkgs.AsList()

	var src bytes.Buffer
	err := tmpl.Execute(&src, data)
	if err != nil {
		log.WithError(err).Error("failed to render template")
		return nil, err
	}

	formatted, err := imports.Process("client.go", src.Bytes(), nil)
	if err != nil {
		log.WithField("code", src.String()).WithError(err).Error("failed to format code")
		return nil, err
	}

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
