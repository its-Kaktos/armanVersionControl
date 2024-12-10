package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:     "avc",
	Aliases: []string{"kvc"},
	Short:   "avc is a version control software",
	Long: "avc is a version control software that is heavily inspired by git." +
		"\nThis is just a hobby project, do NOT use in production." +
		"\navc stands for Arman version control",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured: '%s'\n", err)
		os.Exit(1)
	}
}
