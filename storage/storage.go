package storage

import (
	"errors"
	"os"
)

var (
	MainDir             = ".avc"
	dirPerm os.FileMode = 0777
)

var (
	ErrAlreadyInitialized = errors.New("avc repository is already initialized")
	ErrRepoNotInitialized = errors.New("not an avc repository")
)

// ExistsMainDir will check if the .avc directory exists
func ExistsMainDir() (bool, error) {
	_, err := os.Stat(MainDir)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

// Init will initialize an empty avc repository.
func Init() error {
	ok, err := ExistsMainDir()
	if err != nil {
		return err
	}
	if ok {
		return ErrAlreadyInitialized
	}

	return mkdirAllIfDoesNotExists(MainDir, dirPerm)
}

// mkdirAllIfDoesNotExists will make directories if they do not exist
// in path of name with the provided perm as directory permission.
func mkdirAllIfDoesNotExists(name string, perm os.FileMode) error {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return os.MkdirAll(name, perm)
	}

	return err
}
