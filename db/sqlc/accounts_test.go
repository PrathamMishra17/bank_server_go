// This module actually tests the accounts database functions like create read update delete

package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/PrathamMishra17/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  int32(utils.RandomMoney()),
		Currency: utils.RandomCurrency(),
	}

	account, err := testqueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)

	require.NotEmpty(t, account)

	// MySql does not return the account parameters directly like postgress it returns the lastroweffected

	rowAffected, err := account.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rowAffected)

	// getting the last inserted id
	lastInsertId, err := account.LastInsertId()
	require.NoError(t, err)

	//fetching the account
	fetchedAccount, err := testqueries.GetAccount(context.Background(), int32(lastInsertId))
	require.NoError(t, err)
	require.NotEmpty(t, fetchedAccount)

	return fetchedAccount

}

func TestAccountCreate(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	account2, err := testqueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Createdat, account2.Createdat, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	args := UpdateAccountParams{
		Balance: int32(utils.RandomMoney()),
		ID:      account1.ID,
	}

	account2, err := testqueries.UpdateAccount(context.Background(), args)

	rowsAffected, err := account2.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)

	fetchedAccount, err := testqueries.GetAccount(context.Background(), args.ID)
	require.NoError(t, err)
	require.Equal(t, args.Balance, fetchedAccount.Balance)
	require.Equal(t, account1.Owner, fetchedAccount.Owner)

	require.NoError(t, err)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testqueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	fetchAccount, err := testqueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)                             //we need error in fetching
	require.EqualError(t, err, sql.ErrNoRows.Error()) //confirm the error is due to no rows returned
	require.Empty(t, fetchAccount)
}

func TestAccountList(t *testing.T) {
	var lastaccount Account

	for i := 0; i < 10; i++ {
		lastaccount = createRandomAccount(t)
	}

	args := ListAccountsParams{
		Owner:  lastaccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testqueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, account.Owner, lastaccount.Owner)
	}

}
