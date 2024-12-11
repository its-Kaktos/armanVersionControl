package storage

// FileRepresenter defines an interface for types that have a file representation.
// Currently, those types are Blob, Tree, Commit.
type FileRepresenter interface {
	// FileRepresent returns the content as a slice of bytes.
	FileRepresent() []byte
}
