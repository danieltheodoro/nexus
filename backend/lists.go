package main

import (
	"encoding/json"
	"net/http"
)

func lists(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetLists(w, r)

	case http.MethodPost:
		handleCreateList(w, r)

	case http.MethodPatch:
		handleUpdateList(w, r)

	case http.MethodDelete:
		handleDeleteList(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetLists(w http.ResponseWriter, _ *http.Request) {
	lists, err := getLists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(lists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getLists() ([]List, error) {
	lists := []List{}

	rows, err := db.Query(`
		SELECT
			id,
			project_id,
			name,
			position,
			created_at,
			updated_at
		FROM lists
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var list List

		if err := rows.Scan(
			&list.ID,
			&list.ProjectID,
			&list.Name,
			&list.Position,
			&list.CreatedAt,
			&list.UpdatedAt,
		); err != nil {
			return nil, err
		}

		lists = append(lists, list)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

func handleCreateList(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleUpdateList(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleDeleteList(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
