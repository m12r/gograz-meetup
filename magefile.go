//+build mage

package main

import "github.com/magefile/mage/sh"

var Default = Build

// Build builds the gograz-meetup api proxy server
func Build() error {
	return sh.RunWith(
		globalEnv(),
		"go",
		"build", "-o", "bin/gograz-meetup",
	)
}

// Clean cleans the project from previously built binary
func Clean() error {
	return sh.Rm("bin")
}

func globalEnv() map[string]string {
	return map[string]string{
		"GO111MODULE": "on",
	}
}
