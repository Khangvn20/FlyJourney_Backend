package main

import (
	"context"
	"github.com/Khangvn20/FlyJourney_Backend/internal/controller"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/common/router"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/server"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/service"
	"github.com/Khangvn20/FlyJourney_Backend/internal/infra/repository"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := repository.NewPgxDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	//init repository
	userRepo := repository.NewUserRepository(db)

	//init service
	userService := service.NewUserService(userRepo)
	//init controler
	userController := controller.NewUserController(userService)
	r := router.SetupRouter(userController)

	srv := server.New(r, 3000)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
