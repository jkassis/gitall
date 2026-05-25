package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommandRegistersBuildCommands(t *testing.T) {
	for _, name := range []string{"setup", "build", "buildx", "package", "release", "distro"} {
		cmd, _, err := rootCmd.Find([]string{name})
		if err != nil {
			t.Fatalf("find command %q: %v", name, err)
		}
		if cmd == nil || cmd.Name() != name {
			t.Fatalf("command %q not registered; got %#v", name, cmd)
		}
	}
}

func TestExecCapturesStdoutAndStderr(t *testing.T) {
	stdout, stderr, err := Exec("sh", "-c", "printf out; printf err >&2")
	if err != nil {
		t.Fatal(err)
	}
	if stdout != "out" {
		t.Fatalf("unexpected stdout: %q", stdout)
	}
	if stderr != "err" {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
}

func TestCPCopiesIntoDistDirectory(t *testing.T) {
	root := t.TempDir()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("dist", 0755); err != nil {
		t.Fatal(err)
	}

	src := filepath.Join(root, "artifact.deb")
	if err := os.WriteFile(src, []byte("package"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := CP(src, "/ignored/path/artifact.deb"); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(filepath.Join(root, "dist", "artifact.deb"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(got)) != "package" {
		t.Fatalf("unexpected copied content: %q", string(got))
	}
}
