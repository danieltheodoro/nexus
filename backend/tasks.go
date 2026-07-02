package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTasks(w, r)

	case http.MethodPost:
		handleCreateTask(w, r)

	case http.MethodPatch:
		handleUpdateTask(w, r)

	case http.MethodDelete:
		handleDeleteTask(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetTasks(w http.ResponseWriter, _ *http.Request) {
	tasks, err := getTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCreateTask(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleUpdateTask(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleDeleteTask(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func getTasks() ([]Task, error) {
	tasks := []Task{}

	rows, err := db.Query(`
		SELECT
			id,
			list_id,
			creator_id,
			title,
			description,
			priority,
			position,
			due_date,
			completed_at,
			created_at,
			updated_at
		FROM tasks
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		var description sql.NullString
		var dueDate sql.NullString
		var completedAt sql.NullString

		if err := rows.Scan(
			&task.ID,
			&task.ListID,
			&task.CreatorID,
			&task.Title,
			&description,
			&task.Priority,
			&task.Position,
			&dueDate,
			&completedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}

		task.Description = nullString(description)
		task.DueDate = nullString(dueDate)
		task.CompletedAt = nullString(completedAt)

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
