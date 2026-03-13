package postgres

import (
	"backend/internal/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

type Storage struct {
	db *sql.DB
}

// Передаем строку подключения (DSN) вместо пути к файлу
func NewStorage(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение сразу
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Настройки пула соединений (Postgres это любит)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	// В Postgres не нужны PRAGMA, создаем таблицу сразу
	// Используем TIMESTAMPTZ для корректной работы со временем
	query := `CREATE TABLE IF NOT EXISTS stats (
		id SERIAL PRIMARY KEY,
		base TEXT,
		quote TEXT,
		askprice DOUBLE PRECISION,
		bidprice DOUBLE PRECISION,
		source TEXT,
		timedump TIMESTAMPTZ
	)`

	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(stat models.Stat) error {
	query := `INSERT INTO stats (base, quote, askprice, bidprice, source, timedump)
	          VALUES ($1, $2, $3, $4, $5, $6)`
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

		// postgres is able to parse time from string itself
		if err := rows.Scan(&s.Base, &s.Quote, &s.AskPrice, &s.BidPrice, &s.Source, &s.Timedump); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
