package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mritd/tmux-sysinfo/internal/collector"
	"github.com/mritd/tmux-sysinfo/internal/formatter"
	"github.com/spf13/cobra"
)

// Config holds CLI configuration
type Config struct {
	Enabled   string
	MiniStyle bool
	Lite      bool

	Delimiter         string
	PerCPU            bool
	DiskUsagePath     string
	CPUInterval       time.Duration
	ProgressBarFilled string
	ProgressBarBlank  string

	HostTpl string
	CPUTpl  string
	MemTpl  string
	LoadTpl string
	DiskTpl string
}

var (
	build   string
	commit  string
	version string
)

var cfg Config

var rootCmd = &cobra.Command{
	Use:     "tmux-sysinfo",
	Short:   "Tmux system info plugin",
	Version: fmt.Sprintf("%s %s %s", version, build, commit),
	RunE:    run,
}

func run(cmd *cobra.Command, args []string) error {
	// Parse enabled collectors
	enabled := strings.Split(cfg.Enabled, ",")
	for i := range enabled {
		enabled[i] = strings.ToLower(strings.TrimSpace(enabled[i]))
	}
	names, _ := collector.ParseNames(enabled)

	if len(names) == 0 {
		return fmt.Errorf("no valid collectors specified in --enabled")
	}

	// Setup collector options
	opts := collector.Options{
		PerCPU:        cfg.PerCPU,
		DiskUsagePath: cfg.DiskUsagePath,
		CPUInterval:   cfg.CPUInterval,
		Lite:          cfg.Lite,
	}

	// Collect system info concurrently
	mgr := collector.NewManager(opts)
	info := mgr.Collect(names)

	// Setup templates
	tpls := buildTemplates(cmd)

	// Setup funcmap with progress bar config
	funcMap := formatter.NewFuncMapBuilder().
		WithProgressBar(formatter.ProgressBarConfig{
			Filled: cfg.ProgressBarFilled,
			Blank:  cfg.ProgressBarBlank,
		}).
		Build()

	// Render output
	renderer := formatter.NewRenderer(tpls, funcMap, cfg.Delimiter)
	output, err := renderer.Render(info, names)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return err
	}

	fmt.Println(output)
	return nil
}

func buildTemplates(cmd *cobra.Command) formatter.Templates {
	var tpls formatter.Templates
	if cfg.MiniStyle {
		tpls = formatter.MiniTemplates()
	} else {
		tpls = formatter.DefaultTemplates()
	}

	// Override only if the user explicitly set the flag; otherwise the style
	// selection above already applied the correct Mini or Default template.
	if cmd.Flags().Changed("host-tpl") {
		tpls.Host = cfg.HostTpl
	}
	if cmd.Flags().Changed("cpu-tpl") {
		tpls.CPU = cfg.CPUTpl
	}
	if cmd.Flags().Changed("mem-tpl") {
		tpls.Mem = cfg.MemTpl
	}
	if cmd.Flags().Changed("load-tpl") {
		tpls.Load = cfg.LoadTpl
	}
	if cmd.Flags().Changed("disk-tpl") {
		tpls.Disk = cfg.DiskTpl
	}

	return tpls
}

func init() {
	rootCmd.Flags().SortFlags = false

	// Basic options
	rootCmd.Flags().StringVar(&cfg.Enabled, "enabled", "all", "Which information to collect (host,cpu,mem,load,disk,all)")
	rootCmd.Flags().BoolVar(&cfg.MiniStyle, "mini", false, "Use mini template style")
	rootCmd.Flags().BoolVar(&cfg.Lite, "lite", false, "Skip rarely-used expensive data (swap, load misc, disk IO)")
	rootCmd.Flags().StringVar(&cfg.Delimiter, "delimiter", "|", "Delimiter between info sections")

	// Collector options
	rootCmd.Flags().BoolVar(&cfg.PerCPU, "per-cpu", false, "Get usage percentage for each CPU")
	rootCmd.Flags().StringVar(&cfg.DiskUsagePath, "disk-usage-path", "/", "Path for disk usage statistics")
	rootCmd.Flags().DurationVar(&cfg.CPUInterval, "cpu-interval", time.Second, "CPU sampling interval")

	// Progress bar options
	defaultBar := formatter.DefaultProgressBar()
	rootCmd.Flags().StringVar(&cfg.ProgressBarFilled, "progress-bar-filled", defaultBar.Filled, "Progress bar filled character")
	rootCmd.Flags().StringVar(&cfg.ProgressBarBlank, "progress-bar-blank", defaultBar.Blank, "Progress bar blank character")

	// Template options
	defaults := formatter.DefaultTemplates()
	rootCmd.Flags().StringVar(&cfg.HostTpl, "host-tpl", defaults.Host, "Host info template")
	rootCmd.Flags().StringVar(&cfg.CPUTpl, "cpu-tpl", defaults.CPU, "CPU info template")
	rootCmd.Flags().StringVar(&cfg.MemTpl, "mem-tpl", defaults.Mem, "Memory info template")
	rootCmd.Flags().StringVar(&cfg.LoadTpl, "load-tpl", defaults.Load, "Load info template")
	rootCmd.Flags().StringVar(&cfg.DiskTpl, "disk-tpl", defaults.Disk, "Disk info template")
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
