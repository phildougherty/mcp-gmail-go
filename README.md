# Gmail MCP Server (Go)

A Model Context Protocol (MCP) server implementation in Go that provides Gmail integration with streaming HTTP/SSE transport.

## Features

- **Email Management**: Send, read, search, and organize emails
- **Label Management**: Create and manage Gmail labels 
- **Draft Support**: Create and manage email drafts
- **Contact Integration**: Access Gmail contacts
- **Email Analytics**: Get insights on email patterns
- **OAuth2 Authentication**: Secure Google OAuth2 flow
- **Streaming Transport**: Server-Sent Events (SSE) support
- **Batch Operations**: Efficient bulk email operations

## Tools Available

### Core Email Operations
- `send_email` - Send new emails with optional attachments
- `draft_email` - Create email drafts
- `read_email` - Retrieve email content by ID
- `search_emails` - Search emails using Gmail query syntax
- `modify_email` - Add/remove labels from emails
- `delete_email` - Permanently delete emails

### Label Management
- `list_email_labels` - List all available Gmail labels
- `create_label` - Create new Gmail labels

### Additional Features
- `get_contacts` - Retrieve Gmail contacts
- `email_analytics` - Get email statistics and insights

## Installation

1. Clone the repository:
```bash
git clone https://github.com/phildougherty/mcp-gmail-go.git
cd mcp-gmail-go
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up Google OAuth2 credentials:
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing
   - Enable Gmail API
   - Create OAuth2 credentials (Desktop application)
   - Download the credentials as `gcp-oauth.keys.json`

## Configuration

Place your `gcp-oauth.keys.json` file in:
- Current directory, or
- `~/.gmail-mcp/gcp-oauth.keys.json`

Environment variables:
- `GMAIL_OAUTH_PATH` - Custom path to OAuth keys file
- `GMAIL_CREDENTIALS_PATH` - Custom path to stored credentials

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

List available tools:
```bash
curl http://localhost:8080/mcp/tools
```

### 4. Using Tools
Send an email:
```bash
curl -X POST http://localhost:8080/mcp/tools/send_email \
  -H "Content-Type: application/json" \
  -d '{
    "arguments": {
      "to": ["recipient@example.com"],
      "subject": "Test Email",
      "body": "This is a test email from the Gmail MCP server."
    }
  }'
```

Search emails:
```bash
curl -X POST http://localhost:8080/mcp/tools/search_emails \
  -H "Content-Type: application/json" \
  -d '{
    "arguments": {
      "query": "from:someone@example.com",
      "maxResults": 10
    }
  }'
```

### 5. Server-Sent Events (SSE)
Connect to the SSE endpoint for real-time updates:
```bash
curl -N http://localhost:8080/mcp/sse
```

## Development

### Building
```bash
go build -o gmail-mcp-server main.go
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
│   ├── gmail/             # Gmail API client
│   │   ├── client.go      # OAuth2 and service setup
│   │   └── operations.go  # Gmail operations
│   └── mcp/               # MCP protocol implementation
│       ├── server.go      # HTTP server and routing
│       ├── tools.go       # Tool registry and handlers
│       └── types.go       # Data structures and schemas
├── go.mod
└── README.md
```

## API Reference

### Authentication Flow
1. Start server with `-auth` flag
2. Browser opens for Google OAuth2
3. Grant permissions to Gmail API
4. Credentials saved locally for future use

### Tool Schemas
All tools follow JSON Schema specifications. See `internal/mcp/types.go` for complete schema definitions.

### Error Handling
- HTTP status codes indicate request success/failure
- Tool execution errors returned in response body with `isError: true`
- Detailed error messages provided for debugging

### Security
- OAuth2 tokens stored securely in user home directory
- HTTPS recommended for production deployments
- CORS headers configured for cross-origin requests

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
- Check the Gmail API documentation
- Review MCP protocol specifications# mcp-gmail-go
