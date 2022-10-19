package errutil

import (
	"errors"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	badRequestGRPCCode = codes.InvalidArgument
	badRequestHTTPCode = http.StatusBadRequest
)

type BadRequestError struct {
	message  string
	grpcCode codes.Code
	httpCode int
	fields   ErrorField
}

func (i *BadRequestError) Error() string {
	return i.message
}

func (i *BadRequestError) GetGRPCCode() codes.Code {
	return i.grpcCode
}

func (i *BadRequestError) GetHTTPCode() int {
	return i.httpCode
}

func (i *BadRequestError) GRPCStatus() *status.Status {
	stats := status.New(i.GetGRPCCode(), i.Error())
	errFields := i.GetFields()
	if len(errFields) != 0 {
		br := &errdetails.BadRequest{}
		for attr, msg := range i.GetFields() {
			br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       attr,
				Description: msg,
			})
		}
		stats, _ = stats.WithDetails(br)
	}
	return stats
}

func (i *BadRequestError) GetFields() ErrorField {
	return i.fields
}

func NewBadRequestError(msg string) error {
	return &BadRequestError{
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
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == badRequestGRPCCode
}

func IsBadRequestError(err error) bool {
	if IsBadRequestErrorGRPC(err) {
		return true
	}
	var expectedErr *BadRequestError
	return errors.As(err, &expectedErr)
}
