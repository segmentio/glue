package writer

type Writer interface {
	Write(path string, data []byte) error
}
