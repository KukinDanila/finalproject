package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum)
}

func generateToken(password string) (string, error) {
	hash := passwordHash(password)

	claims := jwt.MapClaims{
		"hash": hash,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(password))
}

func checkToken(tokenString, password string) bool {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(password), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	hash, ok := claims["hash"].(string)
	if !ok {
		return false
	}

	return hash == passwordHash(password)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")

	if req.Password != storedPassword {
		writeJson(w, map[string]string{"error": "неверный пароль"})
		return
	}

	token, err := generateToken(req.Password)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"token": token})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")

		if len(pass) > 0 {
			var jwt string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			valid := checkToken(jwt, pass)

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}
