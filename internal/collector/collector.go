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

// DefaultOptions returns default collector options
func DefaultOptions() Options {
	return Options{
		PerCPU:        false,
		DiskUsagePath: "/",
		CPUInterval:   time.Second,
		Lite:          false,
	}
}

// Collector defines the interface for system info collectors
type Collector interface {
	Name() CollectorName
	Collect() error
}

// Manager manages multiple collectors and runs them concurrently
type Manager struct {
	opts Options
	info *types.Info
	mu   sync.Mutex
}

// NewManager creates a new collector manager
func NewManager(opts Options) *Manager {
	return &Manager{
		opts: opts,
		info: &types.Info{},
	}
}

// Collect runs the specified collectors concurrently
func (m *Manager) Collect(names []CollectorName) *types.Info {
	var wg sync.WaitGroup

	for _, name := range names {
		wg.Add(1)
		go func(n CollectorName) {
			defer wg.Done()
			m.collect(n)
		}(name)
	}

	wg.Wait()
	return m.info
}

func (m *Manager) collect(name CollectorName) {
	switch name {
	case NameHost:
		m.collectHost()
	case NameCPU:
		m.collectCPU()
	case NameMem:
		m.collectMem()
	case NameLoad:
		m.collectLoad()
	case NameDisk:
		m.collectDisk()
	}
}

func (m *Manager) collectHost() {
	info := &types.HostInfo{}
	i, err := host.Info()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	} else {
		info.InfoStat = i
	}

	m.mu.Lock()
	m.info.Host = info
	m.mu.Unlock()
}

func (m *Manager) collectCPU() {
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

	m.mu.Lock()
	m.info.CPU = info
	m.mu.Unlock()
}

func (m *Manager) collectMem() {
	info := &types.MemoryInfo{}

	v, err := mem.VirtualMemory()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Stat = v
	}

	// Skip swap memory in lite mode (rarely used in templates)
	if !m.opts.Lite {
		s, err := mem.SwapMemory()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		} else {
			info.Swap = s
		}
	}

	m.mu.Lock()
	m.info.Mem = info
	m.mu.Unlock()
}

func (m *Manager) collectLoad() {
	info := &types.LoadInfo{}

	a, err := load.Avg()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Stat = a
	}

	// Skip load.Misc() in lite mode - it's expensive (~20ms) and rarely used
	if !m.opts.Lite {
		mi, err := load.Misc()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		} else {
			info.Misc = mi
		}
	}

	m.mu.Lock()
	m.info.Load = info
	m.mu.Unlock()
}

func (m *Manager) collectDisk() {
	info := &types.DiskInfo{}

	u, err := disk.Usage(m.opts.DiskUsagePath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	} else {
		info.Stat = u
	}

	// Skip disk IO counters in lite mode - rarely used in status bar templates
	if !m.opts.Lite {
		i, err := disk.IOCounters()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		} else {
			for _, v := range i {
				info.Counters = append(info.Counters, &v)
			}
		}
	}

	m.mu.Lock()
	m.info.Disk = info
	m.mu.Unlock()
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
