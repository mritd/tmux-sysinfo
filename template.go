package main

const (
	defaultCPUInfoTpl  = `CPU: {{(index .CPU.InfoStats 0).ModelName}} {{index .CPU.Percent 0 | percentage}}`
	defaultMemInfoTpl  = `MEM: {{.Mem.Stat.Used | humanizeBytes}}`
	defaultLoadInfoTpl = `LOAD: {{.Load.Stat.Load1 | percentage}}`

	defaultMiniCPUInfoTpl  = `C: {{index .CPU.Percent 0 | percentage}}`
	defaultMiniMemInfoTpl  = `M: {{.Mem.Stat.Used | humanizeBytes}}`
	defaultMiniLoadInfoTpl = `L: {{.Load.Stat.Load1 | percentage}}`
)
