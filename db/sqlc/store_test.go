package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)

	//run a concurrent transfer transaction using go routene and channel
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//checking the result
	for i := 0; i < n; i++ {
		//checking the error
		err := <-errs
		require.NoError(t, err)

		//checking the result
		result := <-results
		require.NotEmpty(t, result)

		//checking the transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accountFrom.ID, transfer.FromAccountID)
		require.Equal(t, accountTo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check account 'from' entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accountFrom.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//check amount 'to' entries
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accountTo.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check account balance
	}
}
