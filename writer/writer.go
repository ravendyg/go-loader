package writer

import (
	"bytes"
	"encoding/gob"
	"loader/dto"
	"os"
)

var tmpPostfix = ".tmp"
var dataPostfix = ".info"

// FileWriter - writer interface
type FileWriter struct {
	fileName string
	file     *os.File
	dataFile *os.File
}

func encode(d dto.ProcessDescriptor) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(d)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decode(st []byte) (*dto.ProcessDescriptor, error) {
	var m dto.ProcessDescriptor
	b := bytes.Buffer{}
	b.Write(st)
	d := gob.NewDecoder(&b)
	err := d.Decode(m)
	return &m, err
}

// NewFileWriter - create new writer
func NewFileWriter(url string, fileName string) (*FileWriter, error) {
	file, err := os.Create(fileName + tmpPostfix)
	if err != nil {
		return nil, err
	}
	dataFile, err := os.Create(fileName + dataPostfix)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		fileName,
		file,
		dataFile,
	}, nil
}

// WriteData - write data
func (fw *FileWriter) WriteData(chunk *dto.Chunk) (int, error) {
	l, err := fw.file.WriteAt(chunk.Data, chunk.Start+chunk.Offset)
	if err != nil {
		return 0, err
	}

	return l, nil
}

// WriteMetaData - write information about the progress
func (fw *FileWriter) WriteMetaData(descriptor *dto.ProcessDescriptor) (int, error) {
	bt, err := encode(*descriptor)
	l, err := fw.dataFile.Write(bt)
	if err != nil {
		return 0, err
	}

	return l, nil
}

// ReadMeatada - read information about an interrupted process
func (fw *FileWriter) ReadMeatada() (*dto.ProcessDescriptor, error) {
	var b []byte
	_, err := fw.dataFile.Read(b)
	if err != nil {
		return nil, err
	}

	return decode(b)
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
