package cmd

import (
	"armanVersionControl/storage"
	"armanVersionControl/structures"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var prettyPrint = false
var catFileCmd = &cobra.Command{
	Use:   "cat-file {hash} ([-p | --pretty-print])",
	Short: "Display content of object by its hash.",
	Long: `This command will display the content of an object stored in object by its hash.

Arguments:
    hash		The required hash representing the object ID stored in object database.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := args[0]
		c, err := storage.FetchByHash(hash)
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
	catFileCmd.Flags().BoolVarP(&prettyPrint, "pretty-print", "p", false, "Pretty print the content of hash object based on its type.")
	RootCmd.AddCommand(catFileCmd)
}

func computeOriginalContent(o storage.Object) (string, error) {
	if !prettyPrint {
		return string(o.Content), nil
	}

	if structures.IsBlobB(o.Content) {
		blob, err := structures.NewBlobFromB(o.Content)
		if err != nil {
			if errors.Is(err, structures.ErrNotABlob) {
				panic("expected a blob but the content is not a blob")
			}

			return "", err
		}

		return string(blob.Content), nil
	}

	if structures.IsTreeB(o.Content) {
		t, err := structures.NewTreeFromObject(o)
		if err != nil {
			return "", err
		}

		return t.String(), nil
	}

	return "", errors.New("invalid content, content should either be a Blob or a Tree")
}
