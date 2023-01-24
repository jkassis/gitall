package main

import (
	"context"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	// ctx := context.Background()
	// ts := oauth2.StaticTokenSource(
	// 	&oauth2.Token{AccessToken: "token"},
	// )
	// tc := oauth2.NewClient(ctx, ts)
	// client := github.NewClient(tc)

	client, err := GithubClientGet(v)
	if err != nil {
		log.Fatal("could not get github client: %v", err)
	}

	// get the status of requested dirs
	s := GitStatiGet(publicKeys, dirs)
	StatiPrint(s)

	// for each that is in sync
	for _, status := range s.NeedsNothingList {
		// open the repo
		repo, err := git.PlainOpen(status.Dir)
		if err != nil {
			log.Errorf("error opening git repo for %v", status.Dir)
			continue
		}

		// get the github owner and repo of the origin url
		var githubOwner, githubRepo string
		{
			remote, err := repo.Remote("origin")
			if err != nil {
				log.Errorf("could not get origin remote: %v", err)
				continue
			}
			url, err := url.Parse(remote.Config().URLs[0])
			if err != nil {
				log.Errorf("could not get origin url: %v", err)
				continue
			}
			parts := strings.Split(url.Path, "/")
			githubOwner = parts[0]
			githubRepo = parts[1]
		}

		// get latestReleaseTagName
		var latestReleaseTagName string
		{
			ctx := context.Background()
			latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, githubOwner, githubRepo)
			if err != nil {
				log.Errorf("could not get latest release for %s/%s: %v", githubOwner, githubRepo, err)
				continue
			}
			latestReleaseTagName = latestRelease.GetTagName()
		}

		// read the tap formula
		formulaPath := BrewTapRepoLocalPath(v) + "/Formula/" + githubRepo + ".rb"
		formulaData, err := os.ReadFile(formulaPath) // TODO should be the name of the binary / go module
		if err != nil {
			log.Errorf("could not open the Formula at %s: %v", formulaPath, err)
			continue
		}

		// update the release
		r, err := regexp.Compile("releases/download/.*?/")
		if err != nil {
			log.Errorf("could not build regexp for replacement: %v", err)
			continue
		}
		os.Stdout.Write(r.ReplaceAll(formulaData, []byte("releases/download/"+latestReleaseTagName)))

		// release
		/**
		Command takes the path to the tap repo and sub path.
		update the formulae.
		Commit the tap repo and push.
		**/
	}
}
