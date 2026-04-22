package main

import (
	"testing"

	"github.com/mritd/tmux-sysinfo/internal/formatter"
	"github.com/spf13/cobra"
)

// setupTestCmd builds a cobra.Command with the same flag registrations as
// the production rootCmd, binding to the package-level cfg so buildTemplates
// sees test-controlled state. Resets cfg after the test.
func setupTestCmd(t *testing.T) *cobra.Command {
	t.Helper()
	orig := cfg
	t.Cleanup(func() { cfg = orig })

	cmd := &cobra.Command{Use: "test"}
	defaults := formatter.DefaultTemplates()
	cmd.Flags().BoolVar(&cfg.MiniStyle, "mini", false, "")
	cmd.Flags().StringVar(&cfg.HostTpl, "host-tpl", defaults.Host, "")
	cmd.Flags().StringVar(&cfg.CPUTpl, "cpu-tpl", defaults.CPU, "")
	cmd.Flags().StringVar(&cfg.MemTpl, "mem-tpl", defaults.Mem, "")
	cmd.Flags().StringVar(&cfg.LoadTpl, "load-tpl", defaults.Load, "")
	cmd.Flags().StringVar(&cfg.DiskTpl, "disk-tpl", defaults.Disk, "")
	return cmd
}

// Regression test for the flag-detection bug: when the user explicitly passes
// --cpu-tpl=<value that equals DefaultCPUTpl> together with --mini, the
// explicit flag MUST win over the mini preset.
func TestBuildTemplates_RespectsExplicitFlagInMiniMode(t *testing.T) {
	cmd := setupTestCmd(t)
	if err := cmd.ParseFlags([]string{"--mini", "--cpu-tpl=" + formatter.DefaultCPUTpl}); err != nil {
		t.Fatalf("ParseFlags: %v", err)
	}

	tpls := buildTemplates(cmd)

	if tpls.CPU != formatter.DefaultCPUTpl {
		t.Errorf("tpls.CPU = %q, want %q (explicit flag must beat mini preset)", tpls.CPU, formatter.DefaultCPUTpl)
	}
}

// When the user passes --mini with no explicit tpl flag, mini preset applies.
func TestBuildTemplates_DefaultInMiniMode(t *testing.T) {
	cmd := setupTestCmd(t)
	if err := cmd.ParseFlags([]string{"--mini"}); err != nil {
		t.Fatalf("ParseFlags: %v", err)
	}

	tpls := buildTemplates(cmd)

	if tpls.CPU != formatter.MiniCPUTpl {
		t.Errorf("tpls.CPU = %q, want %q", tpls.CPU, formatter.MiniCPUTpl)
	}
	if tpls.Host != formatter.MiniHostTpl {
		t.Errorf("tpls.Host = %q, want %q", tpls.Host, formatter.MiniHostTpl)
	}
}

// No flags: falls through to default templates.
func TestBuildTemplates_DefaultNonMini(t *testing.T) {
	cmd := setupTestCmd(t)
	if err := cmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("ParseFlags: %v", err)
	}

	tpls := buildTemplates(cmd)

	if tpls.CPU != formatter.DefaultCPUTpl {
		t.Errorf("tpls.CPU = %q, want %q", tpls.CPU, formatter.DefaultCPUTpl)
	}
}

// Explicit custom value, no mini: custom value wins.
func TestBuildTemplates_ExplicitNonMini(t *testing.T) {
	cmd := setupTestCmd(t)
	if err := cmd.ParseFlags([]string{"--cpu-tpl=CUSTOM"}); err != nil {
		t.Fatalf("ParseFlags: %v", err)
	}

	tpls := buildTemplates(cmd)

	if tpls.CPU != "CUSTOM" {
		t.Errorf("tpls.CPU = %q, want %q", tpls.CPU, "CUSTOM")
	}
}
