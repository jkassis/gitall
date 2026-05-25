package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestGitStatiGetClassifiesRepos(t *testing.T) {
	repos := newLocalClone(t)

	clean := GitStatiGet(nil, []string{repos.work})
	if _, ok := clean.NeedsNothingList[repos.work]; !ok {
		t.Fatalf("clean repo should need nothing: %#v", clean)
	}

	if err := os.WriteFile(filepath.Join(repos.work, "README.md"), []byte("dirty\n"), 0644); err != nil {
		t.Fatal(err)
	}
	dirty := GitStatiGet(nil, []string{repos.work})
	if _, ok := dirty.NeedsCommitList[repos.work]; !ok {
		t.Fatalf("dirty repo should need commit: %#v", dirty)
	}
	if _, ok := dirty.NeedsNothingList[repos.work]; ok {
		t.Fatalf("dirty repo must not also be classified clean: %#v", dirty)
	}

	repos = newLocalClone(t)
	commitFile(t, repos.origin, "README.md", "new origin commit\n", "advance origin")
	outOfSync := GitStatiGet(nil, []string{repos.work})
	if _, ok := outOfSync.NeedsSyncList[repos.work]; !ok {
		t.Fatalf("repo behind origin should need sync: %#v", outOfSync)
	}

	missing := filepath.Join(t.TempDir(), "missing")
	errors := GitStatiGet(nil, []string{missing})
	if _, ok := errors.RepoErrorList[missing]; !ok {
		t.Fatalf("missing repo should be reported as repo error: %#v", errors)
	}
}

func TestGitWhatWhereGetReportsBranchAndOrigin(t *testing.T) {
	repos := newLocalClone(t)

	result := GitWhatWhereGet(nil, []string{repos.work})
	status, ok := result[repos.work]
	if !ok {
		t.Fatalf("expected whatwhere result for repo: %#v", result)
	}
	if !strings.Contains(status.Detail, "master of "+repos.origin) {
		t.Fatalf("unexpected whatwhere detail: %q", status.Detail)
	}
}

type localClone struct {
	origin string
	work   string
}

func newLocalClone(t *testing.T) localClone {
	t.Helper()

	root := t.TempDir()
	originPath := filepath.Join(root, "origin")
	workPath := filepath.Join(root, "work")

	if _, err := git.PlainInit(originPath, false); err != nil {
		t.Fatal(err)
	}
	commitFile(t, originPath, "README.md", "initial\n", "initial")

	if _, err := git.PlainClone(workPath, false, &git.CloneOptions{URL: originPath}); err != nil {
		t.Fatal(err)
	}

	return localClone{origin: originPath, work: workPath}
}

func commitFile(t *testing.T, repoPath, name, content, message string) plumbing.Hash {
	t.Helper()

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatal(err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}

	fullPath := filepath.Join(repoPath, name)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := wt.Add(name); err != nil {
		t.Fatal(err)
	}
	hash, err := wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return hash
}
