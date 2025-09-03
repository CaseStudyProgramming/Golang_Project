package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // driver postgres
)

func NewPostgresDB(host string, port int, user, password, dbname, sslmode string) *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
	return db
}
