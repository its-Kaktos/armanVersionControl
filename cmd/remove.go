package cmd

import (
	"armanVersionControl/track"
	"fmt"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm name",
	Short: "Will remove the provided path from index.",
	Long: `This command updates the current index and removes the content found in the provided path. 
The index holds a snapshot of the current content of the working tree.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := track.Remove(name); err != nil {
			return err
		}

		fmt.Printf("%v removed from index successfully.\n", name)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
