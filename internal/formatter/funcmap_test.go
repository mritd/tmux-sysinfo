package formatter

import (
	"bytes"
	"testing"
	"text/template"
)

func TestPercentageFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"zero", 0, "0%"},
		{"integer", 50, "50%"},
		{"decimal_round_down", 49.4, "49%"},
		{"decimal_round_up", 49.5, "50%"},
		{"hundred", 100, "100%"},
		{"over_hundred", 150, "150%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := percentageFunc(tt.input)
			if result != tt.expected {
				t.Errorf("percentageFunc(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProgressbarFunc(t *testing.T) {
	builder := NewFuncMapBuilder().WithProgressBar(ProgressBarConfig{
		Filled: "#",
		Blank:  "-",
	})

	tests := []struct {
		name     string
		value    float64
		length   int
		expected string
	}{
		{"zero", 0, 10, "----------"},
		{"half", 50, 10, "#####-----"},
		{"full", 100, 10, "##########"},
		{"over", 150, 10, "##########"},
		{"negative", -10, 10, "----------"},
		{"small_length", 25, 4, "#---"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.progressbarFunc(tt.length, tt.value)
			if result != tt.expected {
				t.Errorf("progressbarFunc(%d, %v) = %q, want %q", tt.length, tt.value, result, tt.expected)
			}
		})
	}
}

func TestFgColorFunc(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		text     interface{}
		expected string
	}{
		{"string", "green", "test", "#[fg=green]test#[fg=default]"},
		{"number", "red", 42, "#[fg=red]42#[fg=default]"},
		{"hex_color", "#ff0000", "text", "#[fg=#ff0000]text#[fg=default]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fgColorFunc(tt.color, tt.text)
			if result != tt.expected {
				t.Errorf("fgColorFunc(%q, %v) = %q, want %q", tt.color, tt.text, result, tt.expected)
			}
		})
	}
}

func TestBgColorFunc(t *testing.T) {
	result := bgColorFunc("blue", "text")
	expected := "#[bg=blue]text#[bg=default]"
	if result != expected {
		t.Errorf("bgColorFunc() = %q, want %q", result, expected)
	}
}

func TestStyleFunc(t *testing.T) {
	tests := []struct {
		name     string
		style    string
		text     interface{}
		expected string
	}{
		{"bold", "bold", "text", "#[bold]text#[default]"},
		{"dim", "dim", "text", "#[dim]text#[default]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := styleFunc(tt.style, tt.text)
			if result != tt.expected {
				t.Errorf("styleFunc(%q, %v) = %q, want %q", tt.style, tt.text, result, tt.expected)
			}
		})
	}
}

func TestColorByThresholdFunc(t *testing.T) {
	tests := []struct {
		name      string
		warn      float64
		crit      float64
		colorLow  string
		colorMid  string
		colorHigh string
		value     float64
		expected  string
	}{
		{"low", 50, 80, "green", "yellow", "red", 30, "#[fg=green]30%#[fg=default]"},
		{"mid", 50, 80, "green", "yellow", "red", 60, "#[fg=yellow]60%#[fg=default]"},
		{"high", 50, 80, "green", "yellow", "red", 90, "#[fg=red]90%#[fg=default]"},
		{"at_warn", 50, 80, "green", "yellow", "red", 50, "#[fg=yellow]50%#[fg=default]"},
		{"at_crit", 50, 80, "green", "yellow", "red", 80, "#[fg=red]80%#[fg=default]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorByThresholdFunc(tt.warn, tt.crit, tt.colorLow, tt.colorMid, tt.colorHigh, tt.value)
			if result != tt.expected {
				t.Errorf("colorByThresholdFunc() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIfGreaterThanFunc(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		threshold float64
		expected  bool
	}{
		{"greater", 80, 50, true},
		{"equal", 50, 50, false},
		{"less", 30, 50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ifGreaterThanFunc(tt.value, tt.threshold)
			if result != tt.expected {
				t.Errorf("ifGreaterThanFunc(%v, %v) = %v, want %v", tt.value, tt.threshold, result, tt.expected)
			}
		})
	}
}

func TestIfLessThanFunc(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		threshold float64
		expected  bool
	}{
		{"less", 30, 50, true},
		{"equal", 50, 50, false},
		{"greater", 80, 50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ifLessThanFunc(tt.value, tt.threshold)
			if result != tt.expected {
				t.Errorf("ifLessThanFunc(%v, %v) = %v, want %v", tt.value, tt.threshold, result, tt.expected)
			}
		})
	}
}

func TestTruncateFunc(t *testing.T) {
	tests := []struct {
		name     string
		maxLen   int
		suffix   string
		input    string
		expected string
	}{
		{"no_truncate", 10, "...", "short", "short"},
		{"truncate", 8, "...", "very long string", "very ..."},
		{"exact_length", 5, "...", "hello", "hello"},
		{"suffix_only", 3, "...", "hello", "..."},
		{"empty", 5, "...", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateFunc(tt.maxLen, tt.suffix, tt.input)
			if result != tt.expected {
				t.Errorf("truncateFunc(%d, %q, %q) = %q, want %q", tt.maxLen, tt.suffix, tt.input, result, tt.expected)
			}
		})
	}
}

