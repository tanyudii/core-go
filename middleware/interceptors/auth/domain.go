package auth

import (
	"context"
	"github.com/tanyudii/core-go/ectx"
	"google.golang.org/grpc"
)

// services
type (
	TokenService interface {
		TokenInfo(ctx context.Context, jwtToken string) (*TokenInfoResponse, error)
	}
	Service interface {
		authenticate(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error)
		authenticateGRPC(ctx context.Context) (context.Context, error)
		authenticateToken(md *ectx.ContextMD, authorization string) error
		authorizedUserType(session *ectx.EContext, info *grpc.UnaryServerInfo) bool
		authorizedPermission(session *ectx.EContext, info *grpc.UnaryServerInfo) error
		authorizedScope(session *ectx.EContext, info *grpc.UnaryServerInfo) error
		authorizedInternalCall(ctx context.Context) (context.Context, bool)
	}
)

// entities
type (
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
