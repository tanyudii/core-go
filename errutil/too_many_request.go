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
	tooManyRequestGRPCCode = codes.ResourceExhausted
	tooManyRequestHTTPCode = http.StatusTooManyRequests
)

type TooManyRequestError struct {
	code     int
	name     string
	message  string
	grpcCode codes.Code
	httpCode int
}

func (i *TooManyRequestError) Error() string {
	return i.message
}

func (i *TooManyRequestError) GetCode() int {
	return i.code
}

func (i *TooManyRequestError) GetName() string {
	return i.name
}

func (i *TooManyRequestError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *TooManyRequestError) GetHTTPCode() int {
	return i.httpCode
}

func (i *TooManyRequestError) GetErrorInfoCustom() *errdetails.ErrorInfo {
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

func (i *TooManyRequestError) GRPCStatus() *status.Status {
	stats := status.New(i.GetGRPCCode(), i.Error())
	//set error info custom
	if customErr := i.GetErrorInfoCustom(); customErr != nil {
		stats, _ = stats.WithDetails(customErr)
	}
	return stats
}

func NewTooManyRequestError(msg string) error {
	return &TooManyRequestError{
		message:  msg,
		grpcCode: tooManyRequestGRPCCode,
		httpCode: tooManyRequestHTTPCode,
	}
}

func NewTooManyRequestErrorWithCode(msg string, code int) error {
	return &TooManyRequestError{
		code:     code,
		message:  msg,
		grpcCode: tooManyRequestGRPCCode,
		httpCode: tooManyRequestHTTPCode,
	}
}

func NewTooManyRequestErrorWithName(msg string, name string) error {
	return &TooManyRequestError{
		name:     name,
		message:  msg,
		grpcCode: tooManyRequestGRPCCode,
		httpCode: tooManyRequestHTTPCode,
	}
}

func IsTooManyRequestErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == tooManyRequestGRPCCode
}

func IsTooManyRequestError(err error) bool {
	if IsTooManyRequestErrorGRPC(err) {
		return true
	}
	var expectedErr *TooManyRequestError
	return errors.As(err, &expectedErr)
}
