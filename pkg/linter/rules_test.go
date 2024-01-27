package linter

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRulesNumber(t *testing.T) {
	jsonRules, err := ParseJsonRules()
	if err != nil {
		t.Fatalf("Failed to parse JSON rules: %v", err)
	}
	if len(ruleFuncs) != len(jsonRules) {
		t.Errorf("Rules number mismatch: %d != %d", len(ruleFuncs), len(jsonRules))
	}
}

func TestMissingRootReadmeRule(t *testing.T) {
	t.Run("No README.md present", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "testdir")
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %v", err)
		}

		linter := &Linter{
			Workdir: tempDir,
		}

		result, err := MissingRootReadmeRule(linter)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.Success != false {
			t.Errorf("Expected test to fail when README.md is missing, got success")
		}
	})

	t.Run("README.md present", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "testdir")
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %v", err)
		}

		// Create a README.md file
		readmePath := path.Join(tempDir, "README.md")
		if err := os.WriteFile(readmePath, []byte("# Readme"), 0644); err != nil {
			t.Fatalf("Failed to create README.md: %v", err)
		}

		linter := &Linter{
			Workdir: tempDir,
		}

		result, err := MissingRootReadmeRule(linter)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.Success != true {
			t.Errorf("Expected test to succeed when README.md is present, got failure")
		}
	})
}

func TestMissingRootTemplatesDirRule(t *testing.T) {
	t.Run("templates directory not present", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "testdir_notemplates")
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %v", err)
		}

		linter := &Linter{
			Workdir: tempDir,
		}

		result, err := MissingRootTemplatesDirRule(linter)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.Success != false {
			t.Errorf("Expected failure when 'templates' directory is missing, got success instead")
		}
	})

	t.Run("templates directory is present", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "testdir_templates")
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %v", err)
		}

		// Create the templates directory within the temporary directory
		templatesDirPath := path.Join(tempDir, "templates")
		err = os.Mkdir(templatesDirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create 'templates' directory: %v", err)
		}

		linter := &Linter{
			Workdir: tempDir,
		}

		result, err := MissingRootTemplatesDirRule(linter)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.Success != true {
			t.Errorf("Expected success when 'templates' directory is present, got failure instead")
		}
	})
}

func TestWrongYamlFileExtensionRule(t *testing.T) {
	setupTemporaryDirectoryWithFiles := func(dirSuffix string, filenames []string) (cleanupFunc func(), tempDir string) {
		tempDir, err := os.MkdirTemp("", fmt.Sprintf("testdir-%s", dirSuffix))
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %v", err)
		}

		for _, filename := range filenames {
			filepath := filepath.Join(tempDir, filename)
			file, err := os.Create(filepath)
			if err != nil {
				t.Fatalf("Failed to create file %s: %v", filepath, err)
			}
			file.Close() // Close the file after creating it
		}

		// Return cleanup function and temporary directory path.
		return func() { os.RemoveAll(tempDir) }, tempDir
	}

	t.Run("no .yaml files", func(t *testing.T) {
		cleanup, tempDir := setupTemporaryDirectoryWithFiles("no", []string{"valid.yml", "another.yml", "not_yaml.txt"})
		defer cleanup()

		linter := &Linter{Workdir: tempDir}

		result, err := WrongYamlFileExtensionRule(linter)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.Success != true {
			t.Errorf("Expected Success to be true when there are no .yaml files, got false")
		}
	})

	t.Run(".yaml files present", func(t *testing.T) {
		cleanup, tempDir := setupTemporaryDirectoryWithFiles("yes", []string{"wrong.yaml", "another.yaml", "valid.yml"})
		defer cleanup()

		linter := &Linter{Workdir: tempDir}

		result, err := WrongYamlFileExtensionRule(linter)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("Expected Success to be false when .yaml files are present, got true")
		}
		if len(result.Metadata["files"].([]string)) == 0 {
			t.Errorf("Expected Metadata files to contain .yaml filenames, got empty list")
		}
		// Add more assertions here to verify content of the "files" slice if needed.
	})
}
