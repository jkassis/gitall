package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// MAIN represents the base command when called without any subcommands
var MAIN = &cobra.Command{
	Use:   "gitall",
	Short: "A CLI for operations on groups of git repos.",
}

func init() {
	CMDStatusInit()
	CMDUpdateTapInit()
}

func main() {
	err := MAIN.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
