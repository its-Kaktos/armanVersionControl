package storage

import (
	"armanVersionControl/storage/objectstore"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

type EntryKind int32

func (e EntryKind) String() string {
	return []string{"KindTree", "KindBlob"}[e]
}

const (
	KindTree EntryKind = iota
	KindBlob
)

const (
	// currentTreeVersion represents the latest (current) version of Tree.
	currentTreeVersion uint16 = 0
	// treeMagicNumber represents the Tree unique identifier.
	treeMagicNumber = 200
)

var (
	// currentTreeSignature represents the latest (current) signature of Tree.
	currentTreeSignature []byte

	// currentTreeHeader represents the first few bytes of the file representation
	// of a Tree. If any file starts with this header, we will know it's a Tree.
	currentTreeHeader []byte
)

var (
	ErrNotATree = errors.New("not a valid Tree")
)

func init() {
	currentTreeSignature = make([]byte, 2)
	// BigEndian is chosen because that is the network byte order
	// and will save few bytes when storing it in the file. Plus
	// that's how git represents numbers in the file as well.
	_, err := binary.Encode(currentTreeSignature, binary.BigEndian, treeMagicNumber+currentTreeVersion)
	if err != nil {
		panic(err)
	}

	// Since Go file-level variables are initialized before the init function,
	// I can't place this line where currentTreeHeader is defined.
	// Otherwise, currentTreeSignature will be nil when currentTreeHeader is initialized.
	// git adds a null byte in the header before the content starts. I don't think I'll need it,
	// but oh well, who cares if there's a null byte in the header even if I don't need it?
	currentTreeHeader = []byte(fmt.Sprintf("%v \u0000", currentTreeSignature))
}

// TreeEntry represents a single entry in a Tree structure.
// Each entry can either be a subdirectory (Tree) or a file (Blob), but not
// both. The Kind field specifies what type this entry holds.
type TreeEntry struct {
	// Kind specifies type of the entry.
	Kind EntryKind
	// tree is accessible when Kind is KindTree
	tree *Tree
	// blob is accessible when Kind is KindBlob
	blob *Blob
	// EntryHash is the hash of the Blob's or the Tree's.
	EntryHash string
	// Name represents the file or directory name.
	Name string
	// CreatedDate represents the creation date time of a Blob or Tree.
	CreatedDate time.Time
	// ModifiedDate represents the last date time of when TreeEntry was changed.
	ModifiedDate time.Time
}

// Tree represents the structure of a directory and its files.
// Each entry in Entries represents a subdirectory (Tree) or a file (Blob).
type Tree struct {
	// Hash represents Tree hash (AKA filename) that is stored in avc object store.
	Hash string
	// Entries contain directory files and subdirectories.
	Entries []*TreeEntry
}

// IsTreeS checks whether the signature is a Tree signature.
func IsTreeS(signature uint16) bool {
	return signature >= 200 && signature <= 299
}

// IsTreeB checks whether the content starts with the correct Tree
// header (AKA signature).
func IsTreeB(content []byte) bool {
	if len(content) < len(currentTreeHeader) {
		return false
	}

	b := content[:len(currentTreeHeader)]
	return slices.Equal(b, currentTreeHeader) ||
		IsTreeS(binary.BigEndian.Uint16(b))
}

// FetchTree retrieves a Tree from the object database using the TreeEntry.EntryHash
// and caches the result to prevent redundant calculations on subsequent calls.
func (te *TreeEntry) FetchTree() (Tree, error) {
	if te.Kind != KindTree {
		return Tree{}, fmt.Errorf("expected kind to be %v but got %v", KindTree, te.Kind)
	}

	if te.tree != nil {
		return *te.tree, nil
	}

	tb, err := objectstore.FetchByHash(te.EntryHash)
	if err != nil {
		return Tree{}, err
	}

	t, err := NewTreeFromObject(tb)
	if err != nil {
		return Tree{}, err
	}

	te.tree = &t
	return t, nil
}

// FetchBlob retrieves a Blob from the object database using the TreeEntry.EntryHash
// and caches the result to prevent redundant calculation on subsequent calls.
func (te *TreeEntry) FetchBlob() (Blob, error) {
	if te.Kind != KindBlob {
		return Blob{}, fmt.Errorf("expected kind to be %v but got %v", KindBlob, te.Kind)
	}

	if te.blob != nil {
		return *te.blob, nil
	}

	o, err := objectstore.FetchByHash(te.EntryHash)
	if err != nil {
		return Blob{}, err
	}

	b, err := NewBlobFromB(o.Content)
	if err != nil {
		return Blob{}, err
	}

	te.blob = &b
	return b, nil
}

// FileRepresent will create a file representation of a Tree in binary format.
func (t *Tree) FileRepresent() ([]byte, error) {
	var buf bytes.Buffer

	buf.Write(currentTreeHeader)

	for _, te := range t.Entries {
		err := binary.Write(&buf, binary.BigEndian, int32(te.Kind))
		if err != nil {
			return nil, err
		}

		err = binary.Write(&buf, binary.BigEndian, int32(len(te.EntryHash)))
		if err != nil {
			return nil, err
		}
		buf.WriteString(te.EntryHash)

		err = binary.Write(&buf, binary.BigEndian, int32(len(te.Name)))
		if err != nil {
			return nil, err
		}
		buf.WriteString(te.Name)
	}

	return buf.Bytes(), nil
}

// NewTreeFromPath creates a new Tree form a path but does not store the result
// in object database. Because result is not stored in the database, there will be
// no hash for the Tree or its TreeEntry.
func NewTreeFromPath(name string) (Tree, error) {
	dir, err := os.ReadDir(name)
	if err != nil {
		return Tree{}, err
	}

	tree := Tree{}
	for _, d := range dir {
		if strings.HasPrefix(d.Name(), ".") {
			// Ignore hidden files and directories
			continue
		}

		te := TreeEntry{Name: d.Name()}

		if d.Type().IsDir() {
			t, err := NewTreeFromPath(path.Join(name, d.Name()))
			if err != nil {
				return Tree{}, err
			}

			te.Kind = KindTree
			te.tree = &t

			tree.Entries = append(tree.Entries, &te)
			continue
		}

		if !d.Type().IsRegular() {
			return Tree{}, fmt.Errorf("directory entries should either be a directory or regualr file which '%v' does not follow", path.Join(name, d.Name()))
		}

		c, err := os.ReadFile(path.Join(name, d.Name()))
		if err != nil {
			return Tree{}, err
		}

		te.Kind = KindBlob
		te.blob = &Blob{Content: c}

		tree.Entries = append(tree.Entries, &te)
	}

	return tree, nil
}

// NewTreeFromObject creates a Tree from objectstore.Object.
func NewTreeFromObject(o objectstore.Object) (Tree, error) {
	if !IsTreeB(o.Content) {
		return Tree{}, ErrNotATree
	}
	t := Tree{Hash: o.Hash}

	r := bytes.NewReader(o.Content[len(currentTreeHeader):])

	//pos := 0
	for r.Len() > 0 {
		te := TreeEntry{}

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

		// Parsing Kind
		intBuf := make([]byte, 4)
		_, err := r.Read(intBuf)
		if err != nil {
			return Tree{}, err
		}
		te.Kind = EntryKind(int32(binary.BigEndian.Uint32(intBuf)))

		// Parsing EntryHash
		buf, err := readBuf()
		if err != nil {
			return Tree{}, err
		}
		te.EntryHash = string(buf)

		// Parsing Name
		buf, err = readBuf()
		if err != nil {
			return Tree{}, err
		}
		te.Name = string(buf)

		t.Entries = append(t.Entries, &te)
	}

	return t, nil
}

// StoreTree will store Tree and all its Entries in the object store
// and return the computed hash for Tree.
func (t *Tree) StoreTree() (string, error) {
	for _, te := range t.Entries {
		if te.Kind == KindTree {
			teTree, err := te.FetchTree()
			if err != nil {
				return "", err
			}

			h, err := teTree.StoreTree()
			if err != nil {
				return "", err
			}

			te.EntryHash = h
			continue
		}

		if te.Kind != KindBlob {
			panic("A new unexpected kind detected.")
		}

		teBlob, err := te.FetchBlob()
		if err != nil {
			return "", err
		}
		h, err := teBlob.StoreBlob()
		if err != nil {
			return "", err
		}

		te.EntryHash = h
	}

	b, err := t.FileRepresent()
	if err != nil {
		return "", err
	}

	h, err := objectstore.Store(b)
	// Reuse the previous object of there is a duplicate error
	var ode *objectstore.ObjectDuplicateError
	if errors.As(err, &ode) {
		return ode.Hash, nil
	}

	return h, err
}

func (t *Tree) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tree hash: %v \n", t.Hash))
	sb.WriteString(fmt.Sprintf("Tree entrieis:\n"))

	for _, te := range t.Entries {
		sb.WriteString(fmt.Sprintf("Kind: %v ", te.Kind))
		sb.WriteString(fmt.Sprintf("Hash: %v ", te.EntryHash))
		sb.WriteString(fmt.Sprintf("name: %v ", te.Name))
		c, err := te.CreatedDate.MarshalText()
		if err != nil {
			panic(err)
		}
		sb.WriteString(fmt.Sprintf("created date: %v ", string(c)))
		m, err := te.ModifiedDate.MarshalText()
		if err != nil {
			panic(err)
		}
		sb.WriteString(fmt.Sprintf("modified date: %v ", string(m)))

		sb.WriteRune('\n')
	}

	return sb.String()
}
