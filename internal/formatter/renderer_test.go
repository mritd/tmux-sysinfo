package formatter

import (
	"strings"
	"testing"

	"github.com/mritd/tmux-sysinfo/internal/collector"
	"github.com/mritd/tmux-sysinfo/internal/types"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

func newTestInfo() *types.Info {
	return &types.Info{
		Host: &types.HostInfo{
			InfoStat: &host.InfoStat{
				OS:            "linux",
				KernelVersion: "5.15.0",
				Hostname:      "testhost",
			},
		},
		CPU: &types.CPUInfo{
			Percent: []float64{25.5},
			InfoStats: []cpu.InfoStat{
				{ModelName: "Test CPU"},
			},
		},
		Mem: &types.MemoryInfo{
			Stat: &mem.VirtualMemoryStat{
				Total:       16 * 1024 * 1024 * 1024, // 16 GB
				Used:        8 * 1024 * 1024 * 1024,  // 8 GB
				UsedPercent: 50.0,
			},
		},
		Load: &types.LoadInfo{
			Stat: &load.AvgStat{
				Load1:  1.5,
				Load5:  2.0,
				Load15: 1.8,
			},
		},
		Disk: &types.DiskInfo{
			Stat: &disk.UsageStat{
				Path:        "/",
				Total:       500 * 1024 * 1024 * 1024, // 500 GB
				Used:        250 * 1024 * 1024 * 1024, // 250 GB
				UsedPercent: 50.0,
			},
		},
	}
}

func TestNewRenderer(t *testing.T) {
	tpls := DefaultTemplates()
	funcMap := NewFuncMapBuilder().Build()
	renderer := NewRenderer(tpls, funcMap, "|")

	// Verify renderer is created with correct config
	if renderer.delimiter != "|" {
		t.Errorf("renderer.delimiter = %q, want %q", renderer.delimiter, "|")
	}

	if renderer.templates.Host != DefaultHostTpl {
		t.Errorf("renderer.templates.Host = %q, want %q", renderer.templates.Host, DefaultHostTpl)
	}
}

func TestRender_SingleCollector(t *testing.T) {
	funcMap := NewFuncMapBuilder().Build()
	info := newTestInfo()

	tests := []struct {
		name     string
		tpls     Templates
		names    []collector.CollectorName
		contains []string
	}{
		{
			name:     "host_only",
			tpls:     DefaultTemplates(),
			names:    []collector.CollectorName{collector.NameHost},
			contains: []string{"HOST:", "linux", "5.15.0"},
		},
		{
			name:     "cpu_only",
			tpls:     DefaultTemplates(),
			names:    []collector.CollectorName{collector.NameCPU},
			contains: []string{"CPU:", "Test CPU", "26%"}, // 25.5 rounds to 26
		},
		{
			name:     "mem_only",
			tpls:     DefaultTemplates(),
			names:    []collector.CollectorName{collector.NameMem},
			contains: []string{"MEM:", "GB"},
		},
		{
			name:     "load_only",
			tpls:     DefaultTemplates(),
			names:    []collector.CollectorName{collector.NameLoad},
			contains: []string{"LOAD:", "2%"}, // 1.5 rounds to 2
		},
		{
			name:     "disk_only",
			tpls:     DefaultTemplates(),
			names:    []collector.CollectorName{collector.NameDisk},
			contains: []string{"DISK:", "50%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(tt.tpls, funcMap, "|")
			output, err := renderer.Render(info, tt.names)

			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			for _, substr := range tt.contains {
				if !strings.Contains(output, substr) {
					t.Errorf("Render() output %q does not contain %q", output, substr)
				}
			}
		})
	}
}

func TestRender_MultipleCollectors(t *testing.T) {
	funcMap := NewFuncMapBuilder().Build()
	tpls := DefaultTemplates()
	renderer := NewRenderer(tpls, funcMap, "|")
	info := newTestInfo()

	names := []collector.CollectorName{collector.NameCPU, collector.NameMem}
	output, err := renderer.Render(info, names)

	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Should contain delimiter
	if !strings.Contains(output, " | ") {
		t.Errorf("Render() output %q does not contain delimiter", output)
	}

	// Should contain both CPU and MEM
	if !strings.Contains(output, "CPU:") {
		t.Errorf("Render() output %q does not contain CPU:", output)
	}
	if !strings.Contains(output, "MEM:") {
		t.Errorf("Render() output %q does not contain MEM:", output)
	}
}

func TestRender_CustomDelimiter(t *testing.T) {
	funcMap := NewFuncMapBuilder().Build()
	tpls := DefaultTemplates()
	renderer := NewRenderer(tpls, funcMap, ":::")
	info := newTestInfo()

	names := []collector.CollectorName{collector.NameCPU, collector.NameMem}
	output, err := renderer.Render(info, names)

	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(output, " ::: ") {
		t.Errorf("Render() output %q does not contain custom delimiter ' ::: '", output)
	}
}

func TestRender_MiniTemplates(t *testing.T) {
	funcMap := NewFuncMapBuilder().Build()
	tpls := MiniTemplates()
	renderer := NewRenderer(tpls, funcMap, "|")
	info := newTestInfo()

	names := []collector.CollectorName{collector.NameHost, collector.NameCPU, collector.NameMem}
	output, err := renderer.Render(info, names)

	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Mini templates use short prefixes
	if !strings.Contains(output, "S:") {
		t.Errorf("Render() mini output %q does not contain 'S:'", output)
	}
	if !strings.Contains(output, "C:") {
		t.Errorf("Render() mini output %q does not contain 'C:'", output)
	}
	if !strings.Contains(output, "M:") {
		t.Errorf("Render() mini output %q does not contain 'M:'", output)
	}
}

func TestRender_EmptyCollectors(t *testing.T) {
	funcMap := NewFuncMapBuilder().Build()
	tpls := DefaultTemplates()
	renderer := NewRenderer(tpls, funcMap, "|")
	info := newTestInfo()

	output, err := renderer.Render(info, []collector.CollectorName{})

	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if output != "" {
		t.Errorf("Render() with empty collectors = %q, want empty string", output)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	funcMap := NewFuncMapBuilder().Build()
	tpls := Templates{
		CPU: `{{.Invalid.Field}}`, // Invalid field
	}
	renderer := NewRenderer(tpls, funcMap, "|")
	info := newTestInfo()

	_, err := renderer.Render(info, []collector.CollectorName{collector.NameCPU})

	if err == nil {
		t.Error("Render() with invalid template should return error")
	}
}

func TestRender_TemplateSyntaxError(t *testing.T) {
	funcMap := NewFuncMapBuilder().Build()
	tpls := Templates{
		CPU: `{{.CPU.Percent | invalid_func}}`, // Unknown function
	}
	renderer := NewRenderer(tpls, funcMap, "|")
	info := newTestInfo()

	_, err := renderer.Render(info, []collector.CollectorName{collector.NameCPU})

	if err == nil {
		t.Error("Render() with template syntax error should return error")
	}
}

func TestDefaultTemplates(t *testing.T) {
	tpls := DefaultTemplates()

	if tpls.Host != DefaultHostTpl {
		t.Errorf("DefaultTemplates().Host = %q, want %q", tpls.Host, DefaultHostTpl)
	}
	if tpls.CPU != DefaultCPUTpl {
		t.Errorf("DefaultTemplates().CPU = %q, want %q", tpls.CPU, DefaultCPUTpl)
	}
	if tpls.Mem != DefaultMemTpl {
		t.Errorf("DefaultTemplates().Mem = %q, want %q", tpls.Mem, DefaultMemTpl)
	}
	if tpls.Load != DefaultLoadTpl {
		t.Errorf("DefaultTemplates().Load = %q, want %q", tpls.Load, DefaultLoadTpl)
	}
	if tpls.Disk != DefaultDiskTpl {
		t.Errorf("DefaultTemplates().Disk = %q, want %q", tpls.Disk, DefaultDiskTpl)
	}
}

func TestMiniTemplates(t *testing.T) {
	tpls := MiniTemplates()

	if tpls.Host != MiniHostTpl {
		t.Errorf("MiniTemplates().Host = %q, want %q", tpls.Host, MiniHostTpl)
	}
	if tpls.CPU != MiniCPUTpl {
		t.Errorf("MiniTemplates().CPU = %q, want %q", tpls.CPU, MiniCPUTpl)
	}
	if tpls.Mem != MiniMemTpl {
		t.Errorf("MiniTemplates().Mem = %q, want %q", tpls.Mem, MiniMemTpl)
	}
	if tpls.Load != MiniLoadTpl {
		t.Errorf("MiniTemplates().Load = %q, want %q", tpls.Load, MiniLoadTpl)
	}
	if tpls.Disk != MiniDiskTpl {
		t.Errorf("MiniTemplates().Disk = %q, want %q", tpls.Disk, MiniDiskTpl)
	}
}

func TestRender_PreservesOrder(t *testing.T) {
	funcMap := NewFuncMapBuilder().Build()
	tpls := Templates{
		Host: "1",
		CPU:  "2",
		Mem:  "3",
		Load: "4",
		Disk: "5",
	}
	renderer := NewRenderer(tpls, funcMap, "-")
	info := newTestInfo()

	// Order: Disk, CPU, Mem
	names := []collector.CollectorName{collector.NameDisk, collector.NameCPU, collector.NameMem}
	output, err := renderer.Render(info, names)

	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := "5 - 2 - 3"
	if output != expected {
		t.Errorf("Render() = %q, want %q (order should be preserved)", output, expected)
	}
}

func TestRender_WithProgressBar(t *testing.T) {
	funcMap := NewFuncMapBuilder().WithProgressBar(ProgressBarConfig{
		Filled: "#",
		Blank:  "-",
	}).Build()

	tpls := Templates{
		Disk: `DISK: {{.Disk.Stat.UsedPercent | progressbar 10}}`,
	}
	renderer := NewRenderer(tpls, funcMap, "|")
	info := newTestInfo() // UsedPercent is 50.0

	output, err := renderer.Render(info, []collector.CollectorName{collector.NameDisk})

	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := "DISK: #####-----"
	if output != expected {
		t.Errorf("Render() = %q, want %q", output, expected)
	}
}

func TestRender_WithColorByThreshold(t *testing.T) {
	funcMap := NewFuncMapBuilder().Build()
	tpls := Templates{
		Disk: `{{.Disk.Stat.UsedPercent | colorByThreshold 30 70 "green" "yellow" "red"}}`,
	}
	renderer := NewRenderer(tpls, funcMap, "|")
	info := newTestInfo() // UsedPercent is 50.0

	output, err := renderer.Render(info, []collector.CollectorName{collector.NameDisk})

	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// 50% is between 30 and 70, so should be yellow
	if !strings.Contains(output, "yellow") {
		t.Errorf("Render() = %q, want to contain 'yellow'", output)
	}
}
