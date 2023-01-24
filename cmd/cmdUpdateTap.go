package main

import (
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

	CMDGitConfig(c, v)
	MAIN.AddCommand(c)
}

func CMDUpdateTap(v *viper.Viper, dirs []string) {
	publicKeys, err := PubKsGet(v)
	if err != nil {
		log.Fatal("could not get publicKeys: %v", err)
	}

	s := GitStatiGet(publicKeys, dirs)
	StatiPrint(s)

	// for _, status := range s.NeedsNothingList {
	// 	// release
	// 	/**
	// 	Tar and gz built objects
	// 	Do not ignore build directory
	// 	Change build dir to dist directory
	// 	Invoke semvar to choose the version
	// 	Commit and tag the repo using the new version
	// 	Push the version
	// 	Run gh release to release the binaries to GitHub
	// 	**/

	// 	/**
	// 	Command takes the path to the tap repo and sub path.
	// 	update the formulae.
	// 	Commit the tap repo and push.
	// 	**/
	// }
}
