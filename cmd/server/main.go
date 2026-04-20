package main

import (
	"log/slog"

	core_server "github.com/Otvetov/kinotower-go/internal/core/server"
)

func main() {
	slog.Info("Starting server...")

	server := core_server.NewServer()
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server error", "error", err)
	}
}
