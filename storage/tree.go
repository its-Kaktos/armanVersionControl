package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

type EntryKind int32

const (
	KindTree = iota
	KindBlob
)

// TODO have to change all these uint16 to binary representation when storing in the file
// TODO calculate total header size and update README.md
// TODO add content size to Tree and Blob
/*
TODO
- How can I calculate the Struct size? Do I need to calculate it? I know that
	when I need to store strings I need to know the size of that string to know
	how far of the file content is for that string.
- Add content size to a Blob and Tree for every Tree entry.
- Should I store a single pointer that points either to a Blob or a Tree to avoid
	having a null followed for Tree followed by a pointer to a Blob?
- Also, when storing TreeEntry in a file, there would be no pointer, only a string that
	stores the sha1 hash of a Tree or a Blob. Should I Change this?
- Add a size for any string in Tree or Blob such as Name in the TreeEntry.
- After all the above is done, create doc for it in the README.md and create new README
	files if necessary.
*/

const (
	// currentTreeVersion represents the latest (current) version of Tree.
	currentTreeVersion uint16 = 0
	// treeMagicNumber represents the Tree unique identifier.
	treeMagicNumber = 200
)

// currentTreeSignature represents the latest (current) signature of Tree.
var currentTreeSignature []byte

// currentTreeHeader represents the first few bytes of the file representation
// of a Tree. If any file starts with this header, we will know it's a Tree.
var currentTreeHeader []byte

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
	// Tree is accessible when Kind is KindTree
	Tree *Tree
	// Blob is accessible when Kind is KindBlob
	Blob *Blob
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
	if len(content) < 5 {
		return false
	}

	b := content[:5]
	return slices.Equal(b, currentTreeHeader) ||
		IsTreeS(binary.BigEndian.Uint16(b))
}

// TODO after implementation, call this func and store the result in NewTreeeFromPath
// TODO add doc
func (t Tree) FileRepresent() ([]byte, error) {
	var buf bytes.Buffer

	buf.Write(currentTreeHeader)

	for _, te := range t.Entries {
		err := binary.Write(&buf, binary.BigEndian, te.Kind)
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

		cd := te.CreatedDate.String()
		err = binary.Write(&buf, binary.BigEndian, int32(len(cd)))
		if err != nil {
			return nil, err
		}
		buf.WriteString(cd)

		md := te.ModifiedDate.String()
		err = binary.Write(&buf, binary.BigEndian, int32(len(md)))
		if err != nil {
			return nil, err
		}
		buf.WriteString(md)

		buf.WriteRune('\n')
	}

	return buf.Bytes(), nil
}

// TODO refacotre to not store all blobs and objects everytime this func is called, after that update the hash-object cmd docs
// TODO add doc
func NewTreeFromPath(name string) (Tree, error) {
	dir, err := os.ReadDir(name)
	if err != nil {
		if os.IsNotExist(err) {
			return Tree{}, nil
		}

		return Tree{}, err
	}

	tree := Tree{}
	for _, d := range dir {
		if strings.HasPrefix(d.Name(), ".") {
			// Ignore hidden files and directories
			continue
		}

		te := TreeEntry{
			CreatedDate:  time.Now(),
			ModifiedDate: time.Now(),
		}

		if d.Type().IsDir() {
			t, err := NewTreeFromPath(path.Join(name, d.Name()))
			if err != nil {
				return Tree{}, err
			}

			te.Kind = KindTree
			te.Tree = &t

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
		te.Name = d.Name()
		te.Blob = &Blob{Content: c}

		tree.Entries = append(tree.Entries, &te)
	}

	return tree, nil
}
