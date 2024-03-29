package dto

// ChunkDescriptor - chunk download state
type ChunkDescriptor struct {
	ID     int
	Start  int64
	End    int64
	Offset int64
}

// Chunk - DTO
type Chunk struct {
	ID     int
	Cursor int64
	Size   int64
	Data   []byte
}

// ProcessDescriptor - file download state
type ProcessDescriptor struct {
	URL              string
	FileName         string
	ChunkDescriptors []ChunkDescriptor
	Size             int64
}
