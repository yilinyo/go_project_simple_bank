package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/token"
	"github.com/yilinyo/project_bank/util"
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

	//User routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	//注册拦截中间件
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouter.GET("/users/:username", server.getUser)
	// Account routes
	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.listAccount)
	authRouter.DELETE("/accounts/:id", server.deleteAccount)
	authRouter.PUT("/accounts/:id", server.updateAccount)

	// Transfer routes
	authRouter.POST("/transfers", server.createTransfer)
	authRouter.GET("/transfers/:id", server.getTransfer)
	authRouter.GET("/transfers", server.listTransfer)
	authRouter.DELETE("/transfers/:id", server.deleteTransfer)
	authRouter.PUT("/transfers/:id", server.updateTransfer)
	authRouter.GET("/transfers/from/:account_id", server.getTransferByFromAccountId)
	authRouter.GET("/transfers/to/:account_id", server.getTransferByToAccountId)

	// Entry routes
	authRouter.POST("/entries", server.createEntry)
	authRouter.GET("/entries/:id", server.getEntry)
	authRouter.GET("/entries", server.listEntry)
	authRouter.DELETE("/entries/:id", server.deleteEntry)
	authRouter.PUT("/entries/:id", server.updateEntry)
	authRouter.GET("/entries/account/:account_id", server.getEntryByAccountId)

	server.router = router

}

//Start runs

func (server *Server) Start(address string) error {
	return server.router.Run(address)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
