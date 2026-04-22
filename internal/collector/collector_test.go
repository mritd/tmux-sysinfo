package collector

import (
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/mritd/tmux-sysinfo/internal/types"
)

// testDiskPath returns a platform-appropriate path for disk tests.
// Uses home directory which exists on all platforms.
func testDiskPath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	// Fallback to temp dir
	return os.TempDir()
}

func TestParseNames(t *testing.T) {
	tests := []struct {
		name          string
		enabled       []string
		expectedNames []CollectorName
		expectedAll   bool
	}{
		{
			name:          "single_host",
			enabled:       []string{"host"},
			expectedNames: []CollectorName{NameHost},
			expectedAll:   false,
		},
		{
			name:          "single_cpu",
			enabled:       []string{"cpu"},
			expectedNames: []CollectorName{NameCPU},
			expectedAll:   false,
		},
		{
			name:          "mem_alias",
			enabled:       []string{"memory"},
			expectedNames: []CollectorName{NameMem},
			expectedAll:   false,
		},
		{
			name:          "multiple",
			enabled:       []string{"cpu", "mem", "disk"},
			expectedNames: []CollectorName{NameCPU, NameMem, NameDisk},
			expectedAll:   false,
		},
		{
			name:          "all",
			enabled:       []string{"all"},
			expectedNames: []CollectorName{NameHost, NameCPU, NameMem, NameLoad, NameDisk},
			expectedAll:   true,
		},
		{
			name:          "all_with_others",
			enabled:       []string{"cpu", "all", "mem"},
			expectedNames: []CollectorName{NameHost, NameCPU, NameMem, NameLoad, NameDisk}, // "all" returns early
			expectedAll:   true,
		},
		{
			name:          "invalid",
			enabled:       []string{"invalid"},
			expectedNames: nil,
			expectedAll:   false,
		},
		{
			name:          "mixed_valid_invalid",
			enabled:       []string{"cpu", "invalid", "mem"},
			expectedNames: []CollectorName{NameCPU, NameMem},
			expectedAll:   false,
		},
		{
			name:          "empty",
			enabled:       []string{},
			expectedNames: nil,
			expectedAll:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			names, isAll := ParseNames(tt.enabled)

			if isAll != tt.expectedAll {
				t.Errorf("ParseNames() isAll = %v, want %v", isAll, tt.expectedAll)
			}

			if len(names) != len(tt.expectedNames) {
				t.Errorf("ParseNames() len = %d, want %d", len(names), len(tt.expectedNames))
				return
			}

			for i, name := range names {
				if name != tt.expectedNames[i] {
					t.Errorf("ParseNames()[%d] = %v, want %v", i, name, tt.expectedNames[i])
				}
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	diskPath := testDiskPath()
	opts := Options{
		PerCPU:        false,
		DiskUsagePath: diskPath,
		CPUInterval:   time.Second,
	}
	mgr := NewManager(opts)

	if mgr == nil {
		t.Fatal("NewManager returned nil")
	}
	if mgr.opts.DiskUsagePath != diskPath {
		t.Errorf("NewManager().opts.DiskUsagePath = %q, want %q", mgr.opts.DiskUsagePath, diskPath)
	}
}

func TestCollectHost(t *testing.T) {
	opts := Options{DiskUsagePath: "/", CPUInterval: time.Second}
	mgr := NewManager(opts)

	info := mgr.Collect([]CollectorName{NameHost})

	if info.Host == nil {
		t.Fatal("Collect(host) returned nil Host")
	}

	// Host.OS should be "darwin" on macOS or "linux" on Linux
	expectedOS := runtime.GOOS
	if info.Host.OS != expectedOS {
		t.Errorf("Host.OS = %q, want %q", info.Host.OS, expectedOS)
	}

	// KernelVersion should not be empty
	if info.Host.KernelVersion == "" {
		t.Error("Host.KernelVersion is empty")
	}
}

func TestCollectCPU(t *testing.T) {
	opts := Options{
		PerCPU:      false,
		CPUInterval: 100 * time.Millisecond, // Short interval for test
	}
	mgr := NewManager(opts)

	info := mgr.Collect([]CollectorName{NameCPU})

	if info.CPU == nil {
		t.Fatal("Collect(cpu) returned nil CPU")
	}

	// Should have at least one CPU info
	if len(info.CPU.InfoStats) == 0 {
		t.Error("CPU.InfoStats is empty")
	}

	// Should have percentage data
	if len(info.CPU.Percent) == 0 {
		t.Error("CPU.Percent is empty")
	}

	// Percentage should be between 0 and 100
	for i, p := range info.CPU.Percent {
		if p < 0 || p > 100 {
			t.Errorf("CPU.Percent[%d] = %v, want 0-100", i, p)
		}
	}
}

func TestCollectCPU_PerCPU(t *testing.T) {
	opts := Options{
		PerCPU:      true,
		CPUInterval: 100 * time.Millisecond,
	}
	mgr := NewManager(opts)

	info := mgr.Collect([]CollectorName{NameCPU})

	if info.CPU == nil {
		t.Fatal("Collect(cpu) returned nil CPU")
	}

	// With PerCPU=true, should have multiple percentage values (one per core)
	// At minimum, should match number of logical CPUs
	if len(info.CPU.Percent) < 1 {
		t.Error("CPU.Percent should have at least 1 entry with PerCPU=true")
	}
}

func TestCollectMem(t *testing.T) {
	opts := Options{DiskUsagePath: "/", CPUInterval: time.Second}
	mgr := NewManager(opts)

	info := mgr.Collect([]CollectorName{NameMem})

	if info.Mem == nil {
		t.Fatal("Collect(mem) returned nil Mem")
	}

	if info.Mem.Stat == nil {
		t.Fatal("Mem.Stat is nil")
	}

	// Total memory should be > 0
	if info.Mem.Stat.Total == 0 {
		t.Error("Mem.Stat.Total is 0")
	}

	// Used should be <= Total
	if info.Mem.Stat.Used > info.Mem.Stat.Total {
		t.Errorf("Mem.Stat.Used (%d) > Total (%d)", info.Mem.Stat.Used, info.Mem.Stat.Total)
	}

	// UsedPercent should be 0-100
	if info.Mem.Stat.UsedPercent < 0 || info.Mem.Stat.UsedPercent > 100 {
		t.Errorf("Mem.Stat.UsedPercent = %v, want 0-100", info.Mem.Stat.UsedPercent)
	}
}

func TestCollectLoad(t *testing.T) {
	opts := Options{DiskUsagePath: "/", CPUInterval: time.Second}
	mgr := NewManager(opts)

	info := mgr.Collect([]CollectorName{NameLoad})

	if info.Load == nil {
		t.Fatal("Collect(load) returned nil Load")
	}

	if info.Load.Stat == nil {
		t.Fatal("Load.Stat is nil")
	}

	// Load values should be >= 0
	if info.Load.Stat.Load1 < 0 {
		t.Errorf("Load.Stat.Load1 = %v, want >= 0", info.Load.Stat.Load1)
	}
	if info.Load.Stat.Load5 < 0 {
		t.Errorf("Load.Stat.Load5 = %v, want >= 0", info.Load.Stat.Load5)
	}
	if info.Load.Stat.Load15 < 0 {
		t.Errorf("Load.Stat.Load15 = %v, want >= 0", info.Load.Stat.Load15)
	}
}

func TestCollectDisk(t *testing.T) {
	diskPath := testDiskPath()
	opts := Options{
		DiskUsagePath: diskPath,
	}
	mgr := NewManager(opts)

	info := mgr.Collect([]CollectorName{NameDisk})

	if info.Disk == nil {
		t.Fatal("Collect(disk) returned nil Disk")
	}

	if info.Disk.Stat == nil {
		t.Fatal("Disk.Stat is nil")
	}

	// Total disk space should be > 0
	if info.Disk.Stat.Total == 0 {
		t.Error("Disk.Stat.Total is 0")
	}

	// Used should be <= Total
	if info.Disk.Stat.Used > info.Disk.Stat.Total {
		t.Errorf("Disk.Stat.Used (%d) > Total (%d)", info.Disk.Stat.Used, info.Disk.Stat.Total)
	}

	// UsedPercent should be 0-100
	if info.Disk.Stat.UsedPercent < 0 || info.Disk.Stat.UsedPercent > 100 {
		t.Errorf("Disk.Stat.UsedPercent = %v, want 0-100", info.Disk.Stat.UsedPercent)
	}

	// Path should be set (may differ from input due to symlink resolution)
	if info.Disk.Stat.Path == "" {
		t.Error("Disk.Stat.Path is empty")
	}
}

func TestCollectAll(t *testing.T) {
	opts := Options{
		PerCPU:        false,
		DiskUsagePath: testDiskPath(),
		CPUInterval:   100 * time.Millisecond,
	}
	mgr := NewManager(opts)

	allNames := []CollectorName{NameHost, NameCPU, NameMem, NameLoad, NameDisk}
	info := mgr.Collect(allNames)

	if info.Host == nil {
		t.Error("Collect(all) Host is nil")
	}
	if info.CPU == nil {
		t.Error("Collect(all) CPU is nil")
	}
	if info.Mem == nil {
		t.Error("Collect(all) Mem is nil")
	}
	if info.Load == nil {
		t.Error("Collect(all) Load is nil")
	}
	if info.Disk == nil {
		t.Error("Collect(all) Disk is nil")
	}
}

func TestCollectConcurrency(t *testing.T) {
	// Test that concurrent collection works correctly
	opts := Options{
		PerCPU:        false,
		DiskUsagePath: testDiskPath(),
		CPUInterval:   100 * time.Millisecond,
	}

	// Run multiple concurrent collections
	var wg sync.WaitGroup
	infos := make([]*types.Info, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			mgr := NewManager(opts)
			// NameCPU intentionally excluded: cpu.Percent blocks for CPUInterval
			// (100ms) so 10 concurrent invocations would add > 1 second to the test.
			infos[idx] = mgr.Collect([]CollectorName{NameHost, NameMem, NameLoad, NameDisk})
		}(i)
	}

	wg.Wait()

	// Verify all results are valid
	for i, info := range infos {
		if info.Host == nil {
			t.Errorf("Concurrent collection %d: Host is nil", i)
		}
		if info.Mem == nil {
			t.Errorf("Concurrent collection %d: Mem is nil", i)
		}
	}
}

func TestCollect_DuplicateNames(t *testing.T) {
	// Regression: if two entries in names map to the same CollectorName
	// (which is possible from ParseNames since "mem" and "memory" both
	// yield NameMem), Collect must still return a valid Info without
	// racing on a shared field.
	opts := Options{
		PerCPU:        false,
		DiskUsagePath: testDiskPath(),
		CPUInterval:   100 * time.Millisecond,
	}
	mgr := NewManager(opts)

	info := mgr.Collect([]CollectorName{NameMem, NameMem, NameMem})

	if info.Mem == nil {
		t.Fatal("Collect with duplicate NameMem returned nil Mem")
	}
	if info.Mem.Stat == nil {
		t.Error("Mem.Stat is nil after duplicate Collect")
	}
}

func TestCollectEmpty(t *testing.T) {
	opts := Options{DiskUsagePath: "/", CPUInterval: time.Second}
	mgr := NewManager(opts)

	info := mgr.Collect([]CollectorName{})

	// All fields should be nil when no collectors specified
	if info.Host != nil {
		t.Error("Collect([]) Host should be nil")
	}
	if info.CPU != nil {
		t.Error("Collect([]) CPU should be nil")
	}
	if info.Mem != nil {
		t.Error("Collect([]) Mem should be nil")
	}
	if info.Load != nil {
		t.Error("Collect([]) Load should be nil")
	}
	if info.Disk != nil {
		t.Error("Collect([]) Disk should be nil")
	}
}

func TestCollectorNameConstants(t *testing.T) {
	// Ensure constants have expected values
	if NameHost != "host" {
		t.Errorf("NameHost = %q, want %q", NameHost, "host")
	}
	if NameCPU != "cpu" {
		t.Errorf("NameCPU = %q, want %q", NameCPU, "cpu")
	}
	if NameMem != "mem" {
		t.Errorf("NameMem = %q, want %q", NameMem, "mem")
	}
	if NameLoad != "load" {
		t.Errorf("NameLoad = %q, want %q", NameLoad, "load")
	}
	if NameDisk != "disk" {
		t.Errorf("NameDisk = %q, want %q", NameDisk, "disk")
	}
}
