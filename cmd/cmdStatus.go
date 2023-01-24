package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CMDStatusInit() {
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

	PrvKFilePathFlag(c, v)
	PrvKPasswordFlag(c, v)
	MAIN.AddCommand(c)
}

func CMDStatus(v *viper.Viper, dirs []string) {
	publicKeys, err := PubKsGet(v)
	if err != nil {
		log.Fatal("could not get publicKeys: %v", err)
	}

	s := GitStatiGet(publicKeys, dirs)
	StatiPrint(s)
}
