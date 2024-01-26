package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/load"
	"os"
)

type LoadInfo struct {
	Stat *load.AvgStat
	Misc *load.MiscStat
}

func loadInfo() *LoadInfo {
	var info LoadInfo
	a, err := load.Avg()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.Stat = a
	}

	m, err := load.Misc()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.Misc = m
	}

	return &info
}
