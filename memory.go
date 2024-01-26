package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
)

func memInfo() *MemoryInfo {
	var info MemoryInfo
	v, err := mem.VirtualMemory()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Stat = v
	}

	s, err := mem.SwapMemory()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Swap = s
	}

	return &info
}
