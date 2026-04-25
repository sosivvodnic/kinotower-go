package main

import (
	_ "github.com/lib/pq"
	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	core_logger "github.com/sosivvodnic/kinotower-go/internal/core/logger"
	core_server "github.com/sosivvodnic/kinotower-go/internal/core/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	// Use Overload to ensure values from .env override existing env vars.
	_ = godotenv.Overload()

	if err := core_logger.Init("logs"); err != nil {
		panic("failed to init logger: " + err.Error())
	}
	db, err := core_database.NewDatabase()
	if err != nil {
		core_logger.Log.Error("Failed to connect to database", "error", err)
		return
	}
	defer db.Close()
	core_logger.Log.Info("Starting server", "addr", ":8080")

	server := core_server.NewServer(*db)

	if err := server.ListenAndServe(); err != nil {
		core_logger.Log.Error("Server stopped", "error", err)
	}
}
