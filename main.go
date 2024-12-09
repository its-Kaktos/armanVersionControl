package main

import (
	"armanVersionControl/storage/objectstore"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Write a string: ")

	r := bufio.NewReader(os.Stdin)

	input, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.Replace(input, "\n", "", -1)

	err = objectstore.Init()
	if err != nil {
		panic(err)
	}

	//writeNew([]byte(input))
	s, err := fetch(input)
	if err != nil {
		if errors.Is(err, objectstore.ErrObjectNotFound) {
			fmt.Println("Object not found")
			return
		}

		panic(err)
	}
	fmt.Println("File content is:")
	fmt.Println(s)

	fmt.Println("+++++++DONE+++++++")
}

func fetch(hash string) (string, error) {
	file, err := objectstore.FetchByHash(hash)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func writeNew(input []byte) {
	err := objectstore.Store(input)

	if err != nil {
		if errors.Is(err, objectstore.ErrObjectAlreadyExists) {
			fmt.Println("Error: an object with the provided content already exists")

			return
		}

		panic(err)
	}
}
