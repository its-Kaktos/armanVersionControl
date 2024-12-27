package track

import (
	"armanVersionControl/storage"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"syscall"
	"time"
)

const (
	// currentIndexVersion represents the latest (current) version of Index.
	currentIndexVersion uint16 = 0
	// indexMagicNumber represents the Index unique identifier.
	indexMagicNumber = 400
)

var (
	// currentIndexSignature represents the latest (current) signature of Index.
	currentIndexSignature []byte

	// currentIndexHeader represents the first few bytes of the file representation
	// of an Index. If any file starts with this header, we will know it's an Index.
	currentIndexHeader []byte
	indexFileName                  = path.Join(storage.MainDir, "index")
	filePerm           os.FileMode = 0770
)

var (
	ErrNotAnIndex        = errors.New("not a valid Index")
	ErrIndexNotFound     = errors.New("index file not found")
	ErrNameAlreadyExists = errors.New("path already exits in index")
)

func init() {
	currentIndexSignature = make([]byte, 2)
	// BigEndian is chosen because that is the network byte order
	// and will save few bytes when storing it in the file. Plus
	// that's how git represents numbers in the file as well.
	_, err := binary.Encode(currentIndexSignature, binary.BigEndian, indexMagicNumber+currentIndexVersion)
	if err != nil {
		panic(err)
	}

	// Since Go file-level variables are initialized before the init function,
	// I can't place this line where currentIndexHeader is defined.
	// Otherwise, currentIndexSignature will be nil when currentIndexHeader is initialized.
	// git adds a null byte in the header before the content starts. I don't think I'll need it,
	// but oh well, who cares if there's a null byte in the header even if I don't need it?
	currentIndexHeader = []byte(fmt.Sprintf("%v \u0000", currentIndexSignature))
}

// IndexEntry represents each entry in Index, which can only be a regular file.
type IndexEntry struct {
	// EntryHash is the hash of the Blob.
	EntryHash string
	// Name represents the file name.
	Name string
	// CreatedDate represents the creation date time of a Blob.
	CreatedDate time.Time
	// ModifiedDate represents the last date time of when IndexEntry was changed.
	ModifiedDate time.Time
}

// Index represents a type to track blobs.
type Index struct {
	// Entries are files that are being tracked.
	Entries []IndexEntry
}

// IsIndexS checks whether the signature is an Index signature.
func isIndexS(signature uint16) bool {
	return signature >= 400 && signature <= 499
}

// IsIndexB checks whether the content starts with the correct Index
// header (AKA signature).
func isIndexB(content []byte) bool {
	if len(content) < len(currentIndexHeader) {
		return false
	}

	b := content[:len(currentIndexHeader)]
	return slices.Equal(b, currentIndexHeader) ||
		isIndexS(binary.BigEndian.Uint16(b))
}

func (index Index) fileRepresent() ([]byte, error) {
	var buf bytes.Buffer

	buf.Write(currentIndexHeader)

	for _, ie := range index.Entries {
		err := binary.Write(&buf, binary.BigEndian, int32(len(ie.EntryHash)))
		if err != nil {
			return nil, err
		}
		buf.WriteString(ie.EntryHash)

		err = binary.Write(&buf, binary.BigEndian, int32(len(ie.Name)))
		if err != nil {
			return nil, err
		}
		buf.WriteString(ie.Name)

		cd, err := ie.CreatedDate.MarshalBinary()
		if err != nil {
			return nil, err
		}
		err = binary.Write(&buf, binary.BigEndian, int32(len(cd)))
		if err != nil {
			return nil, err
		}
		buf.Write(cd)

		md, err := ie.ModifiedDate.MarshalBinary()
		err = binary.Write(&buf, binary.BigEndian, int32(len(md)))
		if err != nil {
			return nil, err
		}
		buf.Write(md)
	}

	return buf.Bytes(), nil
}

