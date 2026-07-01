package main

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)

var db *sql.DB

func databaseExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "database/nexus.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createDatabase() (*sql.DB, error) {
	db, err := openDatabase()
	if err != nil {
		return nil, err
	}

	schemaBytes, err := os.ReadFile("database/schema.sql")
	if err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(string(schemaBytes)); err != nil {
		db.Close()
		return nil, err
	}

	seedBytes, err := os.ReadFile("database/seed.sql")
	if err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(string(seedBytes)); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func ensureDatabase() (*sql.DB, error) {
	if !databaseExists("database/nexus.db") {
		return createDatabase()
	}

	return openDatabase()
}
