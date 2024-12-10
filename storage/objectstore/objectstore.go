package objectstore

import (
	"armanVersionControl/hashing"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
)

// TODO add a init func to initilize the .avc directory
// TODO add a validation for every storage to check if .avc exists
// TODO check for hash collision
// TODO add pager instead of reading whole file?
// TODO create a custom error type instead of all these errors?
// TODO add comment for all functions
// TODO add a global var for smallest hash possible for searching?
// TODO fetch by hash of size 3?
// TODO change file perm? should other users see this? in git all have read access only, why?

var (
	ErrObjectAlreadyExists             = errors.New("object already exists")
	ErrHashIsShort                     = errors.New("provided hash is short, it should be at least 3 characters")
	ErrObjectNotFound                  = errors.New("object not found")
	mainDir                            = ".avc"
	objectDir                          = path.Join(mainDir, "objects")
	dirPerm                os.FileMode = 0777
	filePerm               os.FileMode = 0770
)

func Init() error {
	return mkdirAllIfDoesNotExists(mainDir, dirPerm)
}

func Store(content []byte) (objectName string, error error) {
	hashHexadecimal := hex.EncodeToString(hashing.Sha1(content))

	dir := path.Join(objectDir, hashHexadecimal[:2])
	fileName := hashHexadecimal[3:]

	if err := mkdirAllIfDoesNotExists(dir, dirPerm); err != nil {
		return "", err
	}

	filePath := path.Join(dir, fileName)
	if stat, err := os.Stat(filePath); stat != nil || os.IsExist(err) {
		return "", ErrObjectAlreadyExists
	}

	err := os.WriteFile(filePath, content, filePerm)
	if err != nil {
		return "", err
	}

	return hashHexadecimal, err
}

func mkdirAllIfDoesNotExists(path string, perm os.FileMode) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, perm)
	}

	return nil
}

func FetchByHash(hash string) ([]byte, error) {
	if len(hash) < 3 {
		return nil, ErrHashIsShort
	}

	dir := path.Join(objectDir, hash[:2])
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrObjectNotFound
		}

		return nil, err
	}

	rf, err := os.ReadFile(path.Join(dir, hash[2:]))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrObjectNotFound
		}

		return nil, err
	}

	return rf, nil
}

func FetchAllObjectNames() ([]string, error) {
	dir, err := os.ReadDir(objectDir)
	if err != nil {
		return nil, err
	}

	for _, d := range dir {
		if d.Type() != os.ModeDir {
			return nil, errors.New("expcepted a dir but its not???")
		}

		subDir, err := os.ReadDir(path.Join(objectDir, d.Name()))
		if err != nil {
			return nil, err
		}

		fmt.Println(subDir)
		//for _, sd := range subDir {
		//	TODO How to check if this sd is a file or dir?
		//if _, ok := sd.(os.FileMode); !ok {
		//
		//}
		//i, err := sd.Info()
		//fmt.Printf("%v err: %v", i, err)
		//fmt.Printf("%v", sd)
		//}
	}

	return nil, nil
}
