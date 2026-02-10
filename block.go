package anytype

import (
	"context"
)

type BlockClient interface {
	List(ctx context.Context) ([]Block, error)
	Get(ctx context.Context, blockID string) (*Block, error)
	Create(ctx context.Context, request CreateBlockRequest) (*Block, error)
	Update(ctx context.Context, blockID string, request UpdateBlockRequest) error
	Delete(ctx context.Context, blockID string) error
}

type Block struct {
	ID              string
	Text            *Text
	File            *File
	Property        *Property
	ChildrenIDs     []string `json:"children_ids"`
	Align           string
	VerticalAlign   string `json:"vertical_align"`
	BackgroundColor string `json:"background_color"`
}

type Text struct {
	Content  string
	Style    string
	Markdown string
}

type File struct {
	Name           string
	Hash           string
	Mime           string
	Size           int64
	Type           string
	State          string
	Style          string
	AddedAt        int64  `json:"added_at"`
	TargetObjectID string `json:"target_object_id"`
}

type CreateBlockRequest struct {
	Text           *Text
	File           *File
	Property       *Property
	ParentID       string `json:"parent_id"`
	TargetPosition int    `json:"target_position"`
}

type UpdateBlockRequest struct {
	Text            *Text
	File            *File
	Property        *Property
	Align           string
	VerticalAlign   string `json:"vertical_align"`
	BackgroundColor string `json:"background_color"`
}
