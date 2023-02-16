package errorpresenter

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/tanyudii/core-go/errutil"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)
	if err.Extensions == nil {
		err.Extensions = make(map[string]interface{})
	}

	realErr := err.Unwrap()
	customErr, isCustomErr := realErr.(errutil.CustomError)
	if isCustomErr {
		customExtension := errutil.CustomErrorToMapInterface(customErr)
		badRequestErr, isBadRequestErr := realErr.(*errutil.BadRequestError)
		if isBadRequestErr {
			if fields := badRequestErr.GetFields(); len(fields) != 0 {
				customExtension["fields"] = fields
			}
		}
		err.Extensions["errors"] = customExtension
	}

	return err
}
