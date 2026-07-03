package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func taskLabels(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetTaskLabels(w, r, userID)

	case http.MethodPost:
		handleCreateTaskLabel(w, r, userID)

	case http.MethodDelete:
		handleDeleteTaskLabel(w, r, userID)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetTaskLabels(w http.ResponseWriter, r *http.Request, userID int) {
	taskID, err := strconv.Atoi(r.URL.Query().Get("task_id"))
	if err != nil {
		http.Error(w, "Invalid task_id", http.StatusBadRequest)
		return
	}

	if !taskBelongsToUser(taskID, userID) {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	labels, err := getTaskLabels(taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(labels); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCreateTaskLabel(w http.ResponseWriter, r *http.Request, userID int) {
	var req CreateTaskLabelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !taskBelongsToUser(req.TaskID, userID) {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if !labelBelongsToUser(req.LabelID, userID) {
		http.Error(w, "Label not found", http.StatusNotFound)
		return
	}

	taskLabel := TaskLabel{
		TaskID:  req.TaskID,
		LabelID: req.LabelID,
	}

	if err := createTaskLabel(taskLabel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleDeleteTaskLabel(w http.ResponseWriter, r *http.Request, userID int) {
	taskID, err := strconv.Atoi(r.URL.Query().Get("task_id"))
	if err != nil {
		http.Error(w, "Invalid task_id", http.StatusBadRequest)
		return
	}

	labelID, err := strconv.Atoi(r.URL.Query().Get("label_id"))
	if err != nil {
		http.Error(w, "Invalid label_id", http.StatusBadRequest)
		return
	}

	if !taskBelongsToUser(taskID, userID) {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if !labelBelongsToUser(labelID, userID) {
		http.Error(w, "Label not found", http.StatusNotFound)
		return
	}

	if err := deleteTaskLabel(taskID, labelID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getTaskLabels(taskID int, userID int) ([]Label, error) {
	labels := []Label{}

	rows, err := db.Query(`
		SELECT
			l.id,
			l.project_id,
			l.name,
			l.color,
			l.created_at,
			l.updated_at
		FROM task_labels tl
		JOIN labels l
			ON tl.label_id = l.id
		JOIN projects p
			ON l.project_id = p.id
		WHERE tl.task_id = ?
		AND p.user_id = ?
	`,
		taskID,
		userID,
	)
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

func createTaskLabel(taskLabel TaskLabel) error {
	_, err := db.Exec(`
		INSERT INTO task_labels (
			task_id,
			label_id
		) VALUES (?, ?)
	`,
		taskLabel.TaskID,
		taskLabel.LabelID,
	)

	return err
}

func deleteTaskLabel(taskID int, labelID int) error {
	_, err := db.Exec(`
		DELETE FROM task_labels
		WHERE task_id = ?
		AND label_id = ?
	`,
		taskID,
		labelID,
	)

	return err
}
