package main

import (
	"fmt"

	"go.lsp.dev/protocol"
)

// createDiagnostic creates a diagnostic with the given parameters
func createDiagnostic(
	startLine uint32,
	startChar uint32,
	endLine uint32,
	endChar uint32,
	severity protocol.DiagnosticSeverity,
	source string,
	message string,
) protocol.Diagnostic {
	return protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      startLine,
				Character: startChar,
			},
			End: protocol.Position{
				Line:      endLine,
				Character: endChar,
			},
		},
		Severity: severity,
		Source:   source,
		Message:  message,
	}
}

// createErrorDiagnostic creates an error diagnostic
func createErrorDiagnostic(
	startLine uint32,
	startChar uint32,
	endLine uint32,
	endChar uint32,
	source string,
	message string,
) protocol.Diagnostic {
	return createDiagnostic(
		startLine,
		startChar,
		endLine,
		endChar,
		protocol.DiagnosticSeverityError,
		source,
		message,
	)
}

// createWarningDiagnostic creates a warning diagnostic
func createWarningDiagnostic(
	startLine uint32,
	startChar uint32,
	endLine uint32,
	endChar uint32,
	source string,
	message string,
) protocol.Diagnostic {
	return createDiagnostic(
		startLine,
		startChar,
		endLine,
		endChar,
		protocol.DiagnosticSeverityWarning,
		source,
		message,
	)
}

// createInfoDiagnostic creates an information diagnostic
func createInfoDiagnostic(
	startLine uint32,
	startChar uint32,
	endLine uint32,
	endChar uint32,
	source string,
	message string,
) protocol.Diagnostic {
	return createDiagnostic(
		startLine,
		startChar,
		endLine,
		endChar,
		protocol.DiagnosticSeverityInformation,
		source,
		message,
	)
}

// createHintDiagnostic creates a hint diagnostic
func createHintDiagnostic(
	startLine uint32,
	startChar uint32,
	endLine uint32,
	endChar uint32,
	source string,
	message string,
) protocol.Diagnostic {
	return createDiagnostic(
		startLine,
		startChar,
		endLine,
		endChar,
		protocol.DiagnosticSeverityHint,
		source,
		message,
	)
}

// createDiagnosticForLine creates a diagnostic for a full line
func createDiagnosticForLine(
	line uint32,
	severity protocol.DiagnosticSeverity,
	source string,
	message string,
) protocol.Diagnostic {
	return createDiagnostic(
		line, 0,
		line, 999, // Large column number to cover the whole line
		severity,
		source,
		message,
	)
}

// createDiagnosticForRange creates a diagnostic for a text range
func createDiagnosticForRange(
	rangeObj protocol.Range,
	severity protocol.DiagnosticSeverity,
	source string,
	message string,
) protocol.Diagnostic {
	return protocol.Diagnostic{
		Range:    rangeObj,
		Severity: severity,
		Source:   source,
		Message:  message,
	}
}

// formatCELErrorMessage formats a CEL error message for display
func formatCELErrorMessage(errMsg string) string {
	return fmt.Sprintf("CEL Expression Error: %s", errMsg)
}

// formatYAMLErrorMessage formats a YAML error message for display
func formatYAMLErrorMessage(errMsg string) string {
	return fmt.Sprintf("YAML Syntax Error: %s", errMsg)
}

// formatKroSchemaErrorMessage formats a Kro schema error message for display
func formatKroSchemaErrorMessage(errMsg string) string {
	return fmt.Sprintf("Kro Schema Error: %s", errMsg)
}
