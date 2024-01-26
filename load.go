package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/load"
	"os"
)

func loadInfo() *LoadInfo {
	var info LoadInfo
	a, err := load.Avg()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Stat = a
	}

	m, err := load.Misc()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Misc = m
	}

	return &info
}
