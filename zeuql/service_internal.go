package zeuql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/tanyudii/core-go/graphql/middleware/errorpresenter"
	"github.com/tanyudii/core-go/graphql/middleware/recover"
)

func (s *service) graphQLHandler() gin.HandlerFunc {
	srv := handler.NewDefaultServer(s.schema)
	srv.SetErrorPresenter(errorpresenter.ErrorPresenter)
	srv.SetRecoverFunc(recover.Recover)
	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *service) playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", s.cfg.graphQLPath)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *service) initHealthCheck(r *gin.Engine) {
	r.GET("/_health", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
	})
}
