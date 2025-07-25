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

// InputSchema definitions for calendar tools
var (
	CreateEventSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "Event title/summary",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Event description",
			},
			"location": map[string]interface{}{
				"type":        "string",
				"description": "Event location",
			},
			"startTime": map[string]interface{}{
				"type":        "string",
				"description": "Start time in RFC3339 format (e.g., '2023-12-01T10:00:00Z')",
			},
			"endTime": map[string]interface{}{
				"type":        "string",
				"description": "End time in RFC3339 format (e.g., '2023-12-01T11:00:00Z')",
			},
			"startDate": map[string]interface{}{
				"type":        "string",
				"description": "Start date for all-day events (YYYY-MM-DD format)",
			},
			"endDate": map[string]interface{}{
				"type":        "string",
				"description": "End date for all-day events (YYYY-MM-DD format)",
			},
			"timeZone": map[string]interface{}{
				"type":        "string",
				"description": "Time zone (e.g., 'America/New_York')",
			},
			"allDay": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether this is an all-day event",
			},
			"calendarId": map[string]interface{}{
				"type":        "string",
				"description": "Calendar ID (defaults to primary calendar)",
			},
			"attendees": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of attendee email addresses",
			},
		},
		"required": []string{"summary"},
	}

	GetEventSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"eventId": map[string]interface{}{
				"type":        "string",
				"description": "ID of the event to retrieve",
			},
			"calendarId": map[string]interface{}{
				"type":        "string",
				"description": "Calendar ID (defaults to primary calendar)",
			},
		},
		"required": []string{"eventId"},
	}

	UpdateEventSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"eventId": map[string]interface{}{
				"type":        "string",
				"description": "ID of the event to update",
			},
			"calendarId": map[string]interface{}{
				"type":        "string",
				"description": "Calendar ID (defaults to primary calendar)",
			},
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "Event title/summary",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Event description",
			},
			"location": map[string]interface{}{
				"type":        "string",
				"description": "Event location",
			},
			"startTime": map[string]interface{}{
				"type":        "string",
				"description": "Start time in RFC3339 format",
			},
			"endTime": map[string]interface{}{
				"type":        "string",
				"description": "End time in RFC3339 format",
			},
			"timeZone": map[string]interface{}{
				"type":        "string",
				"description": "Time zone",
			},
		},
		"required": []string{"eventId"},
	}

	DeleteEventSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"eventId": map[string]interface{}{
				"type":        "string",
				"description": "ID of the event to delete",
			},
			"calendarId": map[string]interface{}{
				"type":        "string",
				"description": "Calendar ID (defaults to primary calendar)",
			},
		},
		"required": []string{"eventId"},
	}

	ListEventsSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"calendarId": map[string]interface{}{
				"type":        "string",
				"description": "Calendar ID (defaults to primary calendar)",
			},
			"timeMin": map[string]interface{}{
				"type":        "string",
				"description": "Lower bound for event start time (RFC3339 format)",
			},
			"timeMax": map[string]interface{}{
				"type":        "string",
				"description": "Upper bound for event start time (RFC3339 format)",
			},
			"maxResults": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum number of events to return",
			},
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Free text search terms",
			},
			"orderBy": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"startTime", "updated"},
				"description": "Order of the events",
			},
		},
	}

	ListCalendarsSchema = map[string]interface{}{
		"type":        "object",
		"properties":  map[string]interface{}{},
		"description": "Lists all calendars accessible to the user",
	}

	GetCalendarSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"calendarId": map[string]interface{}{
				"type":        "string",
				"description": "Calendar ID (defaults to primary calendar)",
			},
		},
	}

	CreateCalendarSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "Calendar title/name",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Calendar description",
			},
			"timeZone": map[string]interface{}{
				"type":        "string",
				"description": "Calendar time zone",
			},
		},
		"required": []string{"summary"},
	}

	DeleteCalendarSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"calendarId": map[string]interface{}{
				"type":        "string",
				"description": "Calendar ID to delete",
			},
		},
		"required": []string{"calendarId"},
	}

	FreeBusySchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"timeMin": map[string]interface{}{
				"type":        "string",
				"description": "Lower bound for free/busy query (RFC3339 format)",
			},
			"timeMax": map[string]interface{}{
				"type":        "string",
				"description": "Upper bound for free/busy query (RFC3339 format)",
			},
			"calendarIds": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of calendar IDs to query",
			},
		},
		"required": []string{"timeMin", "timeMax", "calendarIds"},
	}
)