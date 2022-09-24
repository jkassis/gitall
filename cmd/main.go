package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MAIN represents the base command when called without any subcommands
var MAIN = &cobra.Command{
	Use:   "gitall",
	Short: "A CLI for operations on groups of git repos.",
}

func main() {
	err := MAIN.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const (
	SSH_PASSWORD = "ssh_key_password"
	SSH_KEY_PATH = "ssh_key_path"
)

func CMDGitConfig(c *cobra.Command, v *viper.Viper) {
	c.PersistentFlags().StringP(SSH_PASSWORD, "p", "", "ssh key password")
	// c.MarkPersistentFlagRequired(SSH_PASSWORD)
	v.BindPFlag(SSH_PASSWORD, c.PersistentFlags().Lookup(SSH_PASSWORD))

	c.PersistentFlags().StringP(SSH_KEY_PATH, "k", "~/.ssh/id_rsa", "ssh key path")
	v.BindPFlag(SSH_KEY_PATH, c.PersistentFlags().Lookup(SSH_KEY_PATH))
}
