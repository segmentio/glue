package client

import "context"

type Client interface {
	Call(method string, args interface{}, reply interface{}) error
}

type ClientContext interface {
	Call(ctx context.Context, method string, args interface{}, reply interface{}) error
}
