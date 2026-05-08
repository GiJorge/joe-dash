package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type AppLink struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}

var db *sql.DB

func main() {
	// 1. Initialize Database
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	var err error
	db, err = sql.Open("sqlite", "data/dashboard.db")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Ensure table exists
	query := `CREATE TABLE IF NOT EXISTS apps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		url TEXT,
		icon TEXT
	);`
	db.Exec(query)

	// 3. Setup Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/apps", appsHandler)       // Handles GET (all) and POST (add)
	mux.HandleFunc("/api/apps/", resourceHandler) // Handles PUT (edit) and DELETE (remove)



	// Add this at the end of your main function
// Serve static files from the Vue build (dist folder)
fileServer := http.FileServer(http.Dir("./dist"))
mux.Handle("/", fileServer)

log.Println("Server starting on http://0.0.0.0:8080")
log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}

// enableCORS allows the Vue dev server (port 5173) to communicate with Go
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// appsHandler handles list and create
func appsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rows, err := db.Query("SELECT id, name, url, icon FROM apps ORDER BY id ASC")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		apps := []AppLink{}
		for rows.Next() {
			var a AppLink
			rows.Scan(&a.ID, &a.Name, &a.URL, &a.Icon)
			apps = append(apps, a)
		}
		json.NewEncoder(w).Encode(apps)

	} else if r.Method == "POST" {
		var a AppLink
		json.NewDecoder(r.Body).Decode(&a)
		res, _ := db.Exec("INSERT INTO apps (name, url, icon) VALUES (?, ?, ?)", a.Name, a.URL, a.Icon)
		id, _ := res.LastInsertId()
		a.ID = int(id)
		json.NewEncoder(w).Encode(a)
	}
}

// resourceHandler handles update and delete based on ID
func resourceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path like "/api/apps/5"
	parts := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(parts[len(parts)-1])

	if r.Method == "DELETE" {
		db.Exec("DELETE FROM apps WHERE id = ?", id)
		w.WriteHeader(http.StatusNoContent)

	} else if r.Method == "PUT" {
		var a AppLink
		json.NewDecoder(r.Body).Decode(&a)
		_, err := db.Exec("UPDATE apps SET name=?, url=?, icon=? WHERE id=?", a.Name, a.URL, a.Icon, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a.ID = id
		json.NewEncoder(w).Encode(a)
	}
}