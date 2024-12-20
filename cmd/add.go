package cmd

import "github.com/spf13/cobra"

// TODO check if the passed path is in current avc repo
// TODO before adding commit, add config to add author and commiter email

var addCmd = &cobra.Command{
	Use:   "add path",
	Short: "Will add the provided path to the index.",
	Long: `This command updates the current index with the content found in the provided path, to prepare and stage content for the next commit.
The index holds a snapshot of the current content of the working tree.

Notes:
	Adding a file to index one time does not mean the file is being tracked by the avc repository indefinitely. Adding a file or directory to the index
	means that current file or directory content and structure is added to the index and prepared to commit and subsequent changes to a file or a directory
	need to be indexed again.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
