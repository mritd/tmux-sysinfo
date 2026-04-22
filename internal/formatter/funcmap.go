package formatter

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/dustin/go-humanize"
)

// ProgressBarConfig holds progress bar display characters
type ProgressBarConfig struct {
	Filled string
	Blank  string
}

// DefaultProgressBar returns default progress bar config
func DefaultProgressBar() ProgressBarConfig {
	return ProgressBarConfig{
		Filled: "≣",
		Blank:  " ",
	}
}

// FuncMapBuilder builds template.FuncMap with custom configuration
type FuncMapBuilder struct {
	progressBar ProgressBarConfig
}

// NewFuncMapBuilder creates a new FuncMapBuilder
func NewFuncMapBuilder() *FuncMapBuilder {
	return &FuncMapBuilder{
		progressBar: DefaultProgressBar(),
	}
}

// WithProgressBar sets progress bar configuration
func (b *FuncMapBuilder) WithProgressBar(cfg ProgressBarConfig) *FuncMapBuilder {
	b.progressBar = cfg
	return b
}

// Build creates the template.FuncMap
func (b *FuncMapBuilder) Build() template.FuncMap {
	return template.FuncMap{
		// Byte formatting
		"humanizeBytes":  humanize.Bytes,
		"humanizeIBytes": humanize.IBytes,

		// Percentage formatting
		"percentage": percentageFunc,

		// Progress bar
		"progressbar": b.progressbarFunc,

		// tmux color support
		"fgColor": fgColorFunc,
		"bgColor": bgColorFunc,
		"style":   styleFunc,

		// Conditional formatting
		"colorByThreshold": colorByThresholdFunc,
		"ifgt":             ifGreaterThanFunc,
		"iflt":             ifLessThanFunc,

		// String utilities
		"truncate": truncateFunc,
		"padLeft":  padLeftFunc,
		"padRight": padRightFunc,

		// Math utilities
		"div": divFunc,
		"mul": mulFunc,
		"add": addFunc,
		"sub": subFunc,

		// Formatting
		"printf": fmt.Sprintf,
	}
}

// percentageFunc formats a float as percentage string
func percentageFunc(f float64) string {
	return fmt.Sprintf("%.0f%%", f)
}

// progressbarFunc creates a progress bar string
// Usage: {{.Value | progressbar 10}}
func (b *FuncMapBuilder) progressbarFunc(length int, f float64) string {
	progress := int((f / 100) * float64(length))
	if progress > length {
		progress = length
	}
	if progress < 0 {
		progress = 0
	}
	return strings.Repeat(b.progressBar.Filled, progress) + strings.Repeat(b.progressBar.Blank, length-progress)
}

// fgColorFunc wraps text with tmux foreground color
// Usage: {{.Value | fgColor "green"}}
func fgColorFunc(color string, text interface{}) string {
	return fmt.Sprintf("#[fg=%s]%v#[fg=default]", color, text)
}

// bgColorFunc wraps text with tmux background color
// Usage: {{.Value | bgColor "red"}}
func bgColorFunc(color string, text interface{}) string {
	return fmt.Sprintf("#[bg=%s]%v#[bg=default]", color, text)
}

// styleFunc wraps text with tmux style
// Usage: {{.Value | style "bold"}}
// Supported: bold, dim, underscore, blink, reverse, hidden, italics, strikethrough
func styleFunc(styleName string, text interface{}) string {
	return fmt.Sprintf("#[%s]%v#[default]", styleName, text)
}

// colorByThresholdFunc returns colored text based on value thresholds
// Usage: {{.Value | colorByThreshold 50 80 "green" "yellow" "red"}}
// < warn: green, warn-crit: yellow, >= crit: red
func colorByThresholdFunc(warn, crit float64, colorLow, colorMid, colorHigh string, value float64) string {
	var color string
	switch {
	case value < warn:
		color = colorLow
	case value < crit:
		color = colorMid
	default:
		color = colorHigh
	}
	return fmt.Sprintf("#[fg=%s]%.0f%%#[fg=default]", color, value)
}

// ifGreaterThanFunc returns trueVal if value > threshold, else falseVal
// Usage: {{if (ifgt .Value 80)}}...{{end}}
func ifGreaterThanFunc(value, threshold float64) bool {
	return value > threshold
}

// ifLessThanFunc returns true if value < threshold
// Usage: {{if (iflt .Value 20)}}...{{end}}
func ifLessThanFunc(value, threshold float64) bool {
	return value < threshold
}

// truncateFunc truncates string to max length, adding suffix if truncated
// Usage: {{"long string" | truncate 5 "..."}}
func truncateFunc(maxLen int, suffix string, s string) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= len(suffix) {
		return suffix[:maxLen]
	}
	return s[:maxLen-len(suffix)] + suffix
}

// padLeftFunc pads string on the left to reach target length
// Usage: {{"42" | padLeft 5 " "}}
func padLeftFunc(length int, pad string, s string) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(pad, length-len(s)) + s
}

// padRightFunc pads string on the right to reach target length
// Usage: {{"42" | padRight 5 " "}}
func padRightFunc(length int, pad string, s string) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(pad, length-len(s))
}

// Math functions for template calculations
func divFunc(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func mulFunc(a, b float64) float64 {
	return a * b
}

func addFunc(a, b float64) float64 {
	return a + b
}

func subFunc(a, b float64) float64 {
	return a - b
}
