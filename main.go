package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yilinyo/project_bank/mail"
	"github.com/yilinyo/project_bank/worker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v5"
	"github.com/yilinyo/project_bank/api"
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/gapi"
	"github.com/yilinyo/project_bank/pb"
	"github.com/yilinyo/project_bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

//const (
//	dbDriver      = "postgres"
//	dbSource      = "postgresql://root:yilin123@localhost:5432/simple_bank?sslmode=disable"
//	serverAddress = "0.0.0.0:8080"
//)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config:")
	}
	if config.Environment == "prod" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("error opening db:")
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisDistributor(redisOpt)
	runDBMigration(config.MigrationURL, config.DBSource)
	store := db.NewStore(connPool)

	//开启taskProcess 进行消费
	go runTaskProcessor(config, redisOpt, store)

	//选择 开启http 还是 grpc 还是httpGateWay
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)

	//runGinServer(config, store)

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("error new server:")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("error starting server:")
	}
}

func runGrpcServer(config util.Config, store db.Store, d worker.TaskDistributor) {

	server, err := gapi.NewServer(config, store, d)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store, d worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, d)
	if err != nil {
		log.Fatal().Msg("cannot create server:")
	}
	//使用proto定以字段名称进行编解码
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server:")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create HTTPServer listener")
	}

	log.Info().Msgf("start httpGateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start httpGateway server")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing migration:")
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("error running migrate up")
	}

	log.Info().Msg("db migrate up done")

}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {

	emailSender := mail.NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword,
	)
	processor := worker.NewRedisTaskProcessor(redisOpt, store, emailSender)
	log.Info().Msg("task processor started")
	err := processor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create runTask processor")
	}
}
