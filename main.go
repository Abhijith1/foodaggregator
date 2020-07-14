package main

import (
	"fmt"
	"github.com/Abhijith01/foodaggregator/internal/api"
	"github.com/Abhijith01/foodaggregator/internal/config"
	"log"
	"net/http"
)

func main() {
	r := api.NewRouter()
	port := config.ServerPort

	fmt.Printf("Started server at port %s\n", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal("Failed to start the server: \n", err)
	}
}
