package errutil

import (
	"errors"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

const (
	badRequestGRPCCode = codes.InvalidArgument
	badRequestHTTPCode = http.StatusBadRequest
)

type BadRequestError struct {
	code     int
	name     string
	message  string
	grpcCode codes.Code
	httpCode int
	fields   ErrorField
}

func (i *BadRequestError) Error() string {
	return i.message
}

func (i *BadRequestError) GetCode() int {
	return i.code
}

func (i *BadRequestError) GetName() string {
	return i.name
}

func (i *BadRequestError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *BadRequestError) GetHTTPCode() int {
	return i.httpCode
}

func (i *BadRequestError) GetFields() ErrorField {
	return i.fields
}

func (i *BadRequestError) GetBadRequestFields() *errdetails.BadRequest {
	errFields := i.GetFields()
	if len(errFields) == 0 {
		return nil
	}
	br := &errdetails.BadRequest{}
	for attr, msg := range i.GetFields() {
		br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       attr,
			Description: msg,
		})
	}
	return br
}

func (i *BadRequestError) GetErrorInfoCustom() *errdetails.ErrorInfo {
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

func (i *BadRequestError) GRPCStatus() *status.Status {
	stats := status.New(i.GetGRPCCode(), i.Error())

	//set error fields
	if fields := i.GetBadRequestFields(); fields != nil {
		stats, _ = stats.WithDetails(fields)
	}

	//set error info custom
	if customErr := i.GetErrorInfoCustom(); customErr != nil {
		stats, _ = stats.WithDetails(customErr)
	}

	return stats
}

func NewBadRequestError(msg string) error {
	return &BadRequestError{
		message:  msg,
		grpcCode: badRequestGRPCCode,
		httpCode: badRequestHTTPCode,
	}
}

func NewBadRequestErrorWithCode(msg string, code int) error {
	return &BadRequestError{
		code:     code,
		message:  msg,
		grpcCode: badRequestGRPCCode,
		httpCode: badRequestHTTPCode,
	}
}

func NewBadRequestErrorWithName(msg string, name string) error {
	return &BadRequestError{
		name:     name,
		message:  msg,
		grpcCode: badRequestGRPCCode,
		httpCode: badRequestHTTPCode,
	}
}

func NewBadRequestErrorWithFields(msg string, fields ErrorField) error {
	return &BadRequestError{
		message:  msg,
		fields:   fields,
		grpcCode: badRequestGRPCCode,
		httpCode: badRequestHTTPCode,
	}
}

func NewBadRequestErrorUsingFieldsOrNil(fields ErrorField) error {
	if len(fields) != 0 {
		firstErr, otherErr := fields.GetFirstErrorAndOtherTotal()
		wordErr := "error"
		if otherErr == 0 {
			return NewBadRequestErrorWithFields(firstErr, fields)
		}
		if otherErr > 1 {
			wordErr = "errors"
		}
		return NewBadRequestErrorWithFields(fmt.Sprintf("%s. and there are %d %s", firstErr, otherErr, wordErr), fields)
	}
	return nil
}

func IsBadRequestErrorGRPC(err error) bool {
	return GetErrorGRPCCodeFromErrorGRPC(err) == badRequestGRPCCode
}

func IsBadRequestError(err error) bool {
	if IsBadRequestErrorGRPC(err) {
		return true
	}
	var expectedErr *BadRequestError
	return errors.As(err, &expectedErr)
}
