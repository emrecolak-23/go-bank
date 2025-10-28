package main

import (
	"database/sql"
	"log"

	"github.com/emrecolak-23/go-bank/api"
	db "github.com/emrecolak-23/go-bank/db/sqlc"
	"github.com/emrecolak-23/go-bank/utils"
	_ "github.com/lib/pq"
)

func main() {

	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("can not load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not conntect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot Create Server:", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot Start Server:", err)
	}
}
