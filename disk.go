package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"os"
)

func diskInfo(path string) *DiskInfo {
	var info DiskInfo
	u, err := disk.Usage(path)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Stat = u
	}

	i, err := disk.IOCounters()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		for _, v := range i {
			info.Counters = append(info.Counters, &v)
		}
	}

	return &info
}
