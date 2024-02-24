package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"mengzhao/db"
)

func main() {
	db, err := createDB()
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}

func createDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pw := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	return db.CreateDatabase(dbName, user, pw, host)
}
