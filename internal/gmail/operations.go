package gmail

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/phildougherty/mcp-gmail-go/internal/types"
	"google.golang.org/api/gmail/v1"
)


// SendEmail sends an email
func (c *Client) SendEmail(args *types.SendEmailArgs) (string, error) {
	message := c.createRawMessage(args, false)
	
	encodedMessage := base64.URLEncoding.EncodeToString([]byte(message))
	
	gmailMessage := &gmail.Message{
		Raw: encodedMessage,
	}
	
	if args.ThreadID != "" {
		gmailMessage.ThreadId = args.ThreadID
	}
	
	result, err := c.service.Users.Messages.Send("me", gmailMessage).Do()
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}
	
	return result.Id, nil
}

// CreateDraft creates a draft email
func (c *Client) CreateDraft(args *types.SendEmailArgs) (string, error) {
	message := c.createRawMessage(args, true)
	
	encodedMessage := base64.URLEncoding.EncodeToString([]byte(message))
	
	gmailMessage := &gmail.Message{
		Raw: encodedMessage,
	}
	
	if args.ThreadID != "" {
		gmailMessage.ThreadId = args.ThreadID
	}
	
	draft := &gmail.Draft{
		Message: gmailMessage,
	}
	
	result, err := c.service.Users.Drafts.Create("me", draft).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create draft: %w", err)
	}
	
	return result.Id, nil
}

// ReadEmail reads an email by ID
func (c *Client) ReadEmail(messageID string) (*types.EmailMessage, error) {
	message, err := c.service.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	
	headers := make(map[string]string)
	var subject, from, to, date string
	
	if message.Payload != nil && message.Payload.Headers != nil {
		for _, header := range message.Payload.Headers {
			headers[header.Name] = header.Value
			switch strings.ToLower(header.Name) {
			case "subject":
				subject = header.Value
			case "from":
				from = header.Value
			case "to":
				to = header.Value
			case "date":
				date = header.Value
			}
		}
	}
	
	body := c.extractBody(message.Payload)
	
	// Get labels
	labels := make([]string, len(message.LabelIds))
	copy(labels, message.LabelIds)
	
	return &types.EmailMessage{
		ID:       message.Id,
		ThreadID: message.ThreadId,
		Subject:  subject,
		From:     from,
		To:       to,
		Date:     date,
		Body:     body,
		Labels:   labels,
		Headers:  headers,
	}, nil
}

// SearchEmails searches for emails
func (c *Client) SearchEmails(query string, maxResults int) ([]*types.EmailMessage, error) {
	call := c.service.Users.Messages.List("me").Q(query)
	if maxResults > 0 {
		call = call.MaxResults(int64(maxResults))
	}
	
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search emails: %w", err)
	}
	
	var messages []*types.EmailMessage
	for _, msg := range response.Messages {
		detail, err := c.service.Users.Messages.Get("me", msg.Id).Format("metadata").MetadataHeaders("Subject", "From", "Date").Do()
		if err != nil {
			continue // Skip errors for individual messages
		}
		
		var subject, from, date string
		if detail.Payload != nil && detail.Payload.Headers != nil {
			for _, header := range detail.Payload.Headers {
				switch header.Name {
				case "Subject":
					subject = header.Value
				case "From":
					from = header.Value
				case "Date":
					date = header.Value
				}
			}
		}
		
		messages = append(messages, &types.EmailMessage{
			ID:      msg.Id,
			Subject: subject,
			From:    from,
			Date:    date,
		})
	}
	
	return messages, nil
}

// ModifyEmail modifies email labels
func (c *Client) ModifyEmail(messageID string, addLabelIDs, removeLabelIDs []string) error {
	request := &gmail.ModifyMessageRequest{
		AddLabelIds:    addLabelIDs,
		RemoveLabelIds: removeLabelIDs,
	}
	
	_, err := c.service.Users.Messages.Modify("me", messageID, request).Do()
	if err != nil {
		return fmt.Errorf("failed to modify email: %w", err)
	}
	
	return nil
}

// DeleteEmail deletes an email
func (c *Client) DeleteEmail(messageID string) error {
	err := c.service.Users.Messages.Delete("me", messageID).Do()
	if err != nil {
		return fmt.Errorf("failed to delete email: %w", err)
	}
	
	return nil
}

// ListLabels lists all Gmail labels
func (c *Client) ListLabels() ([]*types.GmailLabel, error) {
	response, err := c.service.Users.Labels.List("me").Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list labels: %w", err)
	}
	
	var labels []*types.GmailLabel
	for _, label := range response.Labels {
		labelType := "user"
		if strings.HasPrefix(label.Id, "LABEL_") || strings.HasPrefix(label.Id, "CATEGORY_") {
			labelType = "system"
		}
		
		labels = append(labels, &types.GmailLabel{
			ID:   label.Id,
			Name: label.Name,
			Type: labelType,
		})
	}
	
	return labels, nil
}

// CreateLabel creates a new Gmail label
func (c *Client) CreateLabel(name string) (string, error) {
	label := &gmail.Label{
		Name:                   name,
		MessageListVisibility: "show",
		LabelListVisibility:   "labelShow",
	}
	
	result, err := c.service.Users.Labels.Create("me", label).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create label: %w", err)
	}
	
	return result.Id, nil
}

