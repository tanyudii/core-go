package auth

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Service interface {
	authenticate(c *gin.Context) (newCtx context.Context, err error)
	authenticateToken(c *gin.Context) (context.Context, error)
}
