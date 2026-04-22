package collector

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mritd/tmux-sysinfo/internal/types"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

// CollectorName defines the type for collector names
type CollectorName string

// Available collector names
const (
	NameHost CollectorName = "host"
	NameCPU  CollectorName = "cpu"
	NameMem  CollectorName = "mem"
	NameLoad CollectorName = "load"
	NameDisk CollectorName = "disk"
)

// Options holds configuration for collectors
type Options struct {
	PerCPU        bool
	DiskUsagePath string
	CPUInterval   time.Duration

	// Lite mode skips expensive but rarely-used data collection:
	// - load.Misc() (process counts, etc.)
	// - disk.IOCounters() (read/write stats)
	// - mem.SwapMemory()
	// This can significantly improve performance.
	Lite bool
}

// Collector defines the interface for system info collectors
type Collector interface {
	Name() CollectorName
	Collect() error
}

// Manager drives concurrent collection of selected system info sources.
type Manager struct {
	opts Options
}

// NewManager creates a new collector manager.
func NewManager(opts Options) *Manager {
	return &Manager{opts: opts}
}

// Collect runs the specified collectors concurrently and returns the aggregated
// Info. Duplicate names are collapsed so each collector runs at most once — this
// guarantees each goroutine writes a distinct field of *types.Info, which the
// Go memory model treats as independent memory locations; wg.Wait() provides
// the happens-before for the return.
func (m *Manager) Collect(names []CollectorName) *types.Info {
	info := &types.Info{}
	var wg sync.WaitGroup
	seen := make(map[CollectorName]struct{}, len(names))
	for _, name := range names {
		if _, dup := seen[name]; dup {
			continue
		}
		seen[name] = struct{}{}
		wg.Add(1)
		go func(n CollectorName) {
			defer wg.Done()
			switch n {
			case NameHost:
				info.Host = m.collectHost()
			case NameCPU:
				info.CPU = m.collectCPU()
			case NameMem:
				info.Mem = m.collectMem()
			case NameLoad:
				info.Load = m.collectLoad()
			case NameDisk:
				info.Disk = m.collectDisk()
			}
		}(name)
	}
	wg.Wait()
	return info
}

func (m *Manager) collectHost() *types.HostInfo {
	info := &types.HostInfo{}
	i, err := host.Info()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.InfoStat = i
	}
	return info
}

func (m *Manager) collectCPU() *types.CPUInfo {
	info := &types.CPUInfo{}

	i, err := cpu.Info()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.InfoStats = i
	}

	p, err := cpu.Percent(m.opts.CPUInterval, m.opts.PerCPU)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.Percent = p
	}

	return info
}

func (m *Manager) collectMem() *types.MemoryInfo {
	info := &types.MemoryInfo{}

	v, err := mem.VirtualMemory()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Stat = v
	}

	if !m.opts.Lite {
		s, err := mem.SwapMemory()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		} else {
			info.Swap = s
		}
	}

	return info
}

func (m *Manager) collectLoad() *types.LoadInfo {
	info := &types.LoadInfo{}

	a, err := load.Avg()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Stat = a
	}

	if !m.opts.Lite {
		mi, err := load.Misc()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		} else {
			info.Misc = mi
		}
	}

	return info
}

func (m *Manager) collectDisk() *types.DiskInfo {
	info := &types.DiskInfo{}

	u, err := disk.Usage(m.opts.DiskUsagePath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Stat = u
	}

	if !m.opts.Lite {
		i, err := disk.IOCounters()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		} else {
			for _, v := range i {
				info.Counters = append(info.Counters, v)
			}
		}
	}

	return info
}

// ParseNames converts string slice to CollectorName slice
// Returns (names, isAll) - isAll is true if "all" was specified
func ParseNames(enabled []string) ([]CollectorName, bool) {
	var names []CollectorName

	for _, en := range enabled {
		switch en {
		case "host":
			names = append(names, NameHost)
		case "cpu":
			names = append(names, NameCPU)
		case "mem", "memory":
			names = append(names, NameMem)
		case "load":
			names = append(names, NameLoad)
		case "disk":
			names = append(names, NameDisk)
		case "all":
			return []CollectorName{NameHost, NameCPU, NameMem, NameLoad, NameDisk}, true
		}
	}

	return names, false
}
