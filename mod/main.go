package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type SyncReq struct {
	Dir    string
	Reason string
}

var checkMark = 0x2714
var colorReset = "\033[0m"

// var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"

// var colorBlue = "\033[34m"
// var colorPurple = "\033[35m"
// var colorCyan = "\033[36m"
// var colorWhite = "\033[37m"

func main() {
	var err error
	handle := func(err error) {
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
	}
	var wd string
	wd, err = os.Getwd()
	handle(err)

	syncGoodList := make([]string, 0)
	syncBadList := make([]SyncReq, 0)

	for _, dir := range os.Args[1:] {
		err = os.Chdir(dir)
		handle(err)

		// do something
		cmd := exec.Command("git", "status")
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		handle(err)

		gitStatus := out.String()
		syncGood := true
		if strings.Contains(gitStatus, "untracked files present") {
			syncGood = false
		}

		if syncGood {
			syncGoodList = append(syncGoodList, dir)
		} else {
			syncBadList = append(syncBadList, SyncReq{Dir: dir, Reason: "has untracked files"})
		}

		err = os.Chdir(wd)
		handle(err)
	}

	for _, dir := range syncGoodList {
		fmt.Printf("%s%c%s %s\n", colorGreen, checkMark, colorReset, dir)
	}

	for _, syncReq := range syncBadList {
		fmt.Printf("%sx %s%s (%s)\n", string(colorYellow), string(colorReset), syncReq.Dir, syncReq.Reason)
	}
}
