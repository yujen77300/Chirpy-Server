package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server failed:", err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
