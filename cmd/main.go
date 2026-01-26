package main

import (
	"bytestream/internal/mock"
	"log"
	"net/http"
)

var port = ":8080"

func main() {
	mux := http.NewServeMux()

	mock.RegisterRoutes(mux)
	log.Printf("server listening on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
