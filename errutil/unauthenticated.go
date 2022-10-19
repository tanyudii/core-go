package errutil

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	unauthenticatedGRPCCode = codes.Unauthenticated
	unauthenticatedHTTPCode = http.StatusUnauthorized
)

type UnauthenticatedError struct {
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *UnauthenticatedError) Error() string {
	return i.message
}

func (i *UnauthenticatedError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *UnauthenticatedError) GetHTTPCode() int {
	return i.httpCode
}

func (i *UnauthenticatedError) GRPCStatus() *status.Status {
	return status.New(i.GetGRPCCode(), i.Error())
}

func NewUnauthenticatedError(message string) error {
	return &UnauthenticatedError{
		message:  message,
		grpcCode: unauthenticatedGRPCCode,
		httpCode: unauthenticatedHTTPCode,
	}
}

func IsUnauthenticatedErrorGRPC(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == unauthenticatedGRPCCode
}

func IsUnauthenticatedError(err error) bool {
	if IsUnauthenticatedErrorGRPC(err) {
		return true
	}
	var expectedErr *UnauthenticatedError
	return errors.As(err, &expectedErr)
}
