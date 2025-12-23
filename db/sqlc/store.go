package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

//type TransferTxParams struct {
//	FromAccountID int64 `json:"from_account_id"`
//	ToAccountID   int64 `json:"to_account_id"`
//	// must be positive
//	Amount int64 `json:"amount"`
//}
//
//type TransferTxResult struct {
//	Transfer    Transfer `json:"transfer"`
//	FromAccount Account  `json:"from_account"`
//	ToAccount   Account  `json:"to_account"`
//	FromEntry   Entry    `json:"from_entry"`
//	ToEntry     Entry    `json:"to_entry"`
//}

var txKey = struct {
}{}

//func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
//	var result TransferTxResult
//	err := store.execTx(ctx, func(q *Queries) error {
//		var err error
//
//		//txName := ctx.Value(txKey)
//		//fmt.Println(txName, "Create Transfer")
//		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
//		if err != nil {
//			return err
//		}
//		//fmt.Println(txName, "Create Entry1")
//		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
//			AccountID: arg.FromAccountID,
//			Amount:    -arg.Amount,
//		})
//		if err != nil {
//			return err
//		}
//		//fmt.Println(txName, "Create Entry2")
//		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
//			AccountID: arg.ToAccountID,
//			Amount:    arg.Amount,
//		})
//		if err != nil {
//			return err
//		}
//		//fmt.Println(txName, "get account1")
//		//account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
//		//if err != nil {
//		//	return err
//		//}
//		//fmt.Println(txName, "update account1")
//		//result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
//		//	ID:      arg.FromAccountID,
//		//	Balance: account1.Balance - arg.Amount,
//		//})
//		if arg.FromAccountID < arg.ToAccountID {
//			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
//				ID:     arg.FromAccountID,
//				Amount: -arg.Amount,
//			})
//			if err != nil {
//				return err
//			}
//			//fmt.Println(txName, "get account2")
//			//account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
//			//if err != nil {
//			//	return err
//			//}
//			//fmt.Println(txName, "update account2")
//			//result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
//			//	ID:      arg.ToAccountID,
//			//	Balance: account2.Balance + arg.Amount,
//			//})
//			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
//				ID:     arg.ToAccountID,
//				Amount: arg.Amount,
//			})
//			if err != nil {
//				return err
//			}
//		} else {
//			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
//				ID:     arg.ToAccountID,
//				Amount: arg.Amount,
//			})
//			if err != nil {
//				return err
//			}
//			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
//				ID:     arg.FromAccountID,
//				Amount: -arg.Amount,
//			})
//			if err != nil {
//				return err
//			}
//
//		}
//		return nil
//	})
//	return result, err
//}
