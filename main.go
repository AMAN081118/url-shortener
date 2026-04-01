package main

import (
	"log"
	"net/http"
)

func main() {
	// TODO:
	// create in-memory storage
	storage := NewInMemoryStorage()

	// TODO:
	// pass storage into shortener
	shortener := NewURLShortener(storage)

	server := NewServer(shortener)
	mux := server.routes()

	log.Println("server running on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
