package main

import (
	"fmt"
	"net/http"
	"log"

	"github.com/deFarro/letsdoit_backend/app/middleware"
	"github.com/deFarro/letsdoit_backend/app/router"
	"github.com/deFarro/letsdoit_backend/app/config"
)

func main() {
	config, err := config.GetConfig("config.yml")
	if err != nil {
		log.Fatalf("error while reading config file: %s\n", err)
	}

	router, err := router.NewRouter(config)
	if err != nil {
		log.Fatalf("error while creating router: %s\n", err)
	}

	fmt.Printf("Up and running on localhost:%s\n", config.AppPort)

	http.Handle("/user/login", middleware.Adapt(
		http.HandlerFunc(router.HandleLogin),
		middleware.WithHeaders,
		middleware.WithLogging,
	))

	http.Handle("/user/logout", middleware.Adapt(
		http.HandlerFunc(router.HandleLogout),
		middleware.WithHeaders,
		middleware.WithLogging,
	))

	http.Handle("/todos", middleware.Adapt(
		http.HandlerFunc(router.HandleTodos),
		middleware.WithHeaders,
		middleware.WithLogging,
	))

	http.Handle("/todo", middleware.Adapt(
		http.HandlerFunc(router.HandleEdit),
		middleware.WithHeaders,
		middleware.WithLogging,
	))

	http.Handle("/version", middleware.Adapt(
		http.HandlerFunc(router.HandleVersion),
		middleware.WithHeaders,
	))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

	http.ListenAndServe(":" + config.AppPort, nil)
}
