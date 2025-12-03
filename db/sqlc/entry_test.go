package db

import (
	"context"

	"github.com/stretchr/testify/require"
	"github.com/yilinyo/project_bank/util"

	"testing"
	"time"
)

func TestCreateEntry(t *testing.T) {

	CreateTestEntry(t)

}

func CreateTestEntry(t *testing.T) Entry {

	account := createRandomAccount(t)

	arg := CreateEntryParams{

		AccountID: account.ID,
		Amount:    util.RandomEntryMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.CreatedAt)

	return entry

}

func TestGetEntryByID(t *testing.T) {

	entry := CreateTestEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, entry.Amount, entry2.Amount)
	require.Equal(t, entry.AccountID, entry2.AccountID)
	require.WithinDuration(t, entry.CreatedAt, entry2.CreatedAt, time.Second)

}

func TestGetEntryByAccountID(t *testing.T) {

	entry := CreateTestEntry(t)
	entry2, err := testQueries.GetEntryByAccountId(context.Background(), entry.AccountID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	for _, entry := range entry2 {
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.Amount)
		require.NotZero(t, entry.CreatedAt)
		require.Equal(t, entry.AccountID, entry.AccountID)
	}

}

func TestModifyEntry(t *testing.T) {
	entry1 := CreateTestEntry(t)

	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomEntryMoney(),
	}

	entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

}
func TestDeleteEntry(t *testing.T) {
	entry1 := CreateTestEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	//require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, entry2)
}

func TestListEntry(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateTestEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
