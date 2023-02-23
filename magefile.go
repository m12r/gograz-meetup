//go:build mage
// +build mage

package main

import "github.com/magefile/mage/sh"

var Default = Build

// Build builds the gograz-meetup api proxy server
func Build() error {
	return sh.RunWith(
		globalEnv(),
		"go",
		"build",
		"-trimpath",
		"-ldflags", "-s -w",
		"-o", "bin/gograz-meetup",
	)
}

// Clean cleans the project from previously built binary
func Clean() error {
	return sh.Rm("bin")
}

func globalEnv() map[string]string {
	return map[string]string{
		"CGO_ENABLED": "0",
	}
}