func newIndexFromB(b []byte) (Index, error) {
	if !isIndexB(b) {
		return Index{}, ErrNotAnIndex
	}
	r := bytes.NewReader(b[len(currentIndexHeader):])

	index := Index{}
	readBuf := func() ([]byte, error) {
		countBuf := make([]byte, 4)
		_, err := r.Read(countBuf)
		if err != nil {
			return nil, err
		}

		count := int32(binary.BigEndian.Uint32(countBuf))
		buf := make([]byte, count)
		_, err = r.Read(buf)
		if err != nil {
			return nil, err
		}

		return buf, nil
	}
	for r.Len() > 0 {
		ie := IndexEntry{}

		// Parsing EntryHash
		buf, err := readBuf()
		if err != nil {
			return Index{}, err
		}
		ie.EntryHash = string(buf)

		// Parsing Name
		buf, err = readBuf()
		if err != nil {
			return Index{}, err
		}
		ie.Name = string(buf)

		// Parsing CreatedDate
		buf, err = readBuf()
		if err != nil {
			return Index{}, err
		}
		ie.CreatedDate = time.Now()
		if err = ie.CreatedDate.UnmarshalBinary(buf); err != nil {
			return Index{}, err
		}

		// Parsing ModifiedDate
		buf, err = readBuf()
		if err != nil {
			return Index{}, err
		}
		ie.ModifiedDate = time.Now()
		if err = ie.ModifiedDate.UnmarshalBinary(buf); err != nil {
			return Index{}, err
		}

		index.Entries = append(index.Entries, ie)
	}

	return index, nil
}

// Add will take a name (path) and add it to the current Index to prepare
// the content to be commited. If name is a directory, all subdirectories
// and files in the name will be added to the current Index. If name is an
// empty directory, nothing will be added to Index.
func Add(name string) error {
	s, err := os.Stat(name)
	if err != nil {
		return err
	}

	add := func(n string) error {
		i, err := fetchIndex()
		if err != nil {
			if !errors.Is(err, ErrIndexNotFound) {
				return err
			}
		}
		exists := slices.ContainsFunc(i.Entries, func(ie IndexEntry) bool {
			return ie.Name == n
		})

		if exists {
			return ErrNameAlreadyExists
		}

		rf, err := os.ReadFile(n)
		if err != nil {
			return err
		}

		s, err := os.Stat(n)
		if err != nil {
			return err
		}

		cd := time.Now()
		// For Unix-like systems, we need to use the Sys() method
		// to retrieve platform-specific information.
		if stat, ok := s.Sys().(*syscall.Stat_t); ok {
			cd = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
		}
		e := IndexEntry{
			EntryHash:    storage.ComputeHash(rf),
			Name:         n,
			CreatedDate:  cd,
			ModifiedDate: cd,
		}
		i.Entries = append(i.Entries, e)

		return i.saveIndex()
	}

	if s.Mode().IsRegular() {
		// Read index file and add current value to it.
		return add(name)
	}

	if s.IsDir() {
		return filepath.WalkDir(name, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.Type().IsRegular() {
				// Skip anything that is not a regular file
				return nil
			}

			return add(path)
		})
	}

	return fmt.Errorf("type %v is not supported", s.Mode().Type())
}

// Remove will remove name from Index.
func Remove(name string) error {
	index, err := fetchIndex()
	if err != nil {
		return err
	}

	e := slices.DeleteFunc(index.Entries, func(entry IndexEntry) bool {
		return entry.Name == name
	})

	// Noting was removed from entries because IndexEntry was not found.
	if len(e) == len(index.Entries) {
		return fmt.Errorf("not found")
	}

	index.Entries = e
	return index.saveIndex()
}

// fetchIndex will retrieve Index from the index file stored in
// avc repository.
func fetchIndex() (Index, error) {
	ok, err := storage.ExistsMainDir()
	if err != nil {
		return Index{}, err
	}

	if !ok {
		return Index{}, storage.ErrRepoNotInitialized
	}

	rf, err := os.ReadFile(indexFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return Index{}, ErrIndexNotFound
		}

		return Index{}, err
	}

	return newIndexFromB(rf)
}

// saveIndex persists the current state of the Index to a file.
func (index Index) saveIndex() error {
	ok, err := storage.ExistsMainDir()
	if err != nil {
		return err
	}
	if !ok {
		return storage.ErrRepoNotInitialized
	}

	b, err := index.fileRepresent()
	if err != nil {
		return err
	}

	return os.WriteFile(indexFileName, b, filePerm)
}
