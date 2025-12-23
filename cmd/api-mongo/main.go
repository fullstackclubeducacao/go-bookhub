package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bookhub/internal/config"
	"bookhub/internal/infrastructure/auth"
	"bookhub/internal/infrastructure/database"
	apphttp "bookhub/internal/infrastructure/http"
	"bookhub/internal/infrastructure/http/handler"
	"bookhub/internal/infrastructure/repository"
	"bookhub/internal/usecase"
)

func main() {
	cfg := config.Load()

	mongoDB, err := database.NewMongoConnection(database.MongoConfig{
		URI:         cfg.MongoDB.URI,
		Database:    cfg.MongoDB.Database,
		MaxPoolSize: cfg.MongoDB.MaxPoolSize,
		MinPoolSize: cfg.MongoDB.MinPoolSize,
		MaxIdleTime: cfg.MongoDB.MaxIdleTime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	log.Println("Connected to MongoDB successfully")

	userRepo := repository.NewMongoUserRepository(mongoDB.Database)
	bookRepo := repository.NewMongoBookRepository(mongoDB.Database)
	loanRepo := repository.NewMongoLoanRepository(mongoDB.Database)

	userUseCase := usecase.NewUserUseCase(userRepo)
	bookUseCase := usecase.NewBookUseCase(bookRepo)
	loanUseCase := usecase.NewLoanUseCase(loanRepo, bookRepo, userRepo)

	jwtService := auth.NewJWTService(auth.JWTConfig{
		SecretKey:     cfg.JWT.SecretKey,
		TokenDuration: cfg.JWT.TokenDuration,
		Issuer:        cfg.JWT.Issuer,
	})

	h := handler.NewHandler(userUseCase, bookUseCase, loanUseCase, jwtService)
	router := apphttp.NewRouter(h)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Server starting on port %s (MongoDB backend)", cfg.Server.Port)
		log.Printf("Swagger UI available at http://localhost:%s/docs", cfg.Server.Port)
		log.Printf("OpenAPI spec available at http://localhost:%s/docs/swagger_spec", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
