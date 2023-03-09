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

	newCtx, err := s.authenticateGin(c)
	if err != nil {
		return nil, err
	}

	//skip if graphql mode and request is not contain authorization
	if s.cfg.graphqlMode && newCtx == nil {
		return nil, nil
	}

	session, err := ectx.FromContextWithErr(newCtx)
	if err != nil {
		return nil, err
	}

	//if user authorized with type, will be skip other middleware
	if s.authorizedUserType(session, fullMethod) {
		return nil, nil
	}

	if err = s.authorizedPermission(session, fullMethod); err != nil {
		return nil, errutil.NewUnauthorizedError(err.Error())
	}

	if err = s.authorizedScope(session, fullMethod); err != nil {
		return nil, errutil.NewUnauthorizedError(err.Error())
	}

	return newCtx, nil
}

func (s *service) authenticateGin(c *gin.Context) (context.Context, error) {
	jwtToken := c.GetHeader("Authorization")
	if !s.cfg.graphqlMode && jwtToken == "" {
		return nil, errutil.ErrAuthUnauthenticated
	}

	newCtx, err := s.authenticateToken(c, jwtToken)
	if err != nil {
		return nil, err
	}

	//skip if graphql mode and request is not contain authorization
	if s.cfg.graphqlMode && newCtx == nil {
		return nil, nil
	}

	return newCtx, nil
}

func (s *service) authenticateToken(c *gin.Context, authorization string) (context.Context, error) {
	splitToken := strings.Split(authorization, "Bearer ")
	if len(splitToken) != 2 {
		if s.cfg.graphqlMode {
			return nil, nil
		}
		return nil, errutil.ErrAuthUnauthenticated
	}

	respTokenInfo, err := s.tokenService.TokenInfo(context.Background(), splitToken[1])
	if err != nil {
		return nil, err
	}

	md := ectx.ContextMD{}

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

	md.Set(strings.ToLower(ectx.RequestHeaderKeyScopes), respTokenInfo.Scope)

	reqCtx := ectx.NewEContext(md)
	return ectx.NewContext(c.Request.Context(), reqCtx), nil

}

func (s *service) authorizedUserType(session *ectx.EContext, fullMethod string) bool {
	userType := session.UserType
	if userType == "" {
		return false
	}

	//if trusted user type continue to process request
	for ut := range s.cfg.mapUserTypeTrusted {
		if strings.ToLower(ut) == strings.ToLower(userType) {
			return true
		}
	}

	//skip immediately when route not configured or user type empty
	routeUserTypes, ok := s.cfg.mapUserTypeRoutes[fullMethod]
	if !ok {
		return false
	}

	for _, ut := range routeUserTypes {
		if strings.ToLower(ut) == strings.ToLower(userType) {
			return true
		}
	}

	return false
}

func (s *service) authorizedPermission(session *ectx.EContext, fullMethod string) error {
	//skip immediately when route not configured
	routePermissions, ok := s.cfg.mapPermissionRoutes[fullMethod]
	if !ok {
		return nil
	}

	stringPermissions := session.Permissions
	if stringPermissions == "" {
		return errors.New("user permission is not configured")
	}

	permissions := strings.Split(stringPermissions, ",")
	mapPermission := map[string]bool{}
	for _, p := range permissions {
		mapPermission[p] = true
	}

	for _, rp := range routePermissions {
		if mapPermission[rp] {
			return nil
		}
	}

	return errutil.ErrAuthPermissionNotAllowed
}

func (s *service) authorizedScope(session *ectx.EContext, fullMethod string) error {
	//skip immediately when route not configured
	routeScopes, ok := s.cfg.mapScopeRoutes[fullMethod]
	if !ok {
		return nil
	}

	stringScopes := session.Scopes
	if stringScopes == "" {
		return errors.New("user scope is not configured")
	}

	scopes := strings.Split(stringScopes, ",")
	mapScope := map[string]bool{}
	for _, p := range scopes {
		mapScope[p] = true
	}

	for _, rp := range routeScopes {
		if mapScope[rp] {
			return nil
		}
	}

	return errutil.ErrAuthScopeNotAllowed
}

func (s *service) authorizedInternalCall(ctx context.Context) bool {
	eCtx, ok := ectx.FromContext(ctx)
	if !ok {
		return false
	}
	return eCtx.IsInternal()
}
