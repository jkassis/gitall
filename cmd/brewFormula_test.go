package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-github/v49/github"
)

func TestFormulaNewClassifiesAssetsAndHashesPayloads(t *testing.T) {
	payload := tarGz(t, "gitall", "binary")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer server.Close()

	formula, err := FormulaNew("gitall", "https://example.test/gitall", "v1.2.3", []*github.ReleaseAsset{
		{
			Name:               github.String("gitall-darwin-10.10-arm64.tar.gz"),
			BrowserDownloadURL: github.String(server.URL + "/darwin"),
		},
		{
			Name:               github.String("gitall-linux-amd64.tar.gz"),
			BrowserDownloadURL: github.String(server.URL + "/linux"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if formula.Name != "Gitall" {
		t.Fatalf("unexpected formula name: %q", formula.Name)
	}
	if formula.Version != "1.2.3" {
		t.Fatalf("unexpected formula version: %q", formula.Version)
	}
	if len(formula.MacOSDistros) != 1 || len(formula.LinuxDistros) != 1 {
		t.Fatalf("unexpected distro classification: %#v", formula)
	}

	expectedHashBytes := sha256.Sum256(payload)
	expectedHash := hex.EncodeToString(expectedHashBytes[:])
	if formula.MacOSDistros[0].PayloadSHA256 != expectedHash {
		t.Fatalf("unexpected mac payload hash: %q", formula.MacOSDistros[0].PayloadSHA256)
	}
	if formula.LinuxDistros[0].Arch != AMD || formula.LinuxDistros[0].Bits != Six {
		t.Fatalf("unexpected linux arch classification: %#v", formula.LinuxDistros[0])
	}
}

func TestFormulaRenderIncludesInstallRules(t *testing.T) {
	formula := &Formula{
		Name:    "Gitall",
		Desc:    "CLI to perform git operations on multiple repos at once.",
		HomeURL: "https://example.test/gitall",
		Version: "1.2.3",
		MacOSDistros: []*Distro{
			{
				Platform:      Mac,
				Arch:          ARM,
				Bits:          Six,
				PayloadURL:    "https://example.test/gitall-darwin-arm64.tar.gz",
				PayloadSHA256: "abc123",
				BinaryName:    "gitall-darwin-arm64",
				BinaryRename:  "gitall",
			},
		},
	}

	rendered, err := formula.Render()
	if err != nil {
		t.Fatal(err)
	}
	text := string(rendered)
	for _, want := range []string{
		"class Gitall < Formula",
		`homepage "https://example.test/gitall"`,
		`version "1.2.3"`,
		`url "https://example.test/gitall-darwin-arm64.tar.gz"`,
		`sha256 "abc123"`,
		`bin.install "gitall-darwin-arm64" => "gitall"`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("rendered formula missing %q:\n%s", want, text)
		}
	}
}

func tarGz(t *testing.T, name, content string) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0755, Size: int64(len(content))}); err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(tw, content); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
