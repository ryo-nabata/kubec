package utils

import (
	"os"
	"path/filepath"
	"testing"
	"gopkg.in/yaml.v3"
)

func TestGetKubeConfigPath(t *testing.T) {
	// Clear environment variables
	originalKubeconfig := os.Getenv("KUBECONFIG")
	defer func() {
		if originalKubeconfig != "" {
			os.Setenv("KUBECONFIG", originalKubeconfig)
		} else {
			os.Unsetenv("KUBECONFIG")
		}
	}()

	// Test default path
	os.Unsetenv("KUBECONFIG")
	expectedPath := filepath.Join(GetHomeDir(), ".kube", "config")
	actualPath := GetKubeConfigPath()
	if actualPath != expectedPath {
		t.Errorf("Expected %s, but got %s", expectedPath, actualPath)
	}

	// Test when environment variable is set
	customPath := "/custom/path/config"
	os.Setenv("KUBECONFIG", customPath)
	actualPath = GetKubeConfigPath()
	if actualPath != customPath {
		t.Errorf("Expected %s, but got %s", customPath, actualPath)
	}
}

func TestLoadKubeConfig(t *testing.T) {
	// Create temporary file for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config")
	
	// Test KubeConfig data
	testConfig := KubeConfig{
		ApiVersion:     "v1",
		Kind:           "Config",
		CurrentContext: "test-context",
		Contexts: []Context{
			{
				Name: "test-context",
				Context: ContextInfo{
					Cluster: "test-cluster",
					User:    "test-user",
				},
			},
		},
		Clusters: []Cluster{
			{
				Name: "test-cluster",
				Cluster: ClusterInfo{
					Server: "https://test-server",
				},
			},
		},
		Users: []User{
			{
				Name: "test-user",
				User: UserInfo{
					Token: "test-token",
				},
			},
		},
	}

	// Create YAML file
	data, err := yaml.Marshal(&testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Set environment variable
	os.Setenv("KUBECONFIG", configPath)
	defer os.Unsetenv("KUBECONFIG")

	// Test loadKubeConfig
	config, err := loadKubeConfig()
	if err != nil {
		t.Fatalf("Failed to load kubeconfig: %v", err)
	}

	if config.CurrentContext != "test-context" {
		t.Errorf("Expected current-context 'test-context', but got '%s'", config.CurrentContext)
	}

	if len(config.Contexts) != 1 {
		t.Errorf("Expected 1 context, but got %d", len(config.Contexts))
	}

	if config.Contexts[0].Name != "test-context" {
		t.Errorf("Expected context name 'test-context', but got '%s'", config.Contexts[0].Name)
	}
}

func TestGetContexts(t *testing.T) {
	// Create temporary file for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config")
	
	testConfig := KubeConfig{
		ApiVersion:     "v1",
		Kind:           "Config",
		CurrentContext: "context-a",
		Contexts: []Context{
			{Name: "context-b"},
			{Name: "context-a"},
			{Name: "context-c"},
		},
	}

	data, err := yaml.Marshal(&testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	os.Setenv("KUBECONFIG", configPath)
	defer os.Unsetenv("KUBECONFIG")

	contexts := GetContexts()
	
	// Check if sorted in alphabetical order
	expected := []string{"context-a", "context-b", "context-c"}
	if len(contexts) != len(expected) {
		t.Errorf("Expected %d contexts, but got %d", len(expected), len(contexts))
	}

	for i, context := range contexts {
		if context != expected[i] {
			t.Errorf("Expected context '%s' at index %d, but got '%s'", expected[i], i, context)
		}
	}
}

func TestGetCurrentContext(t *testing.T) {
	// Create temporary file for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config")
	
	testConfig := KubeConfig{
		ApiVersion:     "v1",
		Kind:           "Config",
		CurrentContext: "current-test-context",
		Contexts: []Context{
			{Name: "current-test-context"},
		},
	}

	data, err := yaml.Marshal(&testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	os.Setenv("KUBECONFIG", configPath)
	defer os.Unsetenv("KUBECONFIG")

	currentContext := GetCurrentContext()
	if currentContext != "current-test-context" {
		t.Errorf("Expected current context 'current-test-context', but got '%s'", currentContext)
	}
}

func TestSetCurrentContext(t *testing.T) {
	// Create temporary file for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config")
	
	testConfig := KubeConfig{
		ApiVersion:     "v1",
		Kind:           "Config",
		CurrentContext: "old-context",
		Contexts: []Context{
			{Name: "old-context"},
			{Name: "new-context"},
		},
	}

	data, err := yaml.Marshal(&testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	os.Setenv("KUBECONFIG", configPath)
	defer os.Unsetenv("KUBECONFIG")

	// Switch to existing context
	err = SetCurrentContext("new-context")
	if err != nil {
		t.Errorf("Failed to set current context: %v", err)
	}

	// Check if changes are reflected
	currentContext := GetCurrentContext()
	if currentContext != "new-context" {
		t.Errorf("Expected current context 'new-context', but got '%s'", currentContext)
	}

	// Attempt to switch to non-existent context
	err = SetCurrentContext("non-existent-context")
	if err == nil {
		t.Error("Expected error when setting non-existent context, but got none")
	}
}