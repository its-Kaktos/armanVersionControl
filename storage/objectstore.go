package storage

import (
	"armanVersionControl/hashing"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// TODO change file perm? should other users see this? in git all have read access only, why?
// TODO add pager instead of reading whole file?

// HashCollisionError represents an error for hash collisions.
type HashCollisionError struct {
	Collisions []string
}

func (h *HashCollisionError) Error() string {
	if len(h.Collisions) == 0 {
		return "hash collision detected, but no possible collisions were provided"
	}

	return fmt.Sprintf("hash collision detected. Possible matches:\n%s", strings.Join(h.Collisions, "\n"))
}

// ObjectDuplicateError represents an error for when an object
// with the same hash already exists in the object database.
type ObjectDuplicateError struct {
	// Hash is the object hash.
	Hash string
}

func (o *ObjectDuplicateError) Error() string {
	return fmt.Sprintf("object hash of %v already exists in the object database store", o.Hash)
}

var (
	ErrHashIsShort            = errors.New("provided hash is short, it should be at least 2 characters")
	ErrObjectNotFound         = errors.New("object not found")
	ErrDirectoryIsNotExpected = errors.New("directory is not expected in a directory of object database")
)

var (
	objectDir             = path.Join(MainDir, "objects")
	filePerm  os.FileMode = 0770
)

// Object represents any data in the object database.
type Object struct {
	// Hash is the object hash stored in the object database.
	Hash string
	// Content is the content of the file stored in the object database.
	Content []byte
}

// Store will save content in the object database.
func Store(content []byte) (hash string, e error) {
	ok, err := ExistsMainDir()
	if err != nil {
		return "", err
	}

	if !ok {
		return "", ErrRepoNotInitialized
	}

	hashHex := ComputeHash(content)

	dir := path.Join(objectDir, hashHex[:2])
	fileName := hashHex[2:]
	filePath := path.Join(dir, fileName)

	ok, err = objectExists(path.Join(dir, fileName))
	if err != nil {
		return "", err
	}

	if ok {
		return "", &ObjectDuplicateError{Hash: hashHex}
	}

	if err = mkdirAllIfDoesNotExists(dir, dirPerm); err != nil {
		return "", err
	}

	err = os.WriteFile(filePath, content, filePerm)
	if err != nil {
		return "", err
	}

	return hashHex, err
}

// ComputeHash will compute a hash based on the content and return
// the generated hash.
func ComputeHash(content []byte) string {
	return hex.EncodeToString(hashing.Sha1(content))
}

// FetchByHash will fetch an object from object database by its hash.
func FetchByHash(hash string) (Object, error) {
	ok, err := ExistsMainDir()
	if err != nil {
		return Object{}, err
	}
	if !ok {
		return Object{}, ErrRepoNotInitialized
	}

	if len(hash) < 2 {
		return Object{}, ErrHashIsShort
	}

	dirName := hash[:2]
	dirPath := path.Join(objectDir, dirName)
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			return Object{}, ErrObjectNotFound
		}

		return Object{}, err
	}

	objectsInDir, err := fetchAllFileNamesInDir(dirPath)
	if err != nil {
		return Object{}, err
	}
	if len(objectsInDir) == 0 {
		return Object{}, ErrObjectNotFound
	}
	readFile := func(name string) (Object, error) {
		rf, err := os.ReadFile(name)
		if err != nil {
			if os.IsNotExist(err) {
				return Object{}, ErrObjectNotFound
			}

			return Object{}, err
		}

		return Object{Hash: filepath.Base(name), Content: rf}, nil
	}

	prependToAll := func(co []string, s string) []string {
		var output []string
		for i := range co {
			output = append(output, s+co[i])
		}

		return output
	}

	fileName := hash[2:]
	if fileName == "" {
		// If user provided hash length is 2 and there is only
		// one object in that dirPath, return it.
		if len(objectsInDir) == 1 {
			return readFile(path.Join(dirPath, objectsInDir[0]))
		}

		return Object{}, &HashCollisionError{Collisions: prependToAll(objectsInDir, dirName)}
	}

	var candidates []string
	for _, c := range objectsInDir {
		if strings.Contains(dirName+c, hash) {
			candidates = append(candidates, c)
		}
	}

	if len(candidates) == 0 {
		return Object{}, ErrObjectNotFound
	}

	if len(candidates) > 1 {
		return Object{}, &HashCollisionError{Collisions: prependToAll(candidates, dirName)}
	}

	return readFile(path.Join(dirPath, candidates[0]))
}

// FetchAllObjectNames will fetch all object names from object database.
func FetchAllObjectNames() ([]string, error) {
	ok, err := ExistsMainDir()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrRepoNotInitialized
	}

	dir, err := os.ReadDir(objectDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var output []string
	for _, d := range dir {
		if d.Type().IsRegular() {
			continue
		}

		subDir, err := os.ReadDir(path.Join(objectDir, d.Name()))
		if err != nil {
			return nil, err
		}

		for _, f := range subDir {
			if f.IsDir() {
				return nil, ErrDirectoryIsNotExpected
			}

			output = append(output, d.Name()+f.Name())
		}
	}

	return output, nil
}

// fetchAllFileNamesInDir will fetch all file names in a dir.
func fetchAllFileNamesInDir(dirName string) ([]string, error) {
	dir, err := os.ReadDir(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var output []string
	for _, f := range dir {
		if f.IsDir() {
			return nil, ErrDirectoryIsNotExpected
		}

		output = append(output, f.Name())
	}

	return output, nil
}

// objectExists checks whether name exists
func objectExists(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
