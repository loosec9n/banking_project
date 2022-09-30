package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Store provides all fucntionalities to execute db queries and transactions

type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes the fucntion within a databse transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//creatin queries with New Function
	q := New(tx)
	//
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error : %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input details of the transfer transactions
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"` // from the account transfer happens
	ToAccountID   int64 `json:"to_account_id"`   //to the account transfer happens
	Amount        int64 `json:"amount"`          //money to be transfered
}

// TransferTxResult coantians the output of the TransferTx function
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`     //transfer record
	FromAccount Account  `json:"from_account"` // from the account money is debited from
	ToAccount   Account  `josn:"to_account"`   // to the account the money is credited to
	FromEntry   Entry    `json:"from_entry"`   //entry from the money is debited
	ToEntry     Entry    `json:"to_entry"`     //entry to the account the money is credited

}

// TransferTx performs a money transfer from one account to another
// @create a transfer records
// @add account entries
// @update accounts balance
// in single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//transfer the amount
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//add account entries 'from' account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//add account entries 'to' account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//get account -> update its balance

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.FromAccountID,
		// 	Amount: -arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// // to account -> update its balance
		// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })
		if err != nil {
			return err
		}

		return nil
	})

	return result, err

}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account,
	account2 Account,
	err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}

	return
}
