package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) ListenAndServer(router *gin.Engine) error {
	s.httpServer = &http.Server{
		Addr:         viper.GetString("server.addr"),
		Handler:      router,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
