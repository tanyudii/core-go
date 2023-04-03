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

func GetAcceptLanguage(ctx context.Context, defaultVal ...string) string {
	var acceptLang string
	if eCtx, ok := ectx.FromContext(ctx); ok {
		acceptLang = eCtx.AcceptLanguage
	}
	if acceptLang == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return acceptLang
}

func DuplicateCtx(ctx context.Context) (context.Context, error) {
	eCtx, err := ectx.FromContextWithErr(ctx)
	if err != nil {
		return nil, err
	}
	return ectx.NewContext(context.Background(), eCtx), nil
}

func CreateInternalEContextDummy() *ectx.EContext {
	return &ectx.EContext{
		UserID:         "DummyUserID",
		UserName:       "DummyUserName",
		UserEmail:      "DummyUserEmail",
		UserSerial:     "DummyUserSerial",
		UserType:       "DummyUserType",
		CompanyID:      "DummyCompanyID",
		CompanySerial:  "DummyCompanySerial",
		CompanyName:    "DummyCompanyName",
		Permissions:    "DummyPermissions",
		ClientID:       "DummyClientID",
		ClientName:     "DummyClientName",
		Scopes:         "*",
		IsInternalCall: true,
	}
}

func CreateGRPCContextDummy(ctx context.Context) context.Context {
	eCtxDummy := CreateInternalEContextDummy()
	return ectx.ParseToGrpcCtx(ectx.NewContext(ctx, eCtxDummy))
}
