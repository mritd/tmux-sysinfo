package main

const (
	defaultHostInfoTpl = `HOST: {{.Host.OS}}/{{.Host.KernelVersion}}`
	defaultCPUInfoTpl  = `CPU: {{(index .CPU.InfoStats 0).ModelName}} {{index .CPU.Percent 0 | percentage}}`
	defaultMemInfoTpl  = `MEM: {{.Mem.Stat.Used | humanizeBytes}}`
	defaultDiskInfoTpl = `DISK: {{.Disk.Stat.UsedPercent | percentage}}`
	defaultLoadInfoTpl = `LOAD: {{.Load.Stat.Load1 | percentage}}`

	defaultMiniHostInfoTpl = `S: {{.Host.OS}}`
	defaultMiniCPUInfoTpl  = `C: {{index .CPU.Percent 0 | percentage}}`
	defaultMiniDiskInfoTpl = `D: {{.Disk.Stat.UsedPercent | percentage}}`
	defaultMiniMemInfoTpl  = `M: {{.Mem.Stat.Used | humanizeBytes}}`
	defaultMiniLoadInfoTpl = `L: {{.Load.Stat.Load1 | percentage}}`
)
