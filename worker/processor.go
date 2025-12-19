package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/yilinyo/project_bank/db/sqlc"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	Start() error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Bytes("", task.Payload()).Msg("process task failed")
			}),
			Logger: NewLogger(),
		})
	return &RedisTaskProcessor{
		store:  store,
		server: server,
	}
}

func (r *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}
	user, err := r.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//跳过重试
			return fmt.Errorf("no user row: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("get user wrong: %w", err)
	}
	//todo :real send email
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("processed task")
	return nil
}
func (r *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	//告诉 某个type 的 task 由哪个handler 处理
	mux.HandleFunc(TaskSendVerifyEmail, r.ProcessTaskSendVerifyEmail)

	err := r.server.Start(mux)
	if err != nil {
		return fmt.Errorf("start asynq server err: %w", err)
	}
	return nil
}
