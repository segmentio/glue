package generator

// TemplateData structures input to the template/client.gohtml template.
type TemplateData struct {
	// Package is the name of the output package.
	Package string
	// Service is the name of the service.
	Service string
	// Imports is a list of package paths to import.
	Imports []Import
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

type Import struct {
	Name string
	Path string
}
