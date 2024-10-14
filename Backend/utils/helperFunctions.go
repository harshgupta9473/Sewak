package utils

import (
	"encoding/json"
	"net/http"

	"github.com/harshgupta9473/sevak_backend/models"
)

func IsRoleAlloted(roles []*models.Role, role string) bool {
	for _, v := range roles {
		if v.Name == role {
			return true
		}
	}
	return false
}

type contextKey string


func IsAdmin(role string) bool {
	return role == "admin"
}

func IsRoleMatches(role string, r *http.Request) bool{
	const userContextKey contextKey = "userInfo"
	userInfo, ok := r.Context().Value(userContextKey).(map[string]interface{})
		if !ok || userInfo == nil {
			return false
		}
		return role==userInfo["role"]
}

func WriteJSON(w http.ResponseWriter,status int ,msg any)error{
	w.Header().Set(`Content-Type`,`application/json`)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(msg)
}

func SetAuthTokenCookie(w http.ResponseWriter, tokenValue, cookieName string) {
    http.SetCookie(w, &http.Cookie{
        Name:     cookieName,
        Value:    tokenValue,
        Path:     "/",
        HttpOnly: true,
        // Optionally, add attributes like `Secure` if your app runs over HTTPS.
        // Secure: true,
    })
}