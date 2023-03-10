package auth

import (
	"context"
	"github.com/tanyudii/core-go/ectx"
	"github.com/tanyudii/core-go/errutil"
)

func HasAuth(ctx context.Context) error {
	if _, ok := ectx.FromContext(ctx); ok {
		return nil
	}
	return errutil.ErrAuthUnauthenticated
}

func HasPermission(ctx context.Context, codes []string) error {
	eCtx, ok := ectx.FromContext(ctx)
	if !ok {
		return errutil.ErrAuthPermissionNotAllowed
	}
	return eCtx.HasPermission(codes)
}

func HasScope(ctx context.Context, codes []string) error {
	eCtx, ok := ectx.FromContext(ctx)
	if !ok {
		return errutil.ErrAuthScopeNotAllowed
	}
	return eCtx.HasScope(codes)
}

func HasUserType(ctx context.Context, codes []string) error {
	eCtx, ok := ectx.FromContext(ctx)
	if !ok {
		return errutil.ErrAuthUserTypeNotAllowed
	}
	return eCtx.HasUserType(codes)
}
