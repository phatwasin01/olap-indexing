package main

import (
	"log"

	"github.com/phatwasin01/olap-indexing/api"
	"github.com/phatwasin01/olap-indexing/db"
)

func main() {
	redis := db.RedisClient()
	server, err := api.NewServer(redis)
	if err != nil {
		log.Fatal("Cannot Setup Server", err)
	}
	err = server.Start()
	if err != nil {
		log.Fatal("Cannot Start Server", err)
	}
}
