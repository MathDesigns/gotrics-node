// cmd/node/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gotrics-node/internal/config"
	"gotrics-node/internal/metrics"
	"gotrics-node/internal/sender"
)

func main() {

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	interval := 5 * time.Second

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				metricsData, _ := metrics.CollectSystemMetrics(cfg.NodeID)
				err = sender.SendMetrics(cfg.ServerAddress, metricsData)
				if err != nil {
					log.Printf("Error sending metrics: %v", err)
				}
			}
		}
	}()

	// Wait for termination signal
	<-stop
	log.Println("Node shutting down gracefully")
}
