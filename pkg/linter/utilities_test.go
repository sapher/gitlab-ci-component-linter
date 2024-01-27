package linter

import (
	"os"
	"testing"
)

func TestIsFileExist(t *testing.T) {
	// Setup: Create a temporary file and directory to test with
	tmpFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after the test

	tmpDir, err := os.MkdirTemp("", "exampleDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Clean up after the test

	// Define test cases
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"ExistingFile", tmpFile.Name(), true},
		{"NonExistingFile", "nonexistingfile.txt", false},
		{"ExistingDirectory", tmpDir, false},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsFileExist(tc.path)
			if got != tc.want {
				t.Errorf("IsFileExist(%q) = %v; want %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestIsDirExist(t *testing.T) {
	// Setup: Create a temporary file and directory to test with
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after the test

	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Clean up after the test

	// Define test cases
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"ExistingDirectory", tmpDir, true},
		{"NonExistingDirectory", "nonexistingdirectory", false},
		{"ExistingFile", tmpFile.Name(), false},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsDirExist(tc.path)
			if got != tc.want {
				t.Errorf("IsDirExist(%q) = %v; want %v", tc.path, got, tc.want)
			}
		})
	}
}
