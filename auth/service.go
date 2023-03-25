package auth

import (
	"context"
	"github.com/tanyudii/core-go/ectx"
	"github.com/tanyudii/core-go/errutil"
)

type service struct {
	cfg *Config
}

func NewService(args ...ConfigFunc) Service {
	return &service{
		cfg: generateConfig(args...),
	}
}

func (s *service) IsPublicRoute(ctx context.Context, fullMethod string) (bool, error) {
	routeConfig, err := s.cfg.routeService.GetRouteConfig(ctx, fullMethod)
	if err != nil {
		return false, err
	}
	return routeConfig == nil, nil
}

func (s *service) Authenticate(ctx context.Context, fullMethod string) (context.Context, error) {
	session, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return nil, err
	}

	routeUserTypes := s.cfg.mapUserTypeRoutes[fullMethod]
	routePermissions := s.cfg.mapPermissionRoutes[fullMethod]
	routeScopes := s.cfg.mapScopeRoutes[fullMethod]

	routeConfig, err := s.cfg.routeService.GetRouteConfig(ctx, fullMethod)
	if err != nil {
		return nil, err
	} else if routeConfig != nil {
		routeUserTypes = append(routeUserTypes, routeConfig.GetUserTypes()...)
		routePermissions = append(routePermissions, routeConfig.GetPermissions()...)
		routeScopes = append(routeScopes, routeConfig.GetScopes()...)
	}

	//if user authorized with type, will be skip other middleware
	ok, err := s.authorizedUserType(session, routeUserTypes)
	if err != nil {
		return nil, err
	} else if ok {
		return ctx, nil
	}

	if err = s.authorizedPermission(session, routePermissions); err != nil {
		return nil, errutil.NewUnauthorizedError(err.Error())
	}

	if err = s.authorizedScope(session, routeScopes); err != nil {
		return nil, errutil.NewUnauthorizedError(err.Error())
	}

	return ctx, nil
}

func (s *service) authorizedUserType(session *ectx.EContext, userTypes []string) (bool, error) {
	//if trusted user type continue to process request
	ok, err := session.HasUserTypeByMapCode(s.cfg.mapUserTypeTrusted)
	if err != nil {
		return false, err
	} else if ok {
		return ok, nil
	}
	return session.HasUserType(userTypes)
}

func (s *service) authorizedPermission(session *ectx.EContext, permissions []string) error {
	_, err := session.HasPermission(permissions)
	return err
}

func (s *service) authorizedScope(session *ectx.EContext, scopes []string) error {
	_, err := session.HasScope(scopes)
	return err
}
