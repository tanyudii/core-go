package middleware

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"net/http"
)

func MuxErrorHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	s := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(s.Code())
	newError := runtime.HTTPStatusError{HTTPStatus: httpStatus, Err: err}
	runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, req, &newError)
}
