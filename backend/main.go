package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
)

func databaseExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func createDatabase() error {
	db, err := sql.Open("sqlite", "database/nexus.db")
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	schemaBytes, err := os.ReadFile("database/schema.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(schemaBytes)); err != nil {
		return err
	}

	seedBytes, err := os.ReadFile("database/seed.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(seedBytes)); err != nil {
		return err
	}

	return nil
}

func ensureDatabase() error {
	if !databaseExists("database/nexus.db") {
		return createDatabase()
	}

	db, err := sql.Open("sqlite", "database/nexus.db")
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Ping()
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprint(w, `{"status":"ok"}`)
}

func users(w http.ResponseWriter, _ *http.Request) {
	db, err := sql.Open("sqlite", "database/nexus.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, username, email, first_name, last_name FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprint(w, `[`)
	first := true

	for rows.Next() {
		var id int
		var username, email string
		var firstName, lastName sql.NullString

		if err := rows.Scan(&id, &username, &email, &firstName, &lastName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !first {
			fmt.Fprint(w, `,`)
		}
		first = false

		fmt.Fprintf(w, `{"id":%d,"username":%q,"email":%q}`, id, username, email)
	}

	fmt.Fprint(w, `]`)
}

func main() {
	ensureDatabase()

	http.HandleFunc("/health", health)
	http.HandleFunc("/users", users)
	fmt.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
