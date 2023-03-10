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
	ok, err := eCtx.HasPermission(codes)
	if err != nil {
		return err
	} else if ok {
		return nil
	}
	return errutil.ErrAuthPermissionNotAllowed
}

func HasScope(ctx context.Context, codes []string) error {
	eCtx, ok := ectx.FromContext(ctx)
	if !ok {
		return errutil.ErrAuthScopeNotAllowed
	}
	ok, err := eCtx.HasScope(codes)
	if err != nil {
		return err
	} else if ok {
		return nil
	}
	return errutil.ErrAuthScopeNotAllowed
}

func HasUserType(ctx context.Context, codes []string) error {
	eCtx, ok := ectx.FromContext(ctx)
	if !ok {
		return errutil.ErrAuthUserTypeNotAllowed
	}
	ok, err := eCtx.HasUserType(codes)
	if err != nil {
		return err
	} else if ok {
		return nil
	}
	return errutil.ErrAuthUserTypeNotAllowed
}
