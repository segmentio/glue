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

func (c *Math) Identity(args int) (*int, error) {
	reply := new(int)
	err := c.RPC.Call("Math.Identity", args, reply)
	return reply, err
}
