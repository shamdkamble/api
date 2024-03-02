package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ConversionResult struct to hold the result of conversion
type ConversionResult struct {
	Value float64
}

func main() {
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/convert", convertHandler).Methods("POST")

	// Serve static files using built-in http.FileServer
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}
func convertHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	valueStr := r.FormValue("value")
	unitFrom := r.FormValue("unitFrom")
	unitTo := r.FormValue("unitTo")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		http.Error(w, "Invalid value", http.StatusBadRequest)
		return
	}

	var result float64
	if unitFrom == "kg" && unitTo == "lb" {
		result = value * 2.20462
	} else if unitFrom == "lb" && unitTo == "kg" {
		result = value * 0.453592
	} else if unitFrom == "m" && unitTo == "ft" {
		result = value * 3.28084
	} else if unitFrom == "ft" && unitTo == "m" {
		result = value * 0.3048
	} else {
		result = value // if units are same, no conversion needed
	}

	res := ConversionResult{Value: result}
	tmpl := template.Must(template.ParseFiles("templates/result.html"))
	tmpl.Execute(w, res)
}
