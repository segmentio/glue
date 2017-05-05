package writer

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
)

type FileWriter struct {
	baseDir string
}

func NewFileWriter(dir string) (*FileWriter, error) {
	if dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.WithError(err).Error("failed to create leading directories")
			return nil, err
		}
	}

	return &FileWriter{baseDir: dir}, nil
}

func (fw *FileWriter) Write(path string, data []byte) error {
	f, err := os.Create(filepath.Join(fw.baseDir, path))
	if err != nil {
		log.WithError(err).Error("failed to create file")
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		log.WithError(err).Error("failed to write to file")
		return err
	}

	return nil
}
