package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// main is a script used for creating the schema and inserting test data into a locally running docker postgres database
func main() {
	connInfo := "host=localhost port=5432 user=test password=test dbname=space sslmode=disable"
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatalf("connect to postgres: %s", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("close connection: %s", err)
		}
	}(db)

	if err = executeSQL(db); err != nil {
		log.Fatalf(err.Error())
	}
}

func executeSQL(db *sql.DB) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if !errors.Is(err, sql.ErrTxDone) && err != nil {
			log.Fatalf("rollback transaction: %s", err)
		}
	}(tx)

	sqlFiles := []string{"sql/exoplanet_create.sql", "sql/exoplanet_insert.sql"}

	for _, sqlFile := range sqlFiles {
		sqlBytes, err := os.ReadFile(sqlFile)
		if err != nil {
			return fmt.Errorf("read sql file %s: %w", sqlFile, err)
		}

		if _, err = tx.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("execute %s: %w", sqlFile, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
