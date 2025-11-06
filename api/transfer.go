package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/yilinyo/project_bank/db/sqlc"
)

type createTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
}

type getTransferRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type deleteTransferRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type updateTransferRequest struct {
	Id     int64 `uri:"id" binding:"required"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

type listTransferRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

type getTransferByAccountIdRequest struct {
	AccountId int64 `uri:"account_id" binding:"required"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.FromAccountID == req.ToAccountID {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("from_account_id and to_account_id cannot be the same")))
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.CreateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) getTransfer(ctx *gin.Context) {
	var req getTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.store.GetTransfer(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) listTransfer(ctx *gin.Context) {
	var req listTransferRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListTransferParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	transfers, err := server.store.ListTransfer(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfers)
}

func (server *Server) deleteTransfer(ctx *gin.Context) {
	var req deleteTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteTransfer(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "delete success")
}

func (server *Server) updateTransfer(ctx *gin.Context) {
	var req updateTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err2 := ctx.ShouldBindJSON(&req); err2 != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err2))
		return
	}

	if req.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("amount must be positive")))
		return
	}

	args := db.UpdateTransferParams{
		ID:     req.Id,
		Amount: req.Amount,
	}

	transfer, err := server.store.UpdateTransfer(ctx, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) getTransferByFromAccountId(ctx *gin.Context) {
	var req getTransferByAccountIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfers, err := server.store.GetTransferByFromAccountId(ctx, req.AccountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfers)
}

func (server *Server) getTransferByToAccountId(ctx *gin.Context) {
	var req getTransferByAccountIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfers, err := server.store.GetTransferByToAccountId(ctx, req.AccountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfers)
}

