package auth

import (
	"context"
	"github.com/tanyudii/core-go/ectx"
)

type Service interface {
	IsPublicRoute(ctx context.Context, fullMethod string) (bool, error)
	Authenticate(ctx context.Context, fullMethod string) (context.Context, error)

	authorizedUserType(session *ectx.EContext, userTypes []string) (bool, error)
	authorizedPermission(session *ectx.EContext, permissions []string) error
	authorizedScope(session *ectx.EContext, scopes []string) error

	getRouteConfig(ctx context.Context, fullMethod string) (RouteConfig, error)
}

type (
	TokenService interface {
		TokenInfo(ctx context.Context, jwtToken string) (*TokenInfoResponse, error)
	}
	TokenInfoResponse struct {
		TokenInfo  *TokenInfo
		ClientInfo *ClientInfo
		Scope      string
	}
	TokenInfo struct {
		UserID         string
		UserSerial     string
		UserName       string
		UserEmail      string
		UserType       string
		CompanyID      string
		CompanySerial  string
		CompanyName    string
		Permissions    []string
		IsInternalCall bool
	}
	ClientInfo struct {
		ClientID   string
		ClientName string
	}
)

type (
	RouteService interface {
		GetRouteConfig(ctx context.Context, fullMethod string) (RouteConfig, error)
	}
	RouteConfig interface {
		GetPermissions() []string
		GetScopes() []string
		GetUserTypes() []string
	}
)
