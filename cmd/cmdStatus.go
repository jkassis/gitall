package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/99designs/keyring"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// A general configuration object (feed with flags, conf files, etc.)
	v := viper.New()

	// CLI Command with flag parsing
	c := &cobra.Command{
		Use:   "status",
		Short: "Get the status for multiple git repos",
		// Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			CMDStatus(v, args)
		},
	}

	CMDGitConfig(c, v)
	MAIN.AddCommand(c)
}

type SyncReq struct {
	Dir    string
	Detail string
}

var NL = "\n"
var clrReset = "\033[0m"
var clrRed = "\033[31m"
var clrGreen = "\033[32m"
var clrYellow = "\033[33m"
var clrPurple = "\033[35m"

// var clrBlue = "\033[34m"
// var clrCyan = "\033[36m"
// var clrWhite = "\033[37m"

func CMDStatus(v *viper.Viper, dirs []string) {
	var privateKeyFilePath string
	{
		privateKeyFilePath = v.GetString(SSH_KEY_PATH)
		if privateKeyFilePath == "~/.ssh/id_rsa" {
			dirname, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			privateKeyFilePath = dirname + "/.ssh/id_rsa"
		}
	}

	var privateKeyFilePassword string
	{
		privateKeyFilePassword = v.GetString(SSH_PASSWORD)
		if privateKeyFilePassword == "" {
			// See https://apple.stackexchange.com/questions/265131/recover-ssh-private-key-passphrase-from-keychain/268175#268175
			// do a keychain lookup
			// 	{
			// 		kSecClass: kSecClassGenericPassword,
			// 		kSecAttrAccount: pathToPrivateKey,
			// 		kSecAttrLabel: "SSH: " + pathToPrivateKey,
			// 		kSecAttrService: "OpenSSH",
			// 		kSecAttrNoLegacy: TRUE,
			// 		kSecAttrAccessGroup: "com.apple.ssh.passphrases"
			// }

			{
				ring, _ := keyring.Open(keyring.Config{
					// KeychainName: "icloud",
					// ServiceName: "com.apple.ssh.passphrases",
					KeychainTrustApplication: true,
				})

				keys, err := ring.Keys()
				if err != nil {
					log.Fatal(err)
				}
				if len(keys) == 0 {
					log.Error("No keys found!")
				}
				sort.Strings(keys)
				for _, key := range keys {
					log.Info(key)
					// i, _ := ring.Get(key)
					// log.Info(i.Data)
				}

				// i, _ := ring.Get("Fastly")
				// i, _ := ring.Get(privateKeyFilePath)
			}

			// {
			// 	query := keychain.NewItem()
			// 	query.SetSecClass(keychain.SecClassGenericPassword)
			// 	query.SetAccount("")
			// 	// query.SetLabel("SSH: " + privateKeyFilePath)
			// 	// query.SetService("Fastly")
			// 	// query.SetAccessGroup("com.apple.ssh.passphrases")
			// 	query.SetReturnRef(true)

			// 	// query.SetReturnAttributes(true)
			// 	query.SetMatchLimit(keychain.MatchLimitAll)
			// 	results, err := keychain.QueryItemRef(query)
			// 	if err != nil {
			// 		log.Fatalf("ssh key password not provided and could not lookup in keychain: %v", err)
			// 	} else if len(results) == 0 {
			// 		log.Fatal("key password not provided and could not lookup in keychain: no results")
			// 	} else {
			// 		for _, r := range results {
			// 			fmt.Printf("%#v\n", r)
			// 		}
			// 	}
			// }
		}
	}

	// get publicKeys
	var publicKeys *ssh.PublicKeys
	{
		_, err := os.Stat(privateKeyFilePath)
		if err != nil {
			fmt.Printf(clrRed + "read file %s failed %s\n" + clrReset)
			return
		}
		bytes, err := os.ReadFile(privateKeyFilePath)
		if err != nil {
			log.Fatalf("could not read private key file: %v", err)
		}
		publicKeys, err = ssh.NewPublicKeys("git", bytes, privateKeyFilePassword)
		if err != nil {
			log.Fatalf("could not generate signer keys: %v", err)
		}
	}

	needsSyncList := make([]SyncReq, 0)
	needsCommitList := make([]SyncReq, 0)
	repoErrorList := make([]SyncReq, 0)
	needsNothingList := make([]SyncReq, 0)

	checkErr := func(err error, dir string) bool {
		if err != nil {
			if strings.Contains(err.Error(), "knownhosts") {
				err = fmt.Errorf("problem with known_hosts entry for 'github.com'. try running `ssh-keyscan github.com >> ~/.ssh/known_hosts` on your cli: %v", err)
			}
			repoErrorList = append(repoErrorList, SyncReq{Dir: dir, Detail: err.Error()})
			return true
		}
		return false
	}

REPOS:
	for _, dir := range dirs {
		// open, get worktree, status, and config
		r, err := git.PlainOpen(dir)
		if checkErr(err, dir) {
			continue
		}

		// fetch the origin
		fmt.Printf(clrYellow + " fetching " + dir + " origin" + NL)
		err = r.Fetch(&git.FetchOptions{RemoteName: "origin", Auth: publicKeys, InsecureSkipTLS: true})
		if err != nil {
			if strings.Contains(err.Error(), "already up-to-date") {
				// do nothing
			} else if strings.Contains(err.Error(), "knownhosts") {
				err = fmt.Errorf("problem with known_hosts entry for 'github.com'. try running `ssh-keyscan github.com >> ~/.ssh/known_hosts` on your cli: %v", err)
				repoErrorList = append(repoErrorList, SyncReq{Dir: dir, Detail: err.Error()})
				continue
			}
		}

		// remoteOriginURL := origin.URLs[0]

		// get references for head and remote/origin
		refs, err := r.References()
		if checkErr(err, dir) {
			continue
		}
		refsHeads := make(map[string]string)
		refsOrigin := make(map[string]string)
		err = refs.ForEach(func(ref *plumbing.Reference) error {
			// The HEAD is omitted in a `git show-ref` so we ignore the symbolic
			// references, the HEAD
			if ref.Type() == plumbing.SymbolicReference {
				return nil
			}
			if strings.HasPrefix(string(ref.Name()), "refs/heads/") {
				refsHeads[string(ref.Name()[11:])] = ref.Hash().String()
			}
			if strings.HasPrefix(string(ref.Name()), "refs/remotes/origin/") {
				refsOrigin[string(ref.Name()[20:])] = ref.Hash().String()
			}

			return nil
		})
		if checkErr(err, dir) {
			continue
		}

		// for each head reference
		for headBranch, headHash := range refsHeads {
			originHash, ok := refsOrigin[headBranch]
			if !ok {
				needsSyncList = append(needsSyncList, SyncReq{Dir: dir + " " + headBranch, Detail: clrYellow + "has no origin branch" + clrReset})
				continue REPOS
			}

			if headHash != originHash {
				needsSyncList = append(needsSyncList, SyncReq{Dir: dir + " " + headBranch, Detail: clrYellow + "out of sync with origin" + clrReset})
				continue REPOS
			}
		}

		// now get the current worktree
		w, err := r.Worktree()
		if checkErr(err, dir) {
			continue
		}

		// loop through all status
		s, err := w.Status()
		if checkErr(err, dir) {
			continue
		}
		for _, status := range s {
			if status.Worktree != git.Unmodified {
				needsCommitList = append(needsCommitList, SyncReq{Dir: dir, Detail: clrPurple + "has unstaged changes" + clrReset})
				continue
			}
			if status.Staging != git.Unmodified {
				needsCommitList = append(needsCommitList, SyncReq{Dir: dir, Detail: clrPurple + "has staged changes" + clrReset})
				continue
			}
		}

		needsNothingList = append(needsNothingList, SyncReq{Dir: dir, Detail: clrGreen + "in sync" + clrReset})
	}

	for _, syncReq := range repoErrorList {
		fmt.Printf(clrRed + " x  " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}

	for _, syncReq := range needsNothingList {
		fmt.Printf(clrGreen + " \u2714 " + clrReset + " " + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}

	for _, syncReq := range needsCommitList {
		fmt.Printf(clrPurple + " +  " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}

	for _, syncReq := range needsSyncList {
		fmt.Printf(clrYellow + "<-> " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}
}
