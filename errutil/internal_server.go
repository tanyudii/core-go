package errutil

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	internalServerGRPCCode = codes.Internal
	internalServerHTTPCode = http.StatusInternalServerError
)

type InternalServerError struct {
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *InternalServerError) Error() string {
	return i.message
}

func (i *InternalServerError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *InternalServerError) GetHTTPCode() int {
	return i.httpCode
}

func (i *InternalServerError) GRPCStatus() *status.Status {
	return status.New(i.GetGRPCCode(), i.Error())
}

func NewInternalServerError(msg string) error {
	return &InternalServerError{
		message:  msg,
		grpcCode: internalServerGRPCCode,
		httpCode: internalServerHTTPCode,
	}
}

func IsInternalServerErrorGRPC(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == internalServerGRPCCode
}

func IsInternalServerError(err error) bool {
	if IsInternalServerErrorGRPC(err) {
		return true
	}
	var expectedErr *InternalServerError
	return errors.As(err, &expectedErr)
}
