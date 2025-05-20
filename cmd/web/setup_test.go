package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup code here
	// e.g., initializing configurations, database connections, etc.

	// Run the tests
	os.Exit(m.Run())

	// Teardown code here
	// e.g., closing database connections, cleaning up resources, etc.
}
