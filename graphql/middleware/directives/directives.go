package directives

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/tanyudii/core-go/common"
	"github.com/tanyudii/core-go/ectx"
	"github.com/tanyudii/core-go/errutil"
)

func Auth(ctx context.Context, _ interface{}, next graphql.Resolver) (res interface{}, err error) {
	eCtx, err := ectx.FromContextWithErr(ctx)
	if eCtx.UserID == "" || err != nil {
		return nil, errutil.ErrAuthUnauthenticated
	}
	return next(ctx)
}

func HasPermission(ctx context.Context, obj interface{}, next graphql.Resolver, p string) (res interface{}, err error) {
	if _, err = Auth(ctx, obj, next); err != nil {
		return nil, err
	}

	eCtx, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return nil, err
	}

	permissions := common.ParseStringToSliceBySeparator(p, "|")
	ok, err := eCtx.HasPermission(permissions)
	if err != nil {
		return nil, err
	} else if ok {
		return next(ctx)
	}

	return nil, errutil.ErrAuthPermissionNotAllowed
}

func HasScope(ctx context.Context, obj interface{}, next graphql.Resolver, s string) (res interface{}, err error) {
	if _, err = Auth(ctx, obj, next); err != nil {
		return nil, err
	}

	eCtx, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return nil, err
	}

	scopes := common.ParseStringToSliceBySeparator(s, "|")
	ok, err := eCtx.HasScope(scopes)
	if err != nil {
		return nil, err
	} else if ok {
		return next(ctx)
	}

	return nil, errutil.ErrAuthScopeNotAllowed
}

func HasUserType(ctx context.Context, obj interface{}, next graphql.Resolver, t string) (res interface{}, err error) {
	if _, err = Auth(ctx, obj, next); err != nil {
		return nil, err
	}

	eCtx, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return nil, err
	}

	userTypes := common.ParseStringToSliceBySeparator(t, "|")
	ok, err := eCtx.HasUserType(userTypes)
	if err != nil {
		return nil, err
	} else if ok {
		return next(ctx)
	}

	return nil, errutil.ErrAuthUserTypeNotAllowed
}
