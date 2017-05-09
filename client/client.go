package client

type Client interface {
	Call(method string, args interface{}, reply interface{}) error
}
