package cmd

import (
	"armanVersionControl/storage"
	"fmt"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates an empty avc repository.",
	Long:  "Creates an empty Arman version control repository",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.Init(); err != nil {
			return err
		}

		fmt.Println("An empty avc repository created successfully.")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
