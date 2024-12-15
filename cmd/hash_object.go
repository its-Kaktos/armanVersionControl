package cmd

import (
	"armanVersionControl/storage"
	"armanVersionControl/storage/objectstore"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	filePath string
	content  string
	write    bool
)

// TODO update this Short and Long in the format of `catFileCmd`
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object [-w | --write] { {--file-path | -f} | {--content | -c} }",
	Short: "Computes the object ID value for a specified file path or provided content.",
	Long: `Computes and reports the object ID value for a specified file path or provided content and optionally writes the resulting object into the object database.
If the provided path is a directory path, will generate a tree for the directory and each directory and each file will be a blob in that tree.
Note:
	- Currently this command will NOT reuse currently stored trees or blobs, it will always generate a new tree or blob.
	- If provided path is a directory path, this command will store every file and subdirectory in the object store ignoring the absence of -w flag.
      	Files or directories that start with a '.' AKA the hidden files and directories are ignored.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := computeHashAndWriteIfFlag()
		if err != nil {
			return err
		}

		fmt.Println(s)
		return nil
	},
}

func init() {
	hashObjectCmd.Flags().StringVarP(&filePath, "file-path", "f", "", "Hash content of file located at the given path.")
	hashObjectCmd.Flags().StringVarP(&content, "content", "c", "", "Hash the provided content from initCmd line input.")
	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "When specified will actually store the resulting object into object database.")

	RootCmd.AddCommand(hashObjectCmd)
}

// TODO maybe move these functions to another go file?
func computeHashAndWriteIfFlag() (string, error) {
	if content != "" {
		if write {
			return objectstore.StoreBlob(storage.Blob{Content: []byte(content)})
		}

		return objectstore.ComputeHash([]byte(content)), nil
	}

	if filePath == "" {
		return "", errors.New("filePath can not be empty")
	}

	s, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	if s.IsDir() {
		t, err := computeTree(filePath)
		if err != nil {
			return "", err
		}

		if write {
			return objectstore.StoreTree(&t)
		}

		b, err := t.FileRepresent()
		if err != nil {
			return "", err
		}

		return objectstore.ComputeHash(b), nil
	}

	if !s.Mode().IsRegular() {
		return "", fmt.Errorf("expected a regular file but got %+v", s)
	}

	b, err := computeBlob(filePath)
	if err != nil {
		return "", err
	}

	if write {
		return objectstore.StoreBlob(b)
	}

	return objectstore.ComputeHash(b.FileRepresent()), err
}

func computeBlob(fp string) (storage.Blob, error) {
	if fp == "" {
		return storage.Blob{}, errors.New("fp (file path) can not be empty")
	}

	s, err := os.Stat(filePath)
	if err != nil {
		return storage.Blob{}, err
	}

	if !s.Mode().IsRegular() {
		return storage.Blob{}, fmt.Errorf("expected a regular file but got %+v", s)
	}

	c, err := os.ReadFile(filePath)
	if err != nil {
		return storage.Blob{}, err
	}

	return storage.Blob{Content: c}, nil
}

func computeTree(dirPath string) (storage.Tree, error) {
	if dirPath == "" {
		return storage.Tree{}, errors.New("filepath can not be empty")
	}

	s, err := os.Stat(dirPath)
	if err != nil {
		return storage.Tree{}, err
	}

	if !s.IsDir() {
		return storage.Tree{}, fmt.Errorf("expected a dir but got %+v", s)
	}

	return storage.NewTreeFromPath(dirPath)
}
