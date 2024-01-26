package main

import (
	"bytes"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/template"
)

type Conf struct {
	Enabled string

	Delimiter string
	PerCPU    bool

	CPUTpl  string
	MemTpl  string
	LoadTpl string
}

type Info struct {
	CPU  *CPUInfo
	Mem  *MemoryInfo
	Load *LoadInfo
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

		enabled := strings.Split(conf.Enabled, ",")
		for _, en := range enabled {
			if strings.ToLower(en) == "all" {
				info.CPU = cpuInfo(conf.PerCPU)
				info.Mem = memInfo()
				info.Load = loadInfo()
				tpls = append(tpls, conf.CPUTpl, conf.MemTpl, conf.LoadTpl)
				break
			}
			switch strings.ToLower(en) {
			case "cpu":
				info.CPU = cpuInfo(conf.PerCPU)
				tpls = append(tpls, conf.CPUTpl)
			case "mem", "memory":
				info.Mem = memInfo()
				tpls = append(tpls, conf.MemTpl)
			case "load":
				info.Load = loadInfo()
				tpls = append(tpls, conf.LoadTpl)
			}
		}

		tpl, err := template.New("info").Funcs(template.FuncMap{
			"humanizeBytes":     humanize.Bytes,
			"humanizeIBytes":    humanize.IBytes,
			"humanizeBigBytes":  humanize.BigBytes,
			"humanizeBigIBytes": humanize.BigIBytes,
			"percentage": func(f float64) string {
				return fmt.Sprintf("%.0f%%", f)
			},
		}).Parse(strings.Join(tpls, " "+conf.Delimiter+" "))
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
	rootCmd.Flags().StringVar(&conf.CPUTpl, "cpu-tpl", defaultCPUInfoTpl, "CPU information rendering template")
	rootCmd.Flags().StringVar(&conf.MemTpl, "mem-tpl", defaultMemInfoTpl, "Memory information rendering template")
	rootCmd.Flags().StringVar(&conf.LoadTpl, "load-tpl", defaultLoadInfoTpl, "Load information rendering template")
	rootCmd.Flags().StringVar(&conf.Delimiter, "delimiter", "|", "Delimiter between information areas")
	rootCmd.Flags().BoolVar(&conf.PerCPU, "per-cpu", false, "Get the usage percentage of each CPU")
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
