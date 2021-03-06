package store

import (
	"os"
	"testing"
)

// This file sets up the test fixture.

var tStore DBStore

func TestMain(m *testing.M) {
	st, cleanup := MustGetTempStore()
	tStore = st
	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
}
