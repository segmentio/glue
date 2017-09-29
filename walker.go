package glue

import (
	"errors"
	"fmt"
	"sync"

	"github.com/segmentio/glue/log"
	"github.com/segmentio/glue/provider"
	"github.com/segmentio/glue/writer"
	"golang.org/x/tools/go/loader"
)

// A Walker walks along supplied directions, visits server RPC code, and
// generates RPC client code along the way with the help of others.
type Walker struct {
	// Provider answers RPC-implementation-specific (e.g. stl, gorilla, etc.) questions.
	Provider provider.Provider
	Writer   writer.Writer
}

// Directions tell the Walker where to walk and what to pay attention to along the way.
type Directions struct {
	// Path determines the source code path to walk along. It's required
	Path string
	// Name is the name of the RPC declaration (e.g. `type MathService struct{}`).
	Name string
	// Service is the name of the RPC service. (e.g. `Math` in `Math.Sum`)
	Service string
}

// Walk is the logical entrypoint for Glue. It walks the source code and asks
func (w *Walker) Walk(directions Directions) error {
	var conf loader.Config
	conf.Import(directions.Path)

	prgm, err := conf.Load()
	if err != nil {
		log.Printf("failed to parse Go code: %s", err.Error())
		return err
	}

	var wg sync.WaitGroup
	for _, pkg := range prgm.InitialPackages() {
		wg.Add(1)
		go func(p *loader.PackageInfo) {
			defer wg.Done()
			w.walkPackage(p, directions.Name, directions.Service)
		}(pkg)
	}

	wg.Wait()
	return nil
}

func (w *Walker) walkPackage(pkg *loader.PackageInfo, decl, service string) error {
	visitor := NewVisitor(VisitorConfig{
		Pkg:         pkg,
		Provider:    w.Provider,
		Declaration: decl,
	})
	funcsByRecv := visitor.Go()

	if len(funcsByRecv) == 0 {
		log.Print("could not find RPC declaration")
		return errors.New("not found")
	}

	for _, funcs := range funcsByRecv {
		generator := Generator{
			Provider: w.Provider,
		}

		ident := fmt.Sprintf("%sClient", service)
		src, err := generator.Generate(GenerateInput{
			PackageName: "client",
			Service:     service,
			Funcs:       funcs,
		})
		if err != nil {
			return err
		}

		fname := fmt.Sprintf("generated_%s.go", ident)
		if err := w.Writer.Write(fname, src); err != nil {
			return err
		}

		log.Printf("glue: generated %s", ident)
	}

	return nil
}
