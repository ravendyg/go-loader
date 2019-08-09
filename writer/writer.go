package writer

import (
	"loader/dto"
	"os"
)

var tmpPostfix = ".tmp"
var dataPostfix = ".info"

// FileWriter - writer interface
type FileWriter struct {
	fileName string
	file     *os.File
}

// NewFileWriter - create new writer
func NewFileWriter(url string, fileName string) (*FileWriter, error) {
	file, err := os.Create(fileName + tmpPostfix)
	if err != nil {
		return nil, err
	}

	return &FileWriter{fileName, file}, nil
}

// WriteData - write data
func (fw *FileWriter) WriteData(chunk *dto.Chunk) (int, error) {
	l, err := fw.file.WriteAt(chunk.Data, chunk.Start+chunk.Offset)
	if err != nil {
		return 0, err
	}

	return l, nil
}

// Finish - clean up
func (fw *FileWriter) Finish() error {
	err := fw.file.Close()
	if err != nil {
		return err
	}

	err = os.Rename(fw.fileName+tmpPostfix, fw.fileName)
	if err == nil {
		return err
	}
	return nil
}
