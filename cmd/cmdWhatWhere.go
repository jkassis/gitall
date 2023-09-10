package main

import (
	"fmt"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type WhatWhere map[string]Status

func CMDWhatWhereInit() {
	// A general configuration object (feed with flags, conf files, etc.)
	v := viper.New()

	// CLI Command with flag parsing
	c := &cobra.Command{
		Use:   "whatwhere",
		Short: "List repo-branch of target folders.",
		// Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			CMDWhatWhere(v, args)
		},
	}

	PrvKFilePathFlag(c, v)
	PrvKPasswordFlag(c, v)
	MAIN.AddCommand(c)
}

func CMDWhatWhere(v *viper.Viper, dirs []string) {
	publicKeys, err := PubKsGet(v)
	if err != nil {
		log.Fatalf("could not get publicKeys: %v", err)
	}

	s := GitWhatWhereGet(publicKeys, dirs)
	WhatWherePrint(s)
}

func GitWhatWhereGet(publicKeys *ssh.PublicKeys, dirs []string) WhatWhere {
	s := make(map[string]Status)

	checkErr := func(err error, dir string) bool {
		if err != nil {
			s[dir] = Status{Dir: dir, Detail: clrRed + err.Error() + clrReset}
			return true
		}
		return false
	}

	for _, dir := range dirs {
		// open, get worktree, status, and config
		r, err := git.PlainOpen(dir)
		if checkErr(err, dir) {
			continue
		}

		o, err := r.Remote("origin")
		if err != nil {
			s[dir] = Status{Dir: dir, Detail: fmt.Sprintf(clrRed+"could not get origin %v"+clrReset, err)}
			continue
		}
		h, err := r.Head()
		if err != nil {
			s[dir] = Status{Dir: dir, Detail: fmt.Sprintf(clrRed+"could not get ref for head %v"+clrReset, err)}
			continue
		}
		s[dir] = Status{Dir: dir, Detail: fmt.Sprintf(clrGreen+"%s of %s"+clrReset, h.Name()[11:], o.Config().URLs[0])}
	}
	return s
}

func WhatWherePrint(s WhatWhere) {
	sortedKeys := func(m map[string]Status) []string {
		keys := make([]string, len(m))
		i := 0
		for k := range m {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		return keys
	}

	var keys []string
	keys = sortedKeys(s)
	for _, key := range keys {
		syncReq := s[key]
		fmt.Printf("%20s %s\n", syncReq.Dir, syncReq.Detail)
	}
}
