package cmd

import (
	"armanVersionControl/storage"
	"armanVersionControl/structures"
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

var hashObjectCmd = &cobra.Command{
	Use:   "hash-object [-w | --write] { {--file-path | -f} | {--content | -c} }",
	Short: "Computes the object ID value for a specified file path or provided content.",
	Long: `Computes and reports the object ID value for a specified file path or provided content and optionally writes the resulting object into the object database.
If the provided path is a directory path, will generate a tree for the directory and each directory and each file will be a blob in that tree.

Note:
	- If a file or directory is added to object database more than once, the previous hash and object will be used.
	- Files or directories that start with a '.' AKA the hidden files and directories are ignored.`,
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
			return structures.Blob{Content: []byte(content)}.StoreBlob()
		}

		return storage.ComputeHash([]byte(content)), nil
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
			return t.StoreTree()
		}

		b, err := t.FileRepresent()
		if err != nil {
			return "", err
		}

		return storage.ComputeHash(b), nil
	}

	if !s.Mode().IsRegular() {
		return "", fmt.Errorf("expected a regular file but got %+v", s)
	}

	b, err := computeBlob(filePath)
	if err != nil {
		return "", err
	}

	if write {
		return b.StoreBlob()
	}

	return storage.ComputeHash(b.FileRepresent()), err
}

func computeBlob(fp string) (structures.Blob, error) {
	if fp == "" {
		return structures.Blob{}, errors.New("fp (file path) can not be empty")
	}

	s, err := os.Stat(fp)
	if err != nil {
		return structures.Blob{}, err
	}

	if !s.Mode().IsRegular() {
		return structures.Blob{}, fmt.Errorf("expected a regular file but got %+v", s)
	}

	c, err := os.ReadFile(fp)
	if err != nil {
		return structures.Blob{}, err
	}

	return structures.Blob{Content: c}, nil
}

func computeTree(dirPath string) (structures.Tree, error) {
	if dirPath == "" {
		return structures.Tree{}, errors.New("filepath can not be empty")
	}

	s, err := os.Stat(dirPath)
	if err != nil {
		return structures.Tree{}, err
	}

	if !s.IsDir() {
		return structures.Tree{}, fmt.Errorf("expected a dir but got %+v", s)
	}

	return structures.NewTreeFromPath(dirPath)
}
