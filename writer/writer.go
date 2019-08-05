package writer

import (
	"os"
	"strings"
)

// FileWriter - writer interface
type FileWriter struct {
	file *os.File
}

// NewFileWriter - create new writer
func NewFileWriter(url string) (*FileWriter, error) {
	nameChunks := strings.Split(url, "/")
	fileName := nameChunks[len(nameChunks)-1]
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return &FileWriter{file}, nil
}

func (fw *FileWriter) Write(data []byte, offset int64) (int, error) {
	l, err := fw.file.WriteAt(data, offset)
	if err != nil {
		return 0, err
	}

	return l, nil
}

// Close - close file
func (fw *FileWriter) Close() error {
	if fw.file == nil {
		return nil
	}
	err := fw.file.Close()
	return err
}
