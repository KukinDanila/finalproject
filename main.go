package main

import (
	"finalproject/pkg/db"
	"finalproject/pkg/server"
	"log"
	"os"
)

func main() {
	dbFile := "scheduler.db"

	envDBFile := os.Getenv("TODO_DBFILE")
	if envDBFile != "" {
		dbFile = envDBFile
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
