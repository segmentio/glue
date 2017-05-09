# Glue

Glue generates client code for your Go RPC server. It currently supports
- [net/rpc]
- [gorilla/rpc]

**Status:** Glue is still extremely immature. It needs tests, feedback, and more.
Still, it's functional so try it out and contribute! Just don't be surprised by
breaking changes.


## Installation

`go get github.com/tejasmanohar/glue/cmd/glue`

Then, `glue` should be available at `$GOPATH/bin/glue` (ideally, in your `$PATH`).


## Usage

`glue -name=Service -service=Math [path]` will traverse the provided path (or working
directory if none is provided) and generate clients for RPC methods
(pointed at `Math.*`) declared on `type Service`.

Given the following is in a `*.go` file in your working directory,

```go
package math

//go:generate glue -name Service -service Math
type Service struct{}

type SumArg struct {
	Values []int
}

type SumReply struct {
	Sum int
}

func (s *Service) Sum(arg SumArg, reply *SumReply) error {
	for _, v := range arg.Values {
		reply.Sum += v
	}

	return nil
}
```

`go generate` would output the following to `clients/Service.go`

```go
package client

import (
	"github.com/tejasmanohar/glue/example/stl/math"
)

type Client interface {
	Call(method string, args interface{}, reply interface{}) error
}

func NewMathClient(rpcClient Client) *Math {
	c := new(Math)
	c.RPC = rpcClient
	return c
}

type Math struct {
	RPC Client
}

func (c *Math) Sum(args math.SumArg) (*math.SumReply, error) {
	reply := new(math.SumReply)
	err := c.RPC.Call("Math.Sum", args, reply)
	return reply, err
}
```

## Gorilla

If you use [gorilla/rpc], you're in luck! Just specify `-gorilla`.


## Options

### Output
Glue always outputs code with `client` package. By default, this is in `./client`, but
you can change the output directory via `-out`.

To output code to STDOUT instead of files, supply `-print`.


## FAQ

### How do I use Glue with RPC implementation X?
Glue is modular. If you'd like support for another popular (or interesting, well-maintained)
RPC implementation, open a PR to add a new Glue `provider/`.

Unfortunately, Go doesn't allow dynamic loading of packages so if you'd like Glue
to support an internal or experimental RPC framework, fork Glue and supply another
`provider` in [cmd/glue/main.go](https://github.com/tejasmanohar/glue/blob/master/cmd/glue/main.go).


[net/rpc]: https://golang.org/pkg/net/rpc/
[gorilla/rpc]: github.com/gorilla/rpc
