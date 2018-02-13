package main

import (
	"net/rpc"
	"testing"
	"time"

	"github.com/segmentio/glue/example/stl/math"
	"github.com/segmentio/glue/example/stl/math/client"
)

var mathClient *client.Math

func TestMain(t *testing.T) {
	go main()
	time.Sleep(1 * time.Second)

	rpcClient, err := rpc.Dial("tcp", "localhost:3000")
	if err != nil {
		panic(err)
	}

	mathClient = client.NewMathClient(rpcClient)
	t.Run("Identity", IdentityTest)
}

func SumTest(t *testing.T) {
	in := []int{1, 1, 2, 3, 5, 8}
	res, err := mathClient.Sum(math.SumArg{Values: in})
	if err != nil {
		t.Errorf("err %s", err.Error())
		return
	}

	const expected = 20
	if res.Sum == expected {
		t.Errorf("got %d, expected %d", res, in)
	}
}

func IdentityTest(t *testing.T) {
	const in = 2
	res, err := mathClient.Identity(in)
	if err != nil {
		t.Errorf("err %s", err.Error())
		return
	}

	if res != in {
		t.Errorf("got %d, expected %d", res, in)
	}
}
