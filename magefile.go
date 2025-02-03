//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(Clean)
	fmt.Println("Building...")
	return buildBinary(runtime.GOOS, runtime.GOARCH)
}

// Build the release binaries
func Release() error {
	mg.Deps(Clean)
	fmt.Println("Building...")

	platforms := []string{"linux", "darwin", "windows"}
	archs := []string{"amd64"}

	for _, platform := range platforms {
		for _, arch := range archs {
			if err := buildBinary(platform, arch, platform, arch); err != nil {
				return fmt.Errorf("failed to build for %s/%s: %w", platform, arch, err)
			}
		}
	}

	return nil
}

func buildBinary(platform, arch string, suffixes ...string) error {
	fmt.Println("Building for", platform)

	// Get the current commit hash
	hash, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return fmt.Errorf("failed to get commit hash: %w", err)
	}
	date := time.Now().Format(time.RFC3339)

	fmt.Println("Commit hash:", strings.TrimSpace(string(hash)))
	fmt.Println("Build date:", date)

	buildFlags := fmt.Sprintf("-X main.Commit=%s -X main.Date=%s", hash, date)

	// Store the current GOOS and GOARCH
	defer func() {
		fmt.Println("Restoring GOOS and GOARCH to", runtime.GOOS, runtime.GOARCH)
		os.Setenv("GOOS", runtime.GOOS)
		os.Setenv("GOARCH", runtime.GOARCH)
	}()

	// Set the GOOS and GOARCH to the desired platform
	os.Setenv("GOOS", platform)
	os.Setenv("GOARCH", arch)

	output := "bin/dotmanager"

	if len(suffixes) > 0 {
		output += "-" + strings.Join(suffixes, "-")
	}

	cmd := exec.Command("go", "build", "-o", output, "-ldflags", buildFlags, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("./bin/dotmanager", "/usr/bin/dotmanager")
}

// Clean up after yourself
func Clean() error {
	fmt.Println("Cleaning...")
	return os.RemoveAll("bin")
}
