package main

import (
	"fmt"
	"github.com/Abhijith01/foodaggregator/internal/api"
	"log"
	"net/http"
	"os"
)

func main() {
	r := api.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	fmt.Printf("Started server at port %s\n", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal("Failed to start the server: \n", err)
	}
}