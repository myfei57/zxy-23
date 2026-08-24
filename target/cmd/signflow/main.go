package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"signflow/internal/console"
	"signflow/internal/settings"
)

func main() {
	var (
		addr string
		root string
	)
	flag.StringVar(&addr, "addr", ":8080", "HTTP listen address")
	flag.StringVar(&root, "root", "signflow-data", "data directory for file persistence")
	flag.Parse()

	cfg := settings.Default()
	cfg = cfg.WithRoot(root)

	server, err := console.NewServer(cfg)
	if err != nil {
		log.Fatalf("signflow: initialize console server: %v", err)
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("signflow: listening on %s with data root %s", addr, root)
		errCh <- httpServer.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("signflow: server failed: %v", err)
		}
	case <-stop:
		log.Printf("signflow: shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("signflow: shutdown: %v", err)
		}
	}

	if err := server.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "signflow: close storage: %v\n", err)
	}
}
