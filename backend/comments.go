package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func comments(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetComments(w, r, userID)

	case http.MethodPost:
		handleCreateComment(w, r, userID)

	case http.MethodPatch:
		handleUpdateComment(w, r, userID)

	case http.MethodDelete:
		handleDeleteComment(w, r, userID)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetComments(w http.ResponseWriter, r *http.Request, userID int) {
	idParam := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json")

	if idParam == "" {
		comments, err := getComments(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(comments); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	comment, err := getComment(id, userID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCreateComment(w http.ResponseWriter, r *http.Request, userID int) {
	var req CreateCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !taskBelongsToUser(req.TaskID, userID) {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	comment := Comment{
		TaskID:   req.TaskID,
		AuthorID: userID,
		Content:  req.Content,
	}

	if err := createComment(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleUpdateComment(w http.ResponseWriter, r *http.Request, userID int) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment := Comment{
		Content: req.Content,
	}

	if err := updateComment(id, userID, comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteComment(w http.ResponseWriter, r *http.Request, userID int) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := deleteComment(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getComments(userID int) ([]Comment, error) {
	comments := []Comment{}

	rows, err := db.Query(`
		SELECT
			c.id,
			c.task_id,
			c.author_id,
			c.content,
			c.created_at,
			c.updated_at
		FROM comments c
		JOIN tasks t
			ON c.task_id = t.id
		JOIN lists l
			ON t.list_id = l.id
		JOIN projects p
			ON l.project_id = p.id
		WHERE p.user_id = ?
	`, userID)
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

func getComment(id int, userID int) (Comment, error) {
	var comment Comment

	err := db.QueryRow(`
		SELECT
			c.id,
			c.task_id,
			c.author_id,
			c.content,
			c.created_at,
			c.updated_at
		FROM comments c
		JOIN tasks t
			ON c.task_id = t.id
		JOIN lists l
			ON t.list_id = l.id
		JOIN projects p
			ON l.project_id = p.id
		WHERE c.id = ?
		AND p.user_id = ?
	`, id, userID).Scan(
		&comment.ID,
		&comment.TaskID,
		&comment.AuthorID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func createComment(comment Comment) error {
	_, err := db.Exec(`
		INSERT INTO comments(
			task_id,
			author_id,
			content
		) VALUES (?, ?, ?)
	`,
		comment.TaskID,
		comment.AuthorID,
		comment.Content,
	)

	return err
}

func updateComment(id int, userID int, comment Comment) error {
	_, err := db.Exec(`
		UPDATE comments
		SET
			content = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		AND task_id IN (
			SELECT t.id
			FROM tasks t
			JOIN lists l
				ON t.list_id = l.id
			JOIN projects p
				ON l.project_id = p.id
			WHERE p.user_id = ?
		)
	`,
		comment.Content,
		id,
		userID,
	)

	return err
}

func deleteComment(id int, userID int) error {
	_, err := db.Exec(`
		DELETE FROM comments
		WHERE id = ?
		AND task_id IN (
			SELECT t.id
			FROM tasks t
			JOIN lists l
				ON t.list_id = l.id
			JOIN projects p
				ON l.project_id = p.id
			WHERE p.user_id = ?
		)
	`,
		id,
		userID,
	)

	return err
}
