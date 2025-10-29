package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nwillems/go-run-in-circles/dataaccess"
)

type Runner struct {
	BibNumber      string  `json:"bib_number"`
	Name           string  `json:"name"`
	TotalLaps      int     `json:"total_laps"`
	LastLapTime    string  `json:"last_lap_time"`
	AverageLapTime string  `json:"average_lap_time"`
	FastestLap     string  `json:"fastest_lap"`
	SlowestLap     string  `json:"slowest_lap"`
	AvgPace        string  `json:"avg_pace"`
	TotalDistance  float64 `json:"total_distance"`
	Photo          string  `json:"photo"`
}

func calculateRunnerStats(db *sql.DB, event *dataaccess.Event, rd dataaccess.Runner) Runner {
	var r Runner
	r.BibNumber = rd.BibNumber
	r.Name = rd.Name
	r.Photo = base64.StdEncoding.EncodeToString(rd.Photo)
	lapTimes, err := dataaccess.GetLapTimes(db, r.BibNumber)
	if err != nil {
		return r
	}
	r.TotalLaps = len(lapTimes)
	if r.TotalLaps > 0 {
		var fastest, slowest, total time.Duration
		allTimes := make([]time.Time, 0, len(lapTimes)+1)
		allTimes = append(allTimes, event.StartTimestamp)
		allTimes = append(allTimes, lapTimes...)
		for i := 1; i < len(allTimes); i++ {
			lapDur := allTimes[i].Sub(allTimes[i-1])
			if i == 1 || lapDur < fastest {
				fastest = lapDur
			}
			if i == 1 || lapDur > slowest {
				slowest = lapDur
			}
			total += lapDur
		}
		if len(allTimes) > 1 {
			avg := total / time.Duration(len(allTimes)-1)
			r.AverageLapTime = formatDuration(avg)
			r.FastestLap = formatDuration(fastest)
			r.SlowestLap = formatDuration(slowest)
			r.LastLapTime = formatDuration(allTimes[len(allTimes)-1].Sub(allTimes[len(allTimes)-2]))
			if event.LapDistance > 0 {
				avgPace := avg / time.Duration(event.LapDistance)
				r.AvgPace = formatDuration(avgPace)

				r.TotalDistance = float64(r.TotalLaps) * event.LapDistance
			}
		}
	}
	return r
}

func dashboardHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	event, err := dataaccess.GetLatestEvent(db)
	if err != nil || event.LapDistance <= 0 {
		http.Error(w, "Event lap_distance not set", http.StatusInternalServerError)
		return
	}
	runnersData, err := dataaccess.GetAllRunners(db)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var runners []Runner
	for _, rd := range runnersData {
		runners = append(runners, calculateRunnerStats(db, event, rd))
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(runners)
}
