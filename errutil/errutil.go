package errutil

import (
	"google.golang.org/grpc/codes"
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
