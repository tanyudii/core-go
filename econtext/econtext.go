package ectx

import (
	"context"
	"errors"
	"strings"
)

const (
	RequestHeaderKeyUserID        = "UserID"
	RequestHeaderKeyUserName      = "UserName"
	RequestHeaderKeyUserEmail     = "UserEmail"
	RequestHeaderKeyUserSerial    = "UserSerial"
	RequestHeaderKeyUserType      = "UserType"
	RequestHeaderKeyCompanyID     = "CompanyID"
	RequestHeaderKeyCompanySerial = "CompanySerial"
	RequestHeaderKeyCompanyName   = "CompanyName"
	RequestHeaderKeyPermissions   = "Permissions"
	RequestHeaderKeyScopes        = "Scopes"
	RequestHeaderKeyClientID      = "ClientID"
	RequestHeaderKeyClientName    = "ClientName"
)

type EContext struct {
	UserID        string
	UserName      string
	UserEmail     string
	UserSerial    string
	UserType      string
	CompanyID     string
	CompanySerial string
	CompanyName   string
	Permissions   string //separated by ","
	ClientID      string
	ClientName    string
	Scopes        string //separated by ","
}

func NewEContext(md ContextMD) *EContext {
	return &EContext{
		UserID:        md.Get(strings.ToLower(RequestHeaderKeyUserID)),
		UserName:      md.Get(strings.ToLower(RequestHeaderKeyUserName)),
		UserEmail:     md.Get(strings.ToLower(RequestHeaderKeyUserEmail)),
		UserSerial:    md.Get(strings.ToLower(RequestHeaderKeyUserSerial)),
		UserType:      md.Get(strings.ToLower(RequestHeaderKeyUserType)),
		CompanyID:     md.Get(strings.ToLower(RequestHeaderKeyCompanyID)),
		CompanySerial: md.Get(strings.ToLower(RequestHeaderKeyCompanySerial)),
		CompanyName:   md.Get(strings.ToLower(RequestHeaderKeyCompanyName)),
		Permissions:   md.Get(strings.ToLower(RequestHeaderKeyPermissions)),
		ClientID:      md.Get(strings.ToLower(RequestHeaderKeyClientID)),
		ClientName:    md.Get(strings.ToLower(RequestHeaderKeyClientName)),
		Scopes:        md.Get(strings.ToLower(RequestHeaderKeyScopes)),
	}
}

func (c *EContext) ToContextMD(ctx context.Context) context.Context {
	md := FromIncoming(ctx)
	md.Set(strings.ToLower(RequestHeaderKeyUserID), c.UserID)
	md.Set(strings.ToLower(RequestHeaderKeyUserName), c.UserName)
	md.Set(strings.ToLower(RequestHeaderKeyUserEmail), c.UserEmail)
	md.Set(strings.ToLower(RequestHeaderKeyUserSerial), c.UserSerial)
	md.Set(strings.ToLower(RequestHeaderKeyUserType), c.UserType)
	md.Set(strings.ToLower(RequestHeaderKeyCompanyID), c.CompanyID)
	md.Set(strings.ToLower(RequestHeaderKeyCompanySerial), c.CompanySerial)
	md.Set(strings.ToLower(RequestHeaderKeyCompanyName), c.CompanyName)
	md.Set(strings.ToLower(RequestHeaderKeyPermissions), c.Permissions)
	md.Set(strings.ToLower(RequestHeaderKeyClientID), c.ClientID)
	md.Set(strings.ToLower(RequestHeaderKeyClientName), c.ClientName)
	md.Set(strings.ToLower(RequestHeaderKeyScopes), c.Scopes)
	ctx = NewContext(ctx, c)
	return md.ToIncoming(ctx)
}

var reqCtxKey = "ectx"

func NewContext(ctx context.Context, ectx *EContext) context.Context {
	if ectx == nil {
		return ctx
	}
	return context.WithValue(ctx, reqCtxKey, ectx)
}

func FromContext(ctx context.Context) (*EContext, bool) {
	rc, ok := ctx.Value(reqCtxKey).(*EContext)
	return rc, ok
}

func FromContextWithErr(ctx context.Context) (*EContext, error) {
	val, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("failed to get EContext")
	}
	return val, nil
}
