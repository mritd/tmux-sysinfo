package main

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type Info struct {
	CPU  *CPUInfo
	Mem  *MemoryInfo
	Load *LoadInfo
	Disk *DiskInfo
}

type CPUInfo struct {
	Percent   []float64
	InfoStats []cpu.InfoStat
}

type MemoryInfo struct {
	Stat *mem.VirtualMemoryStat
	Swap *mem.SwapMemoryStat
}

type LoadInfo struct {
	Stat *load.AvgStat
	Misc *load.MiscStat
}

type DiskInfo struct {
	Stat     *disk.UsageStat
	Counters []*disk.IOCountersStat
}
