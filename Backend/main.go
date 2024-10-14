package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/sevak_backend/api/controllers"
	"github.com/harshgupta9473/sevak_backend/api/routes"
	"github.com/harshgupta9473/sevak_backend/database"
	"github.com/harshgupta9473/sevak_backend/models"
)

func main() {
	database.InitDB()
	db := database.GetDB()
	models.InitUser(db)
	router := mux.NewRouter()

	auth := controllers.NewAuthController(db)
	routes.RegisterUserRoutes(router, auth)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Recieved terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)
}
