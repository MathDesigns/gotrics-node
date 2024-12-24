package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Config holds the node configuration
type Config struct {
	ServerURL string `yaml:"server_url"`
	NodeID    string `yaml:"node_id"`
}

// Metrics represents the structure of metrics sent to the server
type Metrics struct {
	NodeID string  `json:"node_id"`
	CPU    float64 `json:"cpu"`
	Memory uint64  `json:"memory"`
}

func main() {
	// Load configuration
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Signal Handling for Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Main loop
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Collect metrics and send to server
				err := sendMetrics(cfg.ServerURL, cfg.NodeID)
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

// loadConfig reads the configuration from a YAML file
func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// sendMetrics collects system metrics and sends them to the server
func sendMetrics(serverURL, nodeID string) error {
	// Collect CPU usage
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}

	// Collect memory usage
	vm, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	// Create metrics payload
	metrics := Metrics{
		NodeID: nodeID,
		CPU:    cpuPercent[0],           // First value is the overall CPU usage
		Memory: vm.Used / (1024 * 1024), // Convert bytes to MB
	}

	// Marshal metrics to JSON
	payload, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	// Send metrics to the server
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Server responded with status: %s", resp.Status)
	}
	return nil
}
