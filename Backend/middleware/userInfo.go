package middleware

import (
	"net/http"

	"github.com/harshgupta9473/sevak_backend/database"
	"github.com/harshgupta9473/sevak_backend/models"
)

type contextKey string

func InfoMiddleware(next http.Handler) http.Handler {
	const userContextKey contextKey = "userInfo"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(userContextKey).(map[string]interface{})
		if !ok || userInfo == nil {
			http.Error(w, "Unauthorized: Could not retrieve user info", http.StatusUnauthorized)
			return
		}
		db := database.GetDB()
		var user models.User

		err := db.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Joins("JOIN roles ON roles.id = user_roles.role_id").
			Where("users.mobile = ? AND users.created_at = ? AND roles.name = ?", userInfo["number"], userInfo["created_at"], userInfo["role"]).
			First(&user).Error

		if err != nil {
			http.Error(w,"unauthorised",http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w,r)
	})
}
