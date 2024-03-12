package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

// ConversionResult struct to hold the result of conversion
type ConversionResult struct {
	Value    float64
	UnitFrom string
	UnitTo   string
	Input    float64
}

// ConversionMap stores the conversion factors for different units
var ConversionMap = map[string]map[string]float64{
	"kg": {"lb": 2.20462, "g": 1000},
	"lb": {"kg": 0.453592, "g": 453.592},
	"g":  {"kg": 0.001, "lb": 0.00220462},
	"m":  {"ft": 3.28084, "cm": 100, "km": 0.001},
	"ft": {"m": 0.3048, "cm": 30.48, "km": 0.0003048},
	"cm": {"m": 0.01, "ft": 0.0328084, "km": 0.00001},
	"km": {"m": 1000, "cm": 100000, "ft": 3280.84},
}

func DbConnection() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected to the database!")

	// Ensure conversion_results table exists
	err = createConversionResultsTable(db)
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// createConversionResultsTable creates the conversion_results table if it doesn't exist
func createConversionResultsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS conversion_results (
        id SERIAL PRIMARY KEY,
        value NUMERIC,
        unit_from VARCHAR(255),
        unit_to VARCHAR(255)
    )`)
	return err
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

	// Check if units are valid
	if _, ok := ConversionMap[unitFrom]; !ok {
		http.Error(w, fmt.Sprintf("Invalid unit from: %s", unitFrom), http.StatusBadRequest)
		return
	}
	if _, ok := ConversionMap[unitTo]; !ok {
		http.Error(w, fmt.Sprintf("Invalid unit to: %s", unitTo), http.StatusBadRequest)
		return
	}

	// Check if conversion factor exists
	conversionFactor, ok := ConversionMap[unitFrom][unitTo]
	if !ok {
		http.Error(w, fmt.Sprintf("Conversion from %s to %s not supported", unitFrom, unitTo), http.StatusBadRequest)
		return
	}

	result := value * conversionFactor
	res := ConversionResult{
		Value:    result,
		UnitFrom: unitFrom,
		UnitTo:   unitTo,
		Input:    value,
	}
	tmpl := template.Must(template.ParseFiles("templates/result.html"))
	tmpl.Execute(w, res)

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	_, err = db.Exec("INSERT INTO conversion_results (value, unit_from, unit_to) VALUES ($1, $2, $3)", result, unitFrom, unitTo)
	if err != nil {
		http.Error(w, "Failed to store conversion result", http.StatusInternalServerError)
		return
	}
}
