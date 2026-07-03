package main

import (
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func nullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return []byte("kore-wa-kimitsu-da")
	}

	return []byte(secret)
}

func generateJWT(user User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(getJWTSecret())
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func validateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method")
			}
			return getJWTSecret(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}

func projectBelongsToUser(projectID int, userID int) bool {
	var exists bool

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM projects
			WHERE id = ?
			AND user_id = ?
		)
	`, projectID, userID).Scan(&exists)

	return err == nil && exists
}

func listBelongsToUser(listID int, userID int) bool {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM lists l
			JOIN projects p
				ON l.project_id = p.id
			WHERE l.id = ?
			AND p.user_id = ?
		)
	`, listID, userID).Scan(&exists)

		return err == nil && exists
}
