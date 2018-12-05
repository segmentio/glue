package writer

import (
	"github.com/segmentio/glue/log"
)

type StdoutWriter struct{}

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

func (s *StdoutWriter) Write(_path string, data []byte) error {
	log.Print(string(data))
	return nil
}
