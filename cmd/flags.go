package main

import (
	"fmt"
	"os"
	"syscall"

	keyring "github.com/99designs/keyring"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

func PasswordPrompt(label string) string {
	var s string
	fmt.Fprint(os.Stderr, label+" ")
	b, _ := term.ReadPassword(int(syscall.Stdin))
	s = string(b)
	fmt.Println()
	return s
}

const SSH_KEY_PATH = "ssh_key_path"

func PrvKFilePathFlag(c *cobra.Command, v *viper.Viper) {
	c.PersistentFlags().StringP(SSH_KEY_PATH, "k", "~/.ssh/id_rsa", "ssh key path")
	v.BindPFlag(SSH_KEY_PATH, c.PersistentFlags().Lookup(SSH_KEY_PATH))
}

func PrvKFilePathGet(v *viper.Viper) (prvKFilePath string, err error) {
	prvKFilePath = v.GetString(SSH_KEY_PATH)
	if prvKFilePath == "~/.ssh/id_rsa" {
		dirname, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		prvKFilePath = dirname + "/.ssh/id_rsa"
	}
	return
}

const SSH_KEY_PASS_PROMPT = "ssh_key_pass_prompt"

func PrvKPasswordFlag(c *cobra.Command, v *viper.Viper) {
	c.PersistentFlags().BoolP(SSH_KEY_PASS_PROMPT, "p", false, "prompt for ssh key password")
	v.BindPFlag(SSH_KEY_PASS_PROMPT, c.PersistentFlags().Lookup(SSH_KEY_PASS_PROMPT))
}

func PrvKPasswordGet(v *viper.Viper, prvKFilePath string) (prvKPassword string, err error) {
	ring, _ := keyring.Open(keyring.Config{ServiceName: "gitall"})
	prvKPasswordItem, err := ring.Get(prvKFilePath)

	// prompt if required or password not found
	prompt := v.GetBool(SSH_KEY_PASS_PROMPT)
	if prompt || err == keyring.ErrKeyNotFound {
		if prompt {
			log.Warnf("prompting by request...")
		} else {
			return "", fmt.Errorf("ssh private key password not found in keychain. use -p to provide it")
		}

		// prompt
		prvKPassword = PasswordPrompt("Enter password for " + prvKFilePath + ": ")

		// add
		err = ring.Set(keyring.Item{Key: prvKFilePath, Data: []byte(prvKPassword)})
		if err != nil {
			return "", fmt.Errorf("could not save password in keychain: %v", err)
		} else {
			log.Warnf("saved password in keychain for %s", prvKFilePath)
		}
	} else if err != nil {
		return "", fmt.Errorf("could not query keychain for ssh private key password: %v", err)
	} else {
		prvKPassword = string(prvKPasswordItem.Data)
		log.Warnf("got password from keychain for %s. use -p to override with prompt", prvKFilePath)
	}
	return
}
