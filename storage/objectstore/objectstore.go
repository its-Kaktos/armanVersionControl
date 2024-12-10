package objectstore

import (
	"armanVersionControl/hashing"
	"encoding/hex"
	"errors"
	"os"
	"path"
)

// TODO when fetch by hash, check if there are multiple version of that hash and report that
// TODO change file perm? should other users see this? in git all have read access only, why?
// TODO add comment for all functions
// TODO add pager instead of reading whole file?

var (
	ErrRepoNotInitialized     = errors.New("not an avc repository")
	ErrAlreadyInitialized     = errors.New("avc repository is already initialized")
	ErrObjectAlreadyExists    = errors.New("object already exists")
	ErrHashIsShort            = errors.New("provided hash is short, it should be at least 3 characters")
	ErrObjectNotFound         = errors.New("object not found")
	ErrDirectoryIsNotExpected = errors.New("directory is not expected in a directory of object database")
)

var (
	mainDir               = ".avc"
	objectDir             = path.Join(mainDir, "objects")
	dirPerm   os.FileMode = 0777
	filePerm  os.FileMode = 0770
)

func Init() error {
	ok, err := existsMainDir()
	if err != nil {
		return err
	}
	if ok {
		return ErrAlreadyInitialized
	}

	return mkdirAllIfDoesNotExists(mainDir, dirPerm)
}

func ComputeHash(content []byte) string {
	return hex.EncodeToString(hashing.Sha1(content))
}

func Store(content []byte) (hash string, e error) {
	ok, err := existsMainDir()
	if err != nil {
		return "", err
	}

	if !ok {
		return "", ErrRepoNotInitialized
	}

	hashHex := ComputeHash(content)

	dir := path.Join(objectDir, hashHex[:2])
	fileName := hashHex[3:]
	filePath := path.Join(dir, fileName)

	ok, err = objectExists(path.Join(dir, fileName))
	if err != nil {
		return "", err
	}

	if ok {
		return "", ErrObjectAlreadyExists
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

func existsMainDir() (bool, error) {
	_, err := os.Stat(mainDir)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func mkdirAllIfDoesNotExists(name string, perm os.FileMode) error {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return os.MkdirAll(name, perm)
	}

	return err
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
	ok, err := existsMainDir()
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

			output = append(output, f.Name())
		}
	}

	return output, nil
}

func objectExists(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
