package storage

import "fmt"

type EntryKind int

const (
	KindTree = iota
	KindBlob
)

const (
	// currentTreeVersion represents the latest (current) version of Tree.
	currentTreeVersion uint16 = 0
	// currentTreeSignature represents the latest (current) signature of Tree.
	currentTreeSignature uint16 = 200 + currentTreeVersion
)

var currentTreeHeader = []byte(fmt.Sprintf("%v \u0000", currentTreeSignature))

// TreeEntry represents a single entry in a Tree structure.
// Each entry can either be a subdirectory (Tree) or a file (Blob), but not
// both. The Kind field specifies what type this entry holds.
type TreeEntry struct {
	// Kind specifies type of the entry.
	Kind EntryKind
	// Tree is a pointer to a Tree if Kind is KindTree, else is nil.
	Tree *Tree
	// Blob is a pointer to a Blob If Kind is KindBlob, else is nil.
	Blob *Blob
}

// Tree represents the structure of a directory and its files.
// Each entry in Entries represents a subdirectory (Tree) or a file (Blob).
type Tree struct {
	Entries []TreeEntry
}
