package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Intent represents a parsed intent file
type Intent struct {
	Name        string                 `json:"name" yaml:"name"`
	Version     string                 `json:"version" yaml:"version"`
	Description string                 `json:"description" yaml:"description"`
	Author      string                 `json:"author" yaml:"author"`
	License     string                 `json:"license" yaml:"license"`
	Tags        []string               `json:"tags" yaml:"tags"`
	Parameters  []Parameter            `json:"parameters" yaml:"parameters"`
	Outputs     []Output               `json:"outputs" yaml:"outputs"`
	Examples    []Example              `json:"examples" yaml:"examples"`
	Script      string                 `json:"script,omitempty" yaml:"script,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
}

// Parameter represents an input parameter for an intent
type Parameter struct {
	Name        string      `json:"name" yaml:"name"`
	Type        string      `json:"type" yaml:"type"`
	Description string      `json:"description" yaml:"description"`
	Required    bool        `json:"required" yaml:"required"`
	Default     interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	Validation  Validation  `json:"validation,omitempty" yaml:"validation,omitempty"`
}

// Output represents an output from an intent
type Output struct {
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"`
	Description string `json:"description" yaml:"description"`
	Format      string `json:"format,omitempty" yaml:"format,omitempty"`
}

// Example represents a usage example
type Example struct {
	Input  map[string]interface{} `json:"input" yaml:"input"`
	Output map[string]interface{} `json:"output" yaml:"output"`
}

// Validation represents parameter validation rules
type Validation struct {
	Min     *float64 `json:"min,omitempty" yaml:"min,omitempty"`
	Max     *float64 `json:"max,omitempty" yaml:"max,omitempty"`
	Pattern string   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Options []string `json:"options,omitempty" yaml:"options,omitempty"`
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
	// First, try to parse as custom ITML format
	if intent, err := parseCustomITMLFormat(content); err == nil {
		return intent, nil
	}
	
	// If custom ITML parsing fails, try JSON
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

// parseCustomITMLFormat parses the custom ITML DSL format
func parseCustomITMLFormat(content []byte) (*Intent, error) {
	lines := strings.Split(string(content), "\n")
	
	intent := &Intent{
		Version:     "1.0.0",
		Author:      "Unknown",
		License:     "MIT",
		Parameters:  []Parameter{},
		Outputs:     []Output{},
		Examples:    []Example{},
		Config:      make(map[string]interface{}),
	}
	
	var currentSection string
	var workflowSteps []string
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Parse intent declaration
		if strings.HasPrefix(line, "intent ") {
			// Extract intent name from quotes
			if strings.Contains(line, `"`) {
				start := strings.Index(line, `"`)
				end := strings.LastIndex(line, `"`)
				if start != -1 && end != -1 && end > start {
					intent.Name = line[start+1:end]
					intent.Description = intent.Name // Use name as description for now
				} else {
					return nil, fmt.Errorf("invalid intent declaration on line %d: %s", i+1, line)
				}
			} else {
				return nil, fmt.Errorf("intent name must be quoted on line %d: %s", i+1, line)
			}
			continue
		}
		
		// Parse section headers
		if strings.HasSuffix(line, ":") {
			currentSection = strings.TrimSuffix(line, ":")
			continue
		}
		
		// Parse inputs section
		if currentSection == "inputs" {
			if strings.HasPrefix(line, "- ") {
				param, err := parseParameter(line[2:])
				if err != nil {
					return nil, fmt.Errorf("invalid parameter on line %d: %w", i+1, err)
				}
				intent.Parameters = append(intent.Parameters, param)
			}
			continue
		}
		
		// Parse workflow section
		if currentSection == "workflow" {
			if strings.HasPrefix(line, "→ ") {
				workflowSteps = append(workflowSteps, line[2:])
			}
			continue
		}
	}
	
	// Convert workflow steps to script
	if len(workflowSteps) > 0 {
		intent.Script = strings.Join(workflowSteps, "\n")
	}
	
	// Add default outputs if none specified
	if len(intent.Outputs) == 0 {
		intent.Outputs = []Output{
			{Name: "result", Type: "string", Description: "Execution result"},
			{Name: "status", Type: "string", Description: "Execution status"},
		}
	}
	
	// Validate the parsed intent
	if err := validateIntent(intent); err != nil {
		return nil, fmt.Errorf("invalid intent format: %w", err)
	}
	
	return intent, nil
}

// parseParameter parses a parameter definition like "name (string) default=\"World\""
func parseParameter(paramStr string) (Parameter, error) {
	param := Parameter{
		Required: false,
	}
	
	// Extract parameter name
	parts := strings.Fields(paramStr)
	if len(parts) == 0 {
		return param, fmt.Errorf("empty parameter definition")
	}
	
	param.Name = parts[0]
	
	// Extract type from parentheses
	for _, part := range parts {
		if strings.HasPrefix(part, "(") && strings.HasSuffix(part, ")") {
			param.Type = strings.Trim(part, "()")
			break
		}
	}
	
	// Extract default value
	if strings.Contains(paramStr, "default=") {
		start := strings.Index(paramStr, "default=")
		if start != -1 {
			valuePart := paramStr[start+8:] // Skip "default="
			if strings.HasPrefix(valuePart, `"`) {
				end := strings.Index(valuePart[1:], `"`)
				if end != -1 {
					param.Default = valuePart[1:end+1]
				}
			} else {
				// Handle unquoted default values
				spaceIndex := strings.Index(valuePart, " ")
				if spaceIndex != -1 {
					param.Default = valuePart[:spaceIndex]
				} else {
					param.Default = valuePart
				}
			}
		}
	}
	
	// Set description
	param.Description = fmt.Sprintf("Parameter %s", param.Name)
	
	return param, nil
}

// parseYAMLFormat parses YAML format (fallback)
func parseYAMLFormat(content []byte) (*Intent, error) {
	var intent Intent
	if err := yaml.Unmarshal(content, &intent); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	
	// Validate the parsed intent
	if err := validateIntent(&intent); err != nil {
		return nil, fmt.Errorf("invalid intent format: %w", err)
	}
	
	return &intent, nil
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
