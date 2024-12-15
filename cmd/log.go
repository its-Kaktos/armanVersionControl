package cmd

import (
	"armanVersionControl/storage"
	"armanVersionControl/storage/objectstore"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

/*
TODO
- Remove printing all objects
- Pass a hash with -p tag to tell print its content
*/
var logCmd = &cobra.Command{
	Use:   "log hash",
	Short: "Prints content of object associated with hash.",
	Long:  "Prints content of object associated with hash, will do further processing to figure out the object type to print it in a readable format.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := args[0]
		c, err := objectstore.FetchByHash(hash)
		if err != nil {
			return err
		}

		s, err := computeOriginalContent(c)
		if err != nil {
			return err
		}

		// the \033[1;32 part is the coloring. Read more at: https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences
		fmt.Println("\033[1;32mHere is the content of provided hash:\033[0m")
		fmt.Println(s)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(logCmd)
}

func computeOriginalContent(b []byte) (string, error) {
	if storage.IsBlobB(b) {
		blob, err := storage.NewBlobFromB(b)
		if err != nil {
			if errors.Is(err, storage.ErrNotABlob) {
				panic("expected a blob but the content is not a blob")
			}

			return "", err
		}

		return string(blob.Content), nil
	}

	// TODO do this for the Tree as well
	//if storage.IsTreeB()

	return "", errors.New("invalid content, content should either be a Blob or a Tree")
}
