package anytype

import (
	"context"
)

// MemberClient provides operations on space members
type MemberClient interface {
	// List retrieves all members of the space
	List(ctx context.Context) (*MemberListResponse, error)
}

// MemberContext provides operations on a specific member
type MemberContext interface {
	// Get retrieves details about this member
	Get(ctx context.Context) (*MemberResponse, error)
}

// Member represents a member of a space
type Member struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GlobalName string `json:"global_name"`
	Identity   string `json:"identity"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	Icon       *Icon  `json:"icon,omitempty"`
}

// MemberListResponse represents the response from List members
type MemberListResponse struct {
	Data []Member `json:"data"`
}

// MemberResponse represents the response from Get member
type MemberResponse struct {
	Member Member `json:"member"`
}

// MemberRole represents a role of a space member
type MemberRole string

// MemberStatus represents the status of a space member
type MemberStatus string

const (
	// Member roles
	MemberRoleViewer       MemberRole = "viewer"
	MemberRoleEditor       MemberRole = "editor"
	MemberRoleOwner        MemberRole = "owner"
	MemberRoleNoPermission MemberRole = "no_permission"

	// Member statuses
	MemberStatusJoining  MemberStatus = "joining"
	MemberStatusActive   MemberStatus = "active"
	MemberStatusRemoved  MemberStatus = "removed"
	MemberStatusDeclined MemberStatus = "declined"
	MemberStatusRemoving MemberStatus = "removing"
	MemberStatusCanceled MemberStatus = "canceled"
)
