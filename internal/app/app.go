package app

import (
	"Authorization/internal/adapter/handler"
	"Authorization/internal/adapter/repository/postgres"
	"Authorization/internal/logger"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func Start() {
	if os.Getenv("RENDER") == "" { // если локально
		_ = godotenv.Load()
	}

	log := logger.SetupLogger()

	db, err := postgres.ConnectDB()
	if err != nil {
		log.Error("DB error", "error", err)
		return
	}
	defer db.Close()

	mux := http.NewServeMux()

	userRepo := postgres.NewPostgreUserRepository(db)
	articleRepo := postgres.NewPostgreArticleRepository(db)
	handler := handler.NewHandler(userRepo, log, articleRepo)

	setupRoutex(mux, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info("Server starting on port " + port) 
	if err := http.ListenAndServe("0.0.0.0:"+port, mux); err != nil {
		log.Error("Error in server", "error", err)
	}
}
