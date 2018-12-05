package generator

import (
	"bytes"
	"go/types"
	"log"

	"github.com/segmentio/glue/provider"

	"golang.org/x/tools/imports"
)

type GenerateInput struct {
	Provider    provider.Provider
	PackageName string
	Service     string
	UseContext  bool

	Funcs []*types.Func
}

func Generate(in GenerateInput) ([]byte, error) {
	data := TemplateData{
		Package:    in.PackageName,
		Service:    in.Service,
		UseContext: in.UseContext,
	}

	resolver := newResolver()
	for _, f := range in.Funcs {
		argT := in.Provider.GetArgType(f)
		replyT := in.Provider.GetReplyType(f)

		data.Methods = append(data.Methods, MethodTemplate{
			Name:      f.Name(),
			ArgType:   resolver.GetTypeString(argT),
			ReplyType: resolver.GetTypeString(replyT),
		})
	}

	data.Imports = resolver.GetImports()

	var src bytes.Buffer
	err := tmpl.Execute(&src, data)
	if err != nil {
		log.Printf("failed to render template: %s", err.Error())
		return nil, err
	}

	formatted, err := imports.Process("client.go", src.Bytes(), nil)
	if err != nil {
		log.Printf("failed to format code: \nCODE:\n%s \nERR:\n%s", src.String(), err.Error())
		return nil, err
	}

	return formatted, err
}
