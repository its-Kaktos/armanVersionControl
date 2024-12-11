package cmd

import (
	"armanVersionControl/storage/objectstore"
	"fmt"
	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file {hash}",
	Short: "Display content of object by its hash.",
	Long: `This command will display the content of an object stored in object by its hash.

Arguments:
    hash		The required hash representing the object ID stored in object database.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := objectstore.FetchByHash(args[0])
		if err != nil {
			return err
		}

		fmt.Println(string(content))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(catFileCmd)
}
