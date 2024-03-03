// server/server.go
package server

import (
	"net/http"

	"workspace/handler"

	"github.com/gorilla/mux"
)

// StartServer starts the HTTP server
func StartServer() {
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", handler.IndexHandler).Methods("GET")
	r.HandleFunc("/convert", handler.ConvertHandler).Methods("POST")

	// Serve static files using built-in http.FileServer
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
