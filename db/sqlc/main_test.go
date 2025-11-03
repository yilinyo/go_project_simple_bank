package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/yilinyo/project_bank/db/util"
	"log"
	"os"
	"testing"
)

//const (
//	dbDriver = "postgres"
//	dbSource = "postgresql://root:yilin123@localhost:5432/simple_bank?sslmode=disable"
//)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("error opening db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())

}
