// internal/core/server/http_server.go
package server

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
    "github.com/Khangvn20/FlyJourney_Backend/internal/controller"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/common/router"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/common/utils"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/service"
    "github.com/Khangvn20/FlyJourney_Backend/internal/infra/repository"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/common/middleware"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

type Server struct {
    Engine *gin.Engine
    Port   int
}

func NewHTTPServer(port int) (*Server, error) {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %v", err)
    }

    // Initialize database
    db, err := repository.NewPgxDatabase()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := db.Ping(ctx); err != nil {
        return nil, fmt.Errorf("database ping failed: %v", err)
    }

    // Initialize repository
    userRepo := repository.NewUserRepository(db)

    // Initialize services
    emailOTPService := service.NewEmailOTPService()
    tokenService := utils.NewTokenService()
    userService := service.NewUserService(userRepo, emailOTPService, tokenService)

    // Initialize controller
    userController := controller.NewUserController(userService)

    // Setup router
    r := gin.Default()
    r.Use(gin.Recovery())
    r.Use(gin.Logger())
    apiV1 := r.Group("/api/v1")
  
   router.UserRoutes(apiV1, userController, middleware.AuthMiddleware(tokenService))

    // Create server
    return &Server{
        Engine: r,
        Port:   port,
    }, nil
}

func (s *Server) Start() error {
    addr := fmt.Sprintf(":%d", s.Port)
    log.Printf("Starting HTTP server at http://localhost%s\n", addr)
    return http.ListenAndServe(addr, s.Engine)
}