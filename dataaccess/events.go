package dataaccess

import (
	"database/sql"
	"time"
)

type Event struct {
	ID             int
	Name           string
	StartTimestamp time.Time
	LapDistance    float64
}

func GetLatestEvent(db *sql.DB) (*Event, error) {
	row := db.QueryRow(`SELECT id, name, start_timestamp, lap_distance FROM events ORDER BY id DESC LIMIT 1`)
	var e Event
	var ts string
	if err := row.Scan(&e.ID, &e.Name, &ts, &e.LapDistance); err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, err
	}
	e.StartTimestamp = t
	return &e, nil
}

func InsertEvent(db *sql.DB, name string, startTimestamp time.Time, lapDistance float64) error {
	_, err := db.Exec(`DELETE FROM events`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO events (name, start_timestamp, lap_distance) VALUES (?, ?, ?)`, name, startTimestamp.Format(time.RFC3339), lapDistance)
	return err
}
