package api

import (
	"paysyncevets/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config utils.Config
	db     *gorm.DB
	router *gin.Engine
}

func NewServer(config utils.Config, db *gorm.DB) (*Server, error) {
	server := &Server{
		config: config,
		db:     db,
	}
	server.setupRoutes()
	return server, nil
}

// set up routes

func (server *Server) setupRoutes() {
	router := gin.Default()
	
	//public routes

	router.POST("/user/createUser", server.createUser)
	router.POST("/user/login", server.userLogin)

	// private routes

	server.router = router
}

// Start runs the http server on a specific address

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// resuable error function for error response

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}