package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func projects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetProjects(w, r)

	case http.MethodPost:
		handleCreateProject(w, r)

	case http.MethodPatch:
		handleUpdateProject(w, r)

	case http.MethodDelete:
		handleDeleteProject(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetProjects(w http.ResponseWriter, _ *http.Request) {
	projects, err := getProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCreateProject(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleUpdateProject(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleDeleteProject(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func getProjects() ([]Project, error) {
	projects := []Project{}

	rows, err := db.Query(`
		SELECT 
			id, 
			user_id, 
			name, 
			description, 
			color,
			is_archived,
			created_at,
			updated_at
		FROM projects
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var project Project
		var description sql.NullString
		var color sql.NullString

		if err := rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Name,
			&description,
			&color,
			&project.IsArchived,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, err
		}

		project.Description = nullString(description)
		project.Color = nullString(color)

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
