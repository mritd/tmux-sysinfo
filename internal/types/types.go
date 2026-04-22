package types

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

// Info holds all collected system information
type Info struct {
	Host *HostInfo
	CPU  *CPUInfo
	Mem  *MemoryInfo
	Load *LoadInfo
	Disk *DiskInfo
}

// HostInfo contains host/OS information
type HostInfo struct {
	*host.InfoStat
}

// CPUInfo contains CPU statistics
type CPUInfo struct {
	Percent   []float64
	InfoStats []cpu.InfoStat
}

// MemoryInfo contains memory statistics
type MemoryInfo struct {
	Stat *mem.VirtualMemoryStat
	Swap *mem.SwapMemoryStat
}

// LoadInfo contains system load statistics
type LoadInfo struct {
	Stat *load.AvgStat
	Misc *load.MiscStat
}

// DiskInfo contains disk statistics
type DiskInfo struct {
	Stat     *disk.UsageStat
	Counters []*disk.IOCountersStat
}
