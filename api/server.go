package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/yilinyo/project_bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Account routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.PUT("/accounts/:id", server.updateAccount)

	// Transfer routes
	router.POST("/transfers", server.createTransfer)
	router.GET("/transfers/:id", server.getTransfer)
	router.GET("/transfers", server.listTransfer)
	router.DELETE("/transfers/:id", server.deleteTransfer)
	router.PUT("/transfers/:id", server.updateTransfer)
	router.GET("/transfers/from/:account_id", server.getTransferByFromAccountId)
	router.GET("/transfers/to/:account_id", server.getTransferByToAccountId)

	// Entry routes
	router.POST("/entries", server.createEntry)
	router.GET("/entries/:id", server.getEntry)
	router.GET("/entries", server.listEntry)
	router.DELETE("/entries/:id", server.deleteEntry)
	router.PUT("/entries/:id", server.updateEntry)
	router.GET("/entries/account/:account_id", server.getEntryByAccountId)

	server.router = router
	return server
}

//Start runs

func (server *Server) Start(address string) error {
	return server.router.Run(address)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
