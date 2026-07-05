package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/arthurgray2k/gocsvtransformer/internal/cli"
)

func main() {
	// Setup context that can be cancelled via Ctrl+C
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	app := cli.NewApp()
	// Skip the binary name (os.Args[0])
	exitCode := app.Run(ctx, os.Args[1:])
	os.Exit(exitCode)
}
