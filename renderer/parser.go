package renderer

import (
	"fmt"
	"regexp"
	"strings"

	ast "github.com/tinybluerobots/jsonforms-parser"
)

// RenderOptions configures the HTML rendering
type RenderOptions struct {
	FullPage      bool           // Generate complete HTML page
	InitialData   map[string]any // Pre-fill form values
	FormID        string         // Form element ID
	FormAction    string         // Form action URL
	FormMethod    string         // Form method (GET/POST)
	PageTitle     string         // Page title (for FullPage mode)
	LinkPlacement string         // Where to place links: "footer" or "header" (default: "footer")
}

// DefaultRenderOptions returns default rendering options
func DefaultRenderOptions() *RenderOptions {
	return &RenderOptions{
		FullPage:      false,
		InitialData:   make(map[string]any),
		FormID:        "jsonform",
		FormAction:    "",
		FormMethod:    "POST",
		PageTitle:     "Form",
		LinkPlacement: "footer",
	}
}

// Render converts a JSON Forms AST to HTML
func Render(tree *ast.AST, options *RenderOptions) (string, error) {
	if options == nil {
		options = DefaultRenderOptions()
	}

	// Ensure defaults
	if options.FormID == "" {
		options.FormID = "jsonform"
	}

	if options.FormMethod == "" {
		options.FormMethod = "POST"
	}

	if options.PageTitle == "" {
		options.PageTitle = "Form"
	}

	if options.LinkPlacement == "" {
		options.LinkPlacement = "footer"
	}

	// Convert schema to map
	var schemaMap map[string]any

	if tree.Schema != nil {
		if sm, ok := tree.Schema.(map[string]any); ok {
			schemaMap = sm
		}
	}

	renderer := &HTMLRenderer{
		generator:    NewHTMLGenerator(),
		schemaHelper: NewSchemaHelper(schemaMap),
		options:      options,
		builder:      &strings.Builder{},
	}

	// Render the form
	formHTML := renderer.renderForm(tree.UISchema)

	if options.FullPage {
		return renderer.generator.WrapFullPage(options.PageTitle, formHTML, GetDefaultCSS()), nil
	}

	return formHTML, nil
}

// HTMLRenderer renders AST to HTML
type HTMLRenderer struct {
	generator    *HTMLGenerator
	schemaHelper *SchemaHelper
	options      *RenderOptions
	builder      *strings.Builder
}

// renderForm renders the complete form
func (r *HTMLRenderer) renderForm(uiSchema ast.UISchemaElement) string {
	r.builder.Reset()

	// Open form
	r.builder.WriteString(r.generator.OpenForm(r.options.FormID, r.options.FormAction, r.options.FormMethod))
	r.generator.Indent()

	// Render links in header if placement is "header"
	if r.options.LinkPlacement == "header" {
		r.renderLinks()
	}

	// Render UI schema
	r.renderElement(uiSchema)

	// Render links in footer if placement is "footer" or not specified
	if r.options.LinkPlacement != "header" {
		r.renderLinks()
	}

	r.generator.Dedent()
	r.builder.WriteString(r.generator.CloseForm())

	return r.builder.String()
}

// renderElement renders a UI schema element
func (r *HTMLRenderer) renderElement(element ast.UISchemaElement) {
	if element == nil {
		return
	}

	switch e := element.(type) {
	case *ast.Control:
		r.renderControl(e)
	case *ast.VerticalLayout:
		r.renderVerticalLayout(e)
	case *ast.HorizontalLayout:
		r.renderHorizontalLayout(e)
	case *ast.Group:
		r.renderGroup(e)
	case *ast.Categorization:
		r.renderCategorization(e)
	case *ast.Category:
		r.renderCategory(e)
	case *ast.Label:
		r.renderLabel(e)
	case *ast.CustomElement:
		r.renderCustomElement(e)
	}
}

