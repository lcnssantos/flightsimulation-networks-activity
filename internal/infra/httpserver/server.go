package httpserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port       int
	gin        *gin.Engine
	controller Controller
}

func NewServer(port int, controller Controller) Server {
	return Server{port: port, controller: controller}
}

func (s *Server) Listen() {
	s.gin = gin.Default()

	s.startRouter()

	s.gin.Run(fmt.Sprintf(":%d", s.port))
}
