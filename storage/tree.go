package storage

import (
	"encoding/binary"
	"fmt"
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

	// TODO fix this comment.
	// Because go file var are initialized first and then init func is run,
	// I cant put this line where the currentTreeHeader is located. Otherwise
	// the treeSignature will be nil when currentTreeHeader is initialized.
	currentTreeHeader = []byte(fmt.Sprintf("%v \u0000", currentTreeSignature))
}

// TreeEntry represents a single entry in a Tree structure.
// Each entry can either be a subdirectory (Tree) or a file (Blob), but not
// both. The Kind field specifies what type this entry holds.
type TreeEntry struct {
	// Kind specifies type of the entry.
	Kind EntryKind
	// EntryHash is the hash of the Blob's or the Tree's.
	EntryHash string
	// NameSize represents the size of the Name in bytes.
	NameSize uint16
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
	Entries []TreeEntry
}
