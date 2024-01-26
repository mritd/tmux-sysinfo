package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"os"
	"time"
)

func cpuInfo(percpu bool) *CPUInfo {
	var info CPUInfo

	i, err := cpu.Info()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.InfoStats = i
	}

	p, err := cpu.Percent(time.Second, percpu)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.Percent = p
	}

	return &info
}
