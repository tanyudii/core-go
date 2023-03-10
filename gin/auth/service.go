package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tanyudii/core-go/ectx"
	"github.com/tanyudii/core-go/errutil"
	"strings"
)

type service struct {
	tokenService TokenService
	cfg          *Config
}

func newService(
	tokenService TokenService,
	args ...ConfigFunc,
) Service {
	return &service{
		tokenService: tokenService,
		cfg:          generateConfig(args...),
	}
}

func (s *service) authenticate(c *gin.Context) (context.Context, error) {
	fullMethod := fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.RequestURI)

	//skip when route is public routes
	if s.cfg.mapPublicRoutes[fullMethod] {
		return nil, nil
	}

	newCtx, err := s.authenticateToken(c)
	if err != nil {
		//skip if graphql mode and error is unauthenticated
		if s.cfg.graphqlMode && errors.Is(err, errutil.ErrAuthUnauthenticated) {
			return nil, nil
		}
		return nil, err
	}

	session, err := ectx.FromContextWithErr(newCtx)
	if err != nil {
		return nil, err
	}

	//if user authorized with type, will be skip other middleware
	ok, err := s.authorizedUserType(session, fullMethod)
	if err != nil {
		return nil, err
	} else if ok {
		return newCtx, nil
	}

	if err = s.authorizedPermission(session, fullMethod); err != nil {
		return nil, errutil.NewUnauthorizedError(err.Error())
	}

	if err = s.authorizedScope(session, fullMethod); err != nil {
		return nil, errutil.NewUnauthorizedError(err.Error())
	}

	return newCtx, nil
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

	md := ectx.ContextMD{}
	md.Set(strings.ToLower(ectx.RequestHeaderKeyScopes), respTokenInfo.Scope)
	md.Set(strings.ToLower(ectx.RequestHeaderKeyAuthorization), token)

	tokenInfo := respTokenInfo.TokenInfo
	if tokenInfo != nil {
		md.Set(strings.ToLower(ectx.RequestHeaderKeyUserID), tokenInfo.UserID)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyUserSerial), tokenInfo.UserSerial)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyUserName), tokenInfo.UserName)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyUserEmail), tokenInfo.UserEmail)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyUserType), tokenInfo.UserType)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyCompanyID), tokenInfo.CompanyID)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyCompanySerial), tokenInfo.CompanySerial)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyCompanyName), tokenInfo.CompanyName)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyPermissions), strings.Join(tokenInfo.Permissions, ","))
	}

	clientInfo := respTokenInfo.ClientInfo
	if clientInfo != nil {
		md.Set(strings.ToLower(ectx.RequestHeaderKeyClientID), clientInfo.ClientID)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyClientName), clientInfo.ClientName)
	}

	reqCtx := ectx.NewEContext(md)
	return ectx.NewContext(c.Request.Context(), reqCtx), nil

}

func (s *service) authorizedUserType(session *ectx.EContext, fullMethod string) (bool, error) {
	//if trusted user type continue to process request
	ok, err := session.HasUserTypeByMapCode(s.cfg.mapUserTypeTrusted)
	if err != nil {
		return false, err
	} else if ok {
		return ok, nil
	}
	return session.HasUserType(s.cfg.mapUserTypeRoutes[fullMethod])
}

func (s *service) authorizedPermission(session *ectx.EContext, fullMethod string) error {
	_, err := session.HasPermission(s.cfg.mapPermissionRoutes[fullMethod])
	return err
}

func (s *service) authorizedScope(session *ectx.EContext, fullMethod string) error {
	_, err := session.HasScope(s.cfg.mapScopeRoutes[fullMethod])
	return err
}

func (s *service) authorizedInternalCall(ctx context.Context) bool {
	eCtx, ok := ectx.FromContext(ctx)
	if !ok {
		return false
	}
	return eCtx.IsInternal()
}
