package dataaccess

import (
	"database/sql"
	"time"
)

type Lap struct {
	ID        int
	BibNumber string
	LapNumber int
	Timestamp time.Time
}

func InsertLap(db *sql.DB, bibNumber string, lapNumber int, timestamp time.Time) error {
	_, err := db.Exec(`INSERT INTO laps (bib_number, lap_number, timestamp) VALUES (?, ?, ?)`, bibNumber, lapNumber, timestamp.Format(time.RFC3339))
	return err
}

func GetLapCount(db *sql.DB, bibNumber string) (int, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM laps WHERE bib_number = ?`, bibNumber).Scan(&count)
	return count, err
}

func GetLapTimes(db *sql.DB, bibNumber string) ([]time.Time, error) {
	rows, err := db.Query(`SELECT timestamp FROM laps WHERE bib_number = ? ORDER BY lap_number`, bibNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var times []time.Time
	for rows.Next() {
		var ts string
		if err := rows.Scan(&ts); err == nil {
			t, err := time.Parse(time.RFC3339, ts)
			if err == nil {
				times = append(times, t)
			}
		}
	}
	return times, nil
}
