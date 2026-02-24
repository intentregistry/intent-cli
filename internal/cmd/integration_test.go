package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchCommand_Integration(t *testing.T) {
	// Test that search command can be created and has expected flags
	cmd := SearchCmd()

	if cmd == nil {
		t.Fatal("SearchCmd() returned nil")
		return
	}

	// Test that JSON flag exists
	if cmd.Flags().Lookup("json") == nil {
		t.Error("Expected 'json' flag not found in search command")
	}
	if cmd.Flags().Lookup("limit") == nil {
		t.Error("Expected 'limit' flag not found in search command")
	}
	if cmd.Flags().Lookup("owner") == nil {
		t.Error("Expected 'owner' flag not found in search command")
	}
	if cmd.Flags().Lookup("sort") == nil {
		t.Error("Expected 'sort' flag not found in search command")
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
	defer func() { _ = os.RemoveAll(tempDir) }()

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
		expectFile  string
	}{
		{
			name:        "init with name",
			args:        []string{"my-test-intent"},
			expectError: false,
			expectFile:  filepath.Join("my-test-intent", "itpkg.json"),
		},
		{
			name:        "init without name (uses directory name)",
			args:        []string{},
			expectError: false,
			expectFile:  "itpkg.json",
		},
		{
			name:        "init with invalid name",
			args:        []string{"invalid name with spaces"},
			expectError: true,
			expectFile:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Remove("itpkg.json")
			_ = os.RemoveAll("my-test-intent")

			cmd := InitCmd()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectFile != "" {
				if _, err := os.Stat(tt.expectFile); err != nil {
					t.Errorf("Expected %s to be created: %v", tt.expectFile, err)
				}
			} else {
				if _, err := os.Stat("itpkg.json"); err == nil {
					t.Error("Expected itpkg.json not to be created")
				}
			}
		})
	}
}

func TestDoctorCommand_Integration(t *testing.T) {
	// Test that doctor command can be created and has expected flags
	cmd := DoctorCmd()

	if cmd == nil {
		t.Fatal("DoctorCmd() returned nil")
		return
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
		"whoami", "search", "version", "completion", "run", "package", "test", "verify",
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
