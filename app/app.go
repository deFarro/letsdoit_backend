package main

import (
	"fmt"
	"net/http"

	"github.com/deFarro/letsdoit_backend/app/middleware"
	"github.com/deFarro/letsdoit_backend/app/router"
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

	http.Handle("/todo", middleware.Adapt(
		http.HandlerFunc(router.HandleEdit),
		middleware.WithHeaders,
		middleware.WithLogging,
	))

	http.ListenAndServe(port, nil)
}
