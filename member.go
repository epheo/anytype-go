package anytype

import (
	"context"
)

type MemberClient interface {
	List(ctx context.Context) (*MemberListResponse, error)
}

type MemberContext interface {
	Get(ctx context.Context) (*MemberResponse, error)
}

type Member struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GlobalName string `json:"global_name"`
	Identity   string `json:"identity"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	Icon       *Icon  `json:"icon,omitempty"`
}

type MemberListResponse struct {
	Data []Member `json:"data"`
}

type MemberResponse struct {
	Member Member `json:"member"`
}

type MemberRole string

type MemberStatus string

const (
	MemberRoleViewer       MemberRole = "viewer"
	MemberRoleEditor       MemberRole = "editor"
	MemberRoleOwner        MemberRole = "owner"
	MemberRoleNoPermission MemberRole = "no_permission"

	MemberStatusJoining  MemberStatus = "joining"
	MemberStatusActive   MemberStatus = "active"
	MemberStatusRemoved  MemberStatus = "removed"
	MemberStatusDeclined MemberStatus = "declined"
	MemberStatusRemoving MemberStatus = "removing"
	MemberStatusCanceled MemberStatus = "canceled"
)
