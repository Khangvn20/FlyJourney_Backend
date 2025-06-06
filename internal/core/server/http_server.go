package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
	Port   int
}

func New(engine *gin.Engine, port int) *Server {
	return &Server{
		Engine: engine,
		Port:   port,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.Port)
	log.Printf("Starting HTTP server at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.Engine)
}
