package gapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/yilinyo/project_bank/db/sqlc"
	util2 "github.com/yilinyo/project_bank/util"
	"github.com/yilinyo/project_bank/worker"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util2.Config{
		TokenSymmetricKey:   util2.RandomStr(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}
