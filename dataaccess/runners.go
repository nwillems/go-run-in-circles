package dataaccess

import (
	"database/sql"
)

type Runner struct {
	BibNumber string
	Name      string
	Photo     []byte
}

func InsertRunner(db *sql.DB, bibNumber, name string, photo []byte) error {
	_, err := db.Exec(`INSERT INTO runners (bib_number, name, photo) VALUES (?, ?, ?)`, bibNumber, name, photo)
	return err
}

func GetAllRunners(db *sql.DB) ([]Runner, error) {
	rows, err := db.Query(`SELECT bib_number, name, photo FROM runners`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var runners []Runner
	for rows.Next() {
		var r Runner
		if err := rows.Scan(&r.BibNumber, &r.Name, &r.Photo); err != nil {
			continue
		}
		runners = append(runners, r)
	}
	return runners, nil
}
