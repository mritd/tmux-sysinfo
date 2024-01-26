package main

import "fmt"

const (
	defaultCPUInfoTpl  = `CPU: {{(index .CPU.InfoStats 0).ModelName}} {{index .CPU.Percent 0 | percentage}}`
	defaultMemInfoTpl  = `MEM: {{.Mem.Stat.Used | humanizeIBytes}}`
	defaultLoadInfoTpl = `LOAD: {{.Load.Stat.Load1 | percentage}}`
)

var defaultTpl = fmt.Sprintf("%s | %s | %s", defaultCPUInfoTpl, defaultMemInfoTpl, defaultLoadInfoTpl)
