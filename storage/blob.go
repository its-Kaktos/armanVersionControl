package storage

import (
	"encoding/binary"
	"fmt"
	"slices"
)

const (
	// currentBlobVersion represents the latest (current) version of Blob.
	currentBlobVersion uint16 = 0
	// blobMagicNumber represents the Blob unique identifier.
	blobMagicNumber uint16 = 100
)

// currentBlobSignature represents the latest (current) signature of Blob.
var currentBlobSignature []byte

// currentBlobHeader represents the first few bytes of the file representation of
// a Blob. If any file starts with this header, we will know it's a Blob.
var currentBlobHeader []byte

func init() {
	currentBlobSignature = make([]byte, 2)
	// BigEndian is chosen because that is the network byte order
	// and will save few bytes when storing it in the file. Plus
	// that's how git represents numbers in the file as well.
	_, err := binary.Encode(currentBlobSignature, binary.BigEndian, treeMagicNumber+currentTreeVersion)
	if err != nil {
		panic(err)
	}

	currentBlobHeader = []byte(fmt.Sprintf("%v \u0000", currentBlobSignature))
}

// Blob is a binary large object which represents the contents of file.
// The signature value of a Blob ranges from 100 to 199.
// When a file's content starts with "100," it indicates that the file
// is a Blob. The remaining two digits represent the version
// of the Blob structure. For example, a Signature value of 121
// indicates that this file is a Blob with a structure version of 21.
type Blob struct {
	Content []byte
}

// IsBlobS checks whether the signature is a Blob signature.
func IsBlobS(signature uint16) bool {
	return signature >= 100 && signature <= 199
}

func IsBlobC(content []byte) bool {
	return slices.Equal(content[:5], currentBlobHeader)
}

func (b Blob) FileRepresent() []byte {
	// Use slices.Clone to ensure the returned slice has a separate underlying array
	// from currentBlobHeader. This prevents modifications to the returned slice
	// from affecting currentBlobHeader.
	return append(slices.Clone(currentBlobHeader), b.Content...)
}
