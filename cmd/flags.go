package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	keyring "github.com/99designs/keyring"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var ringCached keyring.Keyring

func KeyringGet() keyring.Keyring {
	if ringCached == nil {
		var err error
		ringCached, err = keyring.Open(keyring.Config{ServiceName: "gitall"})
		if err != nil {
			log.Fatalf("could not access keyring: %v")
		}
	}
	return ringCached
}

func Prompt(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("could not read prompt response: %v", err)
	}

	// remove the delimeter from the string
	input = strings.TrimSpace(input)
	return input
}

func PromptSecret(label string) string {
	var s string
	fmt.Fprint(os.Stderr, label+" ")
	b, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("could not read secret prompt response: %v", err)
	}
	s = string(b)
	s = strings.TrimSpace(s)
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

func PrvKPasswordGet(v *viper.Viper, prvKFilePath string) (value string, err error) {
	ringKey := prvKFilePath
	ringItem, err := KeyringGet().Get(ringKey)

	// prompt if required or password not found
	prompt := v.GetBool(SSH_KEY_PASS_PROMPT)
	if prompt || err == keyring.ErrKeyNotFound {
		if prompt {
			log.Warnf("prompting by request...")
		} else {
			return "", fmt.Errorf("ssh key password not found in keychain. use -p to provide it")
		}

		// prompt
		value = PromptSecret("Enter ssh key password: ")

		// add
		err = KeyringGet().Set(keyring.Item{Key: ringKey, Data: []byte(value)})
		if err != nil {
			return "", fmt.Errorf("could not save ssh key password in keychain: %v", err)
		} else {
			log.Warnf("saved ssh key password in keychain for %s", prvKFilePath)
		}
	} else if err != nil {
		return "", fmt.Errorf("could not query keychain for ssh key password: %v", err)
	} else {
		value = string(ringItem.Data)
		log.Warnf("got ssh key password from keychain for %s. use -p to override with prompt", prvKFilePath)
	}
	return
}

const GITHUB_PASS_PROMPT = "ghpass"

func GithubPassFlag(c *cobra.Command, v *viper.Viper) {
	c.PersistentFlags().Bool(GITHUB_PASS_PROMPT, false, "prompt for github.com password")
	v.BindPFlag(GITHUB_PASS_PROMPT, c.PersistentFlags().Lookup(GITHUB_PASS_PROMPT))
}

func GitHubPassGet(v *viper.Viper) (value string, err error) {
	const ringKey = "ghpass"
	ringItem, err := KeyringGet().Get(ringKey)

	// prompt if required or password not found
	prompt := v.GetBool(GITHUB_PASS_PROMPT)
	if prompt || err == keyring.ErrKeyNotFound {
		if prompt {
			log.Warnf("prompting by request...")
		} else {
			return "", fmt.Errorf("github.com password not found in keychain. use --ghpass to provide it")
		}

		// prompt
		value = PromptSecret("Enter github.com password: ")

		// add
		err = KeyringGet().Set(keyring.Item{Key: ringKey, Data: []byte(value)})
		if err != nil {
			return "", fmt.Errorf("could not save github.com password in keychain: %v", err)
		} else {
			log.Warnf("saved github.com password in keychain")
		}
	} else if err != nil {
		return "", fmt.Errorf("could not query keychain for github.com password: %v", err)
	} else {
		value = string(ringItem.Data)
		log.Warnf("got github.com password from keychain. use -ghpass to override with prompt")
	}
	return
}

const GITHUB_USER_PROMPT = "ghuser"

func GithubUserFlag(c *cobra.Command, v *viper.Viper) {
	c.PersistentFlags().Bool(GITHUB_USER_PROMPT, false, "prompt for github.com username")
	v.BindPFlag(GITHUB_USER_PROMPT, c.PersistentFlags().Lookup(GITHUB_USER_PROMPT))
}

func GitHubUserGet(v *viper.Viper) (value string, err error) {
	const ringKey = "ghuser"
	ringItem, err := KeyringGet().Get(ringKey)

	// prompt if required or password not found
	prompt := v.GetBool(GITHUB_USER_PROMPT)
	if prompt || err == keyring.ErrKeyNotFound {
		if prompt {
			log.Warnf("prompting by request...")
		} else {
			return "", fmt.Errorf("github.com username not found in keychain. use --ghuser to provide it")
		}

		// prompt
		value = Prompt("Enter username for github.com: ")

		// add
		err = KeyringGet().Set(keyring.Item{Key: ringKey, Data: []byte(value)})
		if err != nil {
			return "", fmt.Errorf("could not save github.com username in keychain: %v", err)
		} else {
			log.Warnf("saved github.com username in keychain")
		}
	} else if err != nil {
		return "", fmt.Errorf("could not query keychain for github.com username: %v", err)
	} else {
		value = string(ringItem.Data)
		log.Warnf("got github.com username from keychain. use -ghuser to override with prompt")
	}
	return
}
