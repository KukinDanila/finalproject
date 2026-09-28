package server

import (
	"finalproject/pkg/api"
	"fmt"
	"net/http"
	"os"
)

func Run() error {
	port := "7540"

	envPort := os.Getenv("TODO_PORT")
	if envPort != "" {
		port = envPort
	}

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init()
	
	fmt.Println("Сервер запущен на порту", port)
	return http.ListenAndServe(":"+port, nil)

}
