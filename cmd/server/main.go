package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Choikyudong/NuriOS/internal/auth"
	"github.com/Choikyudong/NuriOS/internal/config"
	"github.com/Choikyudong/NuriOS/internal/db"
	"github.com/Choikyudong/NuriOS/internal/todo"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("warning: .env not loaded:", err)
	}
	log.Println("DATABASE_URL =", os.Getenv("DATABASE_URL"))

	cfg := config.Load()

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)
	authRepo := auth.NewRepository(pool)
	authService := auth.NewService(authRepo, jwtManager)
	authHandler := auth.NewHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/signup", authHandler.Signup)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	todoRepo := todo.NewRepository(pool)
	todoService := todo.NewService(todoRepo)
	todoHandler := todo.NewHandler(todoService)

	mux.HandleFunc("GET /todos", jwtManager.RequireAuth(todoHandler.List))
	mux.HandleFunc("POST /todos", jwtManager.RequireAuth(todoHandler.Create))
	mux.HandleFunc("PATCH /todos/{id}", jwtManager.RequireAuth(todoHandler.Update))
	mux.HandleFunc("DELETE /todos/{id}", jwtManager.RequireAuth(todoHandler.Delete))

	log.Printf("server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