func TestPadLeftFunc(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		pad      string
		input    string
		expected string
	}{
		{"pad", 5, " ", "42", "   42"},
		{"no_pad", 2, " ", "42", "42"},
		{"longer", 3, " ", "hello", "hello"},
		{"zero_pad", 4, "0", "42", "0042"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padLeftFunc(tt.length, tt.pad, tt.input)
			if result != tt.expected {
				t.Errorf("padLeftFunc(%d, %q, %q) = %q, want %q", tt.length, tt.pad, tt.input, result, tt.expected)
			}
		})
	}
}

func TestPadRightFunc(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		pad      string
		input    string
		expected string
	}{
		{"pad", 5, " ", "42", "42   "},
		{"no_pad", 2, " ", "42", "42"},
		{"longer", 3, " ", "hello", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padRightFunc(tt.length, tt.pad, tt.input)
			if result != tt.expected {
				t.Errorf("padRightFunc(%d, %q, %q) = %q, want %q", tt.length, tt.pad, tt.input, result, tt.expected)
			}
		})
	}
}

func TestMathFuncs(t *testing.T) {
	t.Run("div", func(t *testing.T) {
		if result := divFunc(10, 2); result != 5 {
			t.Errorf("divFunc(10, 2) = %v, want 5", result)
		}
		if result := divFunc(10, 0); result != 0 {
			t.Errorf("divFunc(10, 0) = %v, want 0", result)
		}
	})

	t.Run("mul", func(t *testing.T) {
		if result := mulFunc(3, 4); result != 12 {
			t.Errorf("mulFunc(3, 4) = %v, want 12", result)
		}
	})

	t.Run("add", func(t *testing.T) {
		if result := addFunc(3, 4); result != 7 {
			t.Errorf("addFunc(3, 4) = %v, want 7", result)
		}
	})

	t.Run("sub", func(t *testing.T) {
		if result := subFunc(10, 4); result != 6 {
			t.Errorf("subFunc(10, 4) = %v, want 6", result)
		}
	})
}

func TestFuncMapBuilder(t *testing.T) {
	t.Run("default_progress_bar", func(t *testing.T) {
		builder := NewFuncMapBuilder()
		if builder.progressBar.Filled != "≣" {
			t.Errorf("default filled = %q, want %q", builder.progressBar.Filled, "≣")
		}
	})

	t.Run("custom_progress_bar", func(t *testing.T) {
		builder := NewFuncMapBuilder().WithProgressBar(ProgressBarConfig{
			Filled: "X",
			Blank:  "O",
		})
		result := builder.progressbarFunc(10, 50)
		if result != "XXXXXOOOOO" {
			t.Errorf("progressbarFunc with custom config = %q, want %q", result, "XXXXXOOOOO")
		}
	})

	t.Run("build_returns_funcmap", func(t *testing.T) {
		funcMap := NewFuncMapBuilder().Build()
		if funcMap == nil {
			t.Error("Build() returned nil")
		}

		expectedFuncs := []string{
			"humanizeBytes", "humanizeIBytes", "percentage", "progressbar",
			"fgColor", "bgColor", "style", "colorByThreshold",
			"ifgt", "iflt", "truncate", "padLeft", "padRight",
			"div", "mul", "add", "sub", "printf",
		}
		for _, name := range expectedFuncs {
			if _, ok := funcMap[name]; !ok {
				t.Errorf("FuncMap missing function %q", name)
			}
		}
	})
}

func TestFuncMapInTemplate(t *testing.T) {
	funcMap := NewFuncMapBuilder().WithProgressBar(ProgressBarConfig{
		Filled: "#",
		Blank:  "-",
	}).Build()

	tests := []struct {
		name     string
		tpl      string
		data     interface{}
		expected string
	}{
		{
			name:     "percentage_pipeline",
			tpl:      `{{.Value | percentage}}`,
			data:     map[string]interface{}{"Value": 42.5},
			expected: "42%",
		},
		{
			name:     "progressbar_pipeline",
			tpl:      `{{.Value | progressbar 5}}`,
			data:     map[string]interface{}{"Value": 40.0},
			expected: "##---",
		},
		{
			name:     "color_threshold_pipeline",
			tpl:      `{{.Value | colorByThreshold 50 80 "green" "yellow" "red"}}`,
			data:     map[string]interface{}{"Value": 30.0},
			expected: "#[fg=green]30%#[fg=default]",
		},
		{
			name:     "truncate_pipeline",
			tpl:      `{{.Name | truncate 5 "..."}}`,
			data:     map[string]interface{}{"Name": "very long name"},
			expected: "ve...",
		},
		{
			name:     "chained_functions",
			tpl:      `{{.Value | percentage | fgColor "green"}}`,
			data:     map[string]interface{}{"Value": 50.0},
			expected: "#[fg=green]50%#[fg=default]",
		},
		{
			name:     "conditional",
			tpl:      `{{if (ifgt .Value 50)}}HIGH{{else}}LOW{{end}}`,
			data:     map[string]interface{}{"Value": 80.0},
			expected: "HIGH",
		},
		{
			name:     "math_chain",
			tpl:      `{{div .A .B | printf "%.2f"}}`,
			data:     map[string]interface{}{"A": 10.0, "B": 3.0},
			expected: "3.33",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpl, err := template.New("test").Funcs(funcMap).Parse(tt.tpl)
			if err != nil {
				t.Fatalf("template parse error: %v", err)
			}

			var buf bytes.Buffer
			if err := tpl.Execute(&buf, tt.data); err != nil {
				t.Fatalf("template execute error: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("template output = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}
