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
	unauthorizedGRPCCode = codes.PermissionDenied
	unauthorizedHTTPCode = http.StatusUnauthorized
)

type UnauthorizedError struct {
	code     int
	name     string
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *UnauthorizedError) Error() string {
	return i.message
}

func (i *UnauthorizedError) GetCode() int {
	return i.code
}

func (i *UnauthorizedError) GetName() string {
	return i.name
}

func (i *UnauthorizedError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *UnauthorizedError) GetHTTPCode() int {
	return i.httpCode
}

func (i *UnauthorizedError) GetErrorInfoCustom() *errdetails.ErrorInfo {
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

func (i *UnauthorizedError) GRPCStatus() *status.Status {
	stats := status.New(i.GetGRPCCode(), i.Error())
	//set error info custom
	if customErr := i.GetErrorInfoCustom(); customErr != nil {
		stats, _ = stats.WithDetails(customErr)
	}
	return stats
}

func NewUnauthorizedError(msg string) error {
	return &UnauthorizedError{
		message:  msg,
		grpcCode: unauthorizedGRPCCode,
		httpCode: unauthorizedHTTPCode,
	}
}

func NewUnauthorizedErrorWithCode(msg string, code int) error {
	return &UnauthorizedError{
		code:     code,
		message:  msg,
		grpcCode: unauthorizedGRPCCode,
		httpCode: unauthorizedHTTPCode,
	}
}

func NewUnauthorizedErrorWithName(msg string, name string) error {
	return &UnauthorizedError{
		name:     name,
		message:  msg,
		grpcCode: unauthorizedGRPCCode,
		httpCode: unauthorizedHTTPCode,
	}
}

func IsUnauthorizedErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == unauthorizedGRPCCode
}

func IsUnauthorizedError(err error) bool {
	if IsUnauthenticatedErrorGRPC(err) {
		return true
	}
	var expectedErr *UnauthorizedError
	return errors.As(err, &expectedErr)
}
