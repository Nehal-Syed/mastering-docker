package main

import (
    "log"
    "net/http"
    "os"

    "mastering-docker/internal/config"
    "mastering-docker/internal/database"
    "mastering-docker/internal/routes"
)

func main() {
    // Load config
    cfg := config.LoadConfig()

    // Initialize database
    db, err := database.InitDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Setup routes
    router := routes.SetupRoutes(db, cfg)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    instanceID := os.Getenv("INSTANCE_ID")
    if instanceID == "" {
        instanceID = "single"
    }

    log.Printf("Server %s starting on port %s", instanceID, port)
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}