package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Intent represents a parsed intent file
type Intent struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	License     string                 `json:"license"`
	Tags        []string               `json:"tags"`
	Parameters  []Parameter            `json:"parameters"`
	Outputs     []Output               `json:"outputs"`
	Examples    []Example              `json:"examples"`
	Script      string                 `json:"script,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// Parameter represents an input parameter for an intent
type Parameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Validation  Validation  `json:"validation,omitempty"`
}

// Output represents an output from an intent
type Output struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Format      string `json:"format,omitempty"`
}

// Example represents a usage example
type Example struct {
	Input  map[string]interface{} `json:"input"`
	Output map[string]interface{} `json:"output"`
}

// Validation represents parameter validation rules
type Validation struct {
	Min     *float64 `json:"min,omitempty"`
	Max     *float64 `json:"max,omitempty"`
	Pattern string   `json:"pattern,omitempty"`
	Options []string `json:"options,omitempty"`
}

// ParseITML parses an .itml file and returns an Intent struct
func ParseITML(filename string) (*Intent, error) {
	// Check file extension
	if !strings.HasSuffix(filename, ".itml") {
		return nil, fmt.Errorf("file must have .itml extension")
	}
	
	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Determine file format and parse accordingly
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".itml":
		return parseITMLFormat(content)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// parseITMLFormat parses the .itml format
func parseITMLFormat(content []byte) (*Intent, error) {
	// Try to parse as JSON first
	var intent Intent
	if err := json.Unmarshal(content, &intent); err == nil {
		// Validate the parsed intent
		if err := validateIntent(&intent); err != nil {
			return nil, fmt.Errorf("invalid intent format: %w", err)
		}
		return &intent, nil
	}
	
	// If JSON parsing fails, try YAML
	return parseYAMLFormat(content)
}

// parseYAMLFormat parses YAML format (fallback)
func parseYAMLFormat(content []byte) (*Intent, error) {
	// For now, we'll implement a simple YAML parser
	// In a real implementation, you'd use a YAML library like gopkg.in/yaml.v3
	return nil, fmt.Errorf("YAML parsing not yet implemented, please use JSON format")
}

// validateIntent validates the parsed intent structure
func validateIntent(intent *Intent) error {
	if intent.Name == "" {
		return fmt.Errorf("intent name is required")
	}
	
	if intent.Version == "" {
		return fmt.Errorf("intent version is required")
	}
	
	if intent.Description == "" {
		return fmt.Errorf("intent description is required")
	}
	
	// Validate parameters
	for i, param := range intent.Parameters {
		if param.Name == "" {
			return fmt.Errorf("parameter %d: name is required", i)
		}
		if param.Type == "" {
			return fmt.Errorf("parameter %d: type is required", i)
		}
		if !isValidType(param.Type) {
			return fmt.Errorf("parameter %d: invalid type '%s'", i, param.Type)
		}
	}
	
	// Validate outputs
	for i, output := range intent.Outputs {
		if output.Name == "" {
			return fmt.Errorf("output %d: name is required", i)
		}
		if output.Type == "" {
			return fmt.Errorf("output %d: type is required", i)
		}
		if !isValidType(output.Type) {
			return fmt.Errorf("output %d: invalid type '%s'", i, output.Type)
		}
	}
	
	return nil
}

// isValidType checks if a type is valid
func isValidType(t string) bool {
	validTypes := []string{
		"string", "number", "boolean", "array", "object",
		"integer", "float", "text", "json", "file", "url",
	}
	
	for _, validType := range validTypes {
		if t == validType {
			return true
		}
	}
	return false
}

// GetParameterValue returns the value of a parameter with type conversion
func (i *Intent) GetParameterValue(name string, inputParams map[string]string) (interface{}, error) {
	// Find the parameter definition
	var param *Parameter
	for _, p := range i.Parameters {
		if p.Name == name {
			param = &p
			break
		}
	}
	
	if param == nil {
		return nil, fmt.Errorf("parameter '%s' not defined", name)
	}
	
	// Get the input value
	inputValue, exists := inputParams[name]
	if !exists {
		if param.Required {
			return nil, fmt.Errorf("required parameter '%s' not provided", name)
		}
		return param.Default, nil
	}
	
	// Convert to appropriate type
	return convertToType(inputValue, param.Type)
}

// convertToType converts a string value to the specified type
func convertToType(value, targetType string) (interface{}, error) {
	switch targetType {
	case "string", "text":
		return value, nil
	case "number", "integer", "float":
		// Simple number parsing - in a real implementation, you'd use strconv
		return value, nil // For now, return as string
	case "boolean":
		return strings.ToLower(value) == "true", nil
	case "array":
		// Parse as JSON array
		var result []interface{}
		if err := json.Unmarshal([]byte(value), &result); err != nil {
			return nil, fmt.Errorf("invalid array format: %w", err)
		}
		return result, nil
	case "object", "json":
		// Parse as JSON object
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(value), &result); err != nil {
			return nil, fmt.Errorf("invalid object format: %w", err)
		}
		return result, nil
	default:
		return value, nil
	}
}
