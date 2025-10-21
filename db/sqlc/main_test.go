package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/emrecolak-23/go-bank/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := utils.LoadConfig("../..")

	if err != nil {
		log.Fatal("can not load config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not conntect to db: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
