package errutil

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	permissionDeniedGRPCCode = codes.PermissionDenied
	permissionDeniedHTTPCode = http.StatusUnauthorized
)

type UnauthorizedError struct {
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *UnauthorizedError) Error() string {
	return i.message
}

func (i *UnauthorizedError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *UnauthorizedError) GetHTTPCode() int {
	return i.httpCode
}

func (i *UnauthorizedError) GRPCStatus() *status.Status {
	return status.New(i.GetGRPCCode(), i.Error())
}

func NewUnauthorizedError(msg string) error {
	return &UnauthorizedError{
		message:  msg,
		grpcCode: permissionDeniedGRPCCode,
		httpCode: permissionDeniedHTTPCode,
	}
}

func IsUnauthorizedErrorGRPC(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == permissionDeniedGRPCCode
}

func IsUnauthorizedError(err error) bool {
	if IsUnauthenticatedErrorGRPC(err) {
		return true
	}
	var expectedErr *UnauthorizedError
	return errors.As(err, &expectedErr)
}
