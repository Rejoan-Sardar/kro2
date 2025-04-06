package main

import (
        "strings"

        "go.lsp.dev/protocol"
        "go.lsp.dev/uri"
)

// CompletionItem definitions with kind and detail
var (
        // Top-level Resource Schema Items
        topLevelItems = []protocol.CompletionItem{
                {
                        Label:         "kind",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "Specifies the resource type",
                        Documentation: "Required field that must be set to 'ResourceGraphDefinition'",
                        InsertText:    "kind: ResourceGraphDefinition",
                },
                {
                        Label:         "apiVersion",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "API version for the resource",
                        Documentation: "Required field that must be set to 'kro.run/v1alpha1'",
                        InsertText:    "apiVersion: kro.run/v1alpha1",
                },
                {
                        Label:         "metadata",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "Resource metadata",
                        Documentation: "Contains metadata for the resource such as name and namespace",
                        InsertText:    "metadata:\n  name: ",
                },
                {
                        Label:         "spec",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "Resource specification",
                        Documentation: "Contains the core definition of the resource graph",
                        InsertText:    "spec:\n  ",
                },
        }

        // Main Spec Sections
        specItems = []protocol.CompletionItem{
                {
                        Label:         "resources",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "Resource definitions",
                        Documentation: "Define the resources that make up your infrastructure",
                        InsertText:    "resources:\n  ",
                },
                {
                        Label:         "relations",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "Relation definitions",
                        Documentation: "Define the relations between resources",
                        InsertText:    "relations:\n  ",
                },
                {
                        Label:         "parameters",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "Parameter definitions",
                        Documentation: "Define input parameters for the resource graph",
                        InsertText:    "parameters:\n  ",
                },
        }

        // Resource Properties
        resourceItems = []protocol.CompletionItem{
                {
                        Label:         "resource-name",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "Unique resource name",
                        Documentation: "Define a new resource with a unique identifier",
                        InsertText:    "resource-name:\n  type: \n  template: ",
                },
                {
                        Label:         "type",
                        Kind:          protocol.CompletionItemKindProperty,
                        Detail:        "Resource type",
                        Documentation: "The Kubernetes API type of this resource",
                        InsertText:    "type: ",
                },
                {
                        Label:         "template",
                        Kind:          protocol.CompletionItemKindProperty,
                        Detail:        "Resource template",
                        Documentation: "The YAML template for this resource",
                        InsertText:    "template: |\n    ",
                },
                {
                        Label:         "annotations",
                        Kind:          protocol.CompletionItemKindProperty,
                        Detail:        "Resource annotations",
                        Documentation: "Metadata annotations for this resource",
                        InsertText:    "annotations:\n    ",
                },
        }

        // Relation Properties
        relationItems = []protocol.CompletionItem{
                {
                        Label:         "relation-name",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "Unique relation name",
                        Documentation: "Define a new relation with a unique identifier",
                        InsertText:    "relation-name:\n  from: \n  to: ",
                },
                {
                        Label:         "from",
                        Kind:          protocol.CompletionItemKindProperty,
                        Detail:        "Relation source",
                        Documentation: "The source resource of this relation",
                        InsertText:    "from: ",
                },
                {
                        Label:         "to",
                        Kind:          protocol.CompletionItemKindProperty,
                        Detail:        "Relation target",
                        Documentation: "The target resource of this relation",
                        InsertText:    "to: ",
                },
                {
                        Label:         "condition",
                        Kind:          protocol.CompletionItemKindProperty,
                        Detail:        "Relation condition",
                        Documentation: "A CEL condition that must be true for this relation",
                        InsertText:    "condition: ",
                },
        }

        // Parameter Properties
        parameterItems = []protocol.CompletionItem{
                {
                        Label:         "parameter-name",
                        Kind:          protocol.CompletionItemKindKeyword,
                        Detail:        "Unique parameter name",
                        Documentation: "Define a new input parameter with a unique identifier",
                        InsertText:    "parameter-name:\n  type: \n  default: ",
                },
                {
                        Label:         "type",
                        Kind:          protocol.CompletionItemKindProperty,
                        Detail:        "Parameter type",
                        Documentation: "The data type of this parameter (string, number, boolean, object, array)",
                        InsertText:    "type: ",
                },
                {
                        Label:         "default",
                        Kind:          protocol.CompletionItemKindProperty,
                        Detail:        "Default value",
                        Documentation: "The default value for this parameter if not provided",
                        InsertText:    "default: ",
                },
                {
                        Label:         "description",
                        Kind:          protocol.CompletionItemKindProperty,
                        Detail:        "Parameter description",
                        Documentation: "A description of this parameter's purpose and usage",
                        InsertText:    "description: ",
                },
        }

        // CEL Expression Functions
        celKeywords = []protocol.CompletionItem{
                {
                        Label:         "params",
                        Kind:          protocol.CompletionItemKindVariable,
                        Detail:        "Parameters object",
                        Documentation: "Access input parameters",
                        InsertText:    "params.",
                },
                {
                        Label:         "resources",
                        Kind:          protocol.CompletionItemKindVariable,
                        Detail:        "Resources object",
                        Documentation: "Access defined resources",
                        InsertText:    "resources.",
                },
                {
                        Label:         "relations",
                        Kind:          protocol.CompletionItemKindVariable,
                        Detail:        "Relations object",
                        Documentation: "Access defined relations",
                        InsertText:    "relations.",
                },
                {
                        Label:         "has",
                        Kind:          protocol.CompletionItemKindFunction,
                        Detail:        "CEL has() function",
                        Documentation: "Check if a field exists",
                        InsertText:    "has(",
                },
                {
                        Label:         "size",
                        Kind:          protocol.CompletionItemKindFunction,
                        Detail:        "CEL size() function",
                        Documentation: "Get the size/length of a string, list, or map",
                        InsertText:    "size(",
                },
                {
                        Label:         "string",
                        Kind:          protocol.CompletionItemKindFunction,
                        Detail:        "CEL string() function",
                        Documentation: "Convert a value to string",
                        InsertText:    "string(",
                },
                {
                        Label:         "int",
                        Kind:          protocol.CompletionItemKindFunction,
                        Detail:        "CEL int() function",
                        Documentation: "Convert a value to integer",
                        InsertText:    "int(",
                },
                {
                        Label:         "bool",
                        Kind:          protocol.CompletionItemKindFunction,
                        Detail:        "CEL bool() function",
                        Documentation: "Convert a value to boolean",
                        InsertText:    "bool(",
                },
        }
)

// Note: ValidateAllDocuments is implemented in document_store.go

// getTopLevelCompletionItems returns completion items for top-level fields
func getTopLevelCompletionItems(prefix string) []protocol.CompletionItem {
        // If we're in a spec section, return spec section completions
        if strings.Contains(prefix, "spec:") {
                return specItems
        }
        return topLevelItems
}

// getResourceCompletionItems returns completion items for resource fields
func getResourceCompletionItems(prefix string) []protocol.CompletionItem {
        return resourceItems
}

// getRelationCompletionItems returns completion items for relation fields
func getRelationCompletionItems(prefix string) []protocol.CompletionItem {
        return relationItems
}

// getParameterCompletionItems returns completion items for parameter fields
func getParameterCompletionItems(prefix string) []protocol.CompletionItem {
        return parameterItems
}

// getCELCompletionItems returns completion items for CEL expressions
func getCELCompletionItems(prefix string) []protocol.CompletionItem {
        return celKeywords
}
