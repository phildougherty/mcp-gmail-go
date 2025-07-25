package calendar

import (
	"fmt"
	"time"

	"github.com/phildougherty/mcp-google-calendar-go/internal/types"
	"google.golang.org/api/calendar/v3"
)

// CreateEvent creates a new calendar event
func (c *Client) CreateEvent(args *types.CreateEventArgs) (string, error) {
	event := &calendar.Event{
		Summary:     args.Summary,
		Description: args.Description,
		Location:    args.Location,
	}

	// Set start time
	if args.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, args.StartTime)
		if err != nil {
			return "", fmt.Errorf("invalid start time format: %w", err)
		}
		event.Start = &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: args.TimeZone,
		}
	}

	// Set end time
	if args.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, args.EndTime)
		if err != nil {
			return "", fmt.Errorf("invalid end time format: %w", err)
		}
		event.End = &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: args.TimeZone,
		}
	}

	// Handle all-day events
	if args.AllDay {
		if args.StartDate != "" {
			event.Start = &calendar.EventDateTime{
				Date: args.StartDate,
			}
		}
		if args.EndDate != "" {
			event.End = &calendar.EventDateTime{
				Date: args.EndDate,
			}
		}
	}

	// Add attendees
	if len(args.Attendees) > 0 {
		attendees := make([]*calendar.EventAttendee, len(args.Attendees))
		for i, email := range args.Attendees {
			attendees[i] = &calendar.EventAttendee{
				Email: email,
			}
		}
		event.Attendees = attendees
	}

	// Add reminders
	if len(args.Reminders) > 0 {
		reminders := make([]*calendar.EventReminder, len(args.Reminders))
		for i, reminder := range args.Reminders {
			reminders[i] = &calendar.EventReminder{
				Method:  reminder.Method,
				Minutes: int64(reminder.Minutes),
			}
		}
		event.Reminders = &calendar.EventReminders{
			UseDefault: false,
			Overrides:  reminders,
		}
	}

	calendarID := args.CalendarID
	if calendarID == "" {
		calendarID = "primary"
	}

	result, err := c.service.Events.Insert(calendarID, event).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create event: %w", err)
	}

	return result.Id, nil
}

// GetEvent retrieves a calendar event by ID
func (c *Client) GetEvent(calendarID, eventID string) (*types.CalendarEvent, error) {
	if calendarID == "" {
		calendarID = "primary"
	}

	event, err := c.service.Events.Get(calendarID, eventID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return c.convertToCalendarEvent(event), nil
}

// UpdateEvent updates an existing calendar event
func (c *Client) UpdateEvent(args *types.UpdateEventArgs) error {
	calendarID := args.CalendarID
	if calendarID == "" {
		calendarID = "primary"
	}

	// Get existing event
	event, err := c.service.Events.Get(calendarID, args.EventID).Do()
	if err != nil {
		return fmt.Errorf("failed to get existing event: %w", err)
	}

	// Update fields
	if args.Summary != "" {
		event.Summary = args.Summary
	}
	if args.Description != "" {
		event.Description = args.Description
	}
	if args.Location != "" {
		event.Location = args.Location
	}

	// Update start time
	if args.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, args.StartTime)
		if err != nil {
			return fmt.Errorf("invalid start time format: %w", err)
		}
		event.Start = &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: args.TimeZone,
		}
	}

	// Update end time
	if args.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, args.EndTime)
		if err != nil {
			return fmt.Errorf("invalid end time format: %w", err)
		}
		event.End = &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: args.TimeZone,
		}
	}

	_, err = c.service.Events.Update(calendarID, args.EventID, event).Do()
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// DeleteEvent deletes a calendar event
func (c *Client) DeleteEvent(calendarID, eventID string) error {
	if calendarID == "" {
		calendarID = "primary"
	}

	err := c.service.Events.Delete(calendarID, eventID).Do()
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

// ListEvents lists calendar events
func (c *Client) ListEvents(args *types.ListEventsArgs) ([]*types.CalendarEvent, error) {
	calendarID := args.CalendarID
	if calendarID == "" {
		calendarID = "primary"
	}

	call := c.service.Events.List(calendarID)

	// Set time range
	if args.TimeMin != "" {
		call = call.TimeMin(args.TimeMin)
	}
	if args.TimeMax != "" {
		call = call.TimeMax(args.TimeMax)
	}

	// Set max results
	if args.MaxResults > 0 {
		call = call.MaxResults(int64(args.MaxResults))
	}

	// Set query
	if args.Query != "" {
		call = call.Q(args.Query)
	}

	// Set order
	if args.OrderBy != "" {
		call = call.OrderBy(args.OrderBy)
	}

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	var events []*types.CalendarEvent
	for _, event := range response.Items {
		events = append(events, c.convertToCalendarEvent(event))
	}

	return events, nil
}

// ListCalendars lists available calendars
func (c *Client) ListCalendars() ([]*types.Calendar, error) {
	response, err := c.service.CalendarList.List().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list calendars: %w", err)
	}

	var calendars []*types.Calendar
	for _, cal := range response.Items {
		calendars = append(calendars, &types.Calendar{
			ID:          cal.Id,
			Summary:     cal.Summary,
			Description: cal.Description,
			Primary:     cal.Primary,
			AccessRole:  cal.AccessRole,
			TimeZone:    cal.TimeZone,
		})
	}

	return calendars, nil
}

