package anytype

import (
	"context"
)

type AuthClient interface {
	CreateChallenge(ctx context.Context, appName string) (*CreateChallengeResponse, error)
	CreateApiKey(ctx context.Context, challengeID string, code string) (*CreateApiKeyResponse, error)
}

type CreateChallengeResponse struct {
	ChallengeID string `json:"challenge_id"`
}

type CreateApiKeyResponse struct {
	ApiKey string `json:"api_key"`
}
