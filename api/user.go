package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/db/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	FullName string `json:"full_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
}
type createUserResponse struct {
	Username string `json:"username" binding:"required,alphanum"`
	FullName string `json:"full_name" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
}
type getUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
}

func newUserResponse(user db.User) createUserResponse {
	return createUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
}

func (server *Server) createUser(ctx *gin.Context) {

	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashpwd, err2 := util.HashPassword(req.Password)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf(""))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashpwd,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(pqErr))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)
}
func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)

}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type loginUserResponse struct {
	AccessToken string             `json:"access_token" binding:"required,alphanum"`
	User        createUserResponse `json:"user" binding:"required"`
}

func (server *Server) loginUser(ctx *gin.Context) {

	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !util.CheckPassword(req.Password, user.HashedPassword) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("wrong password and username")))
		return
	}
	token, _, err := server.tokenMaker.CreateToken(user.Username, "admin", server.config.AccessTokenDuration, 1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := loginUserResponse{
		User:        newUserResponse(user),
		AccessToken: token,
	}

	ctx.JSON(http.StatusOK, resp)
	return

}
