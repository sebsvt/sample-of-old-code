package handler

import (
	"errors"
	"strings"
)

func ExtractToken(token string) (string, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(token, bearerPrefix) {
		return "", errors.New("invalid token format")
	}

	// Extract the token without the Bearer prefix
	tokenString := token[len(bearerPrefix):]
	return tokenString, nil
}
