package storage

import (
	"time"
)

// TODO In the future, we need some abstraction on file structures to have some similar functionality.

// currentCommitVersion represents the latest (current) version of Commit.
const currentCommitVersion uint16 = 0

// currentCommitSignature represents the latest (current) signature of Commit.
const currentCommitSignature uint16 = 100 + currentCommitVersion

// Commit represents the structure of a basic commit.
// The signature value of a Commit ranges from 100 to 199.
// When a file's content starts with "100," it indicates that the file
// is a simple commit. The remaining two digits represent the version
// of the commit structure. For example, a Signature value of 121
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
	// Content represents the data of this commit, which can
	// have its own structure.
	Content []byte
}

// IsCommit checks whether the signature is a Commit signature.
func IsCommit(signature uint16) bool {
	return signature >= 100 && signature <= 199
}

// IsRoot will check whether c is a root commit
func (c Commit) IsRoot() bool {
	return c.ParentHash == ""
}

// New will create a new Commit.
func New(parentHash string, author string, authorEmail string, commiter string,
	commiterEmail string, commitDate time.Time, content []byte) *Commit {
	return &Commit{
		ParentHash:    parentHash,
		Author:        author,
		AuthorEmail:   authorEmail,
		Commiter:      commiter,
		CommiterEmail: commiterEmail,
		CommitDate:    commitDate,
		Content:       content,
	}
}
