package renderer

import (
	"fmt"
	"html"
	"strings"
)

// HTMLGenerator provides utilities for generating HTML elements
type HTMLGenerator struct {
	indent int
}

// NewHTMLGenerator creates a new HTML generator
func NewHTMLGenerator() *HTMLGenerator {
	return &HTMLGenerator{indent: 0}
}

// Indent increases the indentation level
func (g *HTMLGenerator) Indent() {
	g.indent++
}

// Dedent decreases the indentation level
func (g *HTMLGenerator) Dedent() {
	if g.indent > 0 {
		g.indent--
	}
}

// GetIndent returns the current indentation string
func (g *HTMLGenerator) GetIndent() string {
	return strings.Repeat("  ", g.indent)
}

// GenerateInput generates an HTML input element
func (g *HTMLGenerator) GenerateInput(name, inputType, value string, attrs map[string]string) string {
	var parts []string

	if inputType == InputTypeCheckbox {
		parts = append(parts, fmt.Sprintf(`<input type="%s" id="%s" name="%s"`, inputType, name, name))

		// For checkboxes, value should be "checked" or empty
		if value == "true" || value == "checked" {
			parts = append(parts, `checked`)
		}
	} else {
		parts = append(parts, fmt.Sprintf(`<input type="%s" id="%s" name="%s"`, inputType, name, name))

		if value != "" {
			parts = append(parts, fmt.Sprintf(`value="%s"`, html.EscapeString(value)))
		}
	}

	// Add additional attributes
	for key, val := range attrs {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, key, html.EscapeString(val)))
	}

	return g.GetIndent() + strings.Join(parts, " ") + ">\n"
}

// GenerateTextarea generates an HTML textarea element
func (g *HTMLGenerator) GenerateTextarea(name, value string, attrs map[string]string) string {
	var parts []string

	parts = append(parts, fmt.Sprintf(`<textarea id="%s" name="%s"`, name, name))

	// Add additional attributes
	for key, val := range attrs {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, key, html.EscapeString(val)))
	}

	parts = append(parts, ">")

	result := g.GetIndent() + strings.Join(parts, " ")
	if value != "" {
		result += html.EscapeString(value)
	}

	result += "</textarea>\n"

	return result
}

// GenerateSelect generates an HTML select element
func (g *HTMLGenerator) GenerateSelect(name string, options []string, value string, attrs map[string]string) string {
	var builder strings.Builder

	var parts []string

	parts = append(parts, fmt.Sprintf(`<select id="%s" name="%s"`, name, name))

	// Add additional attributes
	for key, val := range attrs {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, key, html.EscapeString(val)))
	}

	builder.WriteString(g.GetIndent() + strings.Join(parts, " ") + ">\n")

	g.Indent()
	// Add placeholder option
	builder.WriteString(g.GetIndent() + `<option value="">-- Select --</option>` + "\n")

	// Add options
	for _, opt := range options {
		selected := ""
		if opt == value {
			selected = " selected"
		}

		builder.WriteString(g.GetIndent() + fmt.Sprintf(`<option value="%s"%s>%s</option>`, html.EscapeString(opt), selected, html.EscapeString(opt)) + "\n")
	}

	g.Dedent()

	builder.WriteString(g.GetIndent() + "</select>\n")

	return builder.String()
}

// GenerateLabel generates an HTML label element
func (g *HTMLGenerator) GenerateLabel(text, forID string, required bool) string {
	requiredMark := ""
	if required {
		requiredMark = ` <span class="required">*</span>`
	}

	return g.GetIndent() + fmt.Sprintf(`<label for="%s">%s%s</label>`, forID, html.EscapeString(text), requiredMark) + "\n"
}

// OpenDiv opens a div element with optional class
func (g *HTMLGenerator) OpenDiv(class string) string {
	if class != "" {
		return g.GetIndent() + fmt.Sprintf(`<div class="%s">`, class) + "\n"
	}

	return g.GetIndent() + "<div>\n"
}

// CloseDiv closes a div element
func (g *HTMLGenerator) CloseDiv() string {
	return g.GetIndent() + "</div>\n"
}

