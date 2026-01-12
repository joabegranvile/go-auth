package db

import (
	"database/sql"
	"time"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(conn string) (*Postgres, error) {
	var (
		database *sql.DB
		err      error
	)

	for range 10 {
		database, err = sql.Open("postgres", conn)
		if err == nil {
			err = database.Ping()
		}
		if err == nil {
			return &Postgres{DB: database}, nil
		}
		time.Sleep(2 * time.Second)

	}
	return nil, err
}

func (p *Postgres) Migrate() error {
	_, err := p.DB.Exec(`
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				username TEXT,
				role TEXT
			)
		`)
	return err
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}
