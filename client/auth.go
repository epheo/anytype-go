package client

import (
	"context"
	"net/http"

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

	req, err := ac.client.newRequest(ctx, http.MethodPost, urlPath, requestBody)
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

	req, err := ac.client.newRequest(ctx, http.MethodPost, urlPath, requestBody)
	if err != nil {
		return nil, err
	}

	var result anytype.CreateApiKeyResponse
	if err := ac.client.doRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