// OpenFieldset opens a fieldset element
func (g *HTMLGenerator) OpenFieldset() string {
	return g.GetIndent() + "<fieldset>\n"
}

// CloseFieldset closes a fieldset element
func (g *HTMLGenerator) CloseFieldset() string {
	return g.GetIndent() + "</fieldset>\n"
}

// GenerateLegend generates a legend element
func (g *HTMLGenerator) GenerateLegend(text string) string {
	return g.GetIndent() + fmt.Sprintf("<legend>%s</legend>", html.EscapeString(text)) + "\n"
}

// GenerateParagraph generates a paragraph element
func (g *HTMLGenerator) GenerateParagraph(text string) string {
	return g.GetIndent() + fmt.Sprintf("<p>%s</p>", html.EscapeString(text)) + "\n"
}

// GenerateDescription generates a description/help text
func (g *HTMLGenerator) GenerateDescription(text string) string {
	return g.GetIndent() + fmt.Sprintf(`<small class="description">%s</small>`, html.EscapeString(text)) + "\n"
}

// OpenForm opens a form element
func (g *HTMLGenerator) OpenForm(id, action, method string) string {
	return g.GetIndent() + fmt.Sprintf(`<form id="%s" action="%s" method="%s">`, id, action, method) + "\n"
}

// CloseForm closes a form element
func (g *HTMLGenerator) CloseForm() string {
	return g.GetIndent() + "</form>\n"
}

// GenerateSubmitButton generates a submit button
func (g *HTMLGenerator) GenerateSubmitButton(text string) string {
	return g.GetIndent() + fmt.Sprintf(`<button type="submit">%s</button>`, html.EscapeString(text)) + "\n"
}

// OpenLinksSection opens a links section
func (g *HTMLGenerator) OpenLinksSection() string {
	return g.GetIndent() + `<div class="form-actions">` + "\n"
}

// CloseLinksSection closes a links section
func (g *HTMLGenerator) CloseLinksSection() string {
	return g.GetIndent() + "</div>\n"
}

// GenerateLinkButton generates a button or link for a hypermedia action
func (g *HTMLGenerator) GenerateLinkButton(link HyperSchemaLink, resolvedHref string) string {
	method := strings.ToUpper(link.Method)

	// Submit and create links render as actual submit buttons
	if link.Rel == "submit" || link.Rel == "create" {
		return g.GetIndent() + fmt.Sprintf(`<button type="submit" formaction="%s" formmethod="%s">%s</button>`,
			html.EscapeString(resolvedHref),
			strings.ToLower(method),
			html.EscapeString(link.Title)) + "\n"
	}

	// Determine button class based on rel and method
	var class string

	switch link.Rel {
	case "delete", "remove":
		class = "btn-delete"
	case "cancel":
		class = "btn-link"
	case "prev", "previous":
		class = "btn-prev"
	case "next":
		class = "btn-next"
	default:
		class = "btn-link"
	}

	// GET requests use anchor tags
	if method == "GET" {
		return g.GetIndent() + fmt.Sprintf(`<a href="%s" class="%s">%s</a>`,
			html.EscapeString(resolvedHref),
			class,
			html.EscapeString(link.Title)) + "\n"
	}

	// POST, PUT, PATCH, DELETE use buttons with formaction
	return g.GetIndent() + fmt.Sprintf(`<button type="button" formaction="%s" formmethod="%s" class="%s">%s</button>`,
		html.EscapeString(resolvedHref),
		strings.ToLower(method),
		class,
		html.EscapeString(link.Title)) + "\n"
}

// WrapFullPage wraps HTML content in a complete HTML5 page
func (g *HTMLGenerator) WrapFullPage(title, formHTML, css string) string {
	var builder strings.Builder

	builder.WriteString("<!DOCTYPE html>\n")
	builder.WriteString("<html lang=\"en\">\n")
	builder.WriteString("<head>\n")
	builder.WriteString("  <meta charset=\"UTF-8\">\n")
	builder.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	builder.WriteString(fmt.Sprintf("  <title>%s</title>\n", html.EscapeString(title)))

	if css != "" {
		builder.WriteString("  <style>\n")
		builder.WriteString(css)
		builder.WriteString("  </style>\n")
	}

	builder.WriteString("</head>\n")
	builder.WriteString("<body>\n")
	builder.WriteString(formHTML)
	builder.WriteString("</body>\n")
	builder.WriteString("</html>\n")

	return builder.String()
}

