package errutil

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	tooManyRequestGRPCCode = codes.ResourceExhausted
	tooManyRequestHTTPCode = http.StatusTooManyRequests
)

type TooManyRequestError struct {
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *TooManyRequestError) Error() string {
	return i.message
}

func (i *TooManyRequestError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *TooManyRequestError) GetHTTPCode() int {
	return i.httpCode
}

func (i *TooManyRequestError) GRPCStatus() *status.Status {
	return status.New(i.GetGRPCCode(), i.Error())
}

func NewTooManyRequestError(msg string) error {
	return &TooManyRequestError{
		message:  msg,
		grpcCode: tooManyRequestGRPCCode,
		httpCode: tooManyRequestHTTPCode,
	}
}

func IsTooManyRequestErrorGRPC(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == tooManyRequestGRPCCode
}

func IsTooManyRequestError(err error) bool {
	if IsTooManyRequestErrorGRPC(err) {
		return true
	}
	var expectedErr *TooManyRequestError
	return errors.As(err, &expectedErr)
}
