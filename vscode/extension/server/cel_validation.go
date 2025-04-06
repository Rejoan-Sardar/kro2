package main

import (
	"regexp"
	"strings"

	"go.lsp.dev/protocol"
)

// validateCELExpressions validates CEL expressions in a document and returns diagnostics
func validateCELExpressions(content string) []protocol.Diagnostic {
	var diagnostics []protocol.Diagnostic
	
	// Extract all CEL expressions with their ranges
	celRanges := extractCELRanges(content)
	
	// Get the content of each expression and validate it
	lines := strings.Split(content, "\n")
	for _, rng := range celRanges {
		if int(rng.Start.Line) >= len(lines) {
			continue
		}
		
		line := lines[rng.Start.Line]
		if int(rng.Start.Character) >= len(line) || int(rng.End.Character) > len(line) {
			continue
		}
		
		exprText := line[rng.Start.Character:rng.End.Character]
		expression := extractCELExpression(exprText)
		
		// Validate the expression
		err := validateCELExpression(expression)
		if err != nil {
			diagnostic := createDiagnosticForRange(
				rng,
				protocol.DiagnosticSeverityError,
				"cel",
				formatCELErrorMessage(err.Error()),
			)
			diagnostics = append(diagnostics, diagnostic)
		}
	}
	
	return diagnostics
}

// validateCELExpression validates a single CEL expression
func validateCELExpression(expr string) error {
	// Simple validation for now
	if strings.TrimSpace(expr) == "" {
		return NewCELError("Empty CEL expression")
	}
	
	// Check for unclosed parentheses, brackets, quotes
	if err := checkBalancedDelimiters(expr); err != nil {
		return err
	}
	
	// Check for invalid operators
	if err := checkOperators(expr); err != nil {
		return err
	}
	
	// Check for well-formed ternary expressions
	if err := checkTernaryExpression(expr); err != nil {
		return err
	}
	
	return nil
}

// CELError represents an error in a CEL expression
type CELError struct {
	Message string
}

// NewCELError creates a new CEL error
func NewCELError(msg string) *CELError {
	return &CELError{Message: msg}
}

func (e *CELError) Error() string {
	return e.Message
}

// checkBalancedDelimiters checks if all delimiters are properly balanced
func checkBalancedDelimiters(expr string) error {
	var stack []rune
	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
		'"': '"',
		'\'': '\'',
	}
	
	inSingleQuote := false
	inDoubleQuote := false
	
	for _, ch := range expr {
		// Handle quotes specially
		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}
		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		
		// Skip processing inside quotes
		if inSingleQuote || inDoubleQuote {
			continue
		}
		
		// Check opening delimiters
		if ch == '(' || ch == '[' || ch == '{' {
			stack = append(stack, ch)
			continue
		}
		
		// Check closing delimiters
		if closing, ok := pairs[ch]; ok {
			if len(stack) == 0 {
				return NewCELError("Unexpected closing delimiter: " + string(ch))
			}
			
			if stack[len(stack)-1] != closing {
				return NewCELError("Mismatched delimiter: expected " + string(pairs[ch]) + " but found " + string(stack[len(stack)-1]))
			}
			
			// Pop from stack
			stack = stack[:len(stack)-1]
		}
	}
	
	if inSingleQuote {
		return NewCELError("Unclosed single quote")
	}
	if inDoubleQuote {
		return NewCELError("Unclosed double quote")
	}
	
	if len(stack) > 0 {
		return NewCELError("Unclosed delimiters: " + string(stack))
	}
	
	return nil
}

// checkOperators checks for invalid operators or operator sequences
func checkOperators(expr string) error {
	// Check for invalid operator sequences
	invalidOps := []string{"++", "--", "**", "//", "/*", "*/", "&&&&", "||||"}
	for _, op := range invalidOps {
		if strings.Contains(expr, op) {
			return NewCELError("Invalid operator sequence: " + op)
		}
	}
	
	// Check for trailing operators
	if matched, _ := regexp.MatchString(`[+\-*/&|^%<>]=?\s*$`, expr); matched {
		return NewCELError("Expression ends with an operator")
	}
	
	return nil
}

// checkTernaryExpression checks for well-formed ternary expressions
func checkTernaryExpression(expr string) error {
	// Count question marks and colons outside of strings
	questionCount := 0
	colonCount := 0
	
	inSingleQuote := false
	inDoubleQuote := false
	
	for _, ch := range expr {
		// Handle quotes
		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}
		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		
		// Skip processing inside quotes
		if inSingleQuote || inDoubleQuote {
			continue
		}
		
		if ch == '?' {
			questionCount++
		} else if ch == ':' {
			colonCount++
		}
	}
	
	// For each question mark, there should be a colon
	if questionCount > colonCount {
		return NewCELError("Ternary expression missing colon")
	} else if questionCount < colonCount {
		return NewCELError("Unexpected colon in expression")
	}
	
	return nil
}
