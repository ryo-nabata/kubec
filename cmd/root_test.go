package cmd

import (
	"testing"
)

func TestShouldRunDirectContextSwitch(t *testing.T) {
	// Should return true for normal context names
	normalContexts := []string{
		"my-cluster",
		"production",
		"staging",
		"development",
		"kube-context-1",
		"test-context",
	}
	
	for _, context := range normalContexts {
		if !shouldRunDirectContextSwitch(context) {
			t.Errorf("Expected shouldRunDirectContextSwitch(%s) to return true, but got false", context)
		}
	}
	
	// Should return false for arguments that should be excluded
	excludedArgs := []string{
		"help",
		"version",
		"--help",
		"-h",
		"--version",
		"-v",
		"Help",    // 大文字小文字を区別しない
		"VERSION", // 大文字小文字を区別しない
		"HELP",    // 大文字小文字を区別しない
	}
	
	for _, arg := range excludedArgs {
		if shouldRunDirectContextSwitch(arg) {
			t.Errorf("Expected shouldRunDirectContextSwitch(%s) to return false, but got true", arg)
		}
	}
}

func TestShouldRunDirectContextSwitchCaseInsensitive(t *testing.T) {
	// Confirm case insensitivity
	caseVariations := []string{
		"help",
		"Help",
		"HELP",
		"HeLp",
		"version",
		"Version",
		"VERSION",
		"VeRsIoN",
	}
	
	for _, variation := range caseVariations {
		if shouldRunDirectContextSwitch(variation) {
			t.Errorf("Expected shouldRunDirectContextSwitch(%s) to return false (case insensitive), but got true", variation)
		}
	}
}

func TestShouldRunDirectContextSwitchEdgeCases(t *testing.T) {
	// Edge case testing
	edgeCases := []string{
		"",           // empty string
		" ",          // space
		"  help  ",   // spaces before and after (this is not excluded)
		"help-me",    // contains help but different string
		"version-1",  // contains version but different string
		"my-help",    // contains help but different string
	}
	
	// Anything other than empty string should be treated as normal context name
	for _, testCase := range edgeCases {
		result := shouldRunDirectContextSwitch(testCase)
		if testCase == "" {
			// Behavior for empty string is implementation-dependent, but usually returns true
			if !result {
				t.Errorf("Expected shouldRunDirectContextSwitch(%q) to return true, but got false", testCase)
			}
		} else if testCase == " " || testCase == "  help  " || testCase == "help-me" || testCase == "version-1" || testCase == "my-help" {
			// These should not be excluded
			if !result {
				t.Errorf("Expected shouldRunDirectContextSwitch(%q) to return true, but got false", testCase)
			}
		}
	}
}