package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nwillems/go-run-in-circles/dataaccess"
)

//go:embed web/*
var content embed.FS

var db *sql.DB

func main() {
	// Open SQLite DB (filename from command line or default)
	dbFile := "event.db"
	if len(os.Args) > 1 && os.Args[1] != "" {
		dbFile = os.Args[1]
	}
	var err error
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	// Initialize DB tables
	dataaccess.InitDB(db)

	// Attach handlers from other files
	http.HandleFunc("/laps", func(w http.ResponseWriter, r *http.Request) { lapsHandler(db, w, r) })
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) { dashboardHandler(db, w, r) })
	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) { eventHandler(db, w, r) })
	http.HandleFunc("/runners", func(w http.ResponseWriter, r *http.Request) { runnersHandler(db, w, r) })

	// Serve static files from /web/ at the root URL
	fs := http.FS(content)
	fileServer := http.FileServer(fs)

	// Redirect root to /web/index.html for SPA entry
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			http.ServeFileFS(w, r, content, "web/index.html")
			return
		}

		r.URL.Path = "/web/" + r.URL.Path
		fileServer.ServeHTTP(w, r)
	})

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
