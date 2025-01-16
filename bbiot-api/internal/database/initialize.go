package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func InitialiseDB() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL")+"?sslmode=disable")
	if err != nil {

		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return dbpool, err
}

func AutoMigrate() error {
	iofsDriver, err := iofs.New(EmbeddedFiles, "migrations")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, os.Getenv("DB_URL")+"?sslmode=disable")
	if err != nil {
		return err
	}

	err = migrator.Up()
	switch {
	case errors.Is(err, migrate.ErrNoChange):
		log.Println("No change in migrations")
	case err != nil:
		return err
	}
	log.Println("Migration complete")

	return nil
}
