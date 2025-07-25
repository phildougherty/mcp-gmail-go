package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	OAuth struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURL  string `json:"redirect_url"`
	} `json:"oauth"`
	
	CredentialsPath string `json:"credentials_path,omitempty"`
	OAuthPath       string `json:"oauth_path,omitempty"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	
	// Default paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".gmail-mcp")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	cfg.CredentialsPath = filepath.Join(configDir, "credentials.json")
	cfg.OAuthPath = filepath.Join(configDir, "gcp-oauth.keys.json")
	
	// Override with environment variables if set
	if path := os.Getenv("GMAIL_CREDENTIALS_PATH"); path != "" {
		cfg.CredentialsPath = path
	}
	if path := os.Getenv("GMAIL_OAUTH_PATH"); path != "" {
		cfg.OAuthPath = path
	}
	
	// Load OAuth configuration
	oauthData, err := os.ReadFile(cfg.OAuthPath)
	if err != nil {
		// Try local directory as fallback
		localOAuthPath := "gcp-oauth.keys.json"
		if oauthData, err = os.ReadFile(localOAuthPath); err != nil {
			return nil, fmt.Errorf("OAuth keys file not found. Please place gcp-oauth.keys.json in current directory or %s", configDir)
		}
		// Copy to config directory for future use
		if err := os.WriteFile(cfg.OAuthPath, oauthData, 0600); err != nil {
			return nil, fmt.Errorf("failed to copy OAuth keys to config directory: %w", err)
		}
	}
	
	var oauthFile struct {
		Installed *struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"installed"`
		Web *struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"web"`
	}
	
	if err := json.Unmarshal(oauthData, &oauthFile); err != nil {
		return nil, fmt.Errorf("failed to parse OAuth keys file: %w", err)
	}
	
	if oauthFile.Installed != nil {
		cfg.OAuth.ClientID = oauthFile.Installed.ClientID
		cfg.OAuth.ClientSecret = oauthFile.Installed.ClientSecret
	} else if oauthFile.Web != nil {
		cfg.OAuth.ClientID = oauthFile.Web.ClientID
		cfg.OAuth.ClientSecret = oauthFile.Web.ClientSecret
	} else {
		return nil, fmt.Errorf("invalid OAuth keys file format. File should contain either 'installed' or 'web' credentials")
	}
	
	cfg.OAuth.RedirectURL = "http://localhost:3000/oauth2callback"
	
	return cfg, nil
}