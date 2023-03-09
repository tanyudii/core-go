package auth

import (
	"context"
	"github.com/tanyudii/core-go/ectx"
	"github.com/tanyudii/core-go/errutil"
)

func HasPermission(ctx context.Context, codes []string) error {
	valid := false
	if eCtx, ok := ectx.FromContext(ctx); ok {
		valid = eCtx.HasPermission(codes)
	}
	if valid {
		return nil
	}
	return errutil.ErrAuthPermissionNotAllowed
}

func HasAuth(ctx context.Context) error {
	if _, ok := ectx.FromContext(ctx); ok {
		return nil
	}
	return errutil.ErrAuthUnauthenticated
}
