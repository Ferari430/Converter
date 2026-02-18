package main

import (
	"context"
	"converter/cmd/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	application := app.NewApp()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := application.Start(ctx); err != nil {
		log.Fatal("Failed to start application:", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal")

	cancel()

	if err := application.Stop(); err != nil {
		log.Printf("Error stopping application: %v", err)
	}

	log.Println("Application stopped")
}
