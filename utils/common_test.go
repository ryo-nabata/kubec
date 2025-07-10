package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	homeDir := GetHomeDir()
	
	// Confirm that home directory is not empty
	if homeDir == "" {
		t.Error("Expected non-empty home directory, but got empty string")
	}
	
	// Confirm that home directory exists
	if _, err := os.Stat(homeDir); os.IsNotExist(err) {
		t.Errorf("Home directory %s does not exist", homeDir)
	}
}

func TestFileExists(t *testing.T) {
	// Create temporary file for testing
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "existing.txt")
	nonExistingFile := filepath.Join(tempDir, "non-existing.txt")
	
	// Create file
	err := os.WriteFile(existingFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Test existing file
	if !FileExists(existingFile) {
		t.Errorf("Expected file %s to exist, but FileExists returned false", existingFile)
	}
	
	// Test non-existing file
	if FileExists(nonExistingFile) {
		t.Errorf("Expected file %s to not exist, but FileExists returned true", nonExistingFile)
	}
}

func TestCreateDirectoryIfNotExists(t *testing.T) {
	tempDir := t.TempDir()
	newDir := filepath.Join(tempDir, "new-directory")
	
	// Confirm that directory does not exist
	if FileExists(newDir) {
		t.Errorf("Directory %s should not exist initially", newDir)
	}
	
	// Create directory
	err := CreateDirectoryIfNotExists(newDir)
	if err != nil {
		t.Errorf("Failed to create directory: %v", err)
	}
	
	// Confirm that directory was created
	if !FileExists(newDir) {
		t.Errorf("Directory %s should exist after creation", newDir)
	}
	
	// Confirm that re-running on existing directory does not cause an error
	err = CreateDirectoryIfNotExists(newDir)
	if err != nil {
		t.Errorf("Should not error when directory already exists: %v", err)
	}
}

func TestCreateDirectoryIfNotExistsNestedPath(t *testing.T) {
	tempDir := t.TempDir()
	nestedDir := filepath.Join(tempDir, "level1", "level2", "level3")
	
	// Create deep directory structure
	err := CreateDirectoryIfNotExists(nestedDir)
	if err != nil {
		t.Errorf("Failed to create nested directory: %v", err)
	}
	
	// Confirm that all levels were created
	if !FileExists(nestedDir) {
		t.Errorf("Nested directory %s should exist after creation", nestedDir)
	}
	
	// Confirm that intermediate directories also exist
	level1 := filepath.Join(tempDir, "level1")
	level2 := filepath.Join(tempDir, "level1", "level2")
	
	if !FileExists(level1) {
		t.Errorf("Intermediate directory %s should exist", level1)
	}
	
	if !FileExists(level2) {
		t.Errorf("Intermediate directory %s should exist", level2)
	}
}

func TestGetKubeDirectory(t *testing.T) {
	kubeDir := GetKubeDirectory()
	expectedDir := filepath.Join(GetHomeDir(), ".kube")
	
	if kubeDir != expectedDir {
		t.Errorf("Expected kube directory %s, but got %s", expectedDir, kubeDir)
	}
}

func TestEnsureKubeDirectory(t *testing.T) {
	// To avoid affecting the actual .kube directory,
	// temporarily use a different home directory
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)
	
	kubeDir := GetKubeDirectory()
	
	// Confirm that it does not exist initially
	if FileExists(kubeDir) {
		// Delete if existing directory exists
		os.RemoveAll(kubeDir)
	}
	
	// Create .kube directory
	err := EnsureKubeDirectory()
	if err != nil {
		t.Errorf("Failed to ensure kube directory: %v", err)
	}
	
	// Confirm that directory was created
	if !FileExists(kubeDir) {
		t.Errorf("Kube directory %s should exist after EnsureKubeDirectory", kubeDir)
	}
	
	// Confirm that re-running on existing directory does not cause an error
	err = EnsureKubeDirectory()
	if err != nil {
		t.Errorf("Should not error when kube directory already exists: %v", err)
	}
}