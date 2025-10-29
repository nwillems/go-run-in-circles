package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nwillems/go-run-in-circles/dataaccess"
)

func eventHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		event, err := dataaccess.GetLatestEvent(db)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))
			return
		}
		resp := map[string]interface{}{
			"id":              event.ID,
			"name":            event.Name,
			"start_timestamp": event.StartTimestamp.Format(time.RFC3339),
			"lap_distance":    event.LapDistance,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	case http.MethodPost:
		var req struct {
			Name           string  `json:"name"`
			StartTimestamp string  `json:"start_timestamp"`
			LapDistance    float64 `json:"lap_distance"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		t, err := time.Parse(time.RFC3339, req.StartTimestamp)
		if err != nil {
			http.Error(w, "Invalid timestamp", http.StatusBadRequest)
			return
		}
		err = dataaccess.InsertEvent(db, req.Name, t, req.LapDistance)
		if err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
