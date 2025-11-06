package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/yilinyo/project_bank/db/sqlc"
)

type createEntryRequest struct {
	AccountID int64 `json:"account_id" binding:"required,min=1"`
	Amount    int64 `json:"amount" binding:"required"`
}

type getEntryRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type deleteEntryRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type updateEntryRequest struct {
	Id     int64 `uri:"id" binding:"required"`
	Amount int64 `json:"amount" binding:"required"`
}

type listEntryRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

type getEntryByAccountIdRequest struct {
	AccountId int64 `uri:"account_id" binding:"required"`
}

func (server *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateEntryParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	}

	entry, err := server.store.CreateEntry(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, entry)
}

func (server *Server) getEntry(ctx *gin.Context) {
	var req getEntryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := server.store.GetEntry(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, entry)
}

func (server *Server) listEntry(ctx *gin.Context) {
	var req listEntryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListEntriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	entries, err := server.store.ListEntries(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, entries)
}

func (server *Server) deleteEntry(ctx *gin.Context) {
	var req deleteEntryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteEntry(ctx, req.Id)
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

func (server *Server) updateEntry(ctx *gin.Context) {
	var req updateEntryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err2 := ctx.ShouldBindJSON(&req); err2 != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err2))
		return
	}

	args := db.UpdateEntryParams{
		ID:     req.Id,
		Amount: req.Amount,
	}

	entry, err := server.store.UpdateEntry(ctx, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, entry)
}

func (server *Server) getEntryByAccountId(ctx *gin.Context) {
	var req getEntryByAccountIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entries, err := server.store.GetEntryByAccountId(ctx, req.AccountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, entries)
}

