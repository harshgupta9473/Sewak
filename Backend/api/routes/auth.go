package routes

import (
	"github.com/gorilla/mux"
	"github.com/harshgupta9473/sevak_backend/api/controllers"
)

// RegisterUserRoutes sets up the routes for user-related actions.
func RegisterUserRoutes(router *mux.Router, authController *controllers.AuthController) {
	// Route for handling login
	router.HandleFunc("/api/auth/login", authController.HandleLogin).Methods("POST")

	// Route for handling OTP verification
	router.HandleFunc("/api/auth/verify", authController.HandleVerification).Methods("POST")

	// Add more routes as needed, e.g., for refreshing tokens, user profile updates, etc.
}
