package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/nwillems/go-run-in-circles/dataaccess"
)

// runnersHandler handles POST /runners to register a new runner.
func runnersHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		BibNumber string `json:"bib_number"`
		Name      string `json:"name"`
		Photo     string `json:"photo"` // base64 encoded
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.BibNumber == "" || req.Name == "" || req.Photo == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	photoBytes, err := base64.StdEncoding.DecodeString(req.Photo)
	if err != nil {
		http.Error(w, "Invalid photo encoding", http.StatusBadRequest)
		return
	}
	err = dataaccess.InsertRunner(db, req.BibNumber, req.Name, photoBytes)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
