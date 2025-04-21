package db

import (
	"database/sql"
	"log"

	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	return db
}

func Migrate(db *sql.DB) {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            role TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        );`,
		`CREATE TABLE IF NOT EXISTS pvz (
            id SERIAL PRIMARY KEY,
            city TEXT NOT NULL,
            registered_at TIMESTAMP DEFAULT NOW()
        );`,
		`CREATE TABLE IF NOT EXISTS sessions (
            id SERIAL PRIMARY KEY,
            pvz_id INT REFERENCES pvz(id) ON DELETE CASCADE,
            started_at TIMESTAMP DEFAULT NOW(),
            status TEXT NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS items (
            id SERIAL PRIMARY KEY,
            session_id INT REFERENCES sessions(id) ON DELETE CASCADE,
            type TEXT NOT NULL,
            added_at TIMESTAMP DEFAULT NOW()
        );`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			log.Fatalf("migration error: %v", err)
		}
	}
}