// GetCalendar retrieves a specific calendar
func (c *Client) GetCalendar(calendarID string) (*types.Calendar, error) {
	if calendarID == "" {
		calendarID = "primary"
	}

	cal, err := c.service.Calendars.Get(calendarID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar: %w", err)
	}

	return &types.Calendar{
		ID:          cal.Id,
		Summary:     cal.Summary,
		Description: cal.Description,
		TimeZone:    cal.TimeZone,
	}, nil
}

// CreateCalendar creates a new calendar
func (c *Client) CreateCalendar(args *types.CreateCalendarArgs) (string, error) {
	cal := &calendar.Calendar{
		Summary:     args.Summary,
		Description: args.Description,
		TimeZone:    args.TimeZone,
	}

	result, err := c.service.Calendars.Insert(cal).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create calendar: %w", err)
	}

	return result.Id, nil
}

// DeleteCalendar deletes a calendar
func (c *Client) DeleteCalendar(calendarID string) error {
	err := c.service.Calendars.Delete(calendarID).Do()
	if err != nil {
		return fmt.Errorf("failed to delete calendar: %w", err)
	}

	return nil
}

// GetFreeBusy gets free/busy information
func (c *Client) GetFreeBusy(args *types.FreeBusyArgs) (*types.FreeBusyResponse, error) {
	items := make([]*calendar.FreeBusyRequestItem, len(args.CalendarIDs))
	for i, calID := range args.CalendarIDs {
		items[i] = &calendar.FreeBusyRequestItem{
			Id: calID,
		}
	}

	request := &calendar.FreeBusyRequest{
		TimeMin: args.TimeMin,
		TimeMax: args.TimeMax,
		Items:   items,
	}

	response, err := c.service.Freebusy.Query(request).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get free/busy: %w", err)
	}

	result := &types.FreeBusyResponse{
		TimeMin: response.TimeMin,
		TimeMax: response.TimeMax,
		Calendars: make(map[string]*types.FreeBusyCalendar),
	}

	for calID, cal := range response.Calendars {
		busy := make([]*types.TimePeriod, len(cal.Busy))
		for i, period := range cal.Busy {
			busy[i] = &types.TimePeriod{
				Start: period.Start,
				End:   period.End,
			}
		}

		result.Calendars[calID] = &types.FreeBusyCalendar{
			Busy: busy,
		}
	}

	return result, nil
}

// Helper function to convert Google Calendar event to our type
func (c *Client) convertToCalendarEvent(event *calendar.Event) *types.CalendarEvent {
	calEvent := &types.CalendarEvent{
		ID:          event.Id,
		Summary:     event.Summary,
		Description: event.Description,
		Location:    event.Location,
		Creator:     event.Creator.Email,
		Organizer:   event.Organizer.Email,
		Status:      event.Status,
		HTMLLink:    event.HtmlLink,
		Created:     event.Created,
		Updated:     event.Updated,
	}

	// Start time
	if event.Start != nil {
		if event.Start.DateTime != "" {
			calEvent.StartTime = event.Start.DateTime
			calEvent.StartTimeZone = event.Start.TimeZone
		} else if event.Start.Date != "" {
			calEvent.StartDate = event.Start.Date
			calEvent.AllDay = true
		}
	}

	// End time
	if event.End != nil {
		if event.End.DateTime != "" {
			calEvent.EndTime = event.End.DateTime
			calEvent.EndTimeZone = event.End.TimeZone
		} else if event.End.Date != "" {
			calEvent.EndDate = event.End.Date
			calEvent.AllDay = true
		}
	}

	// Attendees
	if len(event.Attendees) > 0 {
		attendees := make([]*types.EventAttendee, len(event.Attendees))
		for i, attendee := range event.Attendees {
			attendees[i] = &types.EventAttendee{
				Email:          attendee.Email,
				DisplayName:    attendee.DisplayName,
				ResponseStatus: attendee.ResponseStatus,
				Organizer:      attendee.Organizer,
			}
		}
		calEvent.Attendees = attendees
	}

	return calEvent
}