// main.go
package main

import (
    "log" 
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/server"
)

func main() {
    // Initialize HTTP server
    srv, err := server.NewHTTPServer(3000)
    if err != nil {
        log.Fatalf("Failed to initialize server: %v", err)
    }

    // Start server
    if err := srv.Start(); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}