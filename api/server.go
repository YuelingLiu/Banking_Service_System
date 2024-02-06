package api

import (
	db "git/db/sqlc"

	"github.com/gin-gonic/gin"
)

// server serves http requests for banking system
type Server struct {
	store *db.Store 
	router *gin.Engine
}

// NewServer creat a new http server and setup routin 
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/account",server.createAccount)

	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)


	server.router = router
	return server
}

// Start runs the http server on a specific address 
func (server *Server) Start(address string) error{
	return server.router.Run(address)

}

func errorResponse(err error) gin.H {
	 return gin.H{"error":err.Error()}
}