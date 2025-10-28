package token

import (
	"testing"
	"time"

	"github.com/emrecolak-23/go-bank/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {

	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)

}
