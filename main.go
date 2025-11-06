package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/yilinyo/project_bank/api"
	db "github.com/yilinyo/project_bank/db/sqlc"
	"github.com/yilinyo/project_bank/db/util"
)

//const (
//	dbDriver      = "postgres"
//	dbSource      = "postgresql://root:yilin123@localhost:5432/simple_bank?sslmode=disable"
//	serverAddress = "0.0.0.0:8080"
//)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	fmt.Println(config)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("error opening db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("error starting server:", err)
	}

}
