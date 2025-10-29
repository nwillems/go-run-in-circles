package dataaccess

import (
	"database/sql"
	"log"
)

// InitDB creates the required tables if they do not exist.
func InitDB(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS runners (
			bib_number TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			photo BLOB
		);`,
		`CREATE TABLE IF NOT EXISTS laps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bib_number TEXT NOT NULL,
			lap_number INTEGER NOT NULL,
			timestamp DATETIME NOT NULL,
			FOREIGN KEY(bib_number) REFERENCES runners(bib_number)
		);`,
		`CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			start_timestamp DATETIME NOT NULL,
			lap_distance REAL NOT NULL
		);`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
}
