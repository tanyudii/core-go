package auth

import (
	"context"
	"errors"
	"github.com/tanyudii/core-go/ectx"
	"github.com/tanyudii/core-go/errutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type TokenService interface {
	TokenInfo(ctx context.Context, jwtToken string) (*TokenInfoResponse, error)
}

type service struct {
	tokenService TokenService
	cfg          *Config
}

func newService(
	tokenService TokenService,
	args ...ConfigFunc,
) *service {
	return &service{
		tokenService: tokenService,
		cfg:          generateConfig(args...),
	}
}

func (s *service) authenticate(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	//skip when route is public routes
	if s.cfg.mapPublicRoutes[info.FullMethod] {
		return ctx, nil
	}

	if s.authorizedInternalCall(ctx) {
		return ctx, nil
	}

	newCtx, err := s.authenticateGRPC(ctx)
	if err != nil {
		return nil, err
	}

	session, err := ectx.FromContextWithErr(newCtx)
	if err != nil {
		return nil, err
	}

	if err = s.authorizedPermission(session, info); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if err = s.authorizedScope(session, info); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if err = s.authorizedUserType(session, info); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return newCtx, nil
}

func (s *service) authenticateGRPC(ctx context.Context) (context.Context, error) {
	md := ectx.FromIncoming(ctx)
	jwtToken := md.Get(strings.ToLower("authorization"))
	if jwtToken == "" {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err := s.authenticateToken(&md, jwtToken); err != nil {
		return nil, err
	}

	reqCtx := ectx.NewEContext(md)
	return md.ToIncoming(ectx.NewContext(ctx, reqCtx)), nil
}

func (s *service) authenticateToken(md *ectx.ContextMD, authorization string) error {
	splitToken := strings.Split(authorization, "Bearer ")
	if len(splitToken) != 2 {
		return status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	respTokenInfo, err := s.tokenService.TokenInfo(context.Background(), splitToken[1])
	if err != nil {
		return err
	}

	tokenInfo := respTokenInfo.TokenInfo
	if tokenInfo != nil {
		md.Set(strings.ToLower(ectx.RequestHeaderKeyUserID), tokenInfo.UserID)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyUserName), tokenInfo.UserName)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyUserEmail), tokenInfo.UserEmail)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyUserType), tokenInfo.UserType)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyPermissions), strings.Join(tokenInfo.Permissions, ","))
	}

	clientInfo := respTokenInfo.ClientInfo
	if clientInfo != nil {
		md.Set(strings.ToLower(ectx.RequestHeaderKeyClientID), clientInfo.ClientID)
		md.Set(strings.ToLower(ectx.RequestHeaderKeyClientName), clientInfo.ClientName)
	}

	md.Set(strings.ToLower(ectx.RequestHeaderKeyScopes), respTokenInfo.Scope)

	return nil
}

func (s *service) authorizedUserType(session *ectx.EContext, info *grpc.UnaryServerInfo) error {
	//skip immediately when route not configured
	routeUserTypes, ok := s.cfg.mapUserTypeRoutes[info.FullMethod]
	if !ok {
		return nil
	}

	userType := session.UserType
	if userType == "" {
		return errors.New("user type is not configured")
	}

	for _, ut := range routeUserTypes {
		if strings.ToLower(ut) == strings.ToLower(userType) {
			return nil
		}
	}

	return errutil.NewUnauthorizedError("user type is not allowed")
}

func (s *service) authorizedPermission(session *ectx.EContext, info *grpc.UnaryServerInfo) error {
	//skip immediately when route not configured
	routePermissions, ok := s.cfg.mapPermissionRoutes[info.FullMethod]
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

	return errutil.NewUnauthorizedError("user permission is not allowed")
}

func (s *service) authorizedScope(session *ectx.EContext, info *grpc.UnaryServerInfo) error {
	//skip immediately when route not configured
	routeScopes, ok := s.cfg.mapScopeRoutes[info.FullMethod]
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

	return errutil.NewUnauthorizedError("user scope is not allowed")
}

func (s *service) authorizedInternalCall(ctx context.Context) bool {
	eCtx, ok := ectx.FromContext(ctx)
	return ok && eCtx.IsInternalCall
}
