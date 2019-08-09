package dto

// ChunkDescriptor - chunk download state
type ChunkDescriptor struct {
	Start  int64
	End    int64
	Offset int64
}

// Chunk - DTO
type Chunk struct {
	ChunkDescriptor
	Data []byte
}

// ProcessDescriptor - file download state
type ProcessDescriptor struct {
	URL              string
	FileName         string
	ChunkDescriptors []ChunkDescriptor
	Size             int64
}
