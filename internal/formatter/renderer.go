package formatter

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/mritd/tmux-sysinfo/internal/collector"
	"github.com/mritd/tmux-sysinfo/internal/types"
)

// Renderer handles template rendering
type Renderer struct {
	templates Templates
	funcMap   template.FuncMap
	delimiter string
}

// NewRenderer creates a new Renderer
func NewRenderer(tpls Templates, funcMap template.FuncMap, delimiter string) *Renderer {
	return &Renderer{
		templates: tpls,
		funcMap:   funcMap,
		delimiter: delimiter,
	}
}

// Render renders the info using templates for the specified collectors
func (r *Renderer) Render(info *types.Info, names []collector.CollectorName) (string, error) {
	var tpls []string

	for _, name := range names {
		switch name {
		case collector.NameHost:
			tpls = append(tpls, r.templates.Host)
		case collector.NameCPU:
			tpls = append(tpls, r.templates.CPU)
		case collector.NameMem:
			tpls = append(tpls, r.templates.Mem)
		case collector.NameLoad:
			tpls = append(tpls, r.templates.Load)
		case collector.NameDisk:
			tpls = append(tpls, r.templates.Disk)
		}
	}

	tplStr := strings.Join(tpls, " "+r.delimiter+" ")
	tpl, err := template.New("info").Funcs(r.funcMap).Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, info); err != nil {
		return "", err
	}

	return buf.String(), nil
}
