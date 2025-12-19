package work

import (
	"github.com/hibiken/asynq"
	db "github.com/yilinyo/project_bank/db/sqlc"
)

type TaskProcessor interface {
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}
