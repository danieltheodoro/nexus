package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func comments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetComments(w, r)

	case http.MethodPost:
		handleCreateComment(w, r)

	case http.MethodPatch:
		handleUpdateComment(w, r)

	case http.MethodDelete:
		handleDeleteComment(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetComments(w http.ResponseWriter, _ *http.Request) {
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

func handleCreateComment(w http.ResponseWriter, r *http.Request) {
	var req CreateCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment := Comment{
		TaskID:   req.TaskID,
		AuthorID: req.AuthorID,
		Content:  req.Content,
	}

	if err := createComment(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleUpdateComment(w http.ResponseWriter, r *http.Request) {
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

	if err := updateComment(id, comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := deleteComment(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func updateComment(id int, comment Comment) error {
	_, err := db.Exec(`
		UPDATE comments
		SET
			content = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`,
		comment.Content,
		id,
	)

	return err
}

func deleteComment(id int) error {
	_, err := db.Exec(`
						DELETE FROM comments
						WHERE id = ?`,
		id,
	)
	return err
}
