package storage

import (
	"armanVersionControl/hashing"
	"armanVersionControl/storage/objectstore"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use: "hash-object [-w] []",
	Short: "Computes the hash"
}

func Hash(content []byte) []byte {
	return hashing.Sha1(content)
}

func Store(content []byte) (path string, error error) {
	return objectstore.Store(content)
}
