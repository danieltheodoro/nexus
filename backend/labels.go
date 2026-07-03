package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func labels(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetLabels(w, r, userID)

	case http.MethodPost:
		handleCreateLabel(w, r, userID)

	case http.MethodPatch:
		handleUpdateLabel(w, r, userID)

	case http.MethodDelete:
		handleDeleteLabel(w, r, userID)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetLabels(w http.ResponseWriter, r *http.Request, userID int) {
	idParam := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json")

	if idParam == "" {
		labels, err := getLabels(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(labels); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	label, err := getLabel(id, userID)
	if err != nil {
		http.Error(w, "Label not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(label); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCreateLabel(w http.ResponseWriter, r *http.Request, userID int) {
	var req CreateLabelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !projectBelongsToUser(req.ProjectID, userID) {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	label := Label{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Color:     req.Color,
	}

	if err := createLabel(label); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleUpdateLabel(w http.ResponseWriter, r *http.Request, userID int) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateLabelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	label := Label{
		Name:  req.Name,
		Color: req.Color,
	}

	if err := updateLabel(id, userID, label); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteLabel(w http.ResponseWriter, r *http.Request, userID int) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := deleteLabel(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getLabels(userID int) ([]Label, error) {
	labels := []Label{}

	rows, err := db.Query(`
		SELECT
			l.id,
			l.project_id,
			l.name,
			l.color,
			l.created_at,
			l.updated_at
		FROM labels l
		JOIN projects p
			ON l.project_id = p.id
		WHERE p.user_id = ?
	`, userID)
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

func getLabel(id int, userID int) (Label, error) {
	var label Label

	err := db.QueryRow(`
		SELECT
			l.id,
			l.project_id,
			l.name,
			l.color,
			l.created_at,
			l.updated_at
		FROM labels l
		JOIN projects p
			ON l.project_id = p.id
		WHERE l.id = ?
		AND p.user_id = ?
	`, id, userID).Scan(
		&label.ID,
		&label.ProjectID,
		&label.Name,
		&label.Color,
		&label.CreatedAt,
		&label.UpdatedAt,
	)
	if err != nil {
		return Label{}, err
	}

	return label, nil
}

func createLabel(label Label) error {
	_, err := db.Exec(`
		INSERT INTO labels(
			project_id,
			name,
			color
		) VALUES (?, ?, ?)
	`,
		label.ProjectID,
		label.Name,
		label.Color,
	)

	return err
}

func updateLabel(id int, userID int, label Label) error {
	_, err := db.Exec(`
		UPDATE labels
		SET
			name = ?,
			color = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		AND project_id IN (
			SELECT id
			FROM projects
			WHERE user_id = ?
		)
	`,
		label.Name,
		label.Color,
		id,
		userID,
	)

	return err
}

func deleteLabel(id int, userID int) error {
	_, err := db.Exec(`
		DELETE FROM labels
		WHERE id = ?
		AND project_id IN (
			SELECT id
			FROM projects
			WHERE user_id = ?
		)
	`,
		id,
		userID,
	)

	return err
}
