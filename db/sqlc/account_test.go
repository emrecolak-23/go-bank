package db

import (
	"context"
	"testing"
	"time"

	"github.com/emrecolak-23/go-bank/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{

		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	// create account
	account1 := CreateRandomAccount(t)
	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func UpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {

	account1 := CreateRandomAccount(t)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrRecordNotFound)
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	user := createRandomUser(t)
	currencies := []string{"EUR", "USD", "CAD"}

	for _, currency := range currencies {
		arg := CreateAccountParams{
			Owner:    user.Username,
			Balance:  utils.RandomMoney(),
			Currency: currency,
		}

		_, err := testStore.CreateAccount(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListAccountsParams{
		Owner:  user.Username,
		Limit:  2,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	arg2 := ListAccountsParams{
		Owner:  user.Username,
		Limit:  2,
		Offset: 2,
	}

	accounts2, err := testStore.ListAccounts(context.Background(), arg2)
	require.NoError(t, err)
	require.Len(t, accounts2, 1)

	arg3 := ListAccountsParams{
		Owner:  user.Username,
		Limit:  10,
		Offset: 0,
	}

	accounts3, err := testStore.ListAccounts(context.Background(), arg3)
	require.NoError(t, err)
	require.Len(t, accounts3, 3)

	foundCurrencies := make(map[string]bool)
	for _, account := range accounts3 {
		require.Equal(t, user.Username, account.Owner)
		require.False(t, foundCurrencies[account.Currency], "Duplicate currency found")
		foundCurrencies[account.Currency] = true
	}

	for _, currency := range currencies {
		require.True(t, foundCurrencies[currency], "Currency %s not found", currency)
	}
}
