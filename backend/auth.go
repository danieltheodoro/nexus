package main

import (
	"errors"
	"net/http"
	"strings"
)

func getTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("Missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Invalid authorization header")
	}

	return parts[1], nil
}

func authenticate(r *http.Request) (int, error) {

	tokenString, err := getTokenFromRequest(r)
	if err != nil {
		return 0, err
	}

	claims, err := validateJWT(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil

}
