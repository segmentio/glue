package client

import (
	"github.com/segmentio/glue/client"
	"github.com/segmentio/glue/example/gorilla/math"
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

func (c *Math) IdentityMany(args []int) (*[]int, error) {
	reply := new([]int)
	err := c.RPC.Call("Math.IdentityMany", args, reply)
	return reply, err
}

func (c *Math) IdentityManyStruct(args []math.IdentityStruct) (*[]math.IdentityStruct, error) {
	reply := new([]math.IdentityStruct)
	err := c.RPC.Call("Math.IdentityManyStruct", args, reply)
	return reply, err
}
