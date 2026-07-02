package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetUsers(w, r)

	case http.MethodPost:
		handleCreateUser(w, r)

	case http.MethodPatch:
		handleUpdateUser(w, r)

	case http.MethodDelete:
		handleDeleteUser(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

	}
}

func handleGetUsers(w http.ResponseWriter, _ *http.Request) {
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

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		AvatarURL:    req.AvatarURL,
	}

	if err := createUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleUpdateUser(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func handleDeleteUser(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
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

func createUser(user User) error {
	_, err := db.Exec(`
		INSERT INTO users (
			username,
			email,
			password_hash,
			first_name,
			last_name,
			avatar_url
		) VALUES (?, ?, ?, ?, ?, ?)
		`,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.AvatarURL,
	)

	return err
}
