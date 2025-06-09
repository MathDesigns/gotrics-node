package main

import (
	"gotrics-node/internal/config"
	"gotrics-node/internal/hardware"
	"gotrics-node/internal/metrics"
	"gotrics-node/internal/sender"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting Gotrics agent...")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	httpSender := sender.NewSender(cfg.ServerURL, cfg.AuthToken)

	go func() {
		info, err := hardware.CollectInfo()
		if err != nil {
			log.Printf("Error collecting hardware info: %v", err)
			return
		}
		log.Println("Sending initial hardware information...")
		if err := httpSender.SendHardwareInfo(info); err != nil {
			log.Printf("Error sending hardware info: %v", err)
		}
	}()

	metricsTicker := time.NewTicker(cfg.PollInterval)
	defer metricsTicker.Stop()

	hardwareTicker := time.NewTicker(24 * time.Hour)
	defer hardwareTicker.Stop()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-metricsTicker.C:
			log.Println("Collecting metrics...")
			currentMetrics, err := metrics.Collect()
			if err != nil {
				log.Printf("Error collecting metrics: %v", err)
				continue
			}

			go func(m *metrics.Metrics) {
				if err := httpSender.Send(m); err != nil {
					log.Printf("Error sending metrics: %v", err)
				}
			}(currentMetrics)

		case <-hardwareTicker.C:
			go func() {
				info, err := hardware.CollectInfo()
				if err != nil {
					log.Printf("Error collecting hardware info for daily refresh: %v", err)
					return
				}
				log.Println("Performing daily refresh of hardware information...")
				if err := httpSender.SendHardwareInfo(info); err != nil {
					log.Printf("Error sending hardware info for daily refresh: %v", err)
				}
			}()

		case <-shutdown:
			log.Println("Shutdown signal received, stopping agent.")
			return
		}
	}
}
