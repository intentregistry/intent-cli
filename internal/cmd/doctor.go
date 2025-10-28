package cmd

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/intentregistry/intent-cli/internal/httpclient"
	"github.com/spf13/cobra"
)

type CheckResult struct {
	Name        string
	Status      string // "✅", "❌", "⚠️"
	Message     string
	Details     string
	Suggestions []string
}

func DoctorCmd() *cobra.Command {
	var verbose bool
	
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check CLI configuration and connectivity",
		Long: `Run health checks to verify your intent CLI setup.
Checks configuration, API connectivity, authentication, and shell integration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var results []CheckResult
			
			// Check 1: Configuration
			results = append(results, checkConfig())
			
			// Check 2: API Connectivity
			results = append(results, checkAPIConnectivity())
			
			// Check 3: Authentication
			results = append(results, checkAuthentication())
			
			// Check 4: Shell Integration
			results = append(results, checkShellIntegration())
			
			// Check 5: File Permissions
			results = append(results, checkFilePermissions())
			
			// Print results
			printResults(results, verbose)
			
			// Return error if any critical checks failed
			for _, result := range results {
				if result.Status == "❌" {
					return fmt.Errorf("health check failed: %s", result.Message)
				}
			}
			
			return nil
		},
	}
	
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show detailed information for all checks")
	return cmd
}

func checkConfig() CheckResult {
	result := CheckResult{Name: "Configuration"}
	
	cfg := config.Load()
	
	// Check config directory
	configDir := filepath.Join(os.Getenv("HOME"), ".intent")
	if _, err := os.Stat(configDir); err != nil {
		result.Status = "⚠️"
		result.Message = "Config directory not found"
		result.Suggestions = []string{"Run 'intent login' to create configuration"}
		return result
	}
	
	// Check config file
	configFile := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configFile); err != nil {
		result.Status = "⚠️"
		result.Message = "Config file not found"
		result.Suggestions = []string{"Run 'intent login' to create configuration"}
		return result
	}
	
	// Check API URL
	if cfg.APIURL == "" {
		result.Status = "❌"
		result.Message = "API URL not configured"
		result.Suggestions = []string{"Set INTENT_API_URL environment variable or run 'intent login'"}
		return result
	}
	
	result.Status = "✅"
	result.Message = "Configuration looks good"
	result.Details = fmt.Sprintf("API URL: %s", cfg.APIURL)
	return result
}

func checkAPIConnectivity() CheckResult {
	result := CheckResult{Name: "API Connectivity"}
	
	cfg := config.Load()
	if apiURLFlag != "" {
		cfg.APIURL = apiURLFlag
	}
	
	// Parse API URL to get host
	host := strings.TrimPrefix(cfg.APIURL, "https://")
	host = strings.TrimPrefix(host, "http://")
	host = strings.Split(host, "/")[0]
	
	// Test DNS resolution
	_, err := net.LookupHost(host)
	if err != nil {
		result.Status = "❌"
		result.Message = "Cannot resolve API hostname"
		result.Details = fmt.Sprintf("Host: %s, Error: %v", host, err)
		result.Suggestions = []string{"Check your internet connection", "Try: --api-url or INTENT_API_URL"}
		return result
	}
	
	// Test TCP connection
	conn, err := net.DialTimeout("tcp", host+":443", 10*time.Second)
	if err != nil {
		result.Status = "❌"
		result.Message = "Cannot connect to API server"
		result.Details = fmt.Sprintf("Host: %s:443, Error: %v", host, err)
		result.Suggestions = []string{"Check firewall settings", "Try: --api-url or INTENT_API_URL"}
		return result
	}
	conn.Close()
	
	// Test HTTP endpoint
	cl := httpclient.NewWithDebug(cfg, Debug())
	var resp struct {
		Status string `json:"status"`
	}
	err = cl.Get("/v1/health", &resp)
	if err != nil {
		result.Status = "⚠️"
		result.Message = "API endpoint not responding"
		result.Details = err.Error()
		result.Suggestions = []string{"API server may be down", "Check --api-url setting"}
		return result
	}
	
	result.Status = "✅"
	result.Message = "API connectivity is working"
	result.Details = fmt.Sprintf("Connected to %s", cfg.APIURL)
	return result
}

func checkAuthentication() CheckResult {
	result := CheckResult{Name: "Authentication"}
	
	cfg := config.Load()
	if cfg.Token == "" {
		result.Status = "❌"
		result.Message = "No authentication token found"
		result.Suggestions = []string{"Run 'intent login' to authenticate"}
		return result
	}
	
	// Test authentication by calling whoami endpoint
	cl := httpclient.NewWithDebug(cfg, Debug())
	var resp struct {
		User struct {
			Username string `json:"username"`
		} `json:"user"`
	}
	err := cl.Get("/v1/whoami", &resp)
	if err != nil {
		result.Status = "❌"
		result.Message = "Authentication failed"
		result.Details = err.Error()
		result.Suggestions = []string{"Run 'intent login' to refresh your token"}
		return result
	}
	
	result.Status = "✅"
	result.Message = "Authentication is working"
	result.Details = fmt.Sprintf("Logged in as: %s", resp.User.Username)
	return result
}

func checkShellIntegration() CheckResult {
	result := CheckResult{Name: "Shell Integration"}
	
	shell := os.Getenv("SHELL")
	if shell == "" {
		result.Status = "⚠️"
		result.Message = "Shell environment not detected"
		result.Suggestions = []string{"Set SHELL environment variable"}
		return result
	}
	
	// Check if completion is installed
	var completionScript string
	switch {
	case strings.Contains(shell, "bash"):
		completionScript = filepath.Join(os.Getenv("HOME"), ".bashrc")
	case strings.Contains(shell, "zsh"):
		completionScript = filepath.Join(os.Getenv("HOME"), ".zshrc")
	case strings.Contains(shell, "fish"):
		completionScript = filepath.Join(os.Getenv("HOME"), ".config", "fish", "config.fish")
	default:
		result.Status = "⚠️"
		result.Message = "Unsupported shell"
		result.Details = fmt.Sprintf("Shell: %s", shell)
		result.Suggestions = []string{"Shell completion not available for this shell"}
		return result
	}
	
	// Check if completion is configured
	if _, err := os.Stat(completionScript); err == nil {
		content, err := os.ReadFile(completionScript)
		if err == nil && strings.Contains(string(content), "intent completion") {
			result.Status = "✅"
			result.Message = "Shell completion is configured"
			result.Details = fmt.Sprintf("Shell: %s, Script: %s", shell, completionScript)
			return result
		}
	}
	
	result.Status = "⚠️"
	result.Message = "Shell completion not configured"
	result.Details = fmt.Sprintf("Shell: %s", shell)
	result.Suggestions = []string{"Run 'intent completion' to set up shell completion"}
	return result
}

func checkFilePermissions() CheckResult {
	result := CheckResult{Name: "File Permissions"}
	
	// Check config directory permissions
	configDir := filepath.Join(os.Getenv("HOME"), ".intent")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		result.Status = "❌"
		result.Message = "Cannot create config directory"
		result.Details = err.Error()
		result.Suggestions = []string{"Check HOME directory permissions"}
		return result
	}
	
	// Check if we can write to config directory
	testFile := filepath.Join(configDir, ".test-write")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		result.Status = "❌"
		result.Message = "Cannot write to config directory"
		result.Details = err.Error()
		result.Suggestions = []string{"Check config directory permissions"}
		return result
	}
	os.Remove(testFile) // Clean up
	
	// Check current directory permissions
	cwd, err := os.Getwd()
	if err != nil {
		result.Status = "❌"
		result.Message = "Cannot access current directory"
		result.Details = err.Error()
		return result
	}
	
	if err := os.WriteFile(".test-write", []byte("test"), 0644); err != nil {
		result.Status = "⚠️"
		result.Message = "Cannot write to current directory"
		result.Details = err.Error()
		result.Suggestions = []string{"Check current directory permissions for publishing intents"}
		return result
	}
	os.Remove(".test-write") // Clean up
	
	result.Status = "✅"
	result.Message = "File permissions are correct"
	result.Details = fmt.Sprintf("Config dir: %s, Current dir: %s", configDir, cwd)
	return result
}

func printResults(results []CheckResult, verbose bool) {
	fmt.Println("🔍 Intent CLI Health Check")
	fmt.Println(strings.Repeat("=", 40))
	
	allGood := true
	for _, result := range results {
		fmt.Printf("%s %s: %s\n", result.Status, result.Name, result.Message)
		
		if verbose && result.Details != "" {
			fmt.Printf("   Details: %s\n", result.Details)
		}
		
		if len(result.Suggestions) > 0 {
			for _, suggestion := range result.Suggestions {
				fmt.Printf("   💡 %s\n", suggestion)
			}
		}
		
		if result.Status != "✅" {
			allGood = false
		}
		fmt.Println()
	}
	
	if allGood {
		fmt.Println("🎉 All checks passed! Your CLI is ready to use.")
	} else {
		fmt.Println("⚠️  Some issues found. Please address them above.")
	}
}
