package errutil

import "google.golang.org/grpc/codes"

type CustomError interface {
	Error() string
	GetCode() int
	GetName() string
	GetGRPCCode() codes.Code
	GetHTTPCode() int
}

func CustomErrorToMapInterface(err CustomError) map[string]interface{} {
	obj := make(map[string]interface{})
	if code := err.GetCode(); code != 0 {
		obj["code"] = code
	}
	if name := err.GetName(); name != "" {
		obj["name"] = name
	}
	return obj
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
