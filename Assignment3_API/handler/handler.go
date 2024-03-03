// handler/handler.go
package handler

import (
	"html/template"
	"net/http"
	"strconv"
)

// ConversionResult struct to hold the result of conversion
type ConversionResult struct {
	Value float64
}

// IndexHandler handles the index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

// ConvertHandler handles the conversion form submission
func ConvertHandler(w http.ResponseWriter, r *http.Request) {
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
	switch {
	case unitFrom == "kg" && unitTo == "lb":
		result = value * 2.20462
	case unitFrom == "lb" && unitTo == "kg":
		result = value * 0.453592
	case unitFrom == "m" && unitTo == "ft":
		result = value * 3.28084
	case unitFrom == "ft" && unitTo == "m":
		result = value * 0.3048
	// Additional units for weight conversion
	case unitFrom == "g" && unitTo == "kg":
		result = value * 0.001
	case unitFrom == "kg" && unitTo == "g":
		result = value * 1000
	case unitFrom == "lb" && unitTo == "g":
		result = value * 453.592
	case unitFrom == "g" && unitTo == "lb":
		result = value * 0.00220462
	// Additional units for distance conversion
	case unitFrom == "cm" && unitTo == "m":
		result = value * 0.01
	case unitFrom == "m" && unitTo == "cm":
		result = value * 100
	case unitFrom == "km" && unitTo == "m":
		result = value * 1000
	case unitFrom == "m" && unitTo == "km":
		result = value * 0.001
	default:
		result = value // if units are same, no conversion needed
	}

	res := ConversionResult{Value: result}
	tmpl := template.Must(template.ParseFiles("templates/result.html"))
	tmpl.Execute(w, res)
}
