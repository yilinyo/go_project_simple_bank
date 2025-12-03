package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/yilinyo/project_bank/db/sqlc"
	util2 "github.com/yilinyo/project_bank/util"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util2.Config{
		TokenSymmetricKey:   util2.RandomStr(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())

}
