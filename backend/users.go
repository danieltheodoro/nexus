package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func users(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodPost:
		handleCreateUser(w, r)
		return
	}

	userID, err := authenticate(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetUsers(w, r, userID)

	case http.MethodPatch:
		handleUpdateUser(w, r, userID)

	case http.MethodDelete:
		handleDeleteUser(w, r, userID)

	default:
		http.Error(w, "Method NotAllowed", http.StatusMethodNotAllowed)
	}
}

func handleGetUsers(w http.ResponseWriter, r *http.Request, userID int) {
	user, err := getUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	var userID int64

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

	userID, err = createUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"id": userID,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, userID int) {
	id := userID

	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var currentHash string
	var err error

	err = db.QueryRow(
		`SELECT password_hash FROM users WHERE id = ?`,
		id,
	).Scan(&currentHash)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !checkPassword(req.Password, currentHash) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	passwordHash := currentHash

	if req.NewPassword != "" {
		passwordHash, err = hashPassword(req.NewPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	user := User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		AvatarURL:    req.AvatarURL,
		IsActive:     req.IsActive,
	}

	if err := updateUser(id, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, userID int) {
	var err error

	var req DeleteUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var currentHash string

	err = db.QueryRow(
		`SELECT password_hash FROM users WHERE id = ?`,
		userID,
	).Scan(&currentHash)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !checkPassword(req.Password, currentHash) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	if err := deleteUser(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func getUser(userID int) (User, error) {
	var user User
	var firstName sql.NullString
	var lastName sql.NullString
	var avatarURL sql.NullString

	err := db.QueryRow(`
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
		WHERE id = ?
	`, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&firstName,
		&lastName,
		&avatarURL,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	user.FirstName = nullString(firstName)
	user.LastName = nullString(lastName)
	user.AvatarURL = nullString(avatarURL)

	return user, nil
}

func getUserByUsername(username string) (User, error) {
	var user User

	err := db.QueryRow(`
		SELECT
		id,
		username,
		email,
		password_hash
		FROM users
		WHERE username = ?
		`,
		username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
	)

	return user, err
}

func createUser(user User) (int64, error) {
	result, err := db.Exec(`
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
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func updateUser(id int, user User) error {
	_, err := db.Exec(`
		UPDATE users
		SET
			username = ?,
			email = ?,
			password_hash = ?,
			first_name = ?,
			last_name = ?,
			avatar_url = ?,
			is_active = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.AvatarURL,
		user.IsActive,
		id,
	)

	return err
}

func deleteUser(id int) error {
	_, err := db.Exec(`DELETE FROM
							users WHERE
							id = ?`,
		id,
	)
	return err
}
