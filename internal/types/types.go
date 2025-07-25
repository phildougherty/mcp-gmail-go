package types

// EmailMessage represents an email message
type EmailMessage struct {
	ID       string            `json:"id"`
	ThreadID string            `json:"threadId"`
	Subject  string            `json:"subject"`
	From     string            `json:"from"`
	To       string            `json:"to"`
	Date     string            `json:"date"`
	Body     string            `json:"body"`
	Labels   []string          `json:"labels"`
	Headers  map[string]string `json:"headers"`
}

// SendEmailArgs represents arguments for sending an email
type SendEmailArgs struct {
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	HTMLBody    string   `json:"htmlBody,omitempty"`
	CC          []string `json:"cc,omitempty"`
	BCC         []string `json:"bcc,omitempty"`
	ThreadID    string   `json:"threadId,omitempty"`
	InReplyTo   string   `json:"inReplyTo,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

// ReadEmailArgs represents arguments for reading an email
type ReadEmailArgs struct {
	MessageID string `json:"messageId"`
}

// SearchEmailsArgs represents arguments for searching emails
type SearchEmailsArgs struct {
	Query      string `json:"query"`
	MaxResults int    `json:"maxResults,omitempty"`
}

// ModifyEmailArgs represents arguments for modifying an email
type ModifyEmailArgs struct {
	MessageID       string   `json:"messageId"`
	AddLabelIDs     []string `json:"addLabelIds,omitempty"`
	RemoveLabelIDs  []string `json:"removeLabelIds,omitempty"`
}

// DeleteEmailArgs represents arguments for deleting an email
type DeleteEmailArgs struct {
	MessageID string `json:"messageId"`
}

// CreateLabelArgs represents arguments for creating a label
type CreateLabelArgs struct {
	Name                   string `json:"name"`
	MessageListVisibility string `json:"messageListVisibility,omitempty"`
	LabelListVisibility   string `json:"labelListVisibility,omitempty"`
}

// UpdateLabelArgs represents arguments for updating a label
type UpdateLabelArgs struct {
	ID                     string `json:"id"`
	Name                   string `json:"name,omitempty"`
	MessageListVisibility string `json:"messageListVisibility,omitempty"`
	LabelListVisibility   string `json:"labelListVisibility,omitempty"`
}

// DeleteLabelArgs represents arguments for deleting a label
type DeleteLabelArgs struct {
	ID string `json:"id"`
}

// BatchModifyEmailsArgs represents arguments for batch modifying emails
type BatchModifyEmailsArgs struct {
	MessageIDs     []string `json:"messageIds"`
	AddLabelIDs    []string `json:"addLabelIds,omitempty"`
	RemoveLabelIDs []string `json:"removeLabelIds,omitempty"`
	BatchSize      int      `json:"batchSize,omitempty"`
}

// BatchDeleteEmailsArgs represents arguments for batch deleting emails
type BatchDeleteEmailsArgs struct {
	MessageIDs []string `json:"messageIds"`
	BatchSize  int      `json:"batchSize,omitempty"`
}

// DownloadAttachmentArgs represents arguments for downloading an attachment
type DownloadAttachmentArgs struct {
	MessageID    string `json:"messageId"`
	AttachmentID string `json:"attachmentId"`
	Filename     string `json:"filename,omitempty"`
	SavePath     string `json:"savePath,omitempty"`
}

// GetContactsArgs represents arguments for getting contacts
type GetContactsArgs struct {
	MaxResults int    `json:"maxResults,omitempty"`
	Query      string `json:"query,omitempty"`
}

// EmailAnalyticsArgs represents arguments for email analytics
type EmailAnalyticsArgs struct {
	Days      int    `json:"days,omitempty"`
	Query     string `json:"query,omitempty"`
	GroupBy   string `json:"groupBy,omitempty"`
}

// GmailLabel represents a Gmail label
type GmailLabel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Contact represents a contact
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}