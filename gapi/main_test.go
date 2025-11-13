package gapi

import (
	"testing"
	"time"

	db "github.com/emrecolak-23/go-bank/db/sqlc"
	"github.com/emrecolak-23/go-bank/utils"
	"github.com/emrecolak-23/go-bank/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}
