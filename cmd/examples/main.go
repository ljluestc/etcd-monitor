// Package main demonstrates etcd patterns
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/etcd-monitor/taskmaster/pkg/examples"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	endpoints = flag.String("endpoints", "localhost:2379", "etcd endpoints")
	example   = flag.String("example", "all", "Example to run (all, lock, discovery, config, election, etc.)")
	nodeID    = flag.String("node-id", "node-1", "Node ID for election examples")
)

func main() {
	flag.Parse()

	// Create logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Connect to etcd
	client, err := connectEtcd()
	if err != nil {
		logger.Fatal("Failed to connect to etcd", zap.Error(err))
	}
	defer client.Close()

	logger.Info("Connected to etcd", zap.String("endpoints", *endpoints))

	// Run examples
	switch *example {
	case "all":
		examples.RunAllExamples(client, logger)

	case "lock":
		if err := examples.Example1_DistributedLock(client, logger); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "discovery":
		if err := examples.Example2_ServiceDiscovery(client, logger); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "loadbalancer", "lb":
		if err := examples.Example3_LoadBalancing(client, logger); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "config":
		if err := examples.Example4_ConfigurationManagement(client, logger); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "features", "flags":
		if err := examples.Example5_FeatureFlags(client, logger); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "election":
		if err := examples.Example6_LeaderElection(client, logger, *nodeID); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "leader-task":
		if err := examples.Example7_LeaderTask(client, logger, *nodeID); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "transaction", "txn":
		if err := examples.Example8_Transaction(client, logger); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "lease":
		if err := examples.Example9_Leases(client, logger); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "watch":
		if err := examples.Example10_Watch(client, logger); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	case "ha":
		if err := examples.Example11_HighAvailability(client, logger, *nodeID); err != nil {
			logger.Error("Example failed", zap.Error(err))
		}

	default:
		logger.Error("Unknown example", zap.String("example", *example))
		fmt.Println("\nAvailable examples:")
		fmt.Println("  all           - Run all examples")
		fmt.Println("  lock          - Distributed lock")
		fmt.Println("  discovery     - Service discovery")
		fmt.Println("  lb            - Load balancing")
		fmt.Println("  config        - Configuration management")
		fmt.Println("  flags         - Feature flags")
		fmt.Println("  election      - Leader election")
		fmt.Println("  leader-task   - Leader-only tasks")
		fmt.Println("  txn           - Transactions")
		fmt.Println("  lease         - Lease management")
		fmt.Println("  watch         - Watch for changes")
		fmt.Println("  ha            - High availability")
		os.Exit(1)
	}

	logger.Info("Example completed successfully")
}

func connectEtcd() (*clientv3.Client, error) {
	config := clientv3.Config{
		Endpoints:   []string{*endpoints},
		DialTimeout: 5 * time.Second,
	}

	return clientv3.New(config)
}
