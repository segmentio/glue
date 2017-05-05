package writer

import (
	"bufio"
	"os"

	"github.com/apex/log"
)

type StdoutWriter struct {
	Writer *bufio.Writer
}

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{
		Writer: bufio.NewWriter(os.Stdout),
	}
}

func (s *StdoutWriter) Write(_path string, data []byte) error {
	if _, err := s.Writer.Write(data); err != nil {
		log.WithError(err).Error("failed to buffer write")
		return err
	}

	if err := s.Writer.Flush(); err != nil {
		log.WithError(err).Error("failed to flush to stdout")
		return err
	}

	return nil
}
