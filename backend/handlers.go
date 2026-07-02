package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func nullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status": "ok",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func users(w http.ResponseWriter, _ *http.Request) {
	users, err := getUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getUsers() ([]User, error) {
	users := []User{}

	rows, err := db.Query(`
		SELECT 
			id, 
			username, 
			email, 
			first_name, 
			last_name,
			avatar_url,
			is_active,
			created_at,
			updated_at 
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		var firstName sql.NullString
		var lastName sql.NullString
		var avatarURL sql.NullString

		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&firstName,
			&lastName,
			&avatarURL,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		user.FirstName = nullString(firstName)
		user.LastName = nullString(lastName)
		user.AvatarURL = nullString(avatarURL)

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func projects(w http.ResponseWriter, _ *http.Request) {
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

func lists(w http.ResponseWriter, _ *http.Request) {
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

func tasks(w http.ResponseWriter, _ *http.Request) {
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

func labels(w http.ResponseWriter, _ *http.Request) {
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

func comments(w http.ResponseWriter, _ *http.Request) {
	comments, err := getComments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getComments() ([]Comment, error) {
	comments := []Comment{}

	rows, err := db.Query(`
		SELECT
			id,
			task_id,
			author_id,
			content,
			created_at,
			updated_at
		FROM comments
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment

		if err := rows.Scan(
			&comment.ID,
			&comment.TaskID,
			&comment.AuthorID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
