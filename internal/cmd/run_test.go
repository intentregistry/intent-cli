package cmd

import (
	"os"
	"testing"
)

func TestRunCommand_Integration(t *testing.T) {
	// Test that run command can be created and has expected flags
	cmd := RunCmd()
	
	if cmd == nil {
		t.Error("RunCmd() returned nil")
	}
	
	// Test that expected flags exist
	if cmd.Flags().Lookup("inputs") == nil {
		t.Error("Expected 'inputs' flag not found in run command")
	}
	
	if cmd.Flags().Lookup("output-dir") == nil {
		t.Error("Expected 'output-dir' flag not found in run command")
	}
	
	if cmd.Flags().Lookup("verbose") == nil {
		t.Error("Expected 'verbose' flag not found in run command")
	}
	
	// Test command usage
	if cmd.Use != "run FILE.itml [--inputs k=v]" {
		t.Errorf("Expected usage 'run FILE.itml [--inputs k=v]', got '%s'", cmd.Use)
	}
}

func TestRunCommand_FileValidation(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "intent-run-test-*")
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
	
	tests := []struct {
		name        string
		args        []string
		expectError bool
		setupFile   bool
		fileContent string
	}{
		{
			name:        "run with nonexistent file",
			args:        []string{"nonexistent.itml"},
			expectError: true,
			setupFile:   false,
		},
		{
			name:        "run with invalid extension",
			args:        []string{"test.txt"},
			expectError: true,
			setupFile:   true,
			fileContent: `{"name": "test"}`,
		},
		{
			name:        "run with valid .itml file",
			args:        []string{"test.itml"},
			expectError: false,
			setupFile:   true,
			fileContent: `{
				"name": "test",
				"version": "1.0.0",
				"description": "Test intent",
				"parameters": [],
				"outputs": []
			}`,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup file if needed
			if tt.setupFile {
				if err := os.WriteFile(tt.args[0], []byte(tt.fileContent), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				defer os.Remove(tt.args[0])
			}
			
			cmd := RunCmd()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
