package zeuql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func (s *service) initEngine() {
	gin.SetMode(gin.ReleaseMode)
	s.engine = gin.New()
}

func (s *service) graphQLHandler() gin.HandlerFunc {
	srv := handler.NewDefaultServer(s.schema)
	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *service) playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *service) initHealthCheck() {
	s.engine.GET("/_health", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
	})
}
