package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func tasks(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetTasks(w, r, userID)
	case http.MethodPost:
		handleCreateTask(w, r, userID)
	case http.MethodPatch:
		handleUpdateTask(w, r, userID)
	case http.MethodDelete:
		handleDeleteTask(w, r, userID)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetTasks(w http.ResponseWriter, r *http.Request, userID int) {
	idParam := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json")

	if idParam == "" {
		tasks, err := getTasks(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tasks)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	task, err := getTask(id, userID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func handleCreateTask(w http.ResponseWriter, r *http.Request, userID int) {
	var req CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !listBelongsToUser(req.ListID, userID) {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	task := Task{
		ListID:      req.ListID,
		CreatorID:   userID,
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

func handleUpdateTask(w http.ResponseWriter, r *http.Request, userID int) {
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

	if err := updateTask(id, userID, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request, userID int) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := deleteTask(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getTasks(userID int) ([]Task, error) {
	tasks := []Task{}

	rows, err := db.Query(`
		SELECT
			t.id,
			t.list_id,
			t.creator_id,
			t.title,
			t.description,
			t.priority,
			t.position,
			t.due_date,
			t.completed_at,
			t.created_at,
			t.updated_at
		FROM tasks t
		JOIN lists l ON t.list_id = l.id
		JOIN projects p ON l.project_id = p.id
		WHERE p.user_id = ?
	`, userID)
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

func getTask(id int, userID int) (Task, error) {
	var task Task
	var description sql.NullString
	var dueDate sql.NullString
	var completedAt sql.NullString

	err := db.QueryRow(`
		SELECT
			t.id,
			t.list_id,
			t.creator_id,
			t.title,
			t.description,
			t.priority,
			t.position,
			t.due_date,
			t.completed_at,
			t.created_at,
			t.updated_at
		FROM tasks t
		JOIN lists l ON t.list_id = l.id
		JOIN projects p ON l.project_id = p.id
		WHERE t.id = ?
		AND p.user_id = ?
	`, id, userID).Scan(
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
	)
	if err != nil {
		return Task{}, err
	}

	task.Description = nullString(description)
	task.DueDate = nullString(dueDate)
	task.CompletedAt = nullString(completedAt)

	return task, nil
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

func updateTask(id int, userID int, task Task) error {
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
		AND list_id IN (
			SELECT l.id
			FROM lists l
			JOIN projects p ON l.project_id = p.id
			WHERE p.user_id = ?
		)
	`,
		task.Title,
		task.Description,
		task.Priority,
		task.Position,
		task.DueDate,
		task.CompletedAt,
		id,
		userID,
	)

	return err
}

func deleteTask(id int, userID int) error {
	_, err := db.Exec(`
		DELETE FROM tasks
		WHERE id = ?
		AND list_id IN (
			SELECT l.id
			FROM lists l
			JOIN projects p ON l.project_id = p.id
			WHERE p.user_id = ?
		)
	`,
		id,
		userID,
	)

	return err
}
