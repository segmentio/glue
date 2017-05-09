package main

import (
	"net"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/tejasmanohar/glue/example/gorilla/math"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:4000")
	if err != nil {
		panic(err)
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	if err := server.RegisterService(new(math.Service), "Math"); err != nil {
		panic(err)
	}

	if err := http.Serve(inbound, server); err != nil {
		panic(err)
	}
}