// renderControl renders a Control element
func (r *HTMLRenderer) renderControl(control *ast.Control) {
	// Extract field name from scope
	fieldName := r.extractFieldName(control.Scope)
	if fieldName == "" {
		return
	}

	// Get property schema
	propertySchema := r.schemaHelper.GetPropertySchema(control.Scope)

	// Determine input type
	inputType := r.schemaHelper.DeriveInputType(propertySchema)

	// Get initial value
	value := r.getInitialValue(control.Scope)

	// Get label text
	labelText := r.getLabelText(control, propertySchema, fieldName)

	// Check if required
	isRequired := r.schemaHelper.IsRequired(control.Scope)

	// Get validation attributes
	validationAttrs := r.schemaHelper.ExtractValidationAttrs(propertySchema)
	if isRequired && inputType != InputTypeCheckbox {
		validationAttrs["required"] = ""
	}

	// Open form group
	r.builder.WriteString(r.generator.OpenDiv("form-group"))
	r.generator.Indent()

	// Render label (except for checkbox)
	if inputType != "checkbox" {
		r.builder.WriteString(r.generator.GenerateLabel(labelText, fieldName, isRequired))
	}

	// Render input element
	switch inputType {
	case "select":
		enumValues := r.schemaHelper.GetEnumValues(propertySchema)
		r.builder.WriteString(r.generator.GenerateSelect(fieldName, enumValues, value, validationAttrs))

	case "textarea":
		r.builder.WriteString(r.generator.GenerateTextarea(fieldName, value, validationAttrs))

	case "checkbox":
		r.builder.WriteString(r.generator.GenerateInput(fieldName, inputType, value, validationAttrs))
		r.builder.WriteString(r.generator.GenerateLabel(labelText, fieldName, isRequired))

	default:
		r.builder.WriteString(r.generator.GenerateInput(fieldName, inputType, value, validationAttrs))
	}

	// Add description if available
	if description := r.schemaHelper.GetDescription(propertySchema); description != "" {
		r.builder.WriteString(r.generator.GenerateDescription(description))
	}

	r.generator.Dedent()
	r.builder.WriteString(r.generator.CloseDiv())
}

// renderVerticalLayout renders a VerticalLayout element
func (r *HTMLRenderer) renderVerticalLayout(layout *ast.VerticalLayout) {
	r.builder.WriteString(r.generator.OpenDiv("vertical-layout"))
	r.generator.Indent()

	for _, child := range layout.Elements {
		r.renderElement(child)
	}

	r.generator.Dedent()
	r.builder.WriteString(r.generator.CloseDiv())
}

// renderHorizontalLayout renders a HorizontalLayout element
func (r *HTMLRenderer) renderHorizontalLayout(layout *ast.HorizontalLayout) {
	r.builder.WriteString(r.generator.OpenDiv("horizontal-layout"))
	r.generator.Indent()

	for _, child := range layout.Elements {
		r.renderElement(child)
	}

	r.generator.Dedent()
	r.builder.WriteString(r.generator.CloseDiv())
}

// renderGroup renders a Group element
func (r *HTMLRenderer) renderGroup(group *ast.Group) {
	r.builder.WriteString(r.generator.OpenFieldset())
	r.generator.Indent()

	r.builder.WriteString(r.generator.GenerateLegend(group.Label))

	for _, child := range group.Elements {
		r.renderElement(child)
	}

	r.generator.Dedent()
	r.builder.WriteString(r.generator.CloseFieldset())
}

// renderCategorization renders a Categorization element
func (r *HTMLRenderer) renderCategorization(cat *ast.Categorization) {
	// For now, render categories as sequential fieldsets
	// A full implementation could use tabs/accordions with JavaScript
	r.builder.WriteString(r.generator.OpenDiv("categorization"))
	r.generator.Indent()

	for _, child := range cat.Elements {
		r.renderElement(child)
	}

	r.generator.Dedent()
	r.builder.WriteString(r.generator.CloseDiv())
}

// renderCategory renders a Category element
func (r *HTMLRenderer) renderCategory(category *ast.Category) {
	r.builder.WriteString(r.generator.OpenFieldset())
	r.generator.Indent()

	r.builder.WriteString(r.generator.GenerateLegend(category.Label))

	for _, child := range category.Elements {
		r.renderElement(child)
	}

	r.generator.Dedent()
	r.builder.WriteString(r.generator.CloseFieldset())
}

// renderLabel renders a Label element
func (r *HTMLRenderer) renderLabel(label *ast.Label) {
	r.builder.WriteString(r.generator.GenerateParagraph(label.Text))
}

