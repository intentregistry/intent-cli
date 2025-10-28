package cmd

import (
	"os"
	"testing"
)

func TestTestCommand_Integration(t *testing.T) {
	// Test that test command can be created and has expected flags
	cmd := TestCmd()
	
	if cmd == nil {
		t.Error("TestCmd() returned nil")
	}
	
	// Test that expected flags exist
	expectedFlags := []string{"verbose", "format", "timeout", "parallel", "coverage", "output-dir"}
	for _, flagName := range expectedFlags {
		if cmd.Flags().Lookup(flagName) == nil {
			t.Errorf("Expected '%s' flag not found in test command", flagName)
		}
	}
	
	// Test command usage
	if cmd.Use != "test [path]" {
		t.Errorf("Expected usage 'test [path]', got '%s'", cmd.Use)
	}
}

func TestTestCommand_FileDiscovery(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "intent-test-discovery-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	// Create test files
	testFiles := map[string]string{
		"test.itml": `{
			"name": "test-intent",
			"version": "1.0.0",
			"description": "Test intent",
			"parameters": [],
			"outputs": [],
			"examples": [{"input": {}, "output": {"status": "success"}}]
		}`,
		"test.test.json": `{
			"name": "custom-test",
			"type": "unit",
			"description": "Custom test",
			"input": {},
			"expected": {"status": "success"}
		}`,
		"notest.txt": "not a test file",
	}
	
	for filename, content := range testFiles {
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}
	
	// Test discovery
	tests, err := discoverTests(".")
	if err != nil {
		t.Fatalf("Failed to discover tests: %v", err)
	}
	
	// Should find 2 tests (1 from .itml example, 1 from .test.json)
	if len(tests) != 2 {
		t.Errorf("Expected 2 tests, found %d", len(tests))
	}
	
	// Check test names
	testNames := make(map[string]bool)
	for _, test := range tests {
		testNames[test.Name] = true
	}
	
	expectedNames := []string{"test-intent.example", "custom-test"}
	for _, expectedName := range expectedNames {
		if !testNames[expectedName] {
			t.Errorf("Expected test name '%s' not found", expectedName)
		}
	}
}

func TestTestCommand_OutputComparison(t *testing.T) {
	tests := []struct {
		name     string
		actual   map[string]interface{}
		expected map[string]interface{}
		result   bool
	}{
		{
			name:     "exact match",
			actual:   map[string]interface{}{"status": "success", "result": "hello"},
			expected: map[string]interface{}{"status": "success", "result": "hello"},
			result:   true,
		},
		{
			name:     "greeting alias",
			actual:   map[string]interface{}{"result": "hello world"},
			expected: map[string]interface{}{"greeting": "hello world"},
			result:   true,
		},
		{
			name:     "status flexible",
			actual:   map[string]interface{}{"status": "success"},
			expected: map[string]interface{}{"status": "success"},
			result:   true,
		},
		{
			name:     "template partial match",
			actual:   map[string]interface{}{"result": "Hello {{name}}! Welcome to IntentRegistry"},
			expected: map[string]interface{}{"result": "Hello"},
			result:   true,
		},
		{
			name:     "missing field",
			actual:   map[string]interface{}{"status": "success"},
			expected: map[string]interface{}{"status": "success", "result": "hello"},
			result:   false,
		},
		{
			name:     "no match",
			actual:   map[string]interface{}{"result": "hello"},
			expected: map[string]interface{}{"result": "goodbye"},
			result:   false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareOutputs(tt.actual, tt.expected)
			if result != tt.result {
				t.Errorf("Expected %v, got %v", tt.result, result)
			}
		})
	}
}

func TestTestCommand_NoTestsFound(t *testing.T) {
	// Create temporary directory with no tests
	tempDir, err := os.MkdirTemp("", "intent-test-empty-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	// Create a non-test file
	if err := os.WriteFile("readme.txt", []byte("not a test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	
	// Test discovery should find no tests
	tests, err := discoverTests(".")
	if err != nil {
		t.Fatalf("Failed to discover tests: %v", err)
	}
	
	if len(tests) != 0 {
		t.Errorf("Expected 0 tests, found %d", len(tests))
	}
}

func TestTestCommand_InvalidTestFile(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "intent-test-invalid-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	// Create invalid JSON test file
	if err := os.WriteFile("invalid.test.json", []byte("{invalid json}"), 0644); err != nil {
		t.Fatalf("Failed to create invalid test file: %v", err)
	}
	
	// Test discovery should handle invalid files gracefully
	tests, err := discoverTests(".")
	if err != nil {
		// This is expected for invalid JSON
		return
	}
	
	// Should find 0 tests due to invalid JSON
	if len(tests) != 0 {
		t.Errorf("Expected 0 tests due to invalid JSON, found %d", len(tests))
	}
}
