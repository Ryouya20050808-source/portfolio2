package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type AccessLogCount struct {
	PostalCode   string `json:"postal_code"`
	RequestCount int    `json:"request_count"`
}

var db *sql.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}

	log.Println("DB connection successful")
}

func SaveAccessLog(postal string) {
	_, err := db.Exec("INSERT INTO access_logs (postal_code) VALUES (?)", postal)
	if err != nil {
		log.Println("Failed to save access log:", err)
	}
}

func accessLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query(`
		SELECT postal_code, COUNT(*) AS request_count
		FROM access_logs
		GROUP BY postal_code
		ORDER BY request_count DESC
	`)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	var logs []AccessLogCount
	for rows.Next() {
		var logEntry AccessLogCount
		if err := rows.Scan(&logEntry.PostalCode, &logEntry.RequestCount); err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		logs = append(logs, logEntry)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_logs": logs,
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, Go API Server!"})
}

type ApiReceive struct {
	Response struct {
		Location []struct {
			Postal     string `json:"postal"`
			Prefecture string `json:"prefecture"`
			City       string `json:"city"`
			Town       string `json:"town"`
			X          string `json:"x"`
			Y          string `json:"y"`
		} `json:"location"`
	} `json:"response"`
}

type ApiReturn struct {
	Postal_code        string  `json:"postal_code"`
	Hit_count          int     `json:"hit_count"`
	Address            string  `json:"address"`
	Tokyo_sta_distance float64 `json:"tokyo_sta_distance"`
}

func distToTokyo(lon, lat float64) float64 {
	xt := 139.7673068
	yt := 35.6809591
	R := 6371.0

	x1 := (lon - xt) * math.Cos(math.Pi*(lat+yt)/360.0)
	y1 := lat - yt
	d := math.Pi * R / 180.0 * math.Sqrt(x1*x1+y1*y1)
	return math.Round(d*10) / 10
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postal := r.URL.Query().Get("postal")
	if postal == "" {
		http.Error(w, "postal parameter required", http.StatusBadRequest)
		return
	}

	SaveAccessLog(postal)

	resp, err := http.Get("https://geoapi.heartrails.com/api/json?method=searchByPostal&postal=" + postal)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("API request error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var rec ApiReceive
	if err := json.Unmarshal(body, &rec); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("JSON parse error:", err)
		return
	}

	if len(rec.Response.Location) == 0 {
		http.Error(w, "no result", http.StatusNotFound)
		return
	}

	loc := rec.Response.Location[0]
	lon, err1 := strconv.ParseFloat(loc.X, 64)
	lat, err2 := strconv.ParseFloat(loc.Y, 64)
	distance := -1.0
	if err1 == nil && err2 == nil {
		distance = distToTokyo(lon, lat)
	}

	address := loc.Prefecture + loc.City + loc.Town

	json.NewEncoder(w).Encode(ApiReturn{
		Postal_code:        postal,
		Hit_count:          len(rec.Response.Location),
		Address:            address,
		Tokyo_sta_distance: distance,
	})
}

func main() {
	InitDB()

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", apiHandler)
	http.HandleFunc("/address/access_logs", accessLogsHandler)

	log.Println("Server running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
