package main

import (
	"regexp"
	"strings"

	"go.lsp.dev/protocol"
)

// validateKroSchema validates Kro schema in a document and returns diagnostics
func validateKroSchema(content string) []protocol.Diagnostic {
	var diagnostics []protocol.Diagnostic

	// Check for required top-level fields: kind, apiVersion, metadata, spec
	if !strings.Contains(content, "kind:") {
		diagnostics = append(diagnostics, createErrorDiagnostic(
			0, 0, 0, 0,
			"kro-schema",
			formatKroSchemaErrorMessage("Missing required field 'kind'"),
		))
	} else if !strings.Contains(content, "kind: ResourceGraphDefinition") {
		diagnostics = append(diagnostics, createErrorDiagnostic(
			0, 0, 0, 0,
			"kro-schema",
			formatKroSchemaErrorMessage("'kind' must be 'ResourceGraphDefinition'"),
		))
	}

	if !strings.Contains(content, "apiVersion:") {
		diagnostics = append(diagnostics, createErrorDiagnostic(
			0, 0, 0, 0,
			"kro-schema",
			formatKroSchemaErrorMessage("Missing required field 'apiVersion'"),
		))
	} else if !strings.Contains(content, "apiVersion: kro.run/v1alpha1") {
		diagnostics = append(diagnostics, createErrorDiagnostic(
			0, 0, 0, 0,
			"kro-schema",
			formatKroSchemaErrorMessage("'apiVersion' must be 'kro.run/v1alpha1'"),
		))
	}

	if !strings.Contains(content, "metadata:") {
		diagnostics = append(diagnostics, createErrorDiagnostic(
			0, 0, 0, 0,
			"kro-schema",
			formatKroSchemaErrorMessage("Missing required field 'metadata'"),
		))
	}

	if !strings.Contains(content, "spec:") {
		diagnostics = append(diagnostics, createErrorDiagnostic(
			0, 0, 0, 0,
			"kro-schema",
			formatKroSchemaErrorMessage("Missing required field 'spec'"),
		))
	}

	// Check for name in metadata
	metadataPattern := regexp.MustCompile(`metadata:\s*(\n\s+[^:]+:[^\n]*)*`)
	if metadataPattern.MatchString(content) {
		metadataBlock := metadataPattern.FindString(content)
		if !strings.Contains(metadataBlock, "name:") {
			// Find the line number for metadata
			lines := strings.Split(content, "\n")
			lineNum := uint32(0)
			for i, line := range lines {
				if strings.TrimSpace(line) == "metadata:" {
					lineNum = uint32(i)
					break
				}
			}
			diagnostics = append(diagnostics, createErrorDiagnostic(
				lineNum, 0, lineNum, 999,
				"kro-schema",
				formatKroSchemaErrorMessage("'metadata' must contain 'name' field"),
			))
		}
	}

	// Check for required spec fields: resources
	specPattern := regexp.MustCompile(`spec:\s*(\n\s+[^:]+:[^\n]*)*`)
	if specPattern.MatchString(content) {
		specBlock := specPattern.FindString(content)
		if !strings.Contains(specBlock, "resources:") {
			// Find the line number for spec
			lines := strings.Split(content, "\n")
			lineNum := uint32(0)
			for i, line := range lines {
				if strings.TrimSpace(line) == "spec:" {
					lineNum = uint32(i)
					break
				}
			}
			diagnostics = append(diagnostics, createErrorDiagnostic(
				lineNum, 0, lineNum, 999,
				"kro-schema",
				formatKroSchemaErrorMessage("'spec' must contain 'resources' field"),
			))
		}
	}

	return diagnostics
}

// findLineNumber finds the line number for a pattern in content
func findLineNumber(content, pattern string) uint32 {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, pattern) {
			return uint32(i)
		}
	}
	return 0
}

// extractCELRanges extracts ranges for CEL expressions in the document
func extractCELRanges(content string) []protocol.Range {
	var ranges []protocol.Range
	re := regexp.MustCompile(`\{\{([^}]*)\}\}`)
	
	// Find all matches
	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		// Find all matches in this line
		matches := re.FindAllStringIndex(line, -1)
		for _, match := range matches {
			// Create a range for this match
			celRange := protocol.Range{
				Start: protocol.Position{
					Line:      uint32(lineNum),
					Character: uint32(match[0]),
				},
				End: protocol.Position{
					Line:      uint32(lineNum),
					Character: uint32(match[1]),
				},
			}
			ranges = append(ranges, celRange)
		}
	}
	
	return ranges
}

// extractCELExpression extracts the CEL expression from a string with {{ }}
func extractCELExpression(expr string) string {
	// Remove {{ and }} and trim spaces
	expr = strings.TrimPrefix(expr, "{{")
	expr = strings.TrimSuffix(expr, "}}")
	return strings.TrimSpace(expr)
}
