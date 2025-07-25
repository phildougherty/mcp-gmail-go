package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/phildougherty/mcp-google-calendar-go/internal/calendar"
	"github.com/phildougherty/mcp-google-calendar-go/internal/types"
)

type ToolRegistry struct {
	calendarClient *calendar.Client
	tools          map[string]Tool
}

func NewToolRegistry(calendarClient *calendar.Client) *ToolRegistry {
	registry := &ToolRegistry{
		calendarClient: calendarClient,
		tools:          make(map[string]Tool),
	}
	
	registry.registerTools()
	return registry
}

func (r *ToolRegistry) registerTools() {
	r.tools["create_event"] = Tool{
		Name:        "create_event",
		Description: "Creates a new calendar event",
		InputSchema: CreateEventSchema,
	}
	
	r.tools["get_event"] = Tool{
		Name:        "get_event",
		Description: "Retrieves a specific calendar event",
		InputSchema: GetEventSchema,
	}
	
	r.tools["update_event"] = Tool{
		Name:        "update_event",
		Description: "Updates an existing calendar event",
		InputSchema: UpdateEventSchema,
	}
	
	r.tools["delete_event"] = Tool{
		Name:        "delete_event",
		Description: "Deletes a calendar event",
		InputSchema: DeleteEventSchema,
	}
	
	r.tools["list_events"] = Tool{
		Name:        "list_events",
		Description: "Lists calendar events",
		InputSchema: ListEventsSchema,
	}
	
	r.tools["list_calendars"] = Tool{
		Name:        "list_calendars",
		Description: "Lists available calendars",
		InputSchema: ListCalendarsSchema,
	}
	
	r.tools["get_calendar"] = Tool{
		Name:        "get_calendar",
		Description: "Retrieves a specific calendar",
		InputSchema: GetCalendarSchema,
	}
	
	r.tools["create_calendar"] = Tool{
		Name:        "create_calendar",
		Description: "Creates a new calendar",
		InputSchema: CreateCalendarSchema,
	}
	
	r.tools["delete_calendar"] = Tool{
		Name:        "delete_calendar",
		Description: "Deletes a calendar",
		InputSchema: DeleteCalendarSchema,
	}
	
	r.tools["get_freebusy"] = Tool{
		Name:        "get_freebusy",
		Description: "Gets free/busy information for calendars",
		InputSchema: FreeBusySchema,
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
	case "create_event":
		return r.handleCreateEvent(args)
	case "get_event":
		return r.handleGetEvent(args)
	case "update_event":
		return r.handleUpdateEvent(args)
	case "delete_event":
		return r.handleDeleteEvent(args)
	case "list_events":
		return r.handleListEvents(args)
	case "list_calendars":
		return r.handleListCalendars(args)
	case "get_calendar":
		return r.handleGetCalendar(args)
	case "create_calendar":
		return r.handleCreateCalendar(args)
	case "delete_calendar":
		return r.handleDeleteCalendar(args)
	case "get_freebusy":
		return r.handleGetFreeBusy(args)
	default:
		return nil, fmt.Errorf("tool implementation not found: %s", name)
	}
}

func (r *ToolRegistry) handleCreateEvent(args json.RawMessage) (*ToolResult, error) {
	var createArgs types.CreateEventArgs
	if err := json.Unmarshal(args, &createArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	eventID, err := r.calendarClient.CreateEvent(&createArgs)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to create event: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Event created successfully with ID: %s", eventID),
		}},
	}, nil
}

func (r *ToolRegistry) handleGetEvent(args json.RawMessage) (*ToolResult, error) {
	var getArgs struct {
		CalendarID string `json:"calendarId,omitempty"`
		EventID    string `json:"eventId"`
	}
	if err := json.Unmarshal(args, &getArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	event, err := r.calendarClient.GetEvent(getArgs.CalendarID, getArgs.EventID)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to get event: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	eventJSON, _ := json.MarshalIndent(event, "", "  ")
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: string(eventJSON),
		}},
	}, nil
}

func (r *ToolRegistry) handleUpdateEvent(args json.RawMessage) (*ToolResult, error) {
	var updateArgs types.UpdateEventArgs
	if err := json.Unmarshal(args, &updateArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	err := r.calendarClient.UpdateEvent(&updateArgs)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to update event: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Event %s updated successfully", updateArgs.EventID),
		}},
	}, nil
}

func (r *ToolRegistry) handleDeleteEvent(args json.RawMessage) (*ToolResult, error) {
	var deleteArgs types.DeleteEventArgs
	if err := json.Unmarshal(args, &deleteArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	err := r.calendarClient.DeleteEvent(deleteArgs.CalendarID, deleteArgs.EventID)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to delete event: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Event %s deleted successfully", deleteArgs.EventID),
		}},
	}, nil
}

func (r *ToolRegistry) handleListEvents(args json.RawMessage) (*ToolResult, error) {
	var listArgs types.ListEventsArgs
	if err := json.Unmarshal(args, &listArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	if listArgs.MaxResults == 0 {
		listArgs.MaxResults = 10
	}
	
	events, err := r.calendarClient.ListEvents(&listArgs)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to list events: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	eventsJSON, _ := json.MarshalIndent(events, "", "  ")
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: string(eventsJSON),
		}},
	}, nil
}

func (r *ToolRegistry) handleListCalendars(args json.RawMessage) (*ToolResult, error) {
	calendars, err := r.calendarClient.ListCalendars()
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to list calendars: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	calendarsJSON, _ := json.MarshalIndent(calendars, "", "  ")
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: string(calendarsJSON),
		}},
	}, nil
}

func (r *ToolRegistry) handleGetCalendar(args json.RawMessage) (*ToolResult, error) {
	var getArgs struct {
		CalendarID string `json:"calendarId,omitempty"`
	}
	if err := json.Unmarshal(args, &getArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	calendar, err := r.calendarClient.GetCalendar(getArgs.CalendarID)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to get calendar: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	calendarJSON, _ := json.MarshalIndent(calendar, "", "  ")
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: string(calendarJSON),
		}},
	}, nil
}

func (r *ToolRegistry) handleCreateCalendar(args json.RawMessage) (*ToolResult, error) {
	var createArgs types.CreateCalendarArgs
	if err := json.Unmarshal(args, &createArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	calendarID, err := r.calendarClient.CreateCalendar(&createArgs)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to create calendar: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Calendar created successfully with ID: %s", calendarID),
		}},
	}, nil
}

func (r *ToolRegistry) handleDeleteCalendar(args json.RawMessage) (*ToolResult, error) {
	var deleteArgs struct {
		CalendarID string `json:"calendarId"`
	}
	if err := json.Unmarshal(args, &deleteArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	err := r.calendarClient.DeleteCalendar(deleteArgs.CalendarID)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to delete calendar: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: fmt.Sprintf("Calendar %s deleted successfully", deleteArgs.CalendarID),
		}},
	}, nil
}

func (r *ToolRegistry) handleGetFreeBusy(args json.RawMessage) (*ToolResult, error) {
	var freeBusyArgs types.FreeBusyArgs
	if err := json.Unmarshal(args, &freeBusyArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	
	response, err := r.calendarClient.GetFreeBusy(&freeBusyArgs)
	if err != nil {
		return &ToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("Failed to get free/busy: %v", err),
			}},
			IsError: true,
		}, nil
	}
	
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	return &ToolResult{
		Content: []Content{{
			Type: "text",
			Text: string(responseJSON),
		}},
	}, nil
}