package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/utils"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("viper config failed in main test", err)
	}
	//var err error
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("not able to connect to database :", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
