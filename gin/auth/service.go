package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tanyudii/core-go/auth"
	"github.com/tanyudii/core-go/ectx"
	"github.com/tanyudii/core-go/errutil"
	"strings"
)

type service struct {
	authService  auth.Service
	tokenService auth.TokenService
	cfg          *Config
}

func newService(
	authService auth.Service,
	tokenService auth.TokenService,
	args ...ConfigFunc,
) Service {
	return &service{
		authService:  authService,
		tokenService: tokenService,
		cfg:          generateConfig(args...),
	}
}

func (s *service) authenticate(c *gin.Context) (context.Context, error) {
	fullMethod := fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.RequestURI)

	//skip when route is public routes
	ok, err := s.authService.IsPublicRoute(c, fullMethod)
	if err != nil {
		return nil, err
	} else if ok {
		return c, nil
	}

	newCtx, err := s.authenticateToken(c)
	if err != nil {
		//skip if graphql mode and error is unauthenticated
		if s.cfg.graphqlMode && errors.Is(err, errutil.ErrAuthUnauthenticated) {
			return nil, nil
		}
		return nil, err
	}

	return s.authService.Authenticate(newCtx, fullMethod)
}

func (s *service) authenticateToken(c *gin.Context) (context.Context, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, errutil.ErrAuthUnauthenticated
	}

	splitToken := strings.Split(token, "Bearer ")
	if len(splitToken) != 2 {
		return nil, errutil.ErrAuthUnauthenticated
	}

	respTokenInfo, err := s.tokenService.TokenInfo(context.Background(), splitToken[1])
	if err != nil {
		return nil, err
	}

	eCtx, ok := ectx.FromContext(c.Request.Context())
	if !ok {
		eCtx = &ectx.EContext{}
	}

	eCtx.Scopes = respTokenInfo.Scope
	eCtx.Authorization = token

	tokenInfo := respTokenInfo.TokenInfo
	if tokenInfo != nil {
		eCtx.UserID = tokenInfo.UserID
		eCtx.UserSerial = tokenInfo.UserSerial
		eCtx.UserName = tokenInfo.UserName
		eCtx.UserEmail = tokenInfo.UserEmail
		eCtx.UserType = tokenInfo.UserType
		eCtx.CompanyID = tokenInfo.CompanyID
		eCtx.CompanySerial = tokenInfo.CompanySerial
		eCtx.CompanyName = tokenInfo.CompanyName
		eCtx.Permissions = strings.Join(tokenInfo.Permissions, ",")
	}

	clientInfo := respTokenInfo.ClientInfo
	if clientInfo != nil {
		eCtx.ClientID = clientInfo.ClientID
		eCtx.ClientName = clientInfo.ClientName
	}

	return ectx.NewContext(c.Request.Context(), eCtx), nil

}
