package structures

import (
	"armanVersionControl/storage"
	"encoding/binary"
	"errors"
	"fmt"
	"slices"
)

const (
	// currentBlobVersion represents the latest (current) version of Blob.
	currentBlobVersion uint16 = 0
	// blobMagicNumber represents the Blob unique identifier.
	blobMagicNumber uint16 = 100
)

var (
	// currentBlobSignature represents the latest (current) signature of Blob.
	currentBlobSignature []byte

	// currentBlobHeader represents the first few bytes of the file representation of
	// a Blob. If any file starts with this header, we will know it's a Blob.
	currentBlobHeader []byte
)

var (
	ErrNotABlob = errors.New("not a valid Blob")
)

func init() {
	currentBlobSignature = make([]byte, 2)
	// BigEndian is chosen because that is the network byte order
	// and will save few bytes when storing it in the file. Plus
	// that's how git represents numbers in the file as well.
	_, err := binary.Encode(currentBlobSignature, binary.BigEndian, blobMagicNumber+currentBlobVersion)
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
	// Content represents the content to store in Blob.
	Content []byte
}

// IsBlobS checks whether the signature is a Blob signature.
func IsBlobS(signature uint16) bool {
	return signature >= 100 && signature <= 199
}

// IsBlobB checks whether the content starts with the correct Blob
// header (AKA signature).
func IsBlobB(content []byte) bool {
	if len(content) < len(currentBlobHeader) {
		return false
	}

	b := content[:9]
	return slices.Equal(b, currentBlobHeader) ||
		IsBlobS(binary.BigEndian.Uint16(b[:len(currentBlobHeader)]))
}

func (b Blob) FileRepresent() []byte {
	// Use slices.Clone to ensure the returned slice has a separate underlying array
	// from currentBlobHeader. This prevents modifications to the returned slice
	// from affecting currentBlobHeader.
	return append(slices.Clone(currentBlobHeader), b.Content...)
}

func NewBlobFromB(b []byte) (Blob, error) {
	if !IsBlobB(b) {
		return Blob{}, ErrNotABlob
	}

	return Blob{Content: slices.Clone(b[len(currentBlobHeader):])}, nil
}

// StoreBlob will store Blob in the avc object store.
// Returns the hash of Blob when stored in avc repository.
func (b Blob) StoreBlob() (string, error) {
	h, err := storage.Store(b.FileRepresent())

	// If error is ObjectDuplicateError, reuse the previous
	// object hash instead of creating a new object.
	var ode *storage.ObjectDuplicateError
	if errors.As(err, &ode) {
		return ode.Hash, nil
	}

	return h, err
}
