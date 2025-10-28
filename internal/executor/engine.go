package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intentregistry/intent-cli/internal/parser"
)

// ExecuteResult represents the result of an intent execution
type ExecuteResult map[string]interface{}

// Execute executes an intent with the given parameters
func Execute(intent *parser.Intent, inputParams map[string]string, outputDir string) (ExecuteResult, error) {
	// Prepare execution context
	context := &ExecutionContext{
		Intent:     intent,
		Inputs:     inputParams,
		OutputDir:  outputDir,
		Results:    make(ExecuteResult),
	}
	
	// Execute based on intent type
	if intent.Script != "" {
		return executeScriptIntent(context)
	}
	
	// Default execution for intents without scripts
	return executeDefaultIntent(context)
}

// ExecutionContext holds the execution state
type ExecutionContext struct {
	Intent    *parser.Intent
	Inputs    map[string]string
	OutputDir string
	Results   ExecuteResult
}

// executeScriptIntent executes an intent with a custom script
func executeScriptIntent(ctx *ExecutionContext) (ExecuteResult, error) {
	// For now, we'll implement a simple script execution
	// In a real implementation, you might support JavaScript, Python, or other languages
	
	script := ctx.Intent.Script
	if strings.HasPrefix(script, "javascript:") {
		return executeJavaScript(script[11:], ctx)
	}
	
	if strings.HasPrefix(script, "python:") {
		return executePython(script[7:], ctx)
	}
	
	// Default to simple template execution
	return executeTemplate(script, ctx)
}

// executeDefaultIntent executes an intent without a custom script
func executeDefaultIntent(ctx *ExecutionContext) (ExecuteResult, error) {
	// Simple execution that processes inputs and generates outputs
	results := make(ExecuteResult)
	
	// Process each output
	for _, output := range ctx.Intent.Outputs {
		switch output.Name {
		case "result":
			results["result"] = processDefaultResult(ctx)
		case "status":
			results["status"] = "success"
		case "message":
			results["message"] = fmt.Sprintf("Intent '%s' executed successfully", ctx.Intent.Name)
		default:
			// Generate output based on output type
			results[output.Name] = generateOutput(output, ctx)
		}
	}
	
	// Save results to output directory if specified
	if ctx.OutputDir != "" {
		if err := saveResults(results, ctx.OutputDir); err != nil {
			return nil, fmt.Errorf("failed to save results: %w", err)
		}
	}
	
	return results, nil
}

// executeJavaScript executes JavaScript code (placeholder)
func executeJavaScript(code string, ctx *ExecutionContext) (ExecuteResult, error) {
	// In a real implementation, you'd use a JavaScript engine like goja
	// For now, return a placeholder result
	return ExecuteResult{
		"result": "JavaScript execution not yet implemented",
		"status": "error",
	}, fmt.Errorf("JavaScript execution not yet implemented")
}

// executePython executes Python code (placeholder)
func executePython(code string, ctx *ExecutionContext) (ExecuteResult, error) {
	// In a real implementation, you'd use a Python interpreter
	// For now, return a placeholder result
	return ExecuteResult{
		"result": "Python execution not yet implemented",
		"status": "error",
	}, fmt.Errorf("Python execution not yet implemented")
}

