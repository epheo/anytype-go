package anytype

import (
	"context"
)

type AuthClient interface {
	CreateChallenge(ctx context.Context, appName string) (*CreateChallengeResponse, error)
	CreateApiKey(ctx context.Context, challengeID string, code string) (*CreateApiKeyResponse, error)

	// Deprecated: Use CreateChallenge instead
	DisplayCode(ctx context.Context, appName string) (*DisplayCodeResponse, error)

	// Deprecated: Use CreateApiKey instead
	GetToken(ctx context.Context, challengeID string, code string) (*TokenResponse, error)
}

type CreateChallengeResponse struct {
	ChallengeID string `json:"challenge_id"`
}

type CreateApiKeyResponse struct {
	ApiKey string `json:"api_key"`
}

// Deprecated: Use CreateChallengeResponse instead
type DisplayCodeResponse struct {
	ChallengeID string `json:"challenge_id"`
}

// Deprecated: Use CreateApiKeyResponse instead
type TokenResponse struct {
	AppKey string `json:"app_key"`
}
