package rtk

import (
	"context"
	"testing"

	"github.com/schlunsen/claude-agent-sdk-go/types"
)

func newCfg(opts ...Option) *config {
	cfg := &config{
		binary:   "rtk",
		commands: toSet(DefaultCommands),
	}
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

func TestRewrite_SimpleCommand(t *testing.T) {
	cfg := newCfg()
	got, changed := cfg.rewrite("git status")
	if !changed {
		t.Fatalf("expected rewrite to change the command")
	}
	if got != "rtk git status" {
		t.Fatalf("unexpected rewrite: %q", got)
	}
}

func TestRewrite_PreservesLeadingWhitespace(t *testing.T) {
	cfg := newCfg()
	got, changed := cfg.rewrite("   git status")
	if !changed || got != "   rtk git status" {
		t.Fatalf("leading ws not preserved: %q", got)
	}
}

func TestRewrite_AlreadyPrefixed(t *testing.T) {
	cfg := newCfg()
	if _, changed := cfg.rewrite("rtk git status"); changed {
		t.Fatalf("should be no-op when already prefixed")
	}
}

func TestRewrite_AlreadyPrefixedAbsolutePath(t *testing.T) {
	cfg := newCfg(WithBinary("/usr/local/bin/rtk"))
	if _, changed := cfg.rewrite("/usr/local/bin/rtk git status"); changed {
		t.Fatalf("absolute-path rtk prefix should be detected")
	}
}

func TestRewrite_UnknownCommand(t *testing.T) {
	cfg := newCfg()
	if _, changed := cfg.rewrite("make build"); changed {
		t.Fatalf("unknown command should not be rewritten")
	}
}

func TestRewrite_Blocked(t *testing.T) {
	cfg := newCfg(WithBlocked("git"))
	if _, changed := cfg.rewrite("git status"); changed {
		t.Fatalf("blocked command must not be rewritten")
	}
}

func TestRewrite_AddedCommands(t *testing.T) {
	cfg := newCfg(WithCommands(), WithAddedCommands("make"))
	got, changed := cfg.rewrite("make build")
	if !changed || got != "rtk make build" {
		t.Fatalf("added command not wrapped: %q", got)
	}
}

func TestRewrite_UltraCompact(t *testing.T) {
	cfg := newCfg(WithUltraCompact(true))
	got, _ := cfg.rewrite("git status")
	if got != "rtk -u git status" {
		t.Fatalf("ultra-compact flag missing: %q", got)
	}
}

func TestRewrite_PipelineAndSeparators(t *testing.T) {
	cfg := newCfg()
	cases := []struct {
		in, want string
	}{
		{"git status | head -n 20", "rtk git status | rtk head -n 20"},
		{"cd foo && git status", "cd foo && rtk git status"},
		{"npm test; git status", "rtk npm test; rtk git status"},
		{"false || git status", "false || rtk git status"},
	}
	for _, tc := range cases {
		got, _ := cfg.rewrite(tc.in)
		if got != tc.want {
			t.Errorf("rewrite(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRewrite_IgnoresQuotedSeparators(t *testing.T) {
	cfg := newCfg()
	// The '&&' inside the quoted argument must not split the command.
	in := `git commit -m "fix && tidy"`
	got, changed := cfg.rewrite(in)
	if !changed {
		t.Fatalf("expected a rewrite")
	}
	want := `rtk git commit -m "fix && tidy"`
	if got != want {
		t.Fatalf("quoted separators mishandled:\n got:  %s\n want: %s", got, want)
	}
}

func TestRewrite_EnvAssignmentsAndSudo(t *testing.T) {
	cfg := newCfg()
	cases := []struct{ in, want string }{
		{"FOO=1 BAR=baz git status", "FOO=1 BAR=baz rtk git status"},
		{"sudo git status", "sudo rtk git status"},
		{"env FOO=1 git status", "env FOO=1 rtk git status"},
	}
	for _, tc := range cases {
		got, _ := cfg.rewrite(tc.in)
		if got != tc.want {
			t.Errorf("rewrite(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRewrite_EmptyAndWhitespace(t *testing.T) {
	cfg := newCfg()
	if _, changed := cfg.rewrite(""); changed {
		t.Fatalf("empty input should not change")
	}
	if _, changed := cfg.rewrite("   "); changed {
		t.Fatalf("whitespace-only input should not change")
	}
}

func TestCallback_NonBashPassesThrough(t *testing.T) {
	cfg := newCfg()
	resp, err := cfg.callback(context.Background(), map[string]interface{}{
		"tool_name":  "Read",
		"tool_input": map[string]interface{}{"file_path": "/tmp/x"},
	}, nil, types.HookContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("response wrong type: %T", resp)
	}
	if cont, _ := m["continue"].(bool); !cont {
		t.Fatalf("non-bash tool should pass through, got %+v", m)
	}
}

func TestCallback_BashRewritesAndReturnsUpdatedInput(t *testing.T) {
	cfg := newCfg()
	raw := map[string]interface{}{
		"tool_name": "Bash",
		"tool_input": map[string]interface{}{
			"command":     "git status",
			"description": "check status",
		},
	}
	resp, err := cfg.callback(context.Background(), raw, nil, types.HookContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("wrong response type: %T", resp)
	}
	spec, ok := m["hookSpecificOutput"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing hookSpecificOutput: %+v", m)
	}
	if spec["hookEventName"] != "PreToolUse" {
		t.Fatalf("wrong hookEventName: %+v", spec)
	}
	if spec["permissionDecision"] != "allow" {
		t.Fatalf("wrong permissionDecision: %+v", spec)
	}
	updated, ok := spec["updatedInput"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing updatedInput: %+v", spec)
	}
	if updated["command"] != "rtk git status" {
		t.Fatalf("command not rewritten: %+v", updated)
	}
	if updated["description"] != "check status" {
		t.Fatalf("sibling keys dropped: %+v", updated)
	}
	// Original input must not have been mutated.
	if raw["tool_input"].(map[string]interface{})["command"] != "git status" {
		t.Fatalf("original tool_input mutated")
	}
}

func TestCallback_BashUnknownCommandPassesThrough(t *testing.T) {
	cfg := newCfg()
	resp, _ := cfg.callback(context.Background(), map[string]interface{}{
		"tool_name": "Bash",
		"tool_input": map[string]interface{}{
			"command": "make build",
		},
	}, nil, types.HookContext{})
	m := resp.(map[string]interface{})
	if _, ok := m["hookSpecificOutput"]; ok {
		t.Fatalf("expected pass-through, got %+v", m)
	}
}

func TestCallback_AcceptsTypedStructInput(t *testing.T) {
	cfg := newCfg()
	in := &types.PreToolUseHookInput{
		HookEventName: "PreToolUse",
		ToolName:      "Bash",
		ToolInput: map[string]interface{}{
			"command": "git status",
		},
	}
	resp, err := cfg.callback(context.Background(), in, nil, types.HookContext{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	m := resp.(map[string]interface{})
	spec, ok := m["hookSpecificOutput"].(map[string]interface{})
	if !ok {
		t.Fatalf("no hookSpecificOutput for typed input: %+v", m)
	}
	updated := spec["updatedInput"].(map[string]interface{})
	if updated["command"] != "rtk git status" {
		t.Fatalf("typed-input rewrite failed: %+v", updated)
	}
}

func TestHook_MatcherIsBash(t *testing.T) {
	h := Hook()
	if h.Matcher == nil || *h.Matcher != "Bash" {
		t.Fatalf("expected Matcher=\"Bash\", got %+v", h.Matcher)
	}
	if len(h.Hooks) != 1 {
		t.Fatalf("expected exactly one hook callback, got %d", len(h.Hooks))
	}
}

func TestHook_OnlyIfInstalledNoOpWhenMissing(t *testing.T) {
	// Point at a binary guaranteed not to exist on PATH.
	cfg := newCfg(
		WithBinary("rtk-does-not-exist-xyzzy"),
		OnlyIfInstalled(),
	)
	resp, err := cfg.callback(context.Background(), map[string]interface{}{
		"tool_name":  "Bash",
		"tool_input": map[string]interface{}{"command": "git status"},
	}, nil, types.HookContext{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	m := resp.(map[string]interface{})
	if _, wrapped := m["hookSpecificOutput"]; wrapped {
		t.Fatalf("expected pass-through when binary missing, got %+v", m)
	}
	if cont, _ := m["continue"].(bool); !cont {
		t.Fatalf("expected continue:true, got %+v", m)
	}
}
