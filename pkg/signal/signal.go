package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Context creates a main context for application that is canceled when application is closed
func Context() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1)
	}()

	return ctx
}
