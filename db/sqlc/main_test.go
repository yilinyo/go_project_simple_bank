package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yilinyo/project_bank/util"

	"log"
	"os"
	"testing"
)

//const (
//	dbDriver = "postgres"
//	dbSource = "postgresql://root:yilin123@localhost:5432/simple_bank?sslmode=disable"
//)

var testStore Store

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("error opening db:", err)
	}
	testStore = NewStore(connPool)
	os.Exit(m.Run())

}
