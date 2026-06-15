package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"embed" // Add this
    "io/fs" // Add this
	"strconv"
	"strings"
	"github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/mem"
	//"github.com/shirou/gopsutil/v3/host" // Add this
	"time"
	"math"



	_ "modernc.org/sqlite"
)

//go:embed frontend/dist/*
var frontendAssets embed.FS


type AppLink struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
	Status string `json:"status"`
}

var db *sql.DB

func main() {
	// 1. Initialize Database
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}
// Register the route

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

	mux.HandleFunc("/api/stats", getStats)
// The folder is now at frontend/dist relative to main.go
distFS, err := fs.Sub(frontendAssets, "frontend/dist")
if err != nil {
    log.Fatal(err)
}

mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    // Check if the requested file exists in the embedded filesystem
    path := r.URL.Path
    f, err := distFS.Open(path[1:]) // Remove leading slash for Open()

    if err != nil {
        // If file not found, serve index.html (let Vue Router handle the route)
        http.ServeFileFS(w, r, distFS, "index.html")
        return
    }
    f.Close()

    // Otherwise, serve the static file
    http.FileServer(http.FS(distFS)).ServeHTTP(w, r)
})

    log.Println("Server starting on http://0.0.0.0:10000")
    log.Fatal(http.ListenAndServe(":10000", enableCORS(mux)))
}



// 1. Create a custom handler to serve index.html for unknown paths
func spaHandler(fs http.FileSystem) http.Handler {
    fileServer := http.FileServer(fs)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Try to open the requested file
        f, err := fs.Open(r.URL.Path)
        if os.IsNotExist(err) {
            // If file doesn't exist, serve index.html
            http.ServeFile(w, r, "dist/index.html")
            return
        }
        f.Close()
        fileServer.ServeHTTP(w, r)
    })
}


// checkLiveness sends a quick HTTP request to verify if the link is reachable
func checkLiveness(targetURL string) string {
	// Fallback for local servers provided without http:// or https://
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	client := http.Client{
		Timeout: 3 * time.Second, // Don't let slow servers hang the dashboard
	}

	// Using HEAD is faster than GET because it doesn't download the body text/HTML
	resp, err := client.Head(targetURL)
	if err != nil {
		// Fallback to GET just in case the server blocks HEAD requests
		resp, err = client.Get(targetURL)
		if err != nil {
			return "offline"
		}
	}
	defer resp.Body.Close()

	// Consider any 2xx or 3xx status code as online
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "online"
	}
	return "offline"
}







func getStats(w http.ResponseWriter, r *http.Request) {
    // 1. Get RAM Stats
    v, err := mem.VirtualMemory()
    if err != nil {
        http.Error(w, "Failed to get RAM info", http.StatusInternalServerError)
        return
    }

    // 2. Get CPU Stats (Calculated over a 500ms window)
    cpuPercent, err := cpu.Percent(500*time.Millisecond, false)
    if err != nil || len(cpuPercent) == 0 {
        http.Error(w, "Failed to get CPU info", http.StatusInternalServerError)
        return
    }
		//hInfo, _ := host.Info()
    // 3. Prepare the response
    stats := map[string]interface{}{
        "cpu":       math.Round(cpuPercent[0]),
        "ram":       math.Round(v.UsedPercent),
        "ram_gb":    math.Round(float64(v.Used) / 1024 / 1024 / 1024),
        "total_gb":  math.Round(float64(v.Total) / 1024 / 1024 / 1024),
		"temp":   math.Round(getCPUTemp()), // <-- Add this
		"hosttemp":   math.Round(getHostTemp()), // <-- HosT
		"platform": getOSName(), // This will now return "Debian GNU/Linux 13 (trixie)"
      //  "version":  hInfo.PlatformVersion, // e.g., "7.0.5-1-cachyos"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}


func getOSName() string {
    // 1. Check for the mapped host file first (Docker environment)
    data, err := os.ReadFile("/host/etc/os-release")

    // 2. Fallback to the standard path if not in Docker
    if err != nil {
        data, err = os.ReadFile("/etc/os-release")
        if err != nil {
            return "Unknown Linux"
        }
    }

    // 3. Parse the file to find PRETTY_NAME="Debian GNU/Linux 13 (trixie)"
    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "PRETTY_NAME=") {
            name := strings.TrimPrefix(line, "PRETTY_NAME=")
            return strings.Trim(name, "\"") // Removes the quotes
        }
    }

    return "Linux"
}



func getHostTemp() float64 {
    // Path inside the VM that points to the physical host sensor
    data, err := os.ReadFile("/mnt/host_thermal/thermal_zone0/temp")
    if err != nil {
        return 0
    }

    tempRaw := strings.TrimSpace(string(data))
    tempInt, _ := strconv.Atoi(tempRaw)

    return float64(tempInt) / 1000.0
}


func getCPUTemp() float64 {
    // Standard path for CPU temperature on most Linux distros
    data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
    if err != nil {
        // Fallback for some hardware/kernels (like Ryzen on older kernels)
        data, err = os.ReadFile("/sys/class/hwmon/hwmon0/temp1_input")
        if err != nil {
            return 0
        }
    }

    // Linux stores temp in millidegrees (e.g., 55000 = 55°C)
    tempRaw := strings.TrimSpace(string(data))
    tempInt, _ := strconv.Atoi(tempRaw)

    return float64(tempInt) / 1000.0
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
// Check status dynamically before returning to frontend
            a.Status = checkLiveness(a.URL)

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
