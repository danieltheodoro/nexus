package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := Task{
		ListID:      req.ListID,
		CreatorID:   req.CreatorID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Position:    req.Position,
		DueDate:     req.DueDate,
		CompletedAt: req.CompletedAt,
	}

	if err := createTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Position:    req.Position,
		DueDate:     req.DueDate,
		CompletedAt: req.CompletedAt,
	}

	if err := updateTask(id, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := deleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func createTask(task Task) error {
	_, err := db.Exec(`
		INSERT INTO tasks(
			list_id,
			creator_id,
			title,
			description,
			priority,
			position,
			due_date,
			completed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.ListID,
		task.CreatorID,
		task.Title,
		task.Description,
		task.Priority,
		task.Position,
		task.DueDate,
		task.CompletedAt,
	)

	return err
}

func updateTask(id int, task Task) error {
	_, err := db.Exec(`
		UPDATE tasks
		SET
			title = ?,
			description = ?,
			priority = ?,
			position = ?,
			due_date = ?,
			completed_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`,
		task.Title,
		task.Description,
		task.Priority,
		task.Position,
		task.DueDate,
		task.CompletedAt,
		id,
	)

	return err
}

func deleteTask(id int) error {
	_, err := db.Exec(`
						DELETE FROM tasks
						WHERE id = ?`,
		id,
	)
	return err
}
