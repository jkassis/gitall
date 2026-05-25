package main

import (
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestPrvKFilePathGetExpandsDefaultHomePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	v := viper.New()
	c := &cobra.Command{}
	PrvKFilePathFlag(c, v)

	got, err := PrvKFilePathGet(v)
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(home, ".ssh", "id_rsa")
	if got != want {
		t.Fatalf("unexpected default key path: got %q want %q", got, want)
	}
}

func TestPrvKFilePathGetReturnsExplicitPath(t *testing.T) {
	v := viper.New()
	c := &cobra.Command{}
	PrvKFilePathFlag(c, v)
	if err := c.PersistentFlags().Set(SSH_KEY_PATH, "/tmp/custom-key"); err != nil {
		t.Fatal(err)
	}

	got, err := PrvKFilePathGet(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != "/tmp/custom-key" {
		t.Fatalf("unexpected explicit key path: %q", got)
	}
}
