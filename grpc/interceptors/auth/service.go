package auth

import (
	"context"
	"github.com/tanyudii/core-go/auth"
	"github.com/tanyudii/core-go/ectx"
	"github.com/tanyudii/core-go/errutil"
	"google.golang.org/grpc"
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

func (s *service) authenticate(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	//skip when route is public routes
	ok, err := s.authService.IsPublicRoute(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	} else if ok {
		return ctx, nil
	}

	if newCtx, ok := s.authorizedInternalCall(ctx); ok {
		return newCtx, nil
	}

	newCtx, err := s.authenticateToken(ctx)
	if err != nil {
		return nil, err
	}

	return s.authService.Authenticate(newCtx, info.FullMethod)
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

func (s *service) authorizedInternalCall(ctx context.Context) (context.Context, bool) {
	md := ectx.FromIncoming(ctx)
	eCtx := ectx.NewEContext(md)
	return md.ToIncoming(ectx.NewContext(ctx, eCtx)), eCtx.IsInternal()
}
