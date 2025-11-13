package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/emrecolak-23/go-bank/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {

	config, err := utils.LoadConfig("../..")

	if err != nil {
		log.Fatal("can not load config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not conntect to db: ", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
