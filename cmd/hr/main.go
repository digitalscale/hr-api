package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gpb.ru/hr/cmd/hr/app"
)

var version string = "unknown"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := app.NewDefaultCommand(version)
	done := make(chan struct{})
	go func() {
		err := cmd.ExecuteContext(ctx)
		if err != nil {
			log.Printf("[error] %v", err)
		}
		close(done)
		cancel()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigs)

	select {
	case sig := <-sigs:
		log.Printf("[info] receive signal %v", sig)
		cancel()
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	select {
	case <-done:
	case <-shutdownCtx.Done():
	}
}
