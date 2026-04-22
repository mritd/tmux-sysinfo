package formatter

// Default templates for each info type
const (
	DefaultHostTpl = `HOST: {{.Host.OS}}/{{.Host.KernelVersion}}`
	DefaultCPUTpl  = `CPU: {{(index .CPU.InfoStats 0).ModelName}} {{index .CPU.Percent 0 | percentage}}`
	DefaultMemTpl  = `MEM: {{.Mem.Stat.Used | humanizeBytes}}`
	DefaultDiskTpl = `DISK: {{.Disk.Stat.UsedPercent | percentage}}`
	DefaultLoadTpl = `LOAD: {{.Load.Stat.Load1 | percentage}}`

	// Mini style templates
	MiniHostTpl = `S: {{.Host.OS}}`
	MiniCPUTpl  = `C: {{index .CPU.Percent 0 | percentage}}`
	MiniDiskTpl = `D: {{.Disk.Stat.UsedPercent | percentage}}`
	MiniMemTpl  = `M: {{.Mem.Stat.Used | humanizeBytes}}`
	MiniLoadTpl = `L: {{.Load.Stat.Load1 | percentage}}`
)

// Templates holds all template strings
type Templates struct {
	Host string
	CPU  string
	Mem  string
	Load string
	Disk string
}

// DefaultTemplates returns default templates
func DefaultTemplates() Templates {
	return Templates{
		Host: DefaultHostTpl,
		CPU:  DefaultCPUTpl,
		Mem:  DefaultMemTpl,
		Load: DefaultLoadTpl,
		Disk: DefaultDiskTpl,
	}
}

// MiniTemplates returns mini style templates
func MiniTemplates() Templates {
	return Templates{
		Host: MiniHostTpl,
		CPU:  MiniCPUTpl,
		Mem:  MiniMemTpl,
		Load: MiniLoadTpl,
		Disk: MiniDiskTpl,
	}
}
