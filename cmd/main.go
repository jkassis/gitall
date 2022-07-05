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
	Detail string
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

	needsAddList := make([]SyncReq, 0)
	needsCommitList := make([]SyncReq, 0)
	needsGitList := make([]SyncReq, 0)
	needsNothingList := make([]SyncReq, 0)

	for _, dir := range os.Args[1:] {
		// go to the directory
		fileInfo, err := os.Stat(dir)
		if err != nil {
			needsGitList = append(needsGitList, SyncReq{Dir: dir, Detail: err.Error()})
			continue
		}

		if !fileInfo.IsDir() {
			continue
		}

		err = os.Chdir(dir)
		if err != nil {
			needsGitList = append(needsGitList, SyncReq{Dir: dir, Detail: err.Error()})
			continue
		}

		// get git status
		cmd := exec.Command("git", "status")
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			needsGitList = append(needsGitList, SyncReq{Dir: dir, Detail: clrRed + "git status error: " + err.Error() + clrReset})
		} else {
			gitStatus := out.String()

			var remoteOriginURL string
			{
				cmd := exec.Command("git", "config", "--get", "remote.origin.url")
				var out bytes.Buffer
				cmd.Stdout = &out
				err = cmd.Run()
				if err != nil {
					needsGitList = append(needsGitList, SyncReq{Dir: dir, Detail: clrRed + "git config error: " + err.Error() + clrReset})
				} else {
					remoteOriginURL = out.String()
					remoteOriginURL = strings.TrimSpace(remoteOriginURL)
				}
			}

			if strings.Contains(gitStatus, "nothing to commit, working tree clean") {
				if strings.Contains(gitStatus, "Your branch is up to date") {
					needsNothingList = append(needsNothingList, SyncReq{Dir: dir, Detail: clrGreen + "in sync (" + remoteOriginURL + ")" + clrReset})
				} else {
					needsAddList = append(needsAddList, SyncReq{Dir: dir, Detail: clrYellow + "out of sync (" + remoteOriginURL + ")" + clrReset})
				}
			} else if strings.Contains(gitStatus, "Changes not staged for commit") {
				needsCommitList = append(needsCommitList, SyncReq{Dir: dir, Detail: clrPurple + "unstaged changes (" + remoteOriginURL + ")" + clrReset})
			} else if strings.Contains(gitStatus, "untracked files present") {
				needsAddList = append(needsAddList, SyncReq{Dir: dir, Detail: clrPurple + "untracked files (" + remoteOriginURL + ")" + clrReset})
			} else {
				needsGitList = append(needsGitList, SyncReq{Dir: dir, Detail: gitStatus})
			}
		}

		err = os.Chdir(wd)
		handle(err, "Restoring working directory: ")
	}

	for _, syncReq := range needsGitList {
		fmt.Printf(clrRed + "x " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}

	for _, syncReq := range needsNothingList {
		fmt.Printf(clrGreen + string(rune(checkMark)) + clrReset + " " + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}

	for _, syncReq := range needsCommitList {
		fmt.Printf(clrPurple + "! " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}

	for _, syncReq := range needsAddList {
		fmt.Printf(clrYellow + "+ " + clrReset + fmt.Sprintf("%-40s", syncReq.Dir) + " " + syncReq.Detail + NL)
	}
}
