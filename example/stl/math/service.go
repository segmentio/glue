package math

//go:generate glue -name Service -service Math
type Service struct{}

type SumArg struct {
	Values []int
}

type SumReply struct {
	Sum int
}

func (s *Service) Sum(arg SumArg, reply *SumReply) error {
	for _, v := range arg.Values {
		reply.Sum += v
	}

	return nil
}

func (s *Service) Identity(arg int, reply *int) error {
	*reply = arg
	return nil
}
