package rest

import (
	"fmt"
	"github.com/gabriel-98/bingo-backend/internal/application"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RestServer
type RestServer struct {
	port int
	router *gin.Engine
	serviceGroup *application.ServiceGroup
	db *gorm.DB
}

func NewServer(port int, serviceGroup *application.ServiceGroup, db *gorm.DB) *RestServer {
	router := gin.New()
	server := &RestServer{
		port: port,
		router: router,
		serviceGroup: serviceGroup,
		db: db,
	}
	server.loadMiddlewares()
	server.loadEndPoints()
	return server
}

func (server *RestServer) Run() {
	server.router.Run(fmt.Sprintf(":%d", server.port))
}

func (server *RestServer) loadEndPoints() {
	server.loadAuthenticationEndPoints()
}

func (server *RestServer) loadAuthenticationEndPoints() {
	server.router.POST("/auth/signup", server.SignupEndPoint)
	server.router.POST("/auth/login", server.LoginEndPoint)
	server.router.POST("/auth/logout", server.LogoutEndPoint)
	server.router.POST("/auth/refresh-token", server.RefreshTokenEndPoint)
}

func (server *RestServer) loadMiddlewares() {
	server.router.Use(server.QueryExecutorMiddleware)
}

func (server *RestServer) QueryExecutorMiddleware(c *gin.Context) {
	// Set QueryExecutor in the context.
	c.Set("QueryExecutor", server.db)

	c.Next()
}