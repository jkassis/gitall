package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

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

func main() {
	// var err error
	// handle := func(err error, ctx string) {
	// 	if err != nil {
	// 		fmt.Println(ctx, err)
	// 		os.Exit(1)
	// 	}
	// }

	// var wd string
	// wd, err = os.Getwd()
	// handle(err, "Getting working directory: ")

	privateKeyFilePassword := "HDfHrRJwYbcN3"
	privateKeyFilePath := "/Users/jkassis/.ssh/id_rsa"
	_, err := os.Stat(privateKeyFilePath)
	if err != nil {
		fmt.Printf(clrRed + "read file %s failed %s\n" + clrReset)
		return
	}

	// Clone the given repository to the given directory
	bytes, err := ioutil.ReadFile(privateKeyFilePath)
	publicKeys, err := ssh.NewPublicKeys("git", bytes, privateKeyFilePassword)

	// publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFilePath, privateKeyFilePassword)
	// if err != nil {
	// 	fmt.Printf(clrRed+"generate publickeys failed: %s\n", err.Error()+clrReset)
	// 	return
	// }

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
	for _, dir := range os.Args[1:] {

		// open, get worktree, status, and config
		r, err := git.PlainOpen(dir)
		if checkErr(err, dir) {
			continue
		}

		// fetch the origin
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
