package main

import "armanVersionControl/cmd"

// TODO add a friendly string for Kind and other iota types
// TODO add config file to store author and commiter emails and stuff and read from it. also add a default config file when innit is called.

func main() {
	cmd.Execute()
}
