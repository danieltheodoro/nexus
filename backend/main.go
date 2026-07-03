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
	http.HandleFunc("/login", login)
	http.HandleFunc("/users", users)
	http.HandleFunc("/projects", projects)
	http.HandleFunc("/lists", lists)
	http.HandleFunc("/tasks", tasks)
	http.HandleFunc("/labels", labels)
	http.HandleFunc("/comments", comments)
	http.HandleFunc("/task-labels", taskLabels)

	fmt.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