// extractFieldName extracts the field name from a JSON Pointer scope
func (r *HTMLRenderer) extractFieldName(scope string) string {
	// Example: "#/properties/name" -> "name"
	// Example: "#/properties/address/properties/city" -> "address_city"
	parts := strings.Split(strings.TrimPrefix(scope, "#/"), "/")

	var fieldParts []string

	for i := 0; i < len(parts); i++ {
		if parts[i] == "properties" && i+1 < len(parts) {
			fieldParts = append(fieldParts, parts[i+1])
			i++ // Skip the property name we just added
		}
	}

	return strings.Join(fieldParts, "_")
}

// getLabelText gets the label text for a control
func (r *HTMLRenderer) getLabelText(control *ast.Control, propertySchema map[string]any, fieldName string) string {
	// Priority: control.Label > schema.title > fieldName
	if control.Label != nil {
		switch label := control.Label.(type) {
		case string:
			return label
		case bool:
			if !label {
				// Label explicitly hidden
				return ""
			}
		case map[string]any:
			// LabelDescription object
			if text, ok := label["text"].(string); ok {
				return text
			}
		}
	}

	// Try schema title
	if title := r.schemaHelper.GetTitle(propertySchema); title != "" {
		return title
	}

	// Fall back to field name (capitalize and humanize)
	return humanizeFieldName(fieldName)
}

// getInitialValue gets the initial value for a field from InitialData
func (r *HTMLRenderer) getInitialValue(scope string) string {
	if r.options.InitialData == nil {
		return ""
	}

	// Extract field name from scope
	fieldName := r.extractFieldName(scope)

	// Look up in initial data
	if val, ok := r.options.InitialData[fieldName]; ok {
		return fmt.Sprintf("%v", val)
	}

	return ""
}

// humanizeFieldName converts "field_name" to "Field Name"
func humanizeFieldName(name string) string {
	// Replace underscores with spaces
	name = strings.ReplaceAll(name, "_", " ")

	// Capitalize first letter of each word
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}

// renderLinks renders hypermedia links from the schema
func (r *HTMLRenderer) renderLinks() {
	links := r.schemaHelper.GetLinks()
	if len(links) == 0 {
		return
	}

	// Open links section
	r.builder.WriteString(r.generator.OpenLinksSection())
	r.generator.Indent()

	// Render all links
	for _, link := range links {
		// Resolve href template
		resolvedHref := resolveHrefTemplate(link.Href, r.options.InitialData)

		// Generate link button
		r.builder.WriteString(r.generator.GenerateLinkButton(link, resolvedHref))
	}

	r.generator.Dedent()
	r.builder.WriteString(r.generator.CloseLinksSection())
}

// renderCustomElement renders a CustomElement
// By default, it renders the element as a comment with its type and recursively renders any child elements
func (r *HTMLRenderer) renderCustomElement(custom *ast.CustomElement) {
	// Add an HTML comment indicating the custom element type
	elementType := custom.GetType()
	r.builder.WriteString(r.generator.GetIndent() + "<!-- Custom element: " + elementType + " -->\n")

	// If the custom element has child elements, render them in a div
	if len(custom.Elements) > 0 {
		r.builder.WriteString(r.generator.OpenDiv("custom-element custom-" + elementType))
		r.generator.Indent()

		for _, child := range custom.Elements {
			r.renderElement(child)
		}

		r.generator.Dedent()
		r.builder.WriteString(r.generator.CloseDiv())
	}
}

// resolveHrefTemplate resolves URI template variables in href
// Supports simple {variable} syntax from RFC 6570
func resolveHrefTemplate(href string, data map[string]any) string {
	if data == nil {
		return href
	}

	// Pattern to match {variable} templates
	templatePattern := regexp.MustCompile(`\{([^}]+)\}`)

	result := templatePattern.ReplaceAllStringFunc(href, func(match string) string {
		// Extract variable name (remove { and })
		varName := strings.Trim(match, "{}")

		// Look up value in data
		if val, ok := data[varName]; ok {
			return fmt.Sprintf("%v", val)
		}

		// If not found, leave the template as-is
		return match
	})

	return result
}
