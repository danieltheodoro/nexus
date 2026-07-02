package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func lists(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetLists(w, r)

	case http.MethodPost:
		handleCreateList(w, r)

	case http.MethodPatch:
		handleUpdateList(w, r)

	case http.MethodDelete:
		handleDeleteList(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetLists(w http.ResponseWriter, _ *http.Request) {
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

func handleCreateList(w http.ResponseWriter, r *http.Request) {
	var req CreateListRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func handleUpdateList(w http.ResponseWriter, r *http.Request) {
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

	if err := updateList(id, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func handleDeleteList(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := deleteList(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func updateList(id int, list List) error {
	_, err := db.Exec(`
		UPDATE lists
		SET
			name = ?,
			position = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`,
		list.Name,
		list.Position,
		id,
	)

	return err
}

func deleteList(id int) error {
	_, err := db.Exec(`
						DELETE FROM lists
						WHERE id = ?`,
		id,
	)
	return err
}
