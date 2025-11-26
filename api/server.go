package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/db/util"
	"github.com/yilinyo/project_bank/token"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("error creating token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	//register validator 处理注册货币验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validateCurrency)
		if err != nil {
			return nil, fmt.Errorf("error registering validation: %w", err)
		}
	}
	server.SetupRouter()
	return server, nil
}

func (server *Server) SetupRouter() {
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

	//User routes
	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)
	router.POST("/users/login", server.loginUser)
	server.router = router

}

//Start runs

func (server *Server) Start(address string) error {
	return server.router.Run(address)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
