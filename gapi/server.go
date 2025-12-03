package gapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/pb"
	"github.com/yilinyo/project_bank/token"
	"github.com/yilinyo/project_bank/util"
)

// Server pb.UnimplementedSimpleBankServer用于向后兼容新的rpc方法
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
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

	return server, nil
}

//Start runs

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
