package main

import (
	"encoding/json"
	"net/http"
)

func labels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetLabels(w, r)

	case http.MethodPost:
		handleCreateLabel(w, r)

	case http.MethodPatch:
		handleUpdateLabel(w, r)

	case http.MethodDelete:
		handleDeleteLabel(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetLabels(w http.ResponseWriter, _ *http.Request) {
	labels, err := getLabels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(labels); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getLabels() ([]Label, error) {
	labels := []Label{}

	rows, err := db.Query(`
		SELECT
			id,
			project_id,
			name,
			color,
			created_at,
			updated_at
		FROM labels
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var label Label

		if err := rows.Scan(
			&label.ID,
			&label.ProjectID,
			&label.Name,
			&label.Color,
			&label.CreatedAt,
			&label.UpdatedAt,
		); err != nil {
			return nil, err
		}

		labels = append(labels, label)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return labels, nil
}

func handleCreateLabel(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleUpdateLabel(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleDeleteLabel(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
