package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nwillems/go-run-in-circles/dataaccess"
)

func lapsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		BibNumber string `json:"bib_number"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	ts := time.Now()
	if req.Timestamp != "" {
		t, err := time.Parse(time.RFC3339, req.Timestamp)
		if err == nil {
			ts = t
		}
	}
	lapCount, err := dataaccess.GetLapCount(db, req.BibNumber)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = dataaccess.InsertLap(db, req.BibNumber, lapCount+1, ts)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
