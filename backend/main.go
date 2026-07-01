package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	var err error

	db, err = ensureDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/health", health)
	http.HandleFunc("/users", users)

	fmt.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
