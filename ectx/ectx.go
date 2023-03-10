package ectx

import (
	"context"
	"errors"
	"github.com/tanyudii/core-go/common"
	"github.com/tanyudii/core-go/errutil"
	"strconv"
	"strings"
)

const (
	ContextKey = "ectx"

	RequestHeaderKeyUserID         = "UserID"
	RequestHeaderKeyUserName       = "UserName"
	RequestHeaderKeyUserEmail      = "UserEmail"
	RequestHeaderKeyUserSerial     = "UserSerial"
	RequestHeaderKeyUserType       = "UserType"
	RequestHeaderKeyCompanyID      = "CompanyID"
	RequestHeaderKeyCompanySerial  = "CompanySerial"
	RequestHeaderKeyCompanyName    = "CompanyName"
	RequestHeaderKeyPermissions    = "Permissions"
	RequestHeaderKeyScopes         = "Scopes"
	RequestHeaderKeyClientID       = "ClientID"
	RequestHeaderKeyClientName     = "ClientName"
	RequestHeaderKeyIsInternalCall = "IsInternalCall"
	RequestHeaderKeyAuthorization  = "Authorization"
	RequestHeaderKeyRequestID      = "RequestID"
)

type EContext struct {
	UserID         string
	UserName       string
	UserEmail      string
	UserSerial     string
	UserType       string
	CompanyID      string
	CompanySerial  string
	CompanyName    string
	Permissions    string //separated by ","
	ClientID       string
	ClientName     string
	Scopes         string //separated by ","
	IsInternalCall bool
	Authorization  string
	RequestID      string
}

func NewEContext(md ContextMD) *EContext {
	isInternalCall, _ := strconv.ParseBool(md.Get(strings.ToLower(RequestHeaderKeyIsInternalCall)))
	return &EContext{
		UserID:         md.Get(strings.ToLower(RequestHeaderKeyUserID)),
		UserName:       md.Get(strings.ToLower(RequestHeaderKeyUserName)),
		UserEmail:      md.Get(strings.ToLower(RequestHeaderKeyUserEmail)),
		UserSerial:     md.Get(strings.ToLower(RequestHeaderKeyUserSerial)),
		UserType:       md.Get(strings.ToLower(RequestHeaderKeyUserType)),
		CompanyID:      md.Get(strings.ToLower(RequestHeaderKeyCompanyID)),
		CompanySerial:  md.Get(strings.ToLower(RequestHeaderKeyCompanySerial)),
		CompanyName:    md.Get(strings.ToLower(RequestHeaderKeyCompanyName)),
		Permissions:    md.Get(strings.ToLower(RequestHeaderKeyPermissions)),
		ClientID:       md.Get(strings.ToLower(RequestHeaderKeyClientID)),
		ClientName:     md.Get(strings.ToLower(RequestHeaderKeyClientName)),
		Scopes:         md.Get(strings.ToLower(RequestHeaderKeyScopes)),
		IsInternalCall: isInternalCall,
		Authorization:  md.Get(strings.ToLower(RequestHeaderKeyAuthorization)),
		RequestID:      md.Get(strings.ToLower(RequestHeaderKeyRequestID)),
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
	md.Set(strings.ToLower(RequestHeaderKeyIsInternalCall), strconv.FormatBool(c.IsInternalCall))
	md.Set(strings.ToLower(RequestHeaderKeyAuthorization), c.Authorization)
	md.Set(strings.ToLower(RequestHeaderKeyRequestID), c.RequestID)
	ctx = NewContext(ctx, c)
	return md.ToIncoming(ctx)
}

func (c *EContext) IsInternal() bool {
	return c.IsInternalCall &&
		c.UserID != "" && c.UserSerial != "" &&
		c.CompanyID != "" && c.CompanySerial != ""
}

func (c *EContext) HasPermission(codes []string) (bool, error) {
	//skip immediately
	if len(codes) == 0 {
		return false, nil
	}
	permissions := common.ParseStringToSliceBySeparator(c.Permissions, ",")
	if len(permissions) != 0 {
		mapCode := make(map[string]bool)
		for _, code := range codes {
			mapCode[code] = true
		}
		for _, p := range permissions {
			if mapCode[p] {
				return true, nil
			}
		}
	}
	return false, errutil.ErrAuthPermissionNotAllowed
}

func (c *EContext) HasScope(codes []string) (bool, error) {
	//skip immediately
	if len(codes) == 0 {
		return false, nil
	}
	scopes := common.ParseStringToSliceBySeparator(c.Scopes, ",")
	if len(scopes) != 0 {
		mapCode := make(map[string]bool)
		for _, code := range codes {
			mapCode[code] = true
		}
		for _, s := range scopes {
			if mapCode[s] {
				return true, nil
			}
		}
	}
	return false, errutil.ErrAuthScopeNotAllowed
}

func (c *EContext) HasUserType(codes []string) (bool, error) {
	//skip immediately
	if len(codes) == 0 {
		return false, nil
	}
	for _, code := range codes {
		if code == c.UserType {
			return true, nil
		}
	}
	return false, errutil.ErrAuthUserTypeNotAllowed
}

func (c *EContext) HasUserTypeByMapCode(codes map[string]bool) (bool, error) {
	//skip immediately
	if len(codes) == 0 {
		return false, nil
	}
	if c.UserType != "" && codes[c.UserType] {
		return true, nil
	}
	return false, errutil.ErrAuthUserTypeNotAllowed
}

func NewContext(ctx context.Context, eCtx *EContext) context.Context {
	if eCtx == nil {
		return ctx
	}
	return context.WithValue(ctx, ContextKey, eCtx)
}

func FromContext(ctx context.Context) (*EContext, bool) {
	rc, ok := ctx.Value(ContextKey).(*EContext)
	return rc, ok
}

func FromContextWithErr(ctx context.Context) (*EContext, error) {
	val, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("failed to get eCtx")
	}
	return val, nil
}

