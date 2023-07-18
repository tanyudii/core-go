package errutil

import (
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

const (
	metaKeyErrorName = "name"
	metaKeyErrorCode = "code"
)

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

func IsErrorCode(err error, code int) bool {
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

func IsErrorName(err error, name string) bool {
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

func Wrap(prevErr, err error) error {
	return fmt.Errorf("%w: %w", err, prevErr)
}
