package cmd

import (
	"os"
	"testing"
)

func TestSearchCommand_Integration(t *testing.T) {
	// Test that search command can be created and has expected flags
	cmd := SearchCmd()
	
    if cmd == nil {
        t.Fatal("SearchCmd() returned nil")
    }
	
	// Test that JSON flag exists
	if cmd.Flags().Lookup("json") == nil {
		t.Error("Expected 'json' flag not found in search command")
	}
	
	// Test command usage
	if cmd.Use != "search <query>" {
		t.Errorf("Expected usage 'search <query>', got '%s'", cmd.Use)
	}
}

func TestInitCommand_Integration(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "intent-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
    defer func() { _ = os.Chdir(originalDir) }()
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	tests := []struct {
		name        string
		args        []string
		expectError bool
		expectFile  bool
	}{
		{
			name:        "init with name",
			args:        []string{"my-test-intent"},
			expectError: false,
			expectFile:  true,
		},
		{
			name:        "init without name (uses directory name)",
			args:        []string{},
			expectError: false,
			expectFile:  true,
		},
		{
			name:        "init with invalid name",
			args:        []string{"invalid name with spaces"},
			expectError: true,
			expectFile:  false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Remove any existing manifest.yaml
			os.Remove("manifest.yaml")
			
			cmd := InitCmd()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			// Check if manifest.yaml was created
			if _, err := os.Stat("manifest.yaml"); tt.expectFile && err != nil {
				t.Errorf("Expected manifest.yaml to be created: %v", err)
			} else if !tt.expectFile && err == nil {
				t.Error("Expected manifest.yaml not to be created")
			}
		})
	}
}

func TestDoctorCommand_Integration(t *testing.T) {
	// Test that doctor command can be created and has expected flags
	cmd := DoctorCmd()
	
    if cmd == nil {
        t.Fatal("DoctorCmd() returned nil")
    }
	
	// Test that verbose flag exists
	if cmd.Flags().Lookup("verbose") == nil {
		t.Error("Expected 'verbose' flag not found in doctor command")
	}
	
	// Test command usage
	if cmd.Use != "doctor" {
		t.Errorf("Expected usage 'doctor', got '%s'", cmd.Use)
	}
}

func TestRootCommand_Integration(t *testing.T) {
	// Test that all commands are properly registered
	cmd := RootCmd()
	
	expectedCommands := []string{
		"init", "doctor", "login", "publish", "install", 
		"whoami", "search", "version", "completion",
	}
	
	for _, expectedCmd := range expectedCommands {
		found := false
		for _, subCmd := range cmd.Commands() {
			if subCmd.Name() == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command '%s' not found", expectedCmd)
		}
	}
	
	// Test global flags
	if cmd.PersistentFlags().Lookup("debug") == nil {
		t.Error("Expected 'debug' flag not found")
	}
	if cmd.PersistentFlags().Lookup("api-url") == nil {
		t.Error("Expected 'api-url' flag not found")
	}
	if cmd.PersistentFlags().Lookup("telemetry") == nil {
		t.Error("Expected 'telemetry' flag not found")
	}
}

