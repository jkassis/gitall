package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	nfpm "github.com/goreleaser/nfpm/v2"
	_ "github.com/goreleaser/nfpm/v2/apk"
	_ "github.com/goreleaser/nfpm/v2/arch"
	_ "github.com/goreleaser/nfpm/v2/deb"
	"github.com/goreleaser/nfpm/v2/files"
	_ "github.com/goreleaser/nfpm/v2/rpm"
)

func main() {
	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: make [build|release]\n")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		usage()
	}

	var err error
	cmd := os.Args[1]
	switch cmd {
	case "build":
		err = build()
	case "release":
		err = release()
	default:
		usage()
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "success")
	os.Exit(0)
}

// This is complicated...
// Currently using Docker image built from github.com/jkassis/xgo
// might want to try https://github.com/crazy-max/goxx
// or https://github.com/techknowlogick/xgo
//
// Needed environment variables:
//
//	DEPS           - Optional list of C dependency packages to build
//	ARGS           - Optional arguments to pass to C dependency configure scripts
//	OUT            - Optional output prefix to override the package name
//	FLAG_V         - Optional verbosity flag to set on the Go builder
//	FLAG_X         - Optional flag to print the build progress commands
//	FLAG_RACE      - Optional race flag to set on the Go builder
//	FLAG_TAGS      - Optional tag flag to set on the Go builder
//	FLAG_LDFLAGS   - Optional ldflags flag to set on the Go builder
//	FLAG_BUILDMODE - Optional buildmode flag to set on the Go builder
//	FLAG_TRIMPATH  - Optional trimpath flag to set on the Go builder
//	TARGETS        - Comma separated list of build targets to compile for
//	GO_VERSION     - Bootstrapped version of Go to disable uncupported targets
//	EXT_GOPATH     - GOPATH elements mounted from the host filesystem
//
// note that cross architecture build with CGO dependencies can only happen
// with dedicated hardware.
//
// build uses docker for cross platform builds
func build() (err error) {
	var pwd string
	pwd, err = os.Getwd()
	if err != nil {
		return err
	}

	xgoCacheDir := os.Getenv("GOPATH") + `/xgo-cache`
	_, err = os.Stat(xgoCacheDir)
	if os.IsNotExist(err) {
		os.MkdirAll(xgoCacheDir, os.ModePerm)
	}

	cmd := exec.Command(
		"docker",
		"run", "--rm",
		"-v", pwd+"/build:/build", // build volume
		"-v", xgoCacheDir+":/deps-cache", //  xgo cache volume
		"-v", pwd+":/source", // sourcecode
		"-e", "FLAG_V=false",
		"-e", "FLAG_X=false",
		"-e", "FLAG_RACE=false ",
		`-e",  "FLAG_LDFLAGS="-w -s"`,
		"-e", "FLAG_BUILDMODE=default ",
		`-e",  "TARGETS="linux/amd64,linux/arm64,darwin/amd64,darwin/arm64,windows/amd64`,
		"jkassis/xgo:1.19.5",
		"./cmd/")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("docker error: %v", err)
	}

	// change permissions of all executables
	return filepath.Walk(pwd+"/build", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(path, 0555)
	})
}

func release() (err error) {
	type Job [4]string

	doOne := func(job Job) error {
		packager := job[0]
		platform := job[1]
		arch := job[2]
		target := job[3]

		fmt.Printf("using %s packager...\n", packager)

		pkg, err := nfpm.Get(packager)
		if err != nil {
			return err
		}

		config := &nfpm.Config{
			Info: nfpm.Info{
				Name:          "gitall",
				Arch:          arch,
				Platform:      platform,
				Version:       "1.0.1",
				VersionSchema: "semver",
				Section:       "default",
				Priority:      "extra",
				Maintainer:    "Jeremy Kassis<jkassis@gmail.com>",
				Description:   "CLI to perform git operations on multiple repos at once.",
				Homepage:      "https://github.com/jkassis/gitall",
				License:       "CC0_1.0",
				Changelog:     "changelog.yaml",
				Overridables: nfpm.Overridables{
					Contents: files.Contents{
						&files.Content{
							Source:      "./build/github.com/jkassis/gitall-linux-" + arch,
							Destination: "/usr/bin/gitall",
						},
					},
				},
			},
		}

		// allow packager to augment config
		info, err := config.Get(packager)
		if err != nil {
			return err
		}
		info = nfpm.WithDefaults(info)

		// open the target file
		target = path.Join(target, pkg.ConventionalFileName(info))
		f, err := os.Create(target)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := pkg.Package(info, f); err != nil {
			return err
		}

		return nil
	}

	jobs := []Job{
		{"archlinux", "linux", "amd64", "./dist"},
		{"archlinux", "linux", "arm64", "./dist"},
		{"apk", "linux", "amd64", "./dist"},
		{"apk", "linux", "arm64", "./dist"},
		{"deb", "linux", "amd64", "./dist"},
		{"deb", "linux", "arm64", "./dist"},
		{"rpm", "linux", "amd64", "./dist"},
		{"rpm", "linux", "arm64", "./dist"},
		{"rpm", "darwin", "amd64", "./dist"},
		{"rpm", "darwin", "arm64", "./dist"},
	}

	for _, job := range jobs {
		err := doOne(job)
		if err != nil {
			return err
		}
	}

	return nil
}
