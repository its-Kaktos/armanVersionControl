package structures

import (
	"fmt"
	"time"
)

const (
	// currentCommitVersion represents the latest (current) version of Commit.
	currentCommitVersion uint16 = 0
	// currentCommitSignature represents the latest (current) signature of Commit.
	currentCommitSignature uint16 = 300 + currentCommitVersion
)

var currentCommitHeader = []byte(fmt.Sprintf("%v \u0000", currentCommitSignature))

// Commit represents the structure of a basic commit.
// The signature value of a Commit ranges from 300 to 399.
// When a file's content starts with "300", it indicates that the file
// is a simple commit. The remaining two digits represent the version
// of the commit structure. For example, a Signature value of 321
// indicates that this file is a basic commit with a structure version of 21.
type Commit struct {
	// ParentHash represents the hash of the
	// previous commit which this commit is based on.
	// In case of the root commit, ParentHash would be empty.
	ParentHash string
	// Author is the name of the author.
	Author string
	// AuthorEmail is the email of the author.
	AuthorEmail string
	// Commiter is the name of the commiter.
	Commiter string
	// CommiterEmail is the email of the commiter.
	CommiterEmail string
	// CommitDate is the date when this commit was created
	CommitDate time.Time
	// Tree represents the snapshot of the working directory
	// when this commit was created.
	Tree Tree
}

// IsCommit checks whether the signature is a Commit signature.
func IsCommit(signature uint16) bool {
	return signature >= 300 && signature <= 399
}

// IsRoot will check whether c is a root commit
func (c Commit) IsRoot() bool {
	return c.ParentHash == ""
}

// New will create a new Commit.
func New(parentHash string, author string, authorEmail string, commiter string,
	commiterEmail string, commitDate time.Time, tree Tree) *Commit {
	return &Commit{
		ParentHash:    parentHash,
		Author:        author,
		AuthorEmail:   authorEmail,
		Commiter:      commiter,
		CommiterEmail: commiterEmail,
		CommitDate:    commitDate,
		Tree:          tree,
	}
}
