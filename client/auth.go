package client

import (
	"context"
	"log"

	"github.com/epheo/anytype-go"
)

type AuthClientImpl struct {
	client *ClientImpl
}

type CreateChallengeRequest struct {
	AppName string `json:"app_name"`
}

type CreateApiKeyRequest struct {
	ChallengeID string `json:"challenge_id"`
	Code        string `json:"code"`
}

func (ac *AuthClientImpl) CreateChallenge(ctx context.Context, appName string) (*anytype.CreateChallengeResponse, error) {
	urlPath := "/auth/challenges"

	requestBody := CreateChallengeRequest{
		AppName: appName,
	}

	req, err := ac.client.newRequest(ctx, "POST", urlPath, requestBody)
	if err != nil {
		return nil, err
	}

	var result anytype.CreateChallengeResponse
	if err := ac.client.doRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateApiKey completes the authentication flow by providing a code
func (ac *AuthClientImpl) CreateApiKey(ctx context.Context, challengeID string, code string) (*anytype.CreateApiKeyResponse, error) {
	urlPath := "/auth/api_keys"

	requestBody := CreateApiKeyRequest{
		ChallengeID: challengeID,
		Code:        code,
	}

	req, err := ac.client.newRequest(ctx, "POST", urlPath, requestBody)
	if err != nil {
		return nil, err
	}

	var result anytype.CreateApiKeyResponse
	if err := ac.client.doRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Deprecated: Use CreateChallenge instead
func (ac *AuthClientImpl) DisplayCode(ctx context.Context, appName string) (*anytype.DisplayCodeResponse, error) {
	// Log warning to stderr to help users discover they're using deprecated API
	log.Println("Warning: DisplayCode is deprecated, use CreateChallenge instead")

	resp, err := ac.CreateChallenge(ctx, appName)
	if err != nil {
		return nil, err
	}

	return &anytype.DisplayCodeResponse{
		ChallengeID: resp.ChallengeID,
	}, nil
}

// Deprecated: Use CreateApiKey instead
func (ac *AuthClientImpl) GetToken(ctx context.Context, challengeID string, code string) (*anytype.TokenResponse, error) {
	// Log warning to stderr to help users discover they're using deprecated API
	log.Println("Warning: GetToken is deprecated, use CreateApiKey instead")

	resp, err := ac.CreateApiKey(ctx, challengeID, code)
	if err != nil {
		return nil, err
	}

	return &anytype.TokenResponse{
		AppKey: resp.ApiKey,
	}, nil
}
