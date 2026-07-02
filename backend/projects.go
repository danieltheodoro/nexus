package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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

func handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := Project{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
	}

	if err := createProject(project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleUpdateProject(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := Project{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		IsArchived:  req.IsArchived,
	}

	if err := updateProject(id, project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := deleteProject(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

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

func createProject(project Project) error {
	_, err := db.Exec(`
		INSERT INTO projects(
			user_id,
			name,
			description,
			color
		) VALUES (?, ?, ?, ?)
	`,
		project.UserID,
		project.Name,
		project.Description,
		project.Color)

	return err

}

func updateProject(id int, project Project) error {
	_, err := db.Exec(`
		UPDATE projects
		SET
			name = ?,
			description = ?,
			color = ?,
			is_archived = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`,
		project.Name,
		project.Description,
		project.Color,
		project.IsArchived,
		id,
	)

	return err
}

func deleteProject(id int) error {
	_, err := db.Exec(`
						DELETE FROM projects
						WHERE id = ?`,
		id,
	)
	return err

}
