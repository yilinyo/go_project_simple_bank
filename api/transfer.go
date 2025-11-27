package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/token"
)

const (
	FROM_ID = "from_id"
	TO_ID   = "to_id"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
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
	//交易账户必须是 相同货币

	//转账用户要核实当前登录人
	_, valid := server.validAccount(ctx, req.FromAccountID, req.Currency, FROM_ID)
	if !valid {
		return
	}
	//被账用户不做限制
	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency, TO_ID)
	if !valid {
		return
	}
	//todo： 可以进行跨货币交易

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)
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
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string, accountType string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if accountType == FROM_ID {
		payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		if account.Owner != payload.Username {
			err := errors.New("from account does not belong the authenticated user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return db.Account{}, false
		}

	}
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
