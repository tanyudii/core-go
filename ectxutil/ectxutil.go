package ectxutil

import (
	"context"
	"github.com/tanyudii/core-go/ectx"
)

func GetUserID(ctx context.Context) (string, error) {
	eCtx, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return "", err
	}
	return eCtx.UserID, nil
}

func GetUserSerial(ctx context.Context) (string, error) {
	eCtx, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return "", err
	}
	return eCtx.UserSerial, nil
}

func GetUserType(ctx context.Context) (string, error) {
	eCtx, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return "", err
	}
	return eCtx.UserType, nil
}

func GetCompanyID(ctx context.Context) (string, error) {
	eCtx, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return "", err
	}
	return eCtx.CompanyID, nil
}

func GetCompanySerial(ctx context.Context) (string, error) {
	eCtx, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return "", err
	}
	return eCtx.CompanySerial, nil
}
