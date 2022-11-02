package errutil

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

const (
	metaKeyErrorName = "name"
	metaKeyErrorCode = "code"
)

type CustomError interface {
	Error() string
	GetGRPCCode() codes.Code
	GetHTTPCode() int
}

type ErrorField map[string]string

func (f ErrorField) GetFirstErrorAndOtherTotal() (string, int) {
	total := len(f)
	if total > 0 {
		total--
	}
	for k := range f {
		return f[k], total
	}
	return "", 0
}

func GetErrorDetailsFromErrorGRPC(err error) []interface{} {
	e, ok := status.FromError(err)
	if !ok {
		return nil
	}
	return e.Details()
}

func GetErrorGRPCCodeFromErrorGRPC(err error) codes.Code {
	e, ok := status.FromError(err)
	if !ok {
		return codes.Unknown
	}
	return e.Code()
}

func IsGRPCCustomErrorCode(err error, code int) bool {
	details := GetErrorDetailsFromErrorGRPC(err)
	if len(details) == 0 {
		return false
	}
	codeStr := strconv.Itoa(code)
	for _, detail := range details {
		errInfo, valid := detail.(*errdetails.ErrorInfo)
		if !valid {
			continue
		}
		if errInfo.Metadata[metaKeyErrorCode] == codeStr {
			return true
		}
	}
	return false
}

func IsGrpcCustomErrorName(err error, name string) bool {
	details := GetErrorDetailsFromErrorGRPC(err)
	if len(details) == 0 {
		return false
	}
	for _, detail := range details {
		errInfo, valid := detail.(*errdetails.ErrorInfo)
		if !valid {
			continue
		}
		if errInfo.Metadata[metaKeyErrorName] == name {
			return true
		}
	}
	return false
}
