package repository

import (
	"database/sql"
	"backend/internal/models"
	"time"

	_ "modernc.org/sqlite" // driver for sqlite
)

type SqlStorage struct {
	db *sql.DB
}

func NewSqlStorage(path string) (*SqlStorage, error) {
	db, err := sql.Open("sqlite", path)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1) //no more than 1 goroutine write to bd

	_, err = db.Exec("PRAGMA journal_mode=WAL;") //allow read file while writing
	if err != nil {
		return nil, err
	}

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS stats (
        symbol TEXT,
        price REAL,
        source TEXT,
        timedump DATETIME
    )`)
	if err != nil {
		return nil, err
	}
	statement.Exec()

	return &SqlStorage{db: db}, nil
}

func (s *SqlStorage) Save(stat models.Stat) error {
	// query request to sql/db
	query := `INSERT INTO stats (symbol, price, source, timedump) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, stat.Symbol, stat.Price, stat.Source, stat.Timedump)
	return err
}

func (r *SqlStorage) GetStat() ([]models.Stat, error) {
	rows, err := r.db.Query("SELECT symbol, price, source, timedump FROM stats ORDER BY timedump DESC LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.Stat
	for rows.Next() {
		var s models.Stat
		var timeStr string

		if err := rows.Scan(&s.Symbol, &s.Price, &s.Source, &timeStr); err != nil {
			return nil, err
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", timeStr)
		if err != nil {
			// Если формат отличается (например, без зоны), пробуем попроще:
			parsedTime, _ = time.Parse(time.RFC3339, timeStr)
		}

		s.Timedump = parsedTime

		stats = append(stats, s)
	}
	return stats, nil
}