// GetDefaultCSS returns default CSS for forms
func GetDefaultCSS() string {
	return `    * {
      box-sizing: border-box;
    }

    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
      line-height: 1.6;
      color: #333;
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
      background-color: #f5f5f5;
    }

    form {
      background: white;
      padding: 30px;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    .form-group {
      margin-bottom: 20px;
    }

    label {
      display: block;
      margin-bottom: 5px;
      font-weight: 500;
      color: #555;
    }

    .required {
      color: #dc3545;
    }

    input[type="text"],
    input[type="email"],
    input[type="number"],
    input[type="tel"],
    input[type="url"],
    input[type="date"],
    input[type="time"],
    input[type="datetime-local"],
    input[type="color"],
    select,
    textarea {
      width: 100%;
      padding: 8px 12px;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 14px;
      font-family: inherit;
    }

    input[type="checkbox"] {
      width: auto;
      margin-right: 8px;
    }

    textarea {
      min-height: 100px;
      resize: vertical;
    }

    input:focus,
    select:focus,
    textarea:focus {
      outline: none;
      border-color: #4CAF50;
      box-shadow: 0 0 0 2px rgba(76, 175, 80, 0.1);
    }

    .description {
      display: block;
      margin-top: 4px;
      color: #666;
      font-size: 13px;
    }

    fieldset {
      border: 1px solid #ddd;
      border-radius: 4px;
      padding: 15px;
      margin-bottom: 20px;
    }

    legend {
      font-weight: 600;
      color: #444;
      padding: 0 10px;
    }

    .vertical-layout {
      display: flex;
      flex-direction: column;
      gap: 0;
    }

    .horizontal-layout {
      display: flex;
      flex-wrap: wrap;
      gap: 15px;
      margin-bottom: 20px;
    }

    .horizontal-layout > .form-group {
      flex: 1;
      min-width: 200px;
    }

    button[type="submit"] {
      background-color: #4CAF50;
      color: white;
      padding: 10px 24px;
      border: none;
      border-radius: 4px;
      font-size: 16px;
      cursor: pointer;
      margin-top: 10px;
    }

    button[type="submit"]:hover {
      background-color: #45a049;
    }

    button[type="submit"]:active {
      background-color: #3d8b40;
    }

    .form-actions {
      margin-top: 20px;
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
    }

    .btn-link {
      display: inline-block;
      padding: 10px 24px;
      border: 1px solid #ddd;
      background: white;
      color: #333;
      text-decoration: none;
      border-radius: 4px;
      font-size: 16px;
      cursor: pointer;
    }

    .btn-link:hover {
      background-color: #f5f5f5;
    }

    .btn-delete {
      display: inline-block;
      padding: 10px 24px;
      border: 1px solid #dc3545;
      background: #dc3545;
      color: white;
      text-decoration: none;
      border-radius: 4px;
      font-size: 16px;
      cursor: pointer;
    }

    .btn-delete:hover {
      background-color: #c82333;
      border-color: #bd2130;
    }

    .btn-prev {
      display: inline-block;
      padding: 10px 24px;
      border: 1px solid #6c757d;
      background: #6c757d;
      color: white;
      text-decoration: none;
      border-radius: 4px;
      font-size: 16px;
      cursor: pointer;
    }

    .btn-prev:hover {
      background-color: #5a6268;
      border-color: #545b62;
    }

    .btn-next {
      display: inline-block;
      padding: 10px 24px;
      border: 1px solid #007bff;
      background: #007bff;
      color: white;
      text-decoration: none;
      border-radius: 4px;
      font-size: 16px;
      cursor: pointer;
    }

    .btn-next:hover {
      background-color: #0069d9;
      border-color: #0062cc;
    }

    @media (max-width: 600px) {
      body {
        padding: 10px;
      }

      form {
        padding: 20px;
      }

      .horizontal-layout {
        flex-direction: column;
      }

      .horizontal-layout > .form-group {
        width: 100%;
      }
    }
`
}
