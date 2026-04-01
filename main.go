package main

import (
	"log"
	"net/http"
)

func main() {
	shortener := NewURLShortener()
	server := NewServer(shortener)

	mux := server.routes()

	log.Println("server running on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
