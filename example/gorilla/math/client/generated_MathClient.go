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

type MathIFace interface {
	Sum(args math.SumArg) (math.SumReply, error)

	Identity(args int) (int, error)

	IdentityMany(args []int) ([]int, error)

	IdentityManyStruct(args []*math.IdentityStruct) ([]math.IdentityStruct, error)

	MapOfPrimitives(args map[string]string) ([]int, error)
}

type Math struct {
	RPC client.Client
}

func (c *Math) Sum(args math.SumArg) (math.SumReply, error) {
	var reply math.SumReply
	err := c.RPC.Call("Math.Sum", args, reply)
	return reply, err
}

func (c *Math) Identity(args int) (int, error) {
	var reply int
	err := c.RPC.Call("Math.Identity", args, reply)
	return reply, err
}

func (c *Math) IdentityMany(args []int) ([]int, error) {
	var reply []int
	err := c.RPC.Call("Math.IdentityMany", args, reply)
	return reply, err
}

func (c *Math) IdentityManyStruct(args []*math.IdentityStruct) ([]math.IdentityStruct, error) {
	var reply []math.IdentityStruct
	err := c.RPC.Call("Math.IdentityManyStruct", args, reply)
	return reply, err
}

func (c *Math) MapOfPrimitives(args map[string]string) ([]int, error) {
	var reply []int
	err := c.RPC.Call("Math.MapOfPrimitives", args, reply)
	return reply, err
}
