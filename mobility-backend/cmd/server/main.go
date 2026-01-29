package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeng/mobility-backend/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	container, err := app.NewContainer(ctx)
	if err != nil {
		panic("bootstrap failed: " + err.Error())
	}

	container.Start(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	container.Shutdown(ctx)
}
