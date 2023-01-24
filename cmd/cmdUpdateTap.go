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

	PrvKFilePathFlag(c, v)
	PrvKPasswordFlag(c, v)
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
	// 	Command takes the path to the tap repo and sub path.
	// 	update the formulae.
	// 	Commit the tap repo and push.
	// 	**/
	// }
}
