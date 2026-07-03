package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func lists(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetLists(w, r, userID)

	case http.MethodPost:
		handleCreateList(w, r, userID)

	case http.MethodPatch:
		handleUpdateList(w, r, userID)

	case http.MethodDelete:
		handleDeleteList(w, r, userID)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetLists(w http.ResponseWriter, r *http.Request, userID int) {
	idParam := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json")

	if idParam == "" {
		lists, err := getLists(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(lists); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	list, err := getList(id, userID)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCreateList(w http.ResponseWriter, r *http.Request, userID int) {
	var req CreateListRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !projectBelongsToUser(req.ProjectID, userID) {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	list := List{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Position:  req.Position,
	}

	if err := createList(list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleUpdateList(w http.ResponseWriter, r *http.Request, userID int) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateListRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list := List{
		Name:     req.Name,
		Position: req.Position,
	}

	if err := updateList(id, userID, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteList(w http.ResponseWriter, r *http.Request, userID int) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := deleteList(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getLists(userID int) ([]List, error) {
	lists := []List{}

	rows, err := db.Query(`
		SELECT
			l.id,
			l.project_id,
			l.name,
			l.position,
			l.created_at,
			l.updated_at
		FROM lists as l
		JOIN projects AS p
			ON l.project_id = p.id
		WHERE p.user_id = ?		
	`,
		userID)
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

func getList(id int, userID int) (List, error) {
	list := List{}

	err := db.QueryRow(`

		SELECT
			l.id,
			l.project_id,
			l.name,
			l.position,
			l.created_at,
			l.updated_at
		FROM lists l
		JOIN projects p
			ON l.project_id = p.id
		WHERE l.id = ?
		AND p.user_id = ?
	`, id, userID).Scan(
		&list.ID,
		&list.ProjectID,
		&list.Name,
		&list.Position,
		&list.CreatedAt,
		&list.UpdatedAt,
	)
	if err != nil {
		return List{}, err
	}
	return list, nil

}

func createList(list List) error {
	_, err := db.Exec(`
		INSERT INTO lists(
			project_id,
			name,
			position
		) VALUES (?, ?, ?)
	`,
		list.ProjectID,
		list.Name,
		list.Position,
	)

	return err

}

func updateList(id int, userID int, list List) error {

	_, err := db.Exec(`
		UPDATE lists
		SET
			name = ?,
			position = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		AND project_id IN (
			SELECT id
			FROM projects
			WHERE user_id = ?
		)
	`,
		list.Name,
		list.Position,
		id,
		userID,
	)
	return err

}

func deleteList(id int, userID int) error {
	_, err := db.Exec(`
		DELETE FROM lists
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
