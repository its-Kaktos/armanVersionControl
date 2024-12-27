package cmd

import (
	"armanVersionControl/storage"
	"fmt"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Currently prints all objects stored in object database.",
	Long:  "Currently prints all objects stored in object database.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := storage.FetchAllObjectNames()
		if err != nil {
			return err
		}

		for _, o := range all {
			fmt.Println(o)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(logCmd)
}
