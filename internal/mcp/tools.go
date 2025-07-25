package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/phildougherty/mcp-gmail-go/internal/gmail"
	"github.com/phildougherty/mcp-gmail-go/internal/types"
)

type ToolRegistry struct {
	gmailClient *gmail.Client
	tools       map[string]Tool
}

func NewToolRegistry(gmailClient *gmail.Client) *ToolRegistry {
	registry := &ToolRegistry{
		gmailClient: gmailClient,
		tools:       make(map[string]Tool),
	}
	
	registry.registerTools()
	return registry
}

func (r *ToolRegistry) registerTools() {
	r.tools["send_email"] = Tool{
		Name:        "send_email",
		Description: "Sends a new email",
		InputSchema: SendEmailSchema,
	}
	
	r.tools["draft_email"] = Tool{
		Name:        "draft_email",
		Description: "Creates a draft email",
		InputSchema: SendEmailSchema,
	}
	
	r.tools["read_email"] = Tool{
		Name:        "read_email",
		Description: "Retrieves the content of a specific email",
		InputSchema: ReadEmailSchema,
	}
	
	r.tools["search_emails"] = Tool{
		Name:        "search_emails",
		Description: "Searches for emails using Gmail search syntax",
		InputSchema: SearchEmailsSchema,
	}
	
	r.tools["modify_email"] = Tool{
		Name:        "modify_email",
		Description: "Modifies email labels (move to different folders)",
		InputSchema: ModifyEmailSchema,
	}
	
	r.tools["delete_email"] = Tool{
		Name:        "delete_email",
		Description: "Permanently deletes an email",
		InputSchema: DeleteEmailSchema,
	}
	
	r.tools["list_email_labels"] = Tool{
		Name:        "list_email_labels",
		Description: "Retrieves all available Gmail labels",
		InputSchema: ListEmailLabelsSchema,
	}
	
	r.tools["create_label"] = Tool{
		Name:        "create_label",
		Description: "Creates a new Gmail label",
		InputSchema: CreateLabelSchema,
	}
	
	r.tools["get_contacts"] = Tool{
		Name:        "get_contacts",
		Description: "Retrieves Gmail contacts",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"maxResults": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of contacts to return",
				},
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query for contacts",
				},
			},
		},
	}
	
	r.tools["email_analytics"] = Tool{
		Name:        "email_analytics",
		Description: "Get email analytics and statistics",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"days": map[string]interface{}{
					"type":        "integer",
					"description": "Number of days to analyze (default: 30)",
				},
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Gmail search query to filter emails",
				},
				"groupBy": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"sender", "label", "day"},
					"description": "How to group the analytics results",
				},
			},
		},
	}
}

func (r *ToolRegistry) ListTools() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

func (r *ToolRegistry) CallTool(name string, args json.RawMessage) (*ToolResult, error) {
	_, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	
	switch name {
	case "send_email":
		return r.handleSendEmail(args)
	case "draft_email":
		return r.handleDraftEmail(args)
	case "read_email":
		return r.handleReadEmail(args)
	case "search_emails":
		return r.handleSearchEmails(args)
	case "modify_email":
		return r.handleModifyEmail(args)
	case "delete_email":
		return r.handleDeleteEmail(args)
	case "list_email_labels":
		return r.handleListLabels(args)
	case "create_label":
		return r.handleCreateLabel(args)
	case "get_contacts":
		return r.handleGetContacts(args)
	case "email_analytics":
		return r.handleEmailAnalytics(args)
	default:
		return nil, fmt.Errorf("tool implementation not found: %s", name)
	}
}

func (r *ToolRegistry) handleSendEmail(args json.RawMessage) (*ToolResult, error) {
	var sendArgs types.SendEmailArgs
	if err := json.Unmarshal(args, &sendArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	messageID, err := r.gmailClient.SendEmail(&sendArgs)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to send email: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Email sent successfully with ID: %s", messageID),
		}},
	}, nil
}

func (r *ToolRegistry) handleDraftEmail(args json.RawMessage) (*ToolResult, error) {
	var draftArgs types.SendEmailArgs
	if err := json.Unmarshal(args, &draftArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	draftID, err := r.gmailClient.CreateDraft(&draftArgs)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to create draft: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Draft created successfully with ID: %s", draftID),
		}},
	}, nil
}

