package cmd

import (
	"armanVersionControl/track"
	"github.com/spf13/cobra"
)

// TODO update Short and Long.
// TODO check how git status shows doc and show comment like that
// TODO Add Status in track to give slice of tracked and untracked files in the repo
// TODO add comment for which files has been updated or created new?
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Will add the provided path to the index.",
	Long: `This command updates the current index with the content found in the provided path, to prepare and stage content for the next commit.
The index holds a snapshot of the current content of the working tree.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := track.Add(name); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
