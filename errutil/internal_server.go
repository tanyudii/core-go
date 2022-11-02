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
	internalServerGRPCCode = codes.Internal
	internalServerHTTPCode = http.StatusInternalServerError
)

type InternalServerError struct {
	code     int
	name     string
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *InternalServerError) Error() string {
	return i.message
}

func (i *InternalServerError) GetCode() int {
	return i.code
}

func (i *InternalServerError) GetName() string {
	return i.name
}

func (i *InternalServerError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *InternalServerError) GetHTTPCode() int {
	return i.httpCode
}

func (i *InternalServerError) GetErrorInfoCustom() *errdetails.ErrorInfo {
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

func (i *InternalServerError) GRPCStatus() *status.Status {
	stats := status.New(i.GetGRPCCode(), i.Error())
	//set error info custom
	if customErr := i.GetErrorInfoCustom(); customErr != nil {
		stats, _ = stats.WithDetails(customErr)
	}
	return stats
}

func NewInternalServerError(msg string) error {
	return &InternalServerError{
		message:  msg,
		grpcCode: internalServerGRPCCode,
		httpCode: internalServerHTTPCode,
	}
}

func NewInternalServerErrorWithCode(msg string, code int) error {
	return &InternalServerError{
		code:     code,
		message:  msg,
		grpcCode: internalServerGRPCCode,
		httpCode: internalServerHTTPCode,
	}
}

func NewInternalServerErrorWithName(msg string, name string) error {
	return &InternalServerError{
		name:     name,
		message:  msg,
		grpcCode: internalServerGRPCCode,
		httpCode: internalServerHTTPCode,
	}
}

func IsInternalServerErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == internalServerGRPCCode
}

func IsInternalServerError(err error) bool {
	if IsInternalServerErrorGRPC(err) {
		return true
	}
	var expectedErr *InternalServerError
	return errors.As(err, &expectedErr)
}
