package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTransferTx for testing the functionality
func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)

	//run a concurrent transfer transaction using go routene and channel
	n := 5
	amount := int64(10)
	existed := make(map[int]bool)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() { //5 go routenes for creating concurrent transactions
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

		//check 'from' account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountFrom.ID, fromAccount.ID)

		//check 'to' account
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountTo.ID, toAccount.ID)

		//check account balance

		diff1 := accountFrom.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - accountTo.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check the final account balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)

	require.Equal(t, accountFrom.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, accountTo.Balance+int64(n)*amount, updateAccount2.Balance)
}

// TestTransferTxDeadlock will simulate the to and fro transaction for a deadlock
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)

	//run a concurrent transfer transaction using go routene and channel
	n := 10
	amount := int64(10)

	errs := make(chan error)
	//run concurrent transfer transactions
	for i := 0; i < n; i++ {
		fromAccountID := accountFrom.ID
		toAccountID := accountTo.ID

		if i%2 == 1 {
			fromAccountID = accountTo.ID
			toAccountID = accountFrom.ID
		}

		go func() { //5 go routenes for creating concurrent transactions
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	//checking the result

	for i := 0; i < n; i++ {
		//checking the error
		err := <-errs
		require.NoError(t, err)
	}

	//check the final account balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)

	require.Equal(t, accountFrom.Balance, updateAccount1.Balance)
	require.Equal(t, accountTo.Balance, updateAccount2.Balance)
}
