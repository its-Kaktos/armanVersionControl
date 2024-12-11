package main

import (
	"armanVersionControl/cmd"
	"armanVersionControl/storage"
)

func main() {
	storage.Blob{Content: []byte{5, 5}}.FileRepresent()

	cmd.Execute()
}
