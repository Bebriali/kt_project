package sqlite

import (
	"backend/internal/models"
	"database/sql"
	"time"

	_ "modernc.org/sqlite" // driver for sqlite
)

type Storage struct {
	db *sql.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)

	if err != nil {
		return nil, err
	}

	//db.SetMaxOpenConns(1) //no more than 1 goroutine write to bd

	_, err = db.Exec("PRAGMA journal_mode=WAL;") //allow read file while writing
	if err != nil {
		return nil, err
	}

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS stats (
        base TEXT,
		quote TEXT,
		askprice REAL,
        bidprice REAL,
        source TEXT,
        timedump DATETIME
    )`)
	if err != nil {
		return nil, err
	}
	statement.Exec()

	return &Storage{db: db}, nil
}

func (s *Storage) Save(stat models.Stat) error {
	query := `INSERT INTO stats (base, quote, askprice, bidprice, source, timedump) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, stat.Base, stat.Quote, stat.AskPrice, stat.BidPrice, stat.Source, stat.Timedump)
	return err
}

func (r *Storage) GetStat() ([]models.Stat, error) {
	rows, err := r.db.Query("SELECT base, quote, askprice, bidprice, source, timedump FROM stats ORDER BY timedump DESC LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.Stat
	for rows.Next() {
		var s models.Stat
		var timeStr string

		if err := rows.Scan(&s.Base, &s.Quote, &s.AskPrice, &s.BidPrice, &s.Source, &timeStr); err != nil {
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
