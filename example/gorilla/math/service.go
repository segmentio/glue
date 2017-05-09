package math

import (
	"net/http"
)

//go:generate glue -gorilla -name Service -service Math
type Service struct{}

type SumArg struct {
	Values []int
}

type SumReply struct {
	Sum int
}

func (s *Service) Sum(r *http.Request, arg *SumArg, reply *SumReply) error {
	for _, v := range arg.Values {
		reply.Sum += v
	}

	return nil
}

func (s *Service) Identity(r *http.Request, arg *int, reply *int) error {
	*reply = *arg
	return nil
}
