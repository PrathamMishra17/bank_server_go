package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	existed := make(map[int]bool)

	store := NewStore(TestDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("Balance1<-", account1.Balance)
	fmt.Println("Balance2<-", account2.Balance)

	// for individual testing run on same go routine but for multiple testing run
	//multiple go routines
	// say testing 5 transactions

	n := int(5)
	amount := int32(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i <= n; i++ {
		txname := fmt.Sprintf("txname %v", i)
		go func() {

			ctx := context.WithValue(context.Background(), txkey{}, txname)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: int64(account1.ID),
				ToAccountId:   int64(account2.ID),
				Amount:        int64(amount),
			})
			errs <- err
			results <- result
		}()

	}

	for i := 1; i <= n; i++ {

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
		fmt.Println(transfer.ID, " ", transfer.Amount, " ", transfer.CreatedAt)

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

		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//checks the account balance
		fmt.Println("tx1", account1.Balance)
		fmt.Println("tx2", account2.Balance)
		diff1 := account1.Balance - fromAccount.Balance

		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // diff1 will always be multiple of amount transferred

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}
	// check the updated account balance
	updatedAccount1, err := testqueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testqueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("After1", updatedAccount1.Balance)
	fmt.Println("After2", updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int32(n)*(amount), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int32(n)*(amount), updatedAccount2.Balance)

}

func TestTransferTxBidirectional(t *testing.T) {
	n := (2)
	amount := int32(10)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("Balance1", account1.Balance)
	fmt.Println("Balance2", account2.Balance)

	errors := make(chan error)

	store := NewStore(TestDB)
	//we had done the transaction
	for i := 1; i <= n; i++ {
		txname := fmt.Sprintf("txname %v", i)
		fromAccount := account1.ID
		toAccount := account2.ID
		if i%2 == 1 {
			fromAccount = account2.ID
			toAccount = account1.ID
		}
		go func() {
			ctx := context.WithValue(context.Background(), txname, txkey{})
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: int64(fromAccount),
				ToAccountId:   int64(toAccount),
				Amount:        int64(amount),
			})

			errors <- err

		}()
	}

	for i := 1; i <= n; i++ {

		err := <-errors
		require.NoError(t, err)

	}
	// check the updated account balance
	updatedAccount1, err := testqueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testqueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("After1", updatedAccount1.Balance)
	fmt.Println("After2", updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
