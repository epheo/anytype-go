package anytype

import (
	"context"
	"iter"
)

type MemberClient interface {
	List(ctx context.Context, opts ...ListOption) (*Page[Member], error)
	All(ctx context.Context, opts ...ListOption) iter.Seq2[Member, error]
}

type MemberContext interface {
	Get(ctx context.Context) (*MemberResponse, error)
}

type Member struct {
	ID         string       `json:"id"`
	Object     string       `json:"object"`
	Name       string       `json:"name"`
	GlobalName string       `json:"global_name"`
	Identity   string       `json:"identity"`
	Role       MemberRole   `json:"role"`
	Status     MemberStatus `json:"status"`
	Icon       *Icon        `json:"icon,omitempty"`
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
