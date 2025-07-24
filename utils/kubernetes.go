package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

type KubeConfig struct {
	ApiVersion     string                 `yaml:"apiVersion"`
	Kind           string                 `yaml:"kind"`
	CurrentContext string                 `yaml:"current-context"`
	Contexts       []Context              `yaml:"contexts"`
	Clusters       []Cluster              `yaml:"clusters"`
	Users          []User                 `yaml:"users"`
	Preferences    map[string]interface{} `yaml:"preferences,omitempty"`
}

type Context struct {
	Name    string      `yaml:"name"`
	Context ContextInfo `yaml:"context"`
}

type ContextInfo struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace,omitempty"`
}

type Cluster struct {
	Name    string      `yaml:"name"`
	Cluster ClusterInfo `yaml:"cluster"`
}

type ClusterInfo struct {
	Server                   string `yaml:"server"`
	CertificateAuthority     string `yaml:"certificate-authority,omitempty"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
	InsecureSkipTLSVerify    bool   `yaml:"insecure-skip-tls-verify,omitempty"`
}

type User struct {
	Name string   `yaml:"name"`
	User UserInfo `yaml:"user"`
}

type UserInfo struct {
	Token                string                 `yaml:"token,omitempty"`
	ClientCertificate    string                 `yaml:"client-certificate,omitempty"`
	ClientKey            string                 `yaml:"client-key,omitempty"`
	ClientCertificateData string                `yaml:"client-certificate-data,omitempty"`
	ClientKeyData        string                 `yaml:"client-key-data,omitempty"`
	Username             string                 `yaml:"username,omitempty"`
	Password             string                 `yaml:"password,omitempty"`
	Exec                 *ExecConfig            `yaml:"exec,omitempty"`
	AuthProvider         map[string]interface{} `yaml:"auth-provider,omitempty"`
}

type ExecConfig struct {
	APIVersion         string                 `yaml:"apiVersion,omitempty"`
	Command            string                 `yaml:"command,omitempty"`
	Args               []string               `yaml:"args,omitempty"`
	Env                []ExecEnvVar           `yaml:"env,omitempty"`
	InteractiveMode    string                 `yaml:"interactiveMode,omitempty"`
	ProvideClusterInfo bool                   `yaml:"provideClusterInfo,omitempty"`
}

type ExecEnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func GetKubeConfigPath() string {
	// Use KUBECONFIG environment variable if set
	if kubeconfigPath := os.Getenv("KUBECONFIG"); kubeconfigPath != "" {
		return kubeconfigPath
	}
	
	// Use default path
	homeDir := GetHomeDir()
	return filepath.Join(homeDir, ".kube", "config")
}

func loadKubeConfig() (*KubeConfig, error) {
	configPath := GetKubeConfigPath()
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("kubeconfig file not found: %s", configPath)
	}
	
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig file: %v", err)
	}
	
	var config KubeConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig file: %v", err)
	}
	
	return &config, nil
}

func GetContexts() []string {
	config, err := loadKubeConfig()
	if err != nil {
		log.Fatalf("failed to load kubeconfig: %v", err)
	}
	
	var contexts []string
	for _, context := range config.Contexts {
		contexts = append(contexts, context.Name)
	}
	
	// Sort in alphabetical order
	sort.Strings(contexts)
	
	return contexts
}

func GetCurrentContext() string {
	config, err := loadKubeConfig()
	if err != nil {
		log.Fatalf("failed to load kubeconfig: %v", err)
	}
	
	return config.CurrentContext
}

func SetCurrentContext(contextName string) error {
	config, err := loadKubeConfig()
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %v", err)
	}
	
	// Check if the specified context exists
	found := false
	for _, context := range config.Contexts {
		if context.Name == contextName {
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("context '%s' not found", contextName)
	}
	
	// Update current-context
	config.CurrentContext = contextName
	
	// Write back to file
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to prepare kubeconfig write: %v", err)
	}
	
	configPath := GetKubeConfigPath()
	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write kubeconfig: %v", err)
	}
	
	return nil
}