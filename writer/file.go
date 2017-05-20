package writer

import (
	"os"
	"path/filepath"

	"github.com/tejasmanohar/glue/log"
)

type FileWriter struct {
	baseDir string
}

func NewFileWriter(dir string) (*FileWriter, error) {
	if dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Printf("failed to create output directory: %s", err.Error())
			return nil, err
		}
	}

	return &FileWriter{baseDir: dir}, nil
}

func (fw *FileWriter) Write(path string, data []byte) error {
	f, err := os.Create(filepath.Join(fw.baseDir, path))
	if err != nil {
		log.Printf("failed to create file: %s", err.Error())
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		log.Printf("failed to write to file: %s", err.Error())
		return err
	}

	return nil
}
