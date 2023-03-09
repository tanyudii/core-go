package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tanyudii/core-go/ectx"
)

// services
type (
	TokenService interface {
		TokenInfo(ctx context.Context, jwtToken string) (*TokenInfoResponse, error)
	}
	Service interface {
		authenticate(c *gin.Context) (newCtx context.Context, err error)
		authenticateGin(c *gin.Context) (context.Context, error)
		authenticateToken(c *gin.Context, authorization string) (context.Context, error)
		authorizedUserType(session *ectx.EContext, fullMethod string) bool
		authorizedPermission(session *ectx.EContext, fullMethod string) error
		authorizedScope(session *ectx.EContext, fullMethod string) error
		authorizedInternalCall(ctx context.Context) bool
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
