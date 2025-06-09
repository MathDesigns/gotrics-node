package metrics

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct {
	Hostname        string    `json:"hostname"`
	Timestamp       time.Time `json:"timestamp"`
	CPUUsagePercent float64   `json:"cpu_usage_percent"`
	MemoryUsedMB    uint64    `json:"memory_used_mb"`
	DiskUsedPercent float64   `json:"disk_used_percent"`
}

func Collect() (*Metrics, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	cpuPercentages, err := cpu.Percent(time.Second, false)
	if err != nil || len(cpuPercentages) == 0 {
		return nil, err
	}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	m := &Metrics{
		Hostname:        hostInfo.Hostname,
		Timestamp:       time.Now().UTC(),
		CPUUsagePercent: cpuPercentages[0],
		MemoryUsedMB:    vmStat.Used / 1024 / 1024, // Convert bytes to MB
		DiskUsedPercent: diskStat.UsedPercent,
	}

	return m, nil
}
