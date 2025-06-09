package hardware

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type PartitionInfo struct {
	Device     string `json:"device"`
	Mountpoint string `json:"mountpoint"`
	Fstype     string `json:"fstype"`
	TotalGB    uint64 `json:"total_gb"`
}

type HardwareInfo struct {
	Hostname        string          `json:"hostname"`
	Platform        string          `json:"platform"`
	OS              string          `json:"os"`
	PlatformFamily  string          `json:"platform_family"`
	PlatformVersion string          `json:"platform_version"`
	KernelVersion   string          `json:"kernel_version"`
	CPUModel        string          `json:"cpu_model"`
	CPUCores        int32           `json:"cpu_cores"`
	CPULogicalCores int32           `json:"cpu_logical_cores"`
	TotalMemoryGB   uint64          `json:"total_memory_gb"`
	Partitions      []PartitionInfo `json:"partitions"`
}

func CollectInfo() (*HardwareInfo, error) {
	hostStat, err := host.Info()
	if err != nil {
		return nil, err
	}

	cpuStat, err := cpu.Info()
	if err != nil || len(cpuStat) == 0 {
		return nil, err
	}

	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var partitionInfos []PartitionInfo
	for _, p := range parts {
		usage, err := disk.Usage(p.Mountpoint)
		if err == nil {
			partitionInfos = append(partitionInfos, PartitionInfo{
				Device:     p.Device,
				Mountpoint: p.Mountpoint,
				Fstype:     p.Fstype,
				TotalGB:    usage.Total / 1024 / 1024 / 1024,
			})
		}
	}

	info := &HardwareInfo{
		Hostname:        hostStat.Hostname,
		Platform:        hostStat.Platform,
		OS:              hostStat.OS,
		PlatformFamily:  hostStat.PlatformFamily,
		PlatformVersion: hostStat.PlatformVersion,
		KernelVersion:   hostStat.KernelVersion,
		CPUModel:        cpuStat[0].ModelName,
		CPUCores:        cpuStat[0].Cores,
		CPULogicalCores: int32(len(cpuStat)),
		TotalMemoryGB:   memStat.Total / 1024 / 1024 / 1024,
		Partitions:      partitionInfos,
	}

	return info, nil
}
