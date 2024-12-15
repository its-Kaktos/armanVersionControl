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

func computeHash() (string, error) {
	c, err := computeContent()
	if err != nil {
		return "", err
	}

	b := storage.Blob{Content: c}
	if write {
		return objectstore.Store(b.FileRepresent())
	}

	return objectstore.ComputeHash(b.FileRepresent()), nil
}

// TODO fix this
func computeAndStore() (, error) {
	if content != "" {
		return []byte(content), nil
	}

	if filePath == "" {
		return nil, errors.New("content and filepath can not both be empty at the same time")
	}

	s, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		t, err := storage.NewTreeFromPath(filePath)
		if err != nil {
			return nil, err
		}

		c, err := t.FileRepresent()
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	if !s.Mode().IsRegular() {
		return nil, fmt.Errorf("%v should be a directory or regular file path", filePath)
	}

	return os.ReadFile(filePath)
}
