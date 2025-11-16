//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	ast "github.com/tinybluerobots/jsonforms-parser"
	"github.com/tinybluerobots/jsonforms-html/renderer"
)

func renderHTML(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return map[string]interface{}{
			"error": "renderHTML expects 1 argument: JSON string with schema and uischema",
		}
	}

	jsonString := args[0].String()

	// Parse the input JSON
	var input struct {
		Schema   json.RawMessage `json:"schema"`
		UISchema json.RawMessage `json:"uischema"`
	}

	if err := json.Unmarshal([]byte(jsonString), &input); err != nil {
		return map[string]interface{}{
			"error": "Invalid JSON: " + err.Error(),
		}
	}

	// Validate we have both schema and uischema
	if len(input.Schema) == 0 || len(input.UISchema) == 0 {
		return map[string]interface{}{
			"error": "JSON must contain both 'schema' and 'uischema' properties",
		}
	}

	// Parse using jsonforms-parser
	tree, err := ast.Parse(input.UISchema, input.Schema)
	if err != nil {
		return map[string]interface{}{
			"error": "Failed to parse schemas: " + err.Error(),
		}
	}

	// Render to HTML
	html, err := renderer.Render(tree, nil)
	if err != nil {
		return map[string]interface{}{
			"error": "Failed to render HTML: " + err.Error(),
		}
	}

	return map[string]interface{}{
		"html": html,
	}
}

func main() {
	// Make the Go function available to JavaScript
	js.Global().Set("renderHTML", js.FuncOf(renderHTML))

	// Keep the program running
	select {}
}
