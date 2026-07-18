package db

import (
	"context"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(TestDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// for individual testing run on same go routine but for multiple testing run
	//multiple go routines
	// say testing 5 transactions

	n := int(5)
	amount := int32(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: int64(account1.ID),
				ToAccountId:   int64(account2.ID),
				Amount:        int64(amount),
			})
			errs <- err
			results <- result
		}()

	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		transfer := result.Transfer

		//check the transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotEmpty(t, transfer.ID)
		require.NotEmpty(t, transfer.CreatedAt)

		//check the transfer record
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check the entries
		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotEmpty(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.CreatedAt)

		//check the transfer record
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, +amount, toEntry.Amount)
		require.NotEmpty(t, toEntry.ID)
		require.NotEmpty(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO test the updated balance
	}

}
