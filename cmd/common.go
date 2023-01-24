package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/go-github/v49/github"
	"github.com/spf13/viper"
)

var NL = "\n"
var clrReset = "\033[0m"
var clrRed = "\033[31m"
var clrGreen = "\033[32m"
var clrYellow = "\033[33m"
var clrPurple = "\033[35m"

// var clrBlue = "\033[34m"
// var clrCyan = "\033[36m"
// var clrWhite = "\033[37m"

// get publicKeys
func PubKsGet(v *viper.Viper) (publicKeys *ssh.PublicKeys, err error) {
	prvKFilePath, err := PrvKFilePathGet(v)
	if err != nil {
		return nil, fmt.Errorf("could not get ssh key path: %v", err)
	}

	prvKPassword, err := PrvKPasswordGet(v, prvKFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not get ssh key password: %v", err)
	}

	_, err = os.Stat(prvKFilePath)
	if err != nil {
		return nil, fmt.Errorf(clrRed + "read file %s failed %s\n" + clrReset)
	}
	bytes, err := os.ReadFile(prvKFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read private key file: %v", err)
	}
	publicKeys, err = ssh.NewPublicKeys("git", bytes, prvKPassword)
	if err != nil {
		return nil, fmt.Errorf("could not generate signer keys: %v", err)
	}
	return publicKeys, nil
}

func ErrKnownHostsWrap(err error) error {
	if err != nil && strings.Contains(err.Error(), "knownhosts") {
		err = fmt.Errorf("problem with known_hosts entry for 'github.com'. try running `ssh-keyscan github.com >> ~/.ssh/known_hosts` on your cli: %v", err)
	}
	return err
}

type Status struct {
	Dir    string
	Detail string
}

type Stati struct {
	NeedsSyncList    []Status
	NeedsCommitList  []Status
	RepoErrorList    []Status
	NeedsNothingList []Status
}

func GitStatiGet(publicKeys *ssh.PublicKeys, dirs []string) *Stati {
	s := &Stati{
		NeedsSyncList:    make([]Status, 0),
		NeedsCommitList:  make([]Status, 0),
		RepoErrorList:    make([]Status, 0),
		NeedsNothingList: make([]Status, 0),
	}

	checkErr := func(err error, dir string) bool {
		if err != nil {
			err = ErrKnownHostsWrap(err)
			s.RepoErrorList = append(s.RepoErrorList, Status{Dir: dir, Detail: err.Error()})
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
				s.RepoErrorList = append(s.RepoErrorList, Status{Dir: dir, Detail: err.Error()})
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
				s.NeedsSyncList = append(s.NeedsSyncList, Status{Dir: dir + " " + headBranch, Detail: clrYellow + "has no origin branch" + clrReset})
				continue REPOS
			}

			if headHash != originHash {
				s.NeedsSyncList = append(s.NeedsSyncList, Status{Dir: dir + " " + headBranch, Detail: clrYellow + "out of sync with origin" + clrReset})
				continue REPOS
			}
		}

		// now get the current worktree
		w, err := r.Worktree()
		if checkErr(err, dir) {
			continue
		}

		// loop through all status
		stati, err := w.Status()
		if checkErr(err, dir) {
			continue
		}
		for _, status := range stati {
			if status.Worktree != git.Unmodified {
				s.NeedsCommitList = append(s.NeedsCommitList, Status{Dir: dir, Detail: clrPurple + "has unstaged changes" + clrReset})
				continue
			}
			if status.Staging != git.Unmodified {
				s.NeedsCommitList = append(s.NeedsCommitList, Status{Dir: dir, Detail: clrPurple + "has staged changes" + clrReset})
				continue
			}
		}

		s.NeedsNothingList = append(s.NeedsNothingList, Status{Dir: dir, Detail: clrGreen + "in sync" + clrReset})
	}
	return s
}

func StatiPrint(s *Stati) {
	for _, syncReq := range s.RepoErrorList {
		fmt.Printf(clrRed + " x  " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}

	for _, syncReq := range s.NeedsNothingList {
		fmt.Printf(clrGreen + " \u2714 " + clrReset + " " + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}

	for _, syncReq := range s.NeedsCommitList {
		fmt.Printf(clrPurple + " +  " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}

	for _, syncReq := range s.NeedsSyncList {
		fmt.Printf(clrYellow + "<-> " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}
}

func GithubClientGet(v *viper.Viper) (*github.Client, error) {
	// eventually use this... https://github.com/cli/cli/blob/v2.21.2/internal/authflow/flow.go#L37
	githubUser, err := GitHubUserGet(v)
	if err != nil {
		return nil, fmt.Errorf("could not get github.com username: %v", err)
	}
	githubPass, err := GitHubPassGet(v)
	if err != nil {
		return nil, fmt.Errorf("could not get github.com password: %v", err)
	}
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(githubUser),
		Password: strings.TrimSpace(githubPass),
	}
	client := github.NewClient(tp.Client())
	return client, nil
}
