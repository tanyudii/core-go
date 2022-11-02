package errutil

import (
	"errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

const (
	unauthenticatedGRPCCode = codes.Unauthenticated
	unauthenticatedHTTPCode = http.StatusUnauthorized
)

type UnauthenticatedError struct {
	code     int
	name     string
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *UnauthenticatedError) Error() string {
	return i.message
}

func (i *UnauthenticatedError) GetCode() int {
	return i.code
}

func (i *UnauthenticatedError) GetName() string {
	return i.name
}

func (i *UnauthenticatedError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *UnauthenticatedError) GetHTTPCode() int {
	return i.httpCode
}

func (i *UnauthenticatedError) GetErrorInfoCustom() *errdetails.ErrorInfo {
	metaData := make(map[string]string)

	//set error code
	if code := i.GetCode(); code != 0 {
		metaData[metaKeyErrorName] = strconv.Itoa(code)
	}

	//set error name
	if name := i.GetName(); name != "" {
		metaData[metaKeyErrorCode] = name
	}

	return &errdetails.ErrorInfo{
		Metadata: metaData,
	}
}

func (i *UnauthenticatedError) GRPCStatus() *status.Status {
	stats := status.New(i.GetGRPCCode(), i.Error())
	//set error info custom
	if customErr := i.GetErrorInfoCustom(); customErr != nil {
		stats, _ = stats.WithDetails(customErr)
	}
	return stats
}

func NewUnauthenticatedError(message string) error {
	return &UnauthenticatedError{
		message:  message,
		grpcCode: unauthenticatedGRPCCode,
		httpCode: unauthenticatedHTTPCode,
	}
}

func NewUnauthenticatedErrorWithCode(msg string, code int) error {
	return &UnauthenticatedError{
		code:     code,
		message:  msg,
		grpcCode: unauthenticatedGRPCCode,
		httpCode: unauthenticatedHTTPCode,
	}
}

func NewUnauthenticatedErrorWithName(msg string, name string) error {
	return &UnauthenticatedError{
		name:     name,
		message:  msg,
		grpcCode: unauthenticatedGRPCCode,
		httpCode: unauthenticatedHTTPCode,
	}
}

func IsUnauthenticatedErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == unauthenticatedGRPCCode
}

func IsUnauthenticatedError(err error) bool {
	if IsUnauthenticatedErrorGRPC(err) {
		return true
	}
	var expectedErr *UnauthenticatedError
	return errors.As(err, &expectedErr)
}
