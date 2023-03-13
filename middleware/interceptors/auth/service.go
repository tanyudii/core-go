package auth

import (
	"context"
	"github.com/tanyudii/core-go/ectx"
	"github.com/tanyudii/core-go/errutil"
	"google.golang.org/grpc"
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

func (s *service) authenticate(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	//skip when route is public routes
	if s.cfg.mapPublicRoutes[info.FullMethod] {
		return ctx, nil
	}

	if newCtx, ok := s.authorizedInternalCall(ctx); ok {
		return newCtx, nil
	}

	newCtx, err := s.authenticateToken(ctx)
	if err != nil {
		return nil, err
	}

	session, err := ectx.FromContextWithErr(newCtx)
	if err != nil {
		return nil, err
	}

	//if user authorized with type, will be skip other middleware
	ok, err := s.authorizedUserType(session, info)
	if err != nil {
		return nil, err
	} else if ok {
		return newCtx, nil
	}

	if err = s.authorizedPermission(session, info); err != nil {
		return nil, errutil.NewUnauthorizedError(err.Error())
	}

	if err = s.authorizedScope(session, info); err != nil {
		return nil, errutil.NewUnauthorizedError(err.Error())
	}

	return newCtx, nil
}

func (s *service) authenticateToken(ctx context.Context) (context.Context, error) {
	md := ectx.FromIncoming(ctx)
	token := md.Get(strings.ToLower(ectx.RequestHeaderKeyAuthorization))
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
	return md.ToIncoming(ectx.NewContext(ctx, reqCtx)), nil
}

func (s *service) authorizedUserType(session *ectx.EContext, info *grpc.UnaryServerInfo) (bool, error) {
	//if trusted user type continue to process request
	ok, err := session.HasUserTypeByMapCode(s.cfg.mapUserTypeTrusted)
	if err != nil {
		return false, err
	} else if ok {
		return ok, nil
	}
	return session.HasUserType(s.cfg.mapUserTypeRoutes[info.FullMethod])
}

func (s *service) authorizedPermission(session *ectx.EContext, info *grpc.UnaryServerInfo) error {
	_, err := session.HasPermission(s.cfg.mapPermissionRoutes[info.FullMethod])
	return err
}

func (s *service) authorizedScope(session *ectx.EContext, info *grpc.UnaryServerInfo) error {
	_, err := session.HasScope(s.cfg.mapScopeRoutes[info.FullMethod])
	return err
}

func (s *service) authorizedInternalCall(ctx context.Context) (context.Context, bool) {
	md := ectx.FromIncoming(ctx)
	reqCtx := ectx.NewEContext(md)
	return md.ToIncoming(ectx.NewContext(ctx, reqCtx)), reqCtx.IsInternal()
}
