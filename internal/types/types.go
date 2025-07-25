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
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone,omitempty"`
	Organization string `json:"organization,omitempty"`
}

// Calendar types

// CalendarEvent represents a calendar event
type CalendarEvent struct {
	ID            string             `json:"id"`
	Summary       string             `json:"summary"`
	Description   string             `json:"description,omitempty"`
	Location      string             `json:"location,omitempty"`
	StartTime     string             `json:"startTime,omitempty"`
	EndTime       string             `json:"endTime,omitempty"`
	StartDate     string             `json:"startDate,omitempty"`
	EndDate       string             `json:"endDate,omitempty"`
	StartTimeZone string             `json:"startTimeZone,omitempty"`
	EndTimeZone   string             `json:"endTimeZone,omitempty"`
	AllDay        bool               `json:"allDay,omitempty"`
	Creator       string             `json:"creator,omitempty"`
	Organizer     string             `json:"organizer,omitempty"`
	Status        string             `json:"status,omitempty"`
	HTMLLink      string             `json:"htmlLink,omitempty"`
	Created       string             `json:"created,omitempty"`
	Updated       string             `json:"updated,omitempty"`
	Attendees     []*EventAttendee   `json:"attendees,omitempty"`
	Reminders     []*EventReminder   `json:"reminders,omitempty"`
}

// CreateEventArgs represents arguments for creating an event
type CreateEventArgs struct {
	Summary     string            `json:"summary"`
	Description string            `json:"description,omitempty"`
	Location    string            `json:"location,omitempty"`
	StartTime   string            `json:"startTime,omitempty"`
	EndTime     string            `json:"endTime,omitempty"`
	StartDate   string            `json:"startDate,omitempty"`
	EndDate     string            `json:"endDate,omitempty"`
	TimeZone    string            `json:"timeZone,omitempty"`
	AllDay      bool              `json:"allDay,omitempty"`
	CalendarID  string            `json:"calendarId,omitempty"`
	Attendees   []string          `json:"attendees,omitempty"`
	Reminders   []*EventReminder  `json:"reminders,omitempty"`
}

// UpdateEventArgs represents arguments for updating an event
type UpdateEventArgs struct {
	EventID     string           `json:"eventId"`
	CalendarID  string           `json:"calendarId,omitempty"`
	Summary     string           `json:"summary,omitempty"`
	Description string           `json:"description,omitempty"`
	Location    string           `json:"location,omitempty"`
	StartTime   string           `json:"startTime,omitempty"`
	EndTime     string           `json:"endTime,omitempty"`
	TimeZone    string           `json:"timeZone,omitempty"`
	Attendees   []string         `json:"attendees,omitempty"`
	Reminders   []*EventReminder `json:"reminders,omitempty"`
}

// ListEventsArgs represents arguments for listing events
type ListEventsArgs struct {
	CalendarID  string `json:"calendarId,omitempty"`
	TimeMin     string `json:"timeMin,omitempty"`
	TimeMax     string `json:"timeMax,omitempty"`
	MaxResults  int    `json:"maxResults,omitempty"`
	Query       string `json:"query,omitempty"`
	OrderBy     string `json:"orderBy,omitempty"`
}

// DeleteEventArgs represents arguments for deleting an event
type DeleteEventArgs struct {
	CalendarID string `json:"calendarId,omitempty"`
	EventID    string `json:"eventId"`
}

// EventAttendee represents an event attendee
type EventAttendee struct {
	Email          string `json:"email"`
	DisplayName    string `json:"displayName,omitempty"`
	ResponseStatus string `json:"responseStatus,omitempty"`
	Organizer      bool   `json:"organizer,omitempty"`
}

// EventReminder represents an event reminder
type EventReminder struct {
	Method  string `json:"method"`
	Minutes int    `json:"minutes"`
}

// Calendar represents a calendar
type Calendar struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	Primary     bool   `json:"primary,omitempty"`
	AccessRole  string `json:"accessRole,omitempty"`
	TimeZone    string `json:"timeZone,omitempty"`
}

// CreateCalendarArgs represents arguments for creating a calendar
type CreateCalendarArgs struct {
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	TimeZone    string `json:"timeZone,omitempty"`
}

// FreeBusyArgs represents arguments for free/busy query
type FreeBusyArgs struct {
	TimeMin     string   `json:"timeMin"`
	TimeMax     string   `json:"timeMax"`
	CalendarIDs []string `json:"calendarIds"`
}

// FreeBusyResponse represents free/busy response
type FreeBusyResponse struct {
	TimeMin   string                       `json:"timeMin"`
	TimeMax   string                       `json:"timeMax"`
	Calendars map[string]*FreeBusyCalendar `json:"calendars"`
}

// FreeBusyCalendar represents free/busy info for a calendar
type FreeBusyCalendar struct {
	Busy []*TimePeriod `json:"busy"`
}

// TimePeriod represents a time period
type TimePeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}