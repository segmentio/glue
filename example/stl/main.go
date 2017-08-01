package main

import (
	"net"
	"net/rpc"

	"github.com/segmentio/glue/example/stl/math"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:3000")
	if err != nil {
		panic(err)
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	if err := rpc.RegisterName("Math", new(math.Service)); err != nil {
		panic(err)
	}

	rpc.Accept(inbound)
}
