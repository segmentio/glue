package client

import (
	"github.com/segmentio/glue/client"

	"github.com/segmentio/glue/example/stl/math"

	math1 "github.com/segmentio/glue/example/stl/math/math"
)

func NewMathClient(rpcClient client.Client) *Math {
	c := new(Math)
	c.RPC = rpcClient
	return c
}

type Math struct {
	RPC client.Client
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

func (c *Math) Abs(args math1.AbsArg) (*float64, error) {
	reply := new(float64)
	err := c.RPC.Call("Math.Abs", args, reply)
	return reply, err
}
