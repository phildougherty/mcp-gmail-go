# Google Calendar MCP Server (Go)

A Model Context Protocol (MCP) server implementation in Go that provides Google Calendar integration with HTTP transport.

## Features

- **Event Management**: Create, read, update, and delete calendar events
- **Calendar Management**: Manage multiple calendars  
- **Free/Busy Queries**: Check availability across calendars
- **All-day Events**: Support for both timed and all-day events
- **Attendee Management**: Add and manage event attendees
- **Time Zone Support**: Handle events across different time zones
- **OAuth2 Authentication**: Secure Google OAuth2 flow
- **JSON-RPC 2.0**: Standard MCP protocol implementation

## Tools Available

### Event Operations
- `create_event` - Create new calendar events with attendees and reminders
- `get_event` - Retrieve event details by ID
- `update_event` - Modify existing events
- `delete_event` - Remove events from calendar
- `list_events` - Search and filter calendar events

### Calendar Management
- `list_calendars` - List all accessible calendars
- `get_calendar` - Get details for a specific calendar
- `create_calendar` - Create new calendars
- `delete_calendar` - Remove calendars

### Availability
- `get_freebusy` - Query free/busy information across calendars

## Installation

1. Clone the repository:
```bash
git clone https://github.com/phildougherty/mcp-google-calendar-go.git
cd mcp-google-calendar-go
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up Google OAuth2 credentials:
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing
   - Enable Google Calendar API
   - Create OAuth2 credentials (Desktop application)
   - Download the credentials as `gcp-oauth.keys.json`

## Configuration

Place your `gcp-oauth.keys.json` file in:
- Current directory, or
- `~/.calendar-mcp/gcp-oauth.keys.json`

Environment variables:
- `CALENDAR_OAUTH_PATH` - Custom path to OAuth keys file
- `CALENDAR_CREDENTIALS_PATH` - Custom path to stored credentials

## Usage

### 1. Authentication
First, authenticate with Google:
```bash
go run main.go -auth
```
This will open a browser window for OAuth2 authentication.

### 2. Start the Server
```bash
go run main.go -port 8080
```

Optional flags:
- `-port` - Server port (default: 8080)
- `-debug` - Enable debug logging

### 3. Test the Server
Check server health:
```bash
curl http://localhost:8080/health
```

### 4. Using MCP Tools
The server implements JSON-RPC 2.0 over HTTP. Example requests:

Initialize connection:
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {"name": "test-client", "version": "1.0.0"}
    }
  }'
```

List available tools:
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'
```

Create a calendar event:
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "create_event",
      "arguments": {
        "summary": "Team Meeting",
        "description": "Weekly team sync",
        "location": "Conference Room A",
        "startTime": "2024-01-15T10:00:00Z",
        "endTime": "2024-01-15T11:00:00Z",
        "timeZone": "America/New_York",
        "attendees": ["colleague@example.com"]
      }
    }
  }'
```

List upcoming events:
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "list_events",
      "arguments": {
        "timeMin": "2024-01-01T00:00:00Z",
        "timeMax": "2024-12-31T23:59:59Z",
        "maxResults": 10,
        "orderBy": "startTime"
      }
    }
  }'
```

## Development

### Building
```bash
go build -o calendar-mcp-server main.go
```

### Running Tests
```bash
go test ./...
```

### Project Structure
```
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── calendar/          # Google Calendar API client
│   │   ├── client.go      # OAuth2 and service setup
│   │   └── operations.go  # Calendar operations
│   ├── mcp/               # MCP protocol implementation
│   │   ├── server.go      # HTTP server and JSON-RPC
│   │   ├── tools.go       # Tool registry and handlers
│   │   └── types.go       # Data structures and schemas
│   └── types/             # Shared type definitions
│       └── types.go
├── go.mod
└── README.md
```

## API Reference

### Authentication Flow
1. Start server with `-auth` flag
2. Browser opens for Google OAuth2
3. Grant permissions to Calendar API
4. Credentials saved locally for future use

### Tool Schemas
All tools follow JSON Schema specifications. See `internal/mcp/types.go` for complete schema definitions.

### Event Time Formats
- **Timed Events**: Use RFC3339 format (e.g., `2024-01-15T10:00:00Z`)
- **All-day Events**: Use date format (e.g., `2024-01-15`)
- **Time Zones**: Use IANA time zone names (e.g., `America/New_York`)

### Error Handling
- JSON-RPC 2.0 error responses for protocol errors
- Tool execution errors returned in response with `isError: true`
- Detailed error messages provided for debugging

### Security
- OAuth2 tokens stored securely in user home directory
- HTTPS recommended for production deployments
- CORS headers configured for cross-origin requests

## Example Integration

This server can be integrated with MCP clients like Claude Desktop. Add to your MCP configuration:

```json
{
  "mcpServers": {
    "google-calendar": {
      "command": "/path/to/calendar-mcp-server",
      "args": ["-port", "8080"]
    }
  }
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
- Create an issue on GitHub
- Check the Google Calendar API documentation
- Review MCP protocol specifications