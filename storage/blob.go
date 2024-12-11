package storage

import (
	"fmt"
	"slices"
)

const (
	// currentBlobVersion represents the latest (current) version of Blob.
	currentBlobVersion uint16 = 0
	// currentBlobSignature represents the latest (current) signature of Blob.
	currentBlobSignature = 100 + currentBlobVersion
)

var currentBlobHeader = []byte(fmt.Sprintf("%v \u0000", currentBlobSignature))

// Blob is a binary large object which represents the contents of file.
// The signature value of a Blob ranges from 100 to 199.
// When a file's content starts with "100," it indicates that the file
// is a Blob. The remaining two digits represent the version
// of the Blob structure. For example, a Signature value of 121
// indicates that this file is a Blob with a structure version of 21.
type Blob struct {
	Content []byte
}

// IsBlob checks whether the signature is a Blob signature.
func IsBlob(signature uint16) bool {
	return signature >= 100 && signature <= 199
}

func (b Blob) FileRepresent() []byte {
	// Use slices.Clone to ensure the returned slice has a separate underlying array
	// from currentBlobHeader. This prevents modifications to the returned slice
	// from affecting currentBlobHeader.
	return append(slices.Clone(currentBlobHeader), b.Content...)
}
