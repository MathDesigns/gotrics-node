package metrics

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SystemMetrics struct {
	NodeID      string `json:"node_id"`
	CPU         string `json:"cpu"`
	CPUModel    string `json:"cpu_model"`
	NumCores    int    `json:"num_cores"`
	NumThreads  int    `json:"num_threads"`
	UsedMemory  uint64 `json:"used_memory"`
	TotalMemory uint64 `json:"total_memory"`
	Uptime      uint64 `json:"uptime"`
	Platform    string `json:"platform"`
	UsedSpace   uint64 `json:"used_space"`
	TotalSpace  uint64 `json:"total_space"`
}

func CollectSystemMetrics(nodeid string) (*SystemMetrics, error) {
	cpuUsage, _ := cpu.Percent(0, false)
	cpuSpecs, _ := cpu.Info()
	memStats, _ := mem.VirtualMemory()
	uptime, _ := host.Uptime()
	hostInfo, _ := host.Info()
	diskStats, _ := disk.Usage("/")

	if len(cpuSpecs) > 0 {
		cpuInfo := cpuSpecs[0]
		numCores, _ := cpu.Counts(false)
		numThreads, _ := cpu.Counts(true)

		return &SystemMetrics{
			NodeID:      nodeid,
			CPU:         fmt.Sprintf("%f", cpuUsage[0]),
			CPUModel:    cpuInfo.ModelName,
			NumCores:    numCores,
			NumThreads:  numThreads,
			UsedMemory:  memStats.Used,
			TotalMemory: memStats.Total,
			Uptime:      uptime,
			Platform:    hostInfo.Platform,
			UsedSpace:   diskStats.Used,
			TotalSpace:  diskStats.Total,
		}, nil
	}

	return nil, fmt.Errorf("no CPU information available")
}
