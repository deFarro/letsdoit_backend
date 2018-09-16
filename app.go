package main

import (
	"fmt"
	"net/http"

	"github.com/letsdoit_back/app/router"
)

func main() {
	port := ":9090"
	fmt.Printf("Up and running on localhost%v\n", port)

	http.HandleFunc("/user", router.HandleUser)
	http.HandleFunc("/todos", router.HandleTodos)

	http.ListenAndServe(port, nil)
}
