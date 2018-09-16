package main

import (
	"fmt"
	"net/http"

	"github.com/letsdoit_back/app/middleware"
	"github.com/letsdoit_back/app/router"
)

func main() {
	port := ":9090"
	fmt.Printf("Up and running on localhost%v\n", port)

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

	http.ListenAndServe(port, nil)
}
