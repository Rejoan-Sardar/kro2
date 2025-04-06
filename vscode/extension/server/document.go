package main

import (
	"bytes"
	"regexp"
	"strings"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
	"gopkg.in/yaml.v3"
)

// Document represents a text document
type Document struct {
	URI     uri.URI
	Content string
	Lang    string
}

// NewDocument creates a new document
func NewDocument(uri uri.URI, content string, lang string) *Document {
	return &Document{
		URI:     uri,
		Content: content,
		Lang:    lang,
	}
}

// Validate validates the document and returns diagnostics
func (d *Document) Validate() []protocol.Diagnostic {
	var diagnostics []protocol.Diagnostic

	// Only process .kro files
	if !strings.HasSuffix(string(d.URI), ".kro") && 
	   !strings.HasSuffix(string(d.URI), ".yaml") && 
	   !strings.HasSuffix(string(d.URI), ".yml") {
		return diagnostics
	}

	// Preprocess the document to handle CEL expressions in templates
	content := preprocessCELExpressions(d.Content)

	// Validate YAML
	yamlDiagnostics := validateYAML(content)
	diagnostics = append(diagnostics, yamlDiagnostics...)

	// If there are YAML errors, don't proceed with further validation
	if len(yamlDiagnostics) > 0 {
		return diagnostics
	}

	// Validate CEL expressions
	celDiagnostics := validateCELExpressions(d.Content)
	diagnostics = append(diagnostics, celDiagnostics...)

	// Validate Kro schema
	schemaDiagnostics := validateKroSchema(d.Content)
	diagnostics = append(diagnostics, schemaDiagnostics...)

	return diagnostics
}

// validateYAML validates YAML syntax and returns diagnostics
func validateYAML(content string) []protocol.Diagnostic {
	var diagnostics []protocol.Diagnostic
	var yamlData interface{}

	// Skip validation for empty documents
	if strings.TrimSpace(content) == "" {
		return diagnostics
	}

	err := yaml.Unmarshal([]byte(content), &yamlData)
	if err != nil {
		// Try to extract line number from YAML error message
		// Example: "yaml: line 5: mapping values are not allowed in this context"
		lineRegex := regexp.MustCompile(`line (\d+):`)
		match := lineRegex.FindStringSubmatch(err.Error())
		
		startLine := uint32(0)
		if len(match) > 1 {
			startLine = uint32(parseIntSafe(match[1], 0))
			// Adjust for 0-based line numbers in LSP
			if startLine > 0 {
				startLine--
			}
		}
		
		// Create error diagnostic
		diagnostic := createErrorDiagnostic(
			startLine, 0,
			startLine, 999,
			"yaml",
			formatYAMLErrorMessage(err.Error()),
		)
		diagnostics = append(diagnostics, diagnostic)
	}
	
	return diagnostics
}

// parseIntSafe parses an integer with a fallback value
func parseIntSafe(s string, fallback int) int {
	// This is a simple implementation to avoid importing strconv
	// In a real implementation, you'd use strconv.Atoi
	result := 0
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			result = result*10 + int(ch-'0')
		} else {
			return fallback
		}
	}
	return result
}

// preprocessCELExpressions replaces CEL expressions in double curly braces
// with placeholder values to prevent YAML syntax errors during validation
func preprocessCELExpressions(content string) string {
	// Define regex to match CEL expressions in double curly braces
	re := regexp.MustCompile(`\{\{([^}]*)\}\}`)
	
	// Replace with a valid YAML string
	return re.ReplaceAllString(content, "'CEL_EXPRESSION'")
}

// GetCompletionItems returns completion items for the given position
func (d *Document) GetCompletionItems(position protocol.Position) []protocol.CompletionItem {
	// Get context at the position to determine what completions to provide
	contextInfo := d.getContextAtPosition(position)
	
	// Based on context, provide appropriate completion items
	switch contextInfo.context {
	case "top-level":
		return getTopLevelCompletionItems(contextInfo.linePrefix)
	case "resources":
		return getResourceCompletionItems(contextInfo.linePrefix)
	case "relations":
		return getRelationCompletionItems(contextInfo.linePrefix)
	case "parameters":
		return getParameterCompletionItems(contextInfo.linePrefix)
	case "cel":
		return getCELCompletionItems(contextInfo.linePrefix)
	default:
		// Default to top-level items if context is unknown
		return getTopLevelCompletionItems(contextInfo.linePrefix)
	}
}

// contextInfo holds information about the context at a position
type contextInfo struct {
	context    string
	linePrefix string
}

// getContextAtPosition analyzes the document to determine the context
// at the given position for providing relevant completions
func (d *Document) getContextAtPosition(position protocol.Position) contextInfo {
	lines := strings.Split(d.Content, "\n")
	
	// Bounds check
	if int(position.Line) >= len(lines) {
		return contextInfo{context: "top-level", linePrefix: ""}
	}
	
	currentLine := lines[position.Line]
	linePrefix := currentLine
	if int(position.Character) < len(currentLine) {
		linePrefix = currentLine[:position.Character]
	}
	
	// Check for CEL expressions (double curly braces)
	if strings.Contains(linePrefix, "{{") && !strings.Contains(linePrefix, "}}") {
		return contextInfo{context: "cel", linePrefix: linePrefix}
	}
	
	// Walk back through previous lines to determine the context
	// Start with the current line and go backwards
	indent := getIndentation(currentLine)
	section := "top-level"
	
	for i := int(position.Line) - 1; i >= 0; i-- {
		prevLine := lines[i]
		prevIndent := getIndentation(prevLine)
		
		// Check for section headers when we find a less indented line
		if prevIndent < indent {
			if strings.HasPrefix(strings.TrimSpace(prevLine), "resources:") {
				section = "resources"
				break
			} else if strings.HasPrefix(strings.TrimSpace(prevLine), "relations:") {
				section = "relations"
				break
			} else if strings.HasPrefix(strings.TrimSpace(prevLine), "parameters:") {
				section = "parameters"
				break
			} else if strings.HasPrefix(strings.TrimSpace(prevLine), "spec:") {
				section = "top-level"
				break
			}
		}
	}
	
	return contextInfo{context: section, linePrefix: linePrefix}
}

// getIndentation returns the number of leading spaces in a string
func getIndentation(s string) int {
	i := 0
	for i < len(s) && s[i] == ' ' {
		i++
	}
	return i
}

// extractYAMLNode attempts to extract a YAML node from content
func extractYAMLNode(content string) (*yaml.Node, error) {
	var node yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(content)))
	err := decoder.Decode(&node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}
