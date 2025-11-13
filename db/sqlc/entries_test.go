package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateEntry(t *testing.T) Entry {
	account := CreateRandomAccount(t)

	arg := CreateEntriesParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	}

	entry, err := testStore.CreateEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, entry.Amount, account.Balance)
	require.Equal(t, entry.AccountID, account.ID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateEntry(t)

	entry2, err := testStore.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	entry1 := CreateEntry(t)

	err := testStore.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testStore.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrRecordNotFound)
	require.Empty(t, entry2)

}

func TestListEntries(t *testing.T) {

	for i := 0; i < 10; i++ {
		CreateEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testStore.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

}

func TestUpdateEntry(t *testing.T) {

	entry1 := CreateEntry(t)

	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: 10,
	}

	entry2, err := testStore.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, entry1.ID, entry2.ID)
	require.NotEqual(t, entry1.Amount, entry2.Amount)

}
