package collector

import (
	"os"
	"testing"
	"time"
)

// benchDiskPath returns a platform-appropriate path for benchmarks
func benchDiskPath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	return os.TempDir()
}

func BenchmarkCollectHost(b *testing.B) {
	opts := Options{DiskUsagePath: "/", CPUInterval: time.Second}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectHost()
	}
}

func BenchmarkCollectCPU(b *testing.B) {
	opts := Options{
		PerCPU:      false,
		CPUInterval: 0, // No wait for benchmark
	}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectCPU()
	}
}

func BenchmarkCollectCPU_WithInterval(b *testing.B) {
	opts := Options{
		PerCPU:      false,
		CPUInterval: 100 * time.Millisecond,
	}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectCPU()
	}
}

func BenchmarkCollectCPU_PerCPU(b *testing.B) {
	opts := Options{
		PerCPU:      true,
		CPUInterval: 0,
	}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectCPU()
	}
}

func BenchmarkCollectMem(b *testing.B) {
	opts := Options{DiskUsagePath: "/", CPUInterval: time.Second}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectMem()
	}
}

func BenchmarkCollectLoad(b *testing.B) {
	opts := Options{DiskUsagePath: "/", CPUInterval: time.Second}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectLoad()
	}
}

func BenchmarkCollectDisk(b *testing.B) {
	opts := Options{
		DiskUsagePath: benchDiskPath(),
	}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectDisk()
	}
}

// BenchmarkCollectAll_Sequential benchmarks sequential collection
func BenchmarkCollectAll_Sequential(b *testing.B) {
	opts := Options{
		PerCPU:        false,
		DiskUsagePath: benchDiskPath(),
		CPUInterval:   0,
	}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectHost()
		mgr.collectCPU()
		mgr.collectMem()
		mgr.collectLoad()
		mgr.collectDisk()
	}
}

// BenchmarkCollectAll_Concurrent benchmarks concurrent collection
func BenchmarkCollectAll_Concurrent(b *testing.B) {
	opts := Options{
		PerCPU:        false,
		DiskUsagePath: benchDiskPath(),
		CPUInterval:   0,
	}
	mgr := NewManager(opts)
	names := []CollectorName{NameHost, NameCPU, NameMem, NameLoad, NameDisk}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Collect(names)
	}
}

// BenchmarkCollectAll_Concurrent_WithCPUInterval benchmarks with realistic CPU sampling
func BenchmarkCollectAll_Concurrent_WithCPUInterval(b *testing.B) {
	opts := Options{
		PerCPU:        false,
		DiskUsagePath: benchDiskPath(),
		CPUInterval:   100 * time.Millisecond,
	}
	mgr := NewManager(opts)
	names := []CollectorName{NameHost, NameCPU, NameMem, NameLoad, NameDisk}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Collect(names)
	}
}

// BenchmarkParseNames benchmarks name parsing
func BenchmarkParseNames(b *testing.B) {
	enabled := []string{"host", "cpu", "mem", "load", "disk"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseNames(enabled)
	}
}

// BenchmarkParseNames_All benchmarks "all" parsing
func BenchmarkParseNames_All(b *testing.B) {
	enabled := []string{"all"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseNames(enabled)
	}
}

// Lite mode benchmarks - skips expensive rarely-used data
func BenchmarkCollectMem_Lite(b *testing.B) {
	opts := Options{Lite: true}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectMem()
	}
}

func BenchmarkCollectLoad_Lite(b *testing.B) {
	opts := Options{Lite: true}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectLoad()
	}
}

func BenchmarkCollectDisk_Lite(b *testing.B) {
	opts := Options{
		DiskUsagePath: benchDiskPath(),
		Lite:          true,
	}
	mgr := NewManager(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.collectDisk()
	}
}

func BenchmarkCollectAll_Concurrent_Lite(b *testing.B) {
	opts := Options{
		PerCPU:        false,
		DiskUsagePath: benchDiskPath(),
		CPUInterval:   0,
		Lite:          true,
	}
	mgr := NewManager(opts)
	names := []CollectorName{NameHost, NameCPU, NameMem, NameLoad, NameDisk}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Collect(names)
	}
}