// GetContacts gets Gmail contacts (simplified - would need People API for full implementation)
func (c *Client) GetContacts(maxResults int, query string) ([]*types.Contact, error) {
	// This is a simplified implementation - for full contacts, you'd need the People API
	// For now, we'll extract contacts from recent emails
	searchQuery := query
	if searchQuery == "" {
		searchQuery = "in:sent OR in:inbox"
	}
	
	messages, err := c.SearchEmails(searchQuery, maxResults*2) // Get more messages to extract contacts
	if err != nil {
		return nil, err
	}
	
	contactMap := make(map[string]*types.Contact)
	
	for _, msg := range messages {
		if msg.From != "" {
			addr, err := mail.ParseAddress(msg.From)
			if err == nil {
				contactMap[addr.Address] = &types.Contact{
					Name:  addr.Name,
					Email: addr.Address,
				}
			}
		}
	}
	
	var contacts []*types.Contact
	count := 0
	for _, contact := range contactMap {
		if count >= maxResults {
			break
		}
		contacts = append(contacts, contact)
		count++
	}
	
	return contacts, nil
}

// GetEmailAnalytics gets email analytics
func (c *Client) GetEmailAnalytics(days int, query, groupBy string) (string, error) {
	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	
	searchQuery := query
	if searchQuery == "" {
		searchQuery = fmt.Sprintf("after:%s before:%s", 
			startDate.Format("2006/01/02"), 
			endDate.Format("2006/01/02"))
	} else {
		searchQuery += fmt.Sprintf(" after:%s before:%s", 
			startDate.Format("2006/01/02"), 
			endDate.Format("2006/01/02"))
	}
	
	messages, err := c.SearchEmails(searchQuery, 1000) // Get more messages for analytics
	if err != nil {
		return "", err
	}
	
	totalCount := len(messages)
	
	analytics := fmt.Sprintf("Email Analytics (%d days)\n", days)
	analytics += fmt.Sprintf("================\n\n")
	analytics += fmt.Sprintf("Total emails: %d\n", totalCount)
	analytics += fmt.Sprintf("Average per day: %.1f\n\n", float64(totalCount)/float64(days))
	
	// Group by analysis
	switch groupBy {
	case "sender":
		senderCount := make(map[string]int)
		for _, msg := range messages {
			if addr, err := mail.ParseAddress(msg.From); err == nil {
				senderCount[addr.Address]++
			}
		}
		
		analytics += "Top Senders:\n"
		for sender, count := range senderCount {
			if count > 1 {
				analytics += fmt.Sprintf("  %s: %d emails\n", sender, count)
			}
		}
		
	case "day":
		dayCount := make(map[string]int)
		for _, msg := range messages {
			if parsedDate, err := mail.ParseDate(msg.Date); err == nil {
				day := parsedDate.Format("2006-01-02")
				dayCount[day]++
			}
		}
		
		analytics += "Daily Breakdown:\n"
		for day, count := range dayCount {
			analytics += fmt.Sprintf("  %s: %d emails\n", day, count)
		}
	}
	
	return analytics, nil
}

// Helper methods

func (c *Client) createRawMessage(args *types.SendEmailArgs, isDraft bool) string {
	var message strings.Builder
	
	// Headers
	message.WriteString(fmt.Sprintf("To: %s\n", strings.Join(args.To, ", ")))
	
	if args.CC != nil && len(args.CC) > 0 {
		message.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(args.CC, ", ")))
	}
	
	if args.BCC != nil && len(args.BCC) > 0 {
		message.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(args.BCC, ", ")))
	}
	
	message.WriteString(fmt.Sprintf("Subject: %s\n", args.Subject))
	
	if args.InReplyTo != "" {
		message.WriteString(fmt.Sprintf("In-Reply-To: %s\n", args.InReplyTo))
	}
	
	// Content type
	if args.HTMLBody != "" {
		message.WriteString("Content-Type: text/html; charset=utf-8\n")
	} else {
		message.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}
	
	message.WriteString("\n") // Empty line between headers and body
	
	// Body
	if args.HTMLBody != "" {
		message.WriteString(args.HTMLBody)
	} else {
		message.WriteString(args.Body)
	}
	
	return message.String()
}

func (c *Client) extractBody(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// If this part has a body, extract it
	if payload.Body != nil && payload.Body.Data != "" {
		decoded, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			// Prefer plain text
			if payload.MimeType == "text/plain" {
				return string(decoded)
			}
			// Fallback to HTML
			if payload.MimeType == "text/html" {
				return string(decoded)
			}
		}
	}
	
	// If this part has nested parts, recursively extract
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if body := c.extractBody(part); body != "" {
				// Prefer plain text over HTML
				if part.MimeType == "text/plain" {
					return body
				}
			}
		}
		// If no plain text found, try HTML
		for _, part := range payload.Parts {
			if body := c.extractBody(part); body != "" && part.MimeType == "text/html" {
				return body
			}
		}
	}
	
	return ""
}