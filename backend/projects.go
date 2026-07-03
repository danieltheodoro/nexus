package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func projects(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetProjects(w, r, userID)

	case http.MethodPost:
		handleCreateProject(w, r, userID)

	case http.MethodPatch:
		handleUpdateProject(w, r, userID)

	case http.MethodDelete:
		handleDeleteProject(w, r, userID)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetProjects(w http.ResponseWriter, r *http.Request, userID int) {
	idParam := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json")

	if idParam == "" {
		projects, err := getProjects(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(projects)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	project, err := getProject(id, userID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(project)
}

func handleCreateProject(w http.ResponseWriter, r *http.Request, userID int) {
	var req CreateProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := Project{
		UserID:      userID,
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

func handleUpdateProject(w http.ResponseWriter, r *http.Request, userID int) {
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

	if err := updateProject(id, userID, project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteProject(w http.ResponseWriter, r *http.Request, userID int) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := deleteProject(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getProjects(userID int) ([]Project, error) {
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
		WHERE user_id = ?
	`, userID)
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

func getProject(id int, userID int) (Project, error) {
	var project Project
	var description sql.NullString
	var color sql.NullString

	err := db.QueryRow(`
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
		WHERE id = ?
		AND user_id = ?
	`, id, userID).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&description,
		&color,
		&project.IsArchived,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return Project{}, err
	}

	project.Description = nullString(description)
	project.Color = nullString(color)

	return project, nil
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
		project.Color,
	)

	return err
}

func updateProject(id int, userID int, project Project) error {
	_, err := db.Exec(`
		UPDATE projects
		SET
			name = ?,
			description = ?,
			color = ?,
			is_archived = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		AND user_id = ?
	`,
		project.Name,
		project.Description,
		project.Color,
		project.IsArchived,
		id,
		userID,
	)

	return err
}

func deleteProject(id int, userID int) error {
	_, err := db.Exec(`
		DELETE FROM projects
		WHERE id = ?
		AND user_id = ?
	`,
		id,
		userID,
	)

	return err
}
