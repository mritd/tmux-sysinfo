package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/template"
)

type Conf struct {
	Enabled   string
	MiniStyle bool

	Delimiter         string
	PerCPU            bool
	DiskUsagePath     string
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

var conf Conf

var rootCmd = &cobra.Command{
	Use:     "tmux-sysinfo",
	Short:   "Tmux system info plugin",
	Version: fmt.Sprintf("%s %s %s", version, build, commit),
	RunE: func(cmd *cobra.Command, args []string) error {
		var info Info
		var tpls []string

		if conf.MiniStyle {
			if conf.HostTpl == defaultHostInfoTpl {
				conf.HostTpl = defaultMiniHostInfoTpl
			}
			if conf.CPUTpl == defaultCPUInfoTpl {
				conf.CPUTpl = defaultMiniCPUInfoTpl
			}
			if conf.MemTpl == defaultMemInfoTpl {
				conf.MemTpl = defaultMiniMemInfoTpl
			}
			if conf.DiskTpl == defaultDiskInfoTpl {
				conf.DiskTpl = defaultMiniDiskInfoTpl
			}
			if conf.LoadTpl == defaultLoadInfoTpl {
				conf.LoadTpl = defaultMiniLoadInfoTpl
			}
		}

		enabled := strings.Split(conf.Enabled, ",")
		for _, en := range enabled {
			switch strings.ToLower(en) {
			case "host":
				info.Host = hostInfo()
				tpls = append(tpls, conf.HostTpl)
			case "cpu":
				info.CPU = cpuInfo(conf.PerCPU)
				tpls = append(tpls, conf.CPUTpl)
			case "mem", "memory":
				info.Mem = memInfo()
				tpls = append(tpls, conf.MemTpl)
			case "load":
				info.Load = loadInfo()
				tpls = append(tpls, conf.LoadTpl)
			case "disk":
				info.Disk = diskInfo(conf.DiskUsagePath)
				tpls = append(tpls, conf.DiskTpl)
			case "all":
				if strings.ToLower(en) == "all" {
					info = Info{
						Host: hostInfo(),
						CPU:  cpuInfo(conf.PerCPU),
						Mem:  memInfo(),
						Load: loadInfo(),
						Disk: diskInfo(conf.DiskUsagePath),
					}
					tpls = append(tpls, conf.HostTpl, conf.CPUTpl, conf.MemTpl, conf.LoadTpl, conf.DiskTpl)
					break
				}
			}
		}

		tpl, err := template.New("info").Funcs(funcMap).Parse(strings.Join(tpls, " "+conf.Delimiter+" "))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err = tpl.Execute(&buf, info); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return err
		}
		fmt.Println(buf.String())
		return nil
	},
}

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.Flags().StringVar(&conf.Enabled, "enabled", "all", "Which information output is enabled")
	rootCmd.Flags().BoolVar(&conf.MiniStyle, "mini", false, "Use default mini template")
	rootCmd.Flags().StringVar(&conf.HostTpl, "host-tpl", defaultHostInfoTpl, "Host information rendering template")
	rootCmd.Flags().StringVar(&conf.CPUTpl, "cpu-tpl", defaultCPUInfoTpl, "CPU information rendering template")
	rootCmd.Flags().StringVar(&conf.MemTpl, "mem-tpl", defaultMemInfoTpl, "Memory information rendering template")
	rootCmd.Flags().StringVar(&conf.LoadTpl, "load-tpl", defaultLoadInfoTpl, "Load information rendering template")
	rootCmd.Flags().StringVar(&conf.DiskTpl, "disk-tpl", defaultDiskInfoTpl, "Disk information rendering template")
	rootCmd.Flags().StringVar(&conf.Delimiter, "delimiter", "|", "Delimiter between information areas")
	rootCmd.Flags().BoolVar(&conf.PerCPU, "per-cpu", false, "Get the usage percentage of each CPU")
	rootCmd.Flags().StringVar(&conf.DiskUsagePath, "disk-usage-path", "/", "Disk statistics path")
	rootCmd.Flags().StringVar(&conf.ProgressBarFilled, "progress-bar-filled", "≣", "Progress bar completion character")
	rootCmd.Flags().StringVar(&conf.ProgressBarBlank, "progress-bar-blank", " ", "Progress bar blank character")
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
