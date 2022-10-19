package errutil

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	notFoundGRPCCode = codes.NotFound
	notFoundHTTPCode = http.StatusNotFound
)

type NotFoundError struct {
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *NotFoundError) Error() string {
	return i.message
}

func (i *NotFoundError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *NotFoundError) GetHTTPCode() int {
	return i.httpCode
}

func (i *NotFoundError) GRPCStatus() *status.Status {
	return status.New(i.GetGRPCCode(), i.Error())
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{
		message:  msg,
		grpcCode: notFoundGRPCCode,
		httpCode: notFoundHTTPCode,
	}
}

func IsNotFoundErrorGRPC(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == notFoundGRPCCode
}

func IsNotFoundError(err error) bool {
	if IsNotFoundErrorGRPC(err) {
		return true
	}
	var expectedErr *NotFoundError
	return errors.As(err, &expectedErr)
}
