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

func init() {
	CMDStatusInit()
}

func main() {
	err := MAIN.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const (
	SSH_KEY_PATH        = "ssh_key_path"
	SSH_KEY_PASS_PROMPT = "ssh_key_pass_prompt"
)

func CMDGitConfig(c *cobra.Command, v *viper.Viper) {
	c.PersistentFlags().BoolP(SSH_KEY_PASS_PROMPT, "p", false, "prompt for ssh key password")
	v.BindPFlag(SSH_KEY_PASS_PROMPT, c.PersistentFlags().Lookup(SSH_KEY_PASS_PROMPT))

	c.PersistentFlags().StringP(SSH_KEY_PATH, "k", "~/.ssh/id_rsa", "ssh key path")
	v.BindPFlag(SSH_KEY_PATH, c.PersistentFlags().Lookup(SSH_KEY_PATH))
}
