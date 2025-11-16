package renderer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Input type constants
const (
	InputTypeText     = "text"
	InputTypeCheckbox = "checkbox"
	InputTypeNumber   = "number"
)

// SchemaHelper provides utilities for working with JSON Schema
type SchemaHelper struct {
	schema map[string]any
}

// NewSchemaHelper creates a new schema helper
func NewSchemaHelper(schema map[string]any) *SchemaHelper {
	return &SchemaHelper{schema: schema}
}

// GetPropertySchema resolves a JSON Pointer scope to a property schema
func (h *SchemaHelper) GetPropertySchema(scope string) map[string]any {
	if h.schema == nil {
		return nil
	}

	// Parse JSON Pointer: "#/properties/name" or "#/properties/address/properties/city"
	path := strings.TrimPrefix(scope, "#/")
	if path == "" {
		return h.schema
	}

	parts := strings.Split(path, "/")
	current := h.schema

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Navigate into the schema
		if next, ok := current[part]; ok {
			if nextMap, ok := next.(map[string]any); ok {
				current = nextMap
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	return current
}

// DeriveInputType derives the HTML input type from a property schema
func (h *SchemaHelper) DeriveInputType(propertySchema map[string]any) string {
	if propertySchema == nil {
		return InputTypeText
	}

	// Check for enum -> select
	if _, hasEnum := propertySchema["enum"]; hasEnum {
		return "select"
	}

	// Check type
	schemaType, _ := propertySchema["type"].(string)

	switch schemaType {
	case "boolean":
		return InputTypeCheckbox
	case "integer":
		return InputTypeNumber
	case "number":
		return InputTypeNumber
	case "string":
		// Check for format
		format, _ := propertySchema["format"].(string)
		switch format {
		case "email":
			return "email"
		case "uri", "url":
			return "url"
		case "date":
			return "date"
		case "time":
			return "time"
		case "date-time":
			return "datetime-local"
		case "color":
			return "color"
		case "tel", "telephone":
			return "tel"
		}

		// Check for multi-line
		if multiLine, ok := propertySchema["multiLine"].(bool); ok && multiLine {
			return "textarea"
		}

		return InputTypeText
	default:
		return InputTypeText
	}
}

// GetEnumValues extracts enum values from a property schema
func (h *SchemaHelper) GetEnumValues(propertySchema map[string]any) []string {
	if propertySchema == nil {
		return nil
	}

	enumVal, ok := propertySchema["enum"]
	if !ok {
		return nil
	}

	// Try to convert to []any first
	enumSlice, ok := enumVal.([]any)
	if !ok {
		return nil
	}

	var result []string
	for _, val := range enumSlice {
		result = append(result, stringValue(val))
	}

	return result
}

// GetTitle extracts the title from a property schema
func (h *SchemaHelper) GetTitle(propertySchema map[string]any) string {
	if propertySchema == nil {
		return ""
	}

	if title, ok := propertySchema["title"].(string); ok {
		return title
	}

	return ""
}

// GetDescription extracts the description from a property schema
func (h *SchemaHelper) GetDescription(propertySchema map[string]any) string {
	if propertySchema == nil {
		return ""
	}

	if desc, ok := propertySchema["description"].(string); ok {
		return desc
	}

	return ""
}

// IsRequired checks if a field is required based on its scope
func (h *SchemaHelper) IsRequired(scope string) bool {
	if h.schema == nil {
		return false
	}

	// Extract the property name from scope
	parts := strings.Split(strings.TrimPrefix(scope, "#/"), "/")

	var propertyName string

	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "properties" && parts[i] != "" {
			propertyName = parts[i]
			break
		}
	}

	if propertyName == "" {
		return false
	}

	// Check if this property is in the required array
	required, ok := h.schema["required"]
	if !ok {
		return false
	}

	requiredSlice, ok := required.([]any)
	if !ok {
		return false
	}

	for _, req := range requiredSlice {
		if reqStr, ok := req.(string); ok && reqStr == propertyName {
			return true
		}
	}

	return false
}

// ExtractValidationAttrs extracts HTML5 validation attributes from schema
func (h *SchemaHelper) ExtractValidationAttrs(propertySchema map[string]any) map[string]string {
	attrs := make(map[string]string)

	if propertySchema == nil {
		return attrs
	}

	// Min/max for numbers
	if min, ok := propertySchema["minimum"]; ok {
		attrs["min"] = stringValue(min)
	}

	if max, ok := propertySchema["maximum"]; ok {
		attrs["max"] = stringValue(max)
	}

	// Min/max length for strings
	if minLen, ok := propertySchema["minLength"]; ok {
		attrs["minlength"] = stringValue(minLen)
	}

	if maxLen, ok := propertySchema["maxLength"]; ok {
		attrs["maxlength"] = stringValue(maxLen)
	}

	// Pattern
	if pattern, ok := propertySchema["pattern"].(string); ok {
		attrs["pattern"] = pattern
	}

	return attrs
}

// HyperSchemaLink represents a hypermedia link from JSON Hyper-Schema
type HyperSchemaLink struct {
	Href   string
	Rel    string
	Title  string
	Method string
}

// getDefaultLinkTitle returns a default title based on the rel attribute
func getDefaultLinkTitle(rel string) string {
	switch rel {
	case "submit":
		return "Submit"
	case "create":
		return "Create"
	case "prev", "previous":
		return "Previous"
	case "next":
		return "Next"
	case "cancel":
		return "Cancel"
	case "delete", "remove":
		return "Delete"
	case "edit", "update":
		return "Edit"
	case "self":
		return "View"
	default:
		// Capitalize first letter of rel
		if rel != "" {
			return strings.ToUpper(rel[:1]) + rel[1:]
		}

		return "Link"
	}
}

// GetLinks extracts hypermedia links from the schema
func (h *SchemaHelper) GetLinks() []HyperSchemaLink {
	if h.schema == nil {
		return nil
	}

	linksVal, ok := h.schema["links"]
	if !ok {
		return nil
	}

	linksSlice, ok := linksVal.([]any)
	if !ok {
		return nil
	}

	var result []HyperSchemaLink

	for _, linkVal := range linksSlice {
		linkMap, ok := linkVal.(map[string]any)
		if !ok {
			continue
		}

		link := HyperSchemaLink{
			Href:   stringValue(linkMap["href"]),
			Rel:    stringValue(linkMap["rel"]),
			Title:  stringValue(linkMap["title"]),
			Method: strings.ToUpper(stringValue(linkMap["method"])),
		}

		// Default method to GET if not specified
		if link.Method == "" {
			link.Method = "GET"
		}

		// Default title based on rel if not specified
		if link.Title == "" {
			link.Title = getDefaultLinkTitle(link.Rel)
		}

		result = append(result, link)
	}

	return result
}

// GetLinkByRel finds a link by its rel attribute
func (h *SchemaHelper) GetLinkByRel(rel string) *HyperSchemaLink {
	links := h.GetLinks()
	for _, link := range links {
		if link.Rel == rel {
			return &link
		}
	}

	return nil
}

// stringValue converts any value to string
func stringValue(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case float64:
		return fmt.Sprintf("%v", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case bool:
		if v {
			return "true"
		}

		return "false"
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}
