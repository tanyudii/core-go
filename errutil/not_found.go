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
	notFoundGRPCCode = codes.NotFound
	notFoundHTTPCode = http.StatusNotFound
)

type NotFoundError struct {
	code     int
	name     string
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *NotFoundError) Error() string {
	return i.message
}

func (i *NotFoundError) GetCode() int {
	return i.code
}

func (i *NotFoundError) GetName() string {
	return i.name
}

func (i *NotFoundError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *NotFoundError) GetHTTPCode() int {
	return i.httpCode
}

func (i *NotFoundError) GetErrorInfoCustom() *errdetails.ErrorInfo {
	metaData := make(map[string]string)

	//set error code
	if code := i.GetCode(); code != 0 {
		metaData[metaKeyErrorCode] = strconv.Itoa(code)
	}

	//set error name
	if name := i.GetName(); name != "" {
		metaData[metaKeyErrorName] = name
	}

	return &errdetails.ErrorInfo{
		Metadata: metaData,
	}
}

func (i *NotFoundError) GRPCStatus() *status.Status {
	stats := status.New(i.GetGRPCCode(), i.Error())
	//set error info custom
	if customErr := i.GetErrorInfoCustom(); customErr != nil {
		stats, _ = stats.WithDetails(customErr)
	}
	return stats
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{
		message:  msg,
		grpcCode: notFoundGRPCCode,
		httpCode: notFoundHTTPCode,
	}
}

func NewNotFoundErrorWithCode(msg string, code int) error {
	return &NotFoundError{
		code:     code,
		message:  msg,
		grpcCode: notFoundGRPCCode,
		httpCode: notFoundHTTPCode,
	}
}

func NewNotFoundErrorWithName(msg string, name string) error {
	return &NotFoundError{
		name:     name,
		message:  msg,
		grpcCode: notFoundGRPCCode,
		httpCode: notFoundHTTPCode,
	}
}

func IsNotFoundErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == notFoundGRPCCode
}

func IsNotFoundError(err error) bool {
	if IsNotFoundErrorGRPC(err) {
		return true
	}
	var expectedErr *NotFoundError
	return errors.As(err, &expectedErr)
}
