package track

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// TODO when add is called add the path to Tree and save Tree to a file named index.

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

// Add will take a name (path) and add it to the current Index to prepare
// the content to be commited. If name is a directory, all subdirectories
// and files in the name will be added to the current Index. If name is an
// empty directory, nothing will be added to Index.
func Add(name string) error {
	s, err := os.Stat(name)
	if err != nil {
		return err
	}
	
	if s.Mode().IsRegular() {
		// Add to Index
		return 
	}

	if s.IsDir() {
		// TODO handle with or without WalkDir func?
		filepath.WalkDir(name, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			
			d.IsDir() {
				
			}
		})
	}

	return nil
}