func (r *ToolRegistry) handleReadEmail(args json.RawMessage) (*ToolResult, error) {
	var readArgs types.ReadEmailArgs
	if err := json.Unmarshal(args, &readArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	message, err := r.gmailClient.ReadEmail(readArgs.MessageID)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to read email: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	content := fmt.Sprintf("Thread ID: %s\nSubject: %s\nFrom: %s\nTo: %s\nDate: %s\n\n%s",
		message.ThreadID, message.Subject, message.From, message.To, message.Date, message.Body)
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: content,
		}},
	}, nil
}

func (r *ToolRegistry) handleSearchEmails(args json.RawMessage) (*ToolResult, error) {
	var searchArgs types.SearchEmailsArgs
	if err := json.Unmarshal(args, &searchArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	if searchArgs.MaxResults == 0 {
		searchArgs.MaxResults = 10
	}
	
	messages, err := r.gmailClient.SearchEmails(searchArgs.Query, searchArgs.MaxResults)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to search emails: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	var result string
	for _, msg := range messages {
		result += fmt.Sprintf("ID: %s\nSubject: %s\nFrom: %s\nDate: %s\n\n", 
			msg.ID, msg.Subject, msg.From, msg.Date)
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: result,
		}},
	}, nil
}

func (r *ToolRegistry) handleModifyEmail(args json.RawMessage) (*ToolResult, error) {
	var modifyArgs types.ModifyEmailArgs
	if err := json.Unmarshal(args, &modifyArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	err := r.gmailClient.ModifyEmail(modifyArgs.MessageID, modifyArgs.AddLabelIDs, modifyArgs.RemoveLabelIDs)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to modify email: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Email %s labels updated successfully", modifyArgs.MessageID),
		}},
	}, nil
}

func (r *ToolRegistry) handleDeleteEmail(args json.RawMessage) (*ToolResult, error) {
	var deleteArgs types.DeleteEmailArgs
	if err := json.Unmarshal(args, &deleteArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	err := r.gmailClient.DeleteEmail(deleteArgs.MessageID)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to delete email: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Email %s deleted successfully", deleteArgs.MessageID),
		}},
	}, nil
}

func (r *ToolRegistry) handleListLabels(args json.RawMessage) (*ToolResult, error) {
	labels, err := r.gmailClient.ListLabels()
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to list labels: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	var result string
	for _, label := range labels {
		result += fmt.Sprintf("ID: %s\nName: %s\nType: %s\n\n", 
			label.ID, label.Name, label.Type)
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: result,
		}},
	}, nil
}

func (r *ToolRegistry) handleCreateLabel(args json.RawMessage) (*ToolResult, error) {
	var createArgs types.CreateLabelArgs
	if err := json.Unmarshal(args, &createArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	labelID, err := r.gmailClient.CreateLabel(createArgs.Name)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to create label: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Label created successfully with ID: %s", labelID),
		}},
	}, nil
}

func (r *ToolRegistry) handleGetContacts(args json.RawMessage) (*ToolResult, error) {
	var contactArgs types.GetContactsArgs
	if err := json.Unmarshal(args, &contactArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	if contactArgs.MaxResults == 0 {
		contactArgs.MaxResults = 50
	}
	
	contacts, err := r.gmailClient.GetContacts(contactArgs.MaxResults, contactArgs.Query)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to get contacts: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	var result string
	for _, contact := range contacts {
		result += fmt.Sprintf("Name: %s\nEmail: %s\n\n", contact.Name, contact.Email)
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: result,
		}},
	}, nil
}

func (r *ToolRegistry) handleEmailAnalytics(args json.RawMessage) (*ToolResult, error) {
	var analyticsArgs types.EmailAnalyticsArgs
	if err := json.Unmarshal(args, &analyticsArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	if analyticsArgs.Days == 0 {
		analyticsArgs.Days = 30
	}
	
	analytics, err := r.gmailClient.GetEmailAnalytics(analyticsArgs.Days, analyticsArgs.Query, analyticsArgs.GroupBy)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to get email analytics: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: analytics,
		}},
	}, nil
}