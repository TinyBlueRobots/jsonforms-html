# JSON Forms HTML Renderer

A WebAssembly-powered web application that converts JSON Forms schemas into semantic HTML forms. Runs entirely in your browser with no server required.

## 🚀 [Try it Live](https://tinybluerobots.github.io/jsonforms-html/)

## Usage

Paste your JSON in this format:

```json
{
  "schema": {
    "type": "object",
    "properties": {
      "name": { "type": "string", "title": "Full Name" },
      "email": { "type": "string", "format": "email" }
    },
    "required": ["name"]
  },
  "uischema": {
    "type": "VerticalLayout",
    "elements": [
      { "type": "Control", "scope": "#/properties/name" },
      { "type": "Control", "scope": "#/properties/email" }
    ]
  }
}
```

Click "Generate HTML" to see the rendered form.

## Building

```bash
GOOS=js GOARCH=wasm go build -o public/jsonforms.wasm main.go
```

## Dependencies

- [jsonforms-parser](https://github.com/tinybluerobots/jsonforms-parser) - JSON Forms AST parser

## Related

- [JSON Forms](https://jsonforms.io/) - Official JavaScript implementation
