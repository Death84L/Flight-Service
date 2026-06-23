package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(host, port, user, password, dbname string) *Postgres {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect postgres:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Postgres not reachable:", err)
	}

	fmt.Println(" Connected to Postgres")

	return &Postgres{DB: db}
}
