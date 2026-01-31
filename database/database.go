package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open DB
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test Connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected")
	return db, nil
}
