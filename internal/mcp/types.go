package mcp

import "encoding/json"

// Tool represents an MCP tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// ToolCallRequest represents a request to call a tool
type ToolCallRequest struct {
	Arguments json.RawMessage `json:"arguments"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content represents content in a tool result
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// InputSchema definitions for tools
var (
	SendEmailSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"to": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of recipient email addresses",
			},
			"subject": map[string]interface{}{
				"type":        "string",
				"description": "Email subject",
			},
			"body": map[string]interface{}{
				"type":        "string",
				"description": "Email body content",
			},
			"htmlBody": map[string]interface{}{
				"type":        "string",
				"description": "HTML version of the email body",
			},
			"cc": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of CC recipients",
			},
			"bcc": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of BCC recipients",
			},
			"threadId": map[string]interface{}{
				"type":        "string",
				"description": "Thread ID to reply to",
			},
			"inReplyTo": map[string]interface{}{
				"type":        "string",
				"description": "Message ID being replied to",
			},
		},
		"required": []string{"to", "subject", "body"},
	}

	ReadEmailSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"messageId": map[string]interface{}{
				"type":        "string",
				"description": "ID of the email message to retrieve",
			},
		},
		"required": []string{"messageId"},
	}

	SearchEmailsSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Gmail search query (e.g., 'from:example@gmail.com')",
			},
			"maxResults": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum number of results to return",
			},
		},
		"required": []string{"query"},
	}

	ModifyEmailSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"messageId": map[string]interface{}{
				"type":        "string",
				"description": "ID of the email message to modify",
			},
			"addLabelIds": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of label IDs to add to the message",
			},
			"removeLabelIds": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of label IDs to remove from the message",
			},
		},
		"required": []string{"messageId"},
	}

	DeleteEmailSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"messageId": map[string]interface{}{
				"type":        "string",
				"description": "ID of the email message to delete",
			},
		},
		"required": []string{"messageId"},
	}

	CreateLabelSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name for the new label",
			},
			"messageListVisibility": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"show", "hide"},
				"description": "Whether to show or hide the label in the message list",
			},
			"labelListVisibility": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"labelShow", "labelShowIfUnread", "labelHide"},
				"description": "Visibility of the label in the label list",
			},
		},
		"required": []string{"name"},
	}

	ListEmailLabelsSchema = map[string]interface{}{
		"type":        "object",
		"properties":  map[string]interface{}{},
		"description": "Retrieves all available Gmail labels",
	}
)