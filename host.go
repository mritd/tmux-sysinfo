package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
	"os"
)

func hostInfo() *HostInfo {
	var info HostInfo

	i, err := host.Info()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info = HostInfo{i}
	}

	return &info
}
