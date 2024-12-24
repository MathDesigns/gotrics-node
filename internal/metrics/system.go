package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type SystemMetrics struct {
	CPU    float64 `json:"cpu"`
	Memory uint64  `json:"memory"`
}

func CollectSystemMetrics() (*SystemMetrics, error) {
	cpuUsage, _ := cpu.Percent(0, false)
	memStats, _ := mem.VirtualMemory()

	return &SystemMetrics{
		CPU:    cpuUsage[0],
		Memory: memStats.Used,
	}, nil
}
