package server

import (
	"fmt"
	"net/http"
	"workspace/handler"

	"github.com/gorilla/mux"
)

// StartServer starts the HTTP server
func StartServer() {
	r := mux.NewRouter()
	handler.DbConnection()

	// Define routes
	r.HandleFunc("/", handler.IndexHandler).Methods("GET")
	r.HandleFunc("/convert", handler.ConvertHandler).Methods("POST")

	// Serve static files using built-in http.FileServer
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	fmt.Println("Server started at http://localhost:5000")
	http.ListenAndServe(":5000", r)
}
