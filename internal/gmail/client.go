package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/phildougherty/mcp-gmail-go/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Client struct {
	service *gmail.Service
	config  *config.Config
	oauth   *oauth2.Config
}

func NewClient(cfg *config.Config) (*Client, error) {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.OAuth.ClientID,
		ClientSecret: cfg.OAuth.ClientSecret,
		RedirectURL:  cfg.OAuth.RedirectURL,
		Scopes:       []string{gmail.GmailModifyScope},
		Endpoint:     google.Endpoint,
	}

	client := &Client{
		config: cfg,
		oauth:  oauthConfig,
	}

	// Try to load existing credentials
	if err := client.loadCredentials(); err == nil {
		return client, nil
	}

	return client, nil
}

func (c *Client) loadCredentials() error {
	data, err := os.ReadFile(c.config.CredentialsPath)
	if err != nil {
		return err
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return err
	}

	httpClient := c.oauth.Client(context.Background(), &token)
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	c.service = service
	return nil
}

func (c *Client) Authenticate() error {
	authURL := c.oauth.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	
	fmt.Printf("Visit this URL to authorize the application: %s\n", authURL)
	
	// Try to open browser automatically
	if err := c.openBrowser(authURL); err != nil {
		fmt.Printf("Failed to open browser automatically: %v\n", err)
		fmt.Printf("Please manually open the URL above.\n")
	}

	// Start local server to handle callback
	return c.handleAuthCallback()
}

func (c *Client) openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func (c *Client) handleAuthCallback() error {
	codeChan := make(chan string)
	errChan := make(chan error)

	server := &http.Server{Addr: ":3000"}
	
	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no authorization code received")
			return
		}
		
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<h1>Authentication successful!</h1><p>You can close this window.</p>"))
		
		codeChan <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	var authCode string
	select {
	case authCode = <-codeChan:
	case err := <-errChan:
		return err
	case <-time.After(5 * time.Minute):
		return fmt.Errorf("authentication timeout")
	}

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	// Exchange code for token
	token, err := c.oauth.Exchange(context.Background(), authCode)
	if err != nil {
		return fmt.Errorf("failed to exchange authorization code: %w", err)
	}

	// Save credentials
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(c.config.CredentialsPath, data, 0600); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	// Initialize service
	httpClient := c.oauth.Client(context.Background(), token)
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %w", err)
	}

	c.service = service
	return nil
}

func (c *Client) Service() *gmail.Service {
	return c.service
}

func (c *Client) IsAuthenticated() bool {
	return c.service != nil
}