package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
)

type MemoryInfo struct {
	Stat *mem.VirtualMemoryStat
	Swap *mem.SwapMemoryStat
}

func memInfo() *MemoryInfo {
	var info MemoryInfo
	v, err := mem.VirtualMemory()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.Stat = v
	}

	s, err := mem.SwapMemory()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.Swap = s
	}

	return &info
}
