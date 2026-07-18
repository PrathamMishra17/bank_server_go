package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateEntryDummy(t *testing.T, account Account) Entry {
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	}

	result, err := testqueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)

	lastInsertedId, err := result.LastInsertId()
	require.NoError(t, err)

	fetchEntry, err := testqueries.GetEntry(context.Background(), int32(lastInsertedId))
	require.Equal(t, args.AccountID, account.ID)
	require.Equal(t, args.Amount, account.Balance)
	return fetchEntry
}

func TestCreateAccount(t *testing.T) {
	account := createRandomAccount(t)
	CreateEntryDummy(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entryCreated := CreateEntryDummy(t, account)

	entry, err := testqueries.GetEntry(context.Background(), entryCreated.ID)

	require.NoError(t, err)

	require.NotEmpty(t, entry)

	require.Equal(t, entry.ID, entryCreated.ID)
	require.Equal(t, entry.Amount, entryCreated.Amount)
	require.Equal(t, entry.AccountID, entryCreated.AccountID)
	require.Equal(t, entry.CreatedAt, entryCreated.CreatedAt, time.Second)

}

func TestListEntry(t *testing.T) {
	account := createRandomAccount(t)
	//first creating 10 entries
	var lastentry Entry
	for i := 0; i < 10; i++ {
		lastentry = CreateEntryDummy(t, account)
	}

	// fetching the entries

	args := ListEntriesParams{
		AccountID: lastentry.AccountID,
		Limit:     5,
		Offset:    0,
	}

	entries, err := testqueries.ListEntries(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entries)

	// comparing each of the entry
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, lastentry.AccountID)
		require.Equal(t, entry.Amount, lastentry.Amount)
		require.Equal(t, entry.CreatedAt, lastentry.CreatedAt, time.Second)
	}
}
