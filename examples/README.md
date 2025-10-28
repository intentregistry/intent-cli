# Intent Examples

This directory contains example `.itml` files demonstrating the Intent Markup Language format.

## .itml Format

The `.itml` format is a JSON-based specification for defining AI intents. It includes:

### Required Fields
- `name`: Unique identifier for the intent
- `version`: Semantic version (e.g., "1.0.0")
- `description`: Human-readable description
- `parameters`: Array of input parameters
- `outputs`: Array of expected outputs

### Optional Fields
- `author`: Author information
- `license`: License type (e.g., "MIT")
- `tags`: Array of tags for categorization
- `examples`: Array of usage examples
- `script`: Custom execution script (supports templates)
- `config`: Additional configuration

### Parameter Definition
```json
{
  "name": "parameter_name",
  "type": "string|number|boolean|array|object",
  "description": "Parameter description",
  "required": true|false,
  "default": "default_value",
  "validation": {
    "min": 0,
    "max": 100,
    "pattern": "regex_pattern",
    "options": ["option1", "option2"]
  }
}
```

### Output Definition
```json
{
  "name": "output_name",
  "type": "string|number|boolean|array|object",
  "description": "Output description",
  "format": "optional_format_hint"
}
```

## Examples

### hello-world.itml
A simple greeting intent that demonstrates:
- Required and optional parameters
- Parameter validation with options
- Template-based script execution
- Multiple output types

### text-processor.itml
A text processing intent that demonstrates:
- Complex parameter validation
- Multiple operation types
- Different output formats
- Usage examples

## Usage

```bash
# Run a simple greeting
intent run examples/hello-world.itml --inputs name=Alice --inputs language=en

# Process text with verbose output
intent run examples/text-processor.itml --inputs text="Hello World" --inputs operation=uppercase --verbose

# Save results to directory
intent run examples/hello-world.itml --inputs name=Bob --output-dir ./results
```

## Template Variables

In script templates, you can use the following variables:
- `{{parameter_name}}`: Input parameter values
- `{{name}}`: Intent name
- `{{description}}`: Intent description
- `{{version}}`: Intent version

## Supported Types

- `string`: Text values
- `number`: Numeric values (integers and floats)
- `boolean`: True/false values
- `array`: JSON arrays
- `object`: JSON objects
- `text`: Alias for string
- `integer`: Alias for number
- `float`: Alias for number
- `json`: Alias for object
- `file`: File path values
- `url`: URL values
