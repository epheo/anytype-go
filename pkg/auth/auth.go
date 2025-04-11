package auth

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/epheo/anyblog/pkg/anytype"
)

// Common errors
var (
	ErrInvalidCode       = errors.New("invalid authorization code")
	ErrAuthFailed        = errors.New("authentication failed")
	ErrChallengeExpired  = errors.New("challenge expired")
	ErrTokenExpired      = errors.New("token expired")
)

// Constants for authentication
const (
	defaultAPIURL   = "http://localhost:31009"
	tokenExpiryDays = 30
	appName         = "anyblog"
)

// AuthManager handles authentication operations
type AuthManager struct {
	apiURL string
}

// NewAuthManager creates a new AuthManager instance
func NewAuthManager(apiURL string) *AuthManager {
	if apiURL == "" {
		apiURL = defaultAPIURL
	}
	return &AuthManager{apiURL: apiURL}
}

// GetConfiguration loads or creates new auth configuration
func (am *AuthManager) GetConfiguration() (*anytype.AuthConfig, error) {
	// Try to load existing auth configuration
	config, err := loadAuthConfig()
	if err == nil && config.AppKey != "" && !isTokenExpired(config.Timestamp) {
		fmt.Println("Loaded existing authentication from config file")
		return config, nil
	}

	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	fmt.Println("No valid authentication found, starting new authentication process")
	return am.createNewAuthConfig()
}

// createNewAuthConfig performs authentication and creates new config
func (am *AuthManager) createNewAuthConfig() (*anytype.AuthConfig, error) {
	// Get challenge ID
	challengeID, err := am.getChallengeID()
	if err != nil {
		return nil, fmt.Errorf("error getting challenge ID: %w", err)
	}

	// Prompt user for authorization code
	code, err := am.promptForAuthCode()
	if err != nil {
		return nil, err
	}

	// Get session token and app key
	sessionToken, appKey, err := am.getAuthToken(challengeID, code)
	if err != nil {
		return nil, fmt.Errorf("error getting auth token: %w", err)
	}

	// Create and save the new auth configuration
	config := &anytype.AuthConfig{
		ApiURL:       am.apiURL,
		SessionToken: sessionToken,
		AppKey:       appKey,
		Timestamp:    time.Now(),
	}

	if err := saveAuthConfig(config); err != nil {
		fmt.Printf("Warning: Failed to save auth config: %v\n", err)
	} else {
		fmt.Println("Authentication saved to config file")
	}

	return config, nil
}

// promptForAuthCode prompts the user to enter the authorization code
func (am *AuthManager) promptForAuthCode() (string, error) {
	fmt.Println("\nPlease enter the authorization code displayed in Anytype:")
	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}

	code = strings.TrimSpace(code)
	if len(code) < 4 {
		return "", ErrInvalidCode
	}

	return code, nil
}

// getChallengeID gets a challenge ID from the API
func (am *AuthManager) getChallengeID() (string, error) {
	url := fmt.Sprintf("%s/v1/auth/display_code?app_name=%s", am.apiURL, appName)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var challengeResp anytype.ChallengeResponse
	if err := json.Unmarshal(body, &challengeResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	return challengeResp.ChallengeID, nil
}

// getAuthToken gets authentication token using challenge ID and code
func (am *AuthManager) getAuthToken(challengeID, code string) (string, string, error) {
	url := fmt.Sprintf("%s/v1/auth/token?challenge_id=%s&code=%s", am.apiURL, challengeID, code)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response: %w", err)
	}

	var authResp anytype.AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", "", fmt.Errorf("error parsing response: %w", err)
	}

	return authResp.SessionToken, authResp.AppKey, nil
}

// SetEnvironmentVariables sets environment variables for the current process
func SetEnvironmentVariables(config *anytype.AuthConfig) error {
	if config == nil {
		return errors.New("nil config provided")
	}

	os.Setenv("ANYTYPE_API_URL", config.ApiURL)
	os.Setenv("ANYTYPE_SESSION_TOKEN", config.SessionToken)
	os.Setenv("ANYTYPE_APP_KEY", config.AppKey)

	fmt.Println("\nEnvironment variables set for this process:")
	fmt.Printf("ANYTYPE_API_URL=%s\n", config.ApiURL)
	fmt.Printf("ANYTYPE_SESSION_TOKEN=%s\n", config.SessionToken)
	fmt.Printf("ANYTYPE_APP_KEY=%s\n", config.AppKey)

	return nil
}

// isTokenExpired checks if the auth token has expired
func isTokenExpired(timestamp time.Time) bool {
	return time.Since(timestamp) > tokenExpiryDays*24*time.Hour
}
