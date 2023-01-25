package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v49/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	giturls "github.com/whilp/git-urls"
)

func CMDUpdateTapInit() {
	// A general configuration object (feed with flags, conf files, etc.)
	v := viper.New()

	// CLI Command with flag parsing
	c := &cobra.Command{
		Use:   "updatetap",
		Short: "Updates a brew tap with the latest releases for multiple git repos.",
		// Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			CMDUpdateTap(v, args)
		},
	}

	PrvKFilePathFlag(c, v)
	PrvKPasswordFlag(c, v)
	GithubPassFlag(c, v)
	GithubUserFlag(c, v)
	BrewTapRepoLocalPathFlag(c, v)
	MAIN.AddCommand(c)
}

const BREW_TAP_REPO_PATH = "brew_tap_repo_path"

func BrewTapRepoLocalPathFlag(c *cobra.Command, v *viper.Viper) {
	c.PersistentFlags().StringP(BREW_TAP_REPO_PATH, "b", "", "brew tap repo path")
	v.BindPFlag(BREW_TAP_REPO_PATH, c.PersistentFlags().Lookup(BREW_TAP_REPO_PATH))
}

func BrewTapRepoLocalPath(v *viper.Viper) string {
	brewTapRepoPath := v.GetString(BREW_TAP_REPO_PATH)
	if _, err := os.Stat(brewTapRepoPath); os.IsNotExist(err) {
		log.Fatalf("path to brew tap repo does not exist: %s", err.Error())
	}
	return brewTapRepoPath
}

func CMDUpdateTap(v *viper.Viper, dirs []string) {
	// get public keys for git
	publicKeys, err := PubKsGet(v)
	if err != nil {
		log.Fatalf("could not get publicKeys: %v", err)
	}

	client, err := GithubClientGet(v)
	if err != nil {
		log.Fatal("could not get github client: %v", err)
	}

	// get the status of requested dirs
	s := GitStatiGet(publicKeys, dirs)
	StatiPrint(s)

	// commit changes to the tap
	tapRepo, err := git.PlainOpen(BrewTapRepoLocalPath(v))
	if err != nil {
		log.Fatalf("error opening git repo for %s: %v", BrewTapRepoLocalPath(v), err)
	}
	tapWorktree, err := tapRepo.Worktree()
	if err != nil {
		log.Fatalf("error getting worktree for tap repo: %v", err)
	}

	// for each that is in sync
	var commitMessage = ""
	for _, status := range s.NeedsNothingList {
		// open the repo
		repo, err := git.PlainOpen(status.Dir)
		if err != nil {
			log.Errorf("error opening git repo for %v", status.Dir)
			continue
		}

		// get the github owner and repo of the origin url
		var githubURL, githubOwner, githubRepo string
		{
			remote, err := repo.Remote("origin")
			if err != nil {
				log.Errorf("could not get origin remote: %v", err)
				continue
			}
			githubURL = remote.Config().URLs[0]
			url, err := giturls.Parse(githubURL)
			if err != nil {
				log.Errorf("could not get origin url: %v", err)
				continue
			}
			parts := strings.Split(url.Path, "/")
			githubOwner = parts[0]
			githubRepo = strings.TrimSuffix(parts[1], ".git")
		}

		// get latestReleaseTagName
		var latestReleaseTagName string
		var latestReleaseAssets []*github.ReleaseAsset
		{
			ctx := context.Background()
			latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, githubOwner, githubRepo)
			if err != nil {
				log.Errorf("could not get latest release for %s/%s: %v", githubOwner, githubRepo, err)
				continue
			}
			latestReleaseTagName = latestRelease.GetTagName()
			latestReleaseAssets = latestRelease.Assets
		}

		formula, err := FormulaNew(
			githubRepo,
			githubURL,
			latestReleaseTagName,
			latestReleaseAssets)

		if err != nil {
			log.Errorf("could not create new formula: %v", err)
			continue
		}

		formulaData, err := formula.Render()
		if err != nil {
			log.Errorf("could not render formula: %v", err)
			continue
		}
		fmt.Print(string(formulaData))

		// write the tap formula
		formulaPath, err := filepath.Abs(BrewTapRepoLocalPath(v) + "/Formula/" + githubRepo + ".rb")
		if err != nil {
			log.Errorf("could not get Abs path to githubRepo: %v", err)
			continue
		}
		if err != nil {
			log.Errorf("could not generate homebrew formula: %v", err)
			continue
		}

		os.WriteFile(formulaPath, formulaData, 0660)

		// add the file to the commit
		_, err = tapWorktree.Add("Formula/" + githubRepo + ".rb")
		if err != nil {
			log.Errorf("could not add %s to tap worktree: %v", formulaPath, err)
			continue
		}

		// update the commitMessage
		commitMessage += githubRepo + " ==> " + latestReleaseTagName + "\n\r"
	}

	// commit the changes to the tap
	config, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		log.Fatalf("could not load global git config: %v", err)
	}
	commit, err := tapWorktree.Commit(
		commitMessage,
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  config.Author.Name,
				Email: config.Author.Email,
				When:  time.Now(),
			},
		})
	if err != nil {
		log.Fatalf("could not create commit for tap repo: %v", err)
	}

	// Confirm
	fmt.Print("\n\n\n\n")
	fmt.Print("Please Confirm...\n\n")
	fmt.Print(commitMessage + "\n")
	proceedResponse := Prompt("Proceed? [Y/n]: ")
	if proceedResponse == "" {
		proceedResponse = "Y"
	}
	if proceedResponse != "Y" {
		log.Warnf("Cancelling update")
		return
	}

	// Prints the current HEAD to verify that all worked well.
	_, err = tapRepo.CommitObject(commit)
	if err != nil {
		log.Fatalf("could not commit to tap: %v", err)
	}

	err = tapRepo.Push(&git.PushOptions{RemoteName: "origin"})
	if err != nil {
		log.Fatalf("could not push to tap: %v", err)
	}

	log.Warnf("Complete")
}
