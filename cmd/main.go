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

var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"

var NL = "\n"

// var colorBlue = "\033[34m"
var colorPurple = "\033[35m"

// var colorCyan = "\033[36m"
// var colorWhite = "\033[37m"

func main() {
	var err error
	handle := func(err error, ctx string) {
		if err != nil {
			fmt.Println(ctx, err)
			os.Exit(1)
		}
	}
	var wd string
	wd, err = os.Getwd()
	handle(err, "Getting working directory: ")

	syncBadList := make([]SyncReq, 0)
	syncGoodList := make([]string, 0)
	syncActionList := make([]SyncReq, 0)

	for _, dir := range os.Args[1:] {
		// go to the directory
		fileInfo, err := os.Stat(dir)
		if err != nil {
			syncBadList = append(syncBadList, SyncReq{Dir: dir, Reason: err.Error()})
			continue
		}

		if !fileInfo.IsDir() {
			continue
		}

		err = os.Chdir(dir)
		if err != nil {
			syncBadList = append(syncBadList, SyncReq{Dir: dir, Reason: err.Error()})
			continue
		}

		// get git status
		cmd := exec.Command("git", "status")
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			syncBadList = append(syncBadList, SyncReq{Dir: dir, Reason: colorRed + "git status error: " + err.Error() + colorReset})
			continue
		}

		// report
		gitStatus := out.String()
		if strings.Contains(gitStatus, "nothing to commit, working tree clean") {
			syncGoodList = append(syncGoodList, dir)
		} else if strings.Contains(gitStatus, "Changes not staged for commit") {
			syncActionList = append(syncActionList, SyncReq{Dir: dir, Reason: colorYellow + "has unstaged changes" + colorReset})
		} else if strings.Contains(gitStatus, "untracked files present") {
			syncActionList = append(syncActionList, SyncReq{Dir: dir, Reason: colorPurple + "has untracked files" + colorReset})
		} else {
			syncBadList = append(syncBadList, SyncReq{Dir: dir, Reason: gitStatus})
		}

		err = os.Chdir(wd)
		handle(err, "Restoring working directory: ")
	}

	for _, dir := range syncGoodList {
		fmt.Printf(colorGreen + string(rune(checkMark)) + colorReset + " " + fmt.Sprintf("%-40s", dir) + NL)
	}

	for _, syncReq := range syncBadList {
		fmt.Printf(colorRed + "x " + colorReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Reason + NL)
	}

	for _, syncReq := range syncActionList {
		fmt.Printf(colorYellow + "! " + colorReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Reason + NL)
	}
}