func ParseToGrpcCtx(ctx context.Context, internalCall ...bool) context.Context {
	if r, ok := FromContext(ctx); ok {
		isInternalCall := true
		if len(internalCall) > 0 {
			isInternalCall = internalCall[0]
		}
		newCtx := FromIncoming(ctx)
		newCtx.Add(strings.ToLower(RequestHeaderKeyUserID), r.UserID)
		newCtx.Add(strings.ToLower(RequestHeaderKeyUserName), r.UserName)
		newCtx.Add(strings.ToLower(RequestHeaderKeyUserEmail), r.UserEmail)
		newCtx.Add(strings.ToLower(RequestHeaderKeyUserSerial), r.UserSerial)
		newCtx.Add(strings.ToLower(RequestHeaderKeyUserType), r.UserType)
		newCtx.Add(strings.ToLower(RequestHeaderKeyCompanyID), r.CompanyID)
		newCtx.Add(strings.ToLower(RequestHeaderKeyCompanySerial), r.CompanySerial)
		newCtx.Add(strings.ToLower(RequestHeaderKeyCompanyName), r.CompanyName)
		newCtx.Add(strings.ToLower(RequestHeaderKeyPermissions), r.Permissions)
		newCtx.Add(strings.ToLower(RequestHeaderKeyClientID), r.ClientID)
		newCtx.Add(strings.ToLower(RequestHeaderKeyClientName), r.ClientName)
		newCtx.Add(strings.ToLower(RequestHeaderKeyScopes), r.Scopes)
		newCtx.Add(strings.ToLower(RequestHeaderKeyIsInternalCall), strconv.FormatBool(isInternalCall))
		newCtx.Add(strings.ToLower(RequestHeaderKeyAuthorization), r.Authorization)
		newCtx.Add(strings.ToLower(RequestHeaderKeyRequestID), r.RequestID)
		return newCtx.ToOutgoing(ctx)
	}
	return ctx
}
