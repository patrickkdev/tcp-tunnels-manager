package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/patrickkdev/tcptunnel/internal/application"
	"github.com/patrickkdev/tcptunnel/internal/infrastructure/db"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer stop()

	repo := db.NewTunnelRowsRepo()

	manager := application.NewManager(repo)
	manager.Run(ctx, 5*time.Second)
}