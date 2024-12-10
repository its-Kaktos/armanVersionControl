package storage

import (
	"armanVersionControl/cmd"
	"armanVersionControl/hashing"
	"armanVersionControl/storage/objectstore"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "hash-object [-w] {{--file-path | -fp} | {{--content | -c}}",
	Short: "Computes the hash",
}

func init() {
	command.Flags().BoolVarP()
	cmd.RootCmd.AddCommand(command)
}

func Hash(content []byte) []byte {
	return hashing.Sha1(content)
}

func Store(content []byte) (path string, error error) {
	return objectstore.Store(content)
}
