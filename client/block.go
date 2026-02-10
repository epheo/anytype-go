package client

import (
	"context"

	"github.com/epheo/anytype-go"
)

// Block operations are not yet supported by the Anytype API.
// These methods are placeholders for future API support.
type BlockClientImpl struct {
	client   *ClientImpl
	spaceID  string
	objectID string
}

func (bc *BlockClientImpl) List(ctx context.Context) ([]anytype.Block, error) {
	// TODO: Implement
	return nil, nil
}

func (bc *BlockClientImpl) Get(ctx context.Context, blockID string) (*anytype.Block, error) {
	// TODO: Implement
	return nil, nil
}

func (bc *BlockClientImpl) Create(ctx context.Context, request anytype.CreateBlockRequest) (*anytype.Block, error) {
	// TODO: Implement
	return nil, nil
}

func (bc *BlockClientImpl) Update(ctx context.Context, blockID string, request anytype.UpdateBlockRequest) error {
	// TODO: Implement
	return nil
}

func (bc *BlockClientImpl) Delete(ctx context.Context, blockID string) error {
	// TODO: Implement
	return nil
}
