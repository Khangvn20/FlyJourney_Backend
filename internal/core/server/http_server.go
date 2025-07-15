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
    "github.com/Khangvn20/FlyJourney_Backend/internal/infra/config"
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
    r:=gin.Default()
     r.Use(func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        allowedOrigins := []string{"http://localhost:5173", "http://localhost:5555"}

        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }
        
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    })

    // Initialize database
    db, err := repository.NewPgxDatabase()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }
    redisConfig := config.NewRedisConfig()
    redisClient, err := config.NewRedisClient(redisConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %v", err)
    }
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    pong, err := redisClient.Ping(ctx).Result()
    if err != nil {
        return nil, fmt.Errorf("redis ping failed: %v", err)
    }
    log.Printf("Redis connection successful: %s", pong)

    defer cancel()
    if err := db.Ping(ctx); err != nil {
        return nil, fmt.Errorf("database ping failed: %v", err)
    }
   //Intitiallize config

    // Initialize repository
    userRepo := repository.NewUserRepository(db)
    flightRepo := repository.NewFlightRepository(db.GetPool())
    // Initialize services
    redisService := service.NewRedisService(redisClient)
    
    emailOTPService := service.NewEmailOTPService()
     tokenService := utils.NewTokenService(redisService)
    userService := service.NewUserService(userRepo, emailOTPService, tokenService)
    flightService := service.NewFlightService(flightRepo)


    // Initialize controller
    userController := controller.NewUserController(userService)
    flightController := controller.NewFlightController(flightService)


    // Setup router
    r.Use(gin.Recovery())
    r.Use(gin.Logger())
    apiV1 := r.Group("/api/v1")
    router.AuthRoutes(apiV1, userController, middleware.AuthMiddleware(tokenService))
    router.UserRoutes(apiV1, userController, middleware.AuthMiddleware(tokenService))
    router.FlightRoutes(apiV1, flightController, middleware.AuthMiddleware(tokenService))
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