package mocks

import (
	"context"

	"github.com/epheo/anytype-go"
)

// MockTypeService implements the anytype.TypeClient interface for testing
type MockTypeService struct {
	CurrentTypeKey   string
	ListFunc         func(ctx context.Context) ([]anytype.Type, error)
	GetFunc          func(ctx context.Context, typeKey string) (*anytype.Type, error)
	GetKeyByNameFunc func(ctx context.Context, name string) (string, error)
	CreateFunc       func(ctx context.Context, request anytype.CreateTypeRequest) (*anytype.TypeResponse, error)
}

// NewMockTypeService creates a new instance of MockTypeService with default implementations
func NewMockTypeService() *MockTypeService {
	return &MockTypeService{
		ListFunc: func(ctx context.Context) ([]anytype.Type, error) {
			return []anytype.Type{
				{
					ID:     "type_page_id",
					Key:    "page",
					Name:   "Page",
					Object: "type",
					Layout: "basic",
				},
				{
					ID:     "type_collection_id",
					Key:    "collection",
					Name:   "Collection",
					Object: "type",
					Layout: "set",
				},
			}, nil
		},
		GetFunc: func(ctx context.Context, typeKey string) (*anytype.Type, error) {
			return &anytype.Type{
				ID:     "type_" + typeKey + "_id",
				Key:    typeKey,
				Name:   "Page",
				Object: "type",
				Layout: "basic",
			}, nil
		},
		GetKeyByNameFunc: func(ctx context.Context, name string) (string, error) {
			if name == "Page" {
				return "page", nil
			}
			return "unknown", nil
		},
		CreateFunc: func(ctx context.Context, request anytype.CreateTypeRequest) (*anytype.TypeResponse, error) {
			return &anytype.TypeResponse{
				Type: anytype.Type{
					Key:    "mock-type-" + request.Name,
					Name:   request.Name,
					Icon:   request.Icon,
					Layout: request.Layout,
				},
			}, nil
		},
	}
}

// List calls the mock implementation
func (s *MockTypeService) List(ctx context.Context) ([]anytype.Type, error) {
	return s.ListFunc(ctx)
}

// Get calls the mock implementation
func (s *MockTypeService) Get(ctx context.Context, typeKey string) (*anytype.Type, error) {
	return s.GetFunc(ctx, typeKey)
}

// GetKeyByName calls the mock implementation
func (s *MockTypeService) GetKeyByName(ctx context.Context, name string) (string, error) {
	return s.GetKeyByNameFunc(ctx, name)
}

// Create calls the mock implementation
func (s *MockTypeService) Create(ctx context.Context, request anytype.CreateTypeRequest) (*anytype.TypeResponse, error) {
	return s.CreateFunc(ctx, request)
}

// Type returns a mock type context for a specific type
func (s *MockTypeService) Type(typeID string) anytype.TypeContext {
	return NewMockTypeContextService(typeID)
}

// MockTypeContextService implements the anytype.TypeContext interface for testing
type MockTypeContextService struct {
	TypeID           string
	GetFunc          func(ctx context.Context) (*anytype.TypeResponse, error)
	TemplatesService *MockTemplatesService
}

// NewMockTypeContextService creates a new instance of MockTypeContextService with default implementations
func NewMockTypeContextService(typeID string) *MockTypeContextService {
	return &MockTypeContextService{
		TypeID: typeID,
		GetFunc: func(ctx context.Context) (*anytype.TypeResponse, error) {
			return &anytype.TypeResponse{
				Type: anytype.Type{
					ID:     "type_" + typeID + "_id",
					Key:    typeID,
					Name:   "Page",
					Object: "type",
					Layout: "basic",
				},
			}, nil
		},
		TemplatesService: NewMockTemplatesService(),
	}
}

// Get calls the mock implementation
func (s *MockTypeContextService) Get(ctx context.Context) (*anytype.TypeResponse, error) {
	return s.GetFunc(ctx)
}

// Templates returns a mock templates service
func (s *MockTypeContextService) Templates() anytype.TemplateClient {
	return s.TemplatesService
}

// Template returns a mock template context for a specific template
func (s *MockTypeContextService) Template(templateID string) anytype.TemplateContext {
	return NewMockTemplateService(templateID)
}

// Templates returns a mock templates service
func (s *MockTypeService) Templates() anytype.TemplateClient {
	return NewMockTemplatesService()
}

// Template returns a mock template context for a specific template
func (s *MockTypeService) Template(templateID string) anytype.TemplateContext {
	return NewMockTemplateService(templateID)
}

// MockTemplatesService implements the anytype.TemplateClient interface for testing
type MockTemplatesService struct {
	ListFunc func(ctx context.Context) ([]anytype.Template, error)
	GetFunc  func(ctx context.Context, templateID string) (*anytype.Template, error)
}

// NewMockTemplatesService creates a new instance of MockTemplatesService with default implementations
func NewMockTemplatesService() *MockTemplatesService {
	return &MockTemplatesService{
		ListFunc: func(ctx context.Context) ([]anytype.Template, error) {
			return []anytype.Template{
				{
					ID:   "mock-template-id",
					Name: "Mock Template",
					Icon: &anytype.Icon{
						Format: anytype.IconFormatEmoji,
						Emoji:  "📄",
					},
				},
			}, nil
		},
		GetFunc: func(ctx context.Context, templateID string) (*anytype.Template, error) {
			return &anytype.Template{
				ID:   templateID,
				Name: "Mock Template",
				Icon: &anytype.Icon{
					Format: anytype.IconFormatEmoji,
					Emoji:  "📄",
				},
			}, nil
		},
	}
}

// List calls the mock implementation
func (s *MockTemplatesService) List(ctx context.Context) ([]anytype.Template, error) {
	return s.ListFunc(ctx)
}

// Get calls the mock implementation
func (s *MockTemplatesService) Get(ctx context.Context, templateID string) (*anytype.Template, error) {
	return s.GetFunc(ctx, templateID)
}

// MockTemplateService implements the anytype.TemplateContext interface for testing
type MockTemplateService struct {
	TemplateID string
	GetFunc    func(ctx context.Context) (*anytype.TemplateResponse, error)
}

// NewMockTemplateService creates a new instance of MockTemplateService with default implementations
func NewMockTemplateService(templateID string) *MockTemplateService {
	return &MockTemplateService{
		TemplateID: templateID,
		GetFunc: func(ctx context.Context) (*anytype.TemplateResponse, error) {
			return &anytype.TemplateResponse{
				Template: anytype.Template{
					ID:   templateID,
					Name: "Mock Template",
					Icon: &anytype.Icon{
						Format: anytype.IconFormatEmoji,
						Emoji:  "📄",
					},
				},
			}, nil
		},
	}
}

// Get calls the mock implementation
func (s *MockTemplateService) Get(ctx context.Context) (*anytype.TemplateResponse, error) {
	return s.GetFunc(ctx)
}
