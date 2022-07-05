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

var NL = "\n"
var checkMark = 0x2714
var clrReset = "\033[0m"
var clrRed = "\033[31m"
var clrGreen = "\033[32m"
var clrYellow = "\033[33m"
var clrPurple = "\033[35m"

// var clrBlue = "\033[34m"
// var clrCyan = "\033[36m"
// var clrWhite = "\033[37m"

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
	syncGoodList := make([]SyncReq, 0)
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
			syncBadList = append(syncBadList, SyncReq{Dir: dir, Reason: clrRed + "git status error: " + err.Error() + clrReset})
		} else {
			gitStatus := out.String()

			var remoteOriginURL string
			{
				cmd := exec.Command("git", "config", "--get", "remote.origin.url")
				var out bytes.Buffer
				cmd.Stdout = &out
				err = cmd.Run()
				if err != nil {
					syncBadList = append(syncBadList, SyncReq{Dir: dir, Reason: clrRed + "git config error: " + err.Error() + clrReset})
				} else {
					remoteOriginURL = out.String()
					remoteOriginURL = strings.TrimSpace(remoteOriginURL)
				}
			}

			if strings.Contains(gitStatus, "nothing to commit, working tree clean") {
				syncGoodList = append(syncGoodList, SyncReq{Dir: dir, Reason: clrGreen + "synced (" + remoteOriginURL + ")" + clrReset})
			} else if strings.Contains(gitStatus, "Changes not staged for commit") {
				syncActionList = append(syncActionList, SyncReq{Dir: dir, Reason: clrYellow + "has unstaged changes (" + remoteOriginURL + ")" + clrReset})
			} else if strings.Contains(gitStatus, "untracked files present") {
				syncActionList = append(syncActionList, SyncReq{Dir: dir, Reason: clrPurple + "has untracked files (" + remoteOriginURL + ")" + clrReset})
			} else {
				syncBadList = append(syncBadList, SyncReq{Dir: dir, Reason: gitStatus})
			}
		}

		err = os.Chdir(wd)
		handle(err, "Restoring working directory: ")
	}

	for _, syncReq := range syncGoodList {
		fmt.Printf(clrGreen + string(rune(checkMark)) + clrReset + " " + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Reason + NL)
	}

	for _, syncReq := range syncBadList {
		fmt.Printf(clrRed + "x " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Reason + NL)
	}

	for _, syncReq := range syncActionList {
		fmt.Printf(clrYellow + "! " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Reason + NL)
	}
}
