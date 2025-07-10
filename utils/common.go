package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	return homeDir
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CreateDirectoryIfNotExists(path string) error {
	if !FileExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}
	return nil
}

func GetKubeDirectory() string {
	return filepath.Join(GetHomeDir(), ".kube")
}

func EnsureKubeDirectory() error {
	kubeDir := GetKubeDirectory()
	return CreateDirectoryIfNotExists(kubeDir)
}