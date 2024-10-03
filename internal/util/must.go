package util

import (
	"os"

	"golang.org/x/sys/unix"
)

var oldHostname string

func init() {
	if n, err := os.Hostname(); err == nil {
		oldHostname = n
	} else {
		panic(err)
	}
}

func Must(err error) {
	resetHostname()
	if err != nil {
		panic(err)
	}
}

func resetHostname() {
	if err := unix.Sethostname([]byte(oldHostname)); err != nil {
		panic(err)
	}
}
