package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanyudii/core-go/ectx"
	"time"
)

func GinRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		eCtx, ok := ectx.FromContext(ctx)
		if !ok {
			eCtx = &ectx.EContext{}
		}

		if eCtx.RequestID == "" {
			requestID := c.GetHeader(ectx.RequestHeaderKeyRequestID)
			if requestID == "" {
				requestID = fmt.Sprintf("%s-%d", uuid.NewString(), time.Now().Unix())
			}
			eCtx.RequestID = requestID
		}

		newCtx := ectx.NewContext(ctx, eCtx)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
