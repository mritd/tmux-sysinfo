package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"humanizeBytes":  humanize.Bytes,
	"humanizeIBytes": humanize.IBytes,
	"percentage":     percentageFunc,
	"progressbar":    progressbarFunc,
}

func percentageFunc(f float64) string {
	return fmt.Sprintf("%.0f%%", f)
}

func progressbarFunc(f float64, l int) string {
	progress := int((f / 100) * float64(l))
	return strings.Repeat(conf.ProgressBarFilled, progress) + strings.Repeat(conf.ProgressBarBlank, l-progress)
}