// executeTemplate executes a simple template
func executeTemplate(template string, ctx *ExecutionContext) (ExecuteResult, error) {
	// Simple template replacement
	result := template
	
	// Create a map of all parameter values (inputs + defaults)
	paramValues := make(map[string]string)
	
	// First, add default values from intent definition
	for _, param := range ctx.Intent.Parameters {
		if param.Default != nil {
			// Convert default value to string
			if defaultStr, ok := param.Default.(string); ok {
				paramValues[param.Name] = defaultStr
			} else {
				// Convert other types to string
				paramValues[param.Name] = fmt.Sprintf("%v", param.Default)
			}
		}
	}
	
	// Then, override with provided input values
	for key, value := range ctx.Inputs {
		paramValues[key] = value
	}
	
	// Replace parameter placeholders
	for key, value := range paramValues {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	
	// Replace intent metadata
	result = strings.ReplaceAll(result, "{{name}}", ctx.Intent.Name)
	result = strings.ReplaceAll(result, "{{description}}", ctx.Intent.Description)
	result = strings.ReplaceAll(result, "{{version}}", ctx.Intent.Version)
	
	return ExecuteResult{
		"result": result,
		"status": "success",
	}, nil
}

// processDefaultResult processes inputs to generate a default result
func processDefaultResult(ctx *ExecutionContext) string {
	if len(ctx.Inputs) == 0 {
		return fmt.Sprintf("Intent '%s' executed with no inputs", ctx.Intent.Name)
	}
	
	var parts []string
	for key, value := range ctx.Inputs {
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	
	return fmt.Sprintf("Intent '%s' processed inputs: %s", ctx.Intent.Name, strings.Join(parts, ", "))
}

// generateOutput generates output based on the output definition
func generateOutput(output parser.Output, ctx *ExecutionContext) interface{} {
	switch output.Type {
	case "string", "text":
		return fmt.Sprintf("Generated %s output", output.Name)
	case "number", "integer", "float":
		return 42
	case "boolean":
		return true
	case "array":
		return []string{"item1", "item2", "item3"}
	case "object", "json":
		return map[string]interface{}{
			"output": output.Name,
			"type":   output.Type,
			"status": "generated",
		}
	default:
		return fmt.Sprintf("Unknown output type: %s", output.Type)
	}
}

// saveResults saves execution results to the output directory
func saveResults(results ExecuteResult, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Save results as JSON
	resultsFile := filepath.Join(outputDir, "results.json")
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	
	if err := os.WriteFile(resultsFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write results file: %w", err)
	}
	
	// Save individual outputs if they're strings
	for name, value := range results {
		if str, ok := value.(string); ok && len(str) > 0 {
			outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.txt", name))
			if err := os.WriteFile(outputFile, []byte(str), 0644); err != nil {
				// Log error but don't fail the entire operation
				fmt.Printf("Warning: failed to save %s: %v\n", name, err)
			}
		}
	}
	
	return nil
}

// ValidateInputs validates input parameters against intent definition
func ValidateInputs(intent *parser.Intent, inputParams map[string]string) error {
	// Check required parameters
	for _, param := range intent.Parameters {
		if param.Required {
			if _, exists := inputParams[param.Name]; !exists {
				return fmt.Errorf("required parameter '%s' not provided", param.Name)
			}
		}
	}
	
	// Validate parameter values
	for name, value := range inputParams {
		// Find parameter definition
		var param *parser.Parameter
		for _, p := range intent.Parameters {
			if p.Name == name {
				param = &p
				break
			}
		}
		
		if param == nil {
			return fmt.Errorf("unknown parameter '%s'", name)
		}
		
		// Validate value
		if err := validateParameterValue(value, param); err != nil {
			return fmt.Errorf("invalid value for parameter '%s': %w", name, err)
		}
	}
	
	return nil
}

// validateParameterValue validates a parameter value against its definition
func validateParameterValue(value string, param *parser.Parameter) error {
	// Basic validation based on type
	switch param.Type {
	case "string", "text":
		// String validation
		if param.Validation.Pattern != "" {
			// In a real implementation, you'd use regex
			// For now, just check if it's not empty
			if value == "" {
				return fmt.Errorf("value cannot be empty")
			}
		}
	case "number", "integer", "float":
		// Number validation
		if param.Validation.Min != nil || param.Validation.Max != nil {
			// In a real implementation, you'd parse the number and check bounds
			// For now, just check if it's not empty
			if value == "" {
				return fmt.Errorf("value cannot be empty")
			}
		}
	case "boolean":
		// Boolean validation
		lowerValue := strings.ToLower(value)
		if lowerValue != "true" && lowerValue != "false" {
			return fmt.Errorf("boolean value must be 'true' or 'false'")
		}
	case "array":
		// Array validation - check if it's valid JSON array
		var arr []interface{}
		if err := json.Unmarshal([]byte(value), &arr); err != nil {
			return fmt.Errorf("invalid array format")
		}
	case "object", "json":
		// Object validation - check if it's valid JSON object
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(value), &obj); err != nil {
			return fmt.Errorf("invalid object format")
		}
	}
	
	// Check options if specified
	if len(param.Validation.Options) > 0 {
		found := false
		for _, option := range param.Validation.Options {
			if value == option {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("value must be one of: %s", strings.Join(param.Validation.Options, ", "))
		}
	}
	
	return nil
}
