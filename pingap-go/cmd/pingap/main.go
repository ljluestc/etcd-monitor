package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/etcd-monitor/pingap-go/pkg/config"
	"github.com/etcd-monitor/pingap-go/pkg/proxy"
)

var (
	configFile = flag.String("config", "config.toml", "Configuration file path")
	version    = flag.Bool("version", false, "Show version information")
)

const (
	Version = "1.0.0"
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("pingap-go version %s\n", Version)
		return
	}

	// Load configuration
	loader := config.NewLoader(*configFile)
	cfg, err := loader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded configuration from %s\n", *configFile)
	fmt.Printf("Starting %s v%s\n", cfg.Basic.Name, Version)

	// Create proxy server
	server, err := proxy.NewServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh

	fmt.Println("\nShutting down gracefully...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracePeriod)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Shutdown error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped")
}
