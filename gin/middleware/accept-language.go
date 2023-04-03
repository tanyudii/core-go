package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tanyudii/core-go/constant"
	"github.com/tanyudii/core-go/ectx"
)

func GinAcceptLanguage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		eCtx, ok := ectx.FromContext(ctx)
		if !ok {
			eCtx = &ectx.EContext{}
		}

		if eCtx.AcceptLanguage == "" {
			eCtx.AcceptLanguage = c.GetHeader(constant.KeyAcceptLanguage)
		}

		newCtx := ectx.NewContext(ctx, eCtx)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
