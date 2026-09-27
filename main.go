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

	//now, _ := time.Parse("20060102", "20240126")
	//result, err := api.NextDate(now, "20240113", "d 7")
	//fmt.Println("Результат:", result, "Ошибка:", err)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
