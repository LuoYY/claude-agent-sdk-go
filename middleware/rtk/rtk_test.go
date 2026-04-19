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
	cfg := newCfg(WithAddedCommands("make"))
	got, changed := cfg.rewrite("make build")
	if !changed || got != "rtk make build" {
		t.Fatalf("added command not wrapped: %q", got)
	}
}

func TestWithCommands_NoArgsIsNoop(t *testing.T) {
	// Passing zero variadic args must NOT wipe the default allow-list,
	// otherwise rtk.Hook(rtk.WithCommands()) would be a silent footgun
	// that disables the middleware entirely.
	cfg := newCfg(WithCommands())
	got, changed := cfg.rewrite("git status")
	if !changed || got != "rtk git status" {
		t.Fatalf("WithCommands() wiped the allow-list: changed=%v got=%q", changed, got)
	}
}

func TestWithCommands_ExplicitEmptySliceWipes(t *testing.T) {
	// Explicitly passing an empty (non-nil) slice is the documented way
	// to request a fresh allow-list.
	cfg := newCfg(WithCommands([]string{}...))
	if _, changed := cfg.rewrite("git status"); changed {
		t.Fatalf("explicit empty slice should wipe allow-list")
	}
}

func TestWithCommands_ReplaceSet(t *testing.T) {
	cfg := newCfg(WithCommands("only-this"))
	if _, changed := cfg.rewrite("git status"); changed {
		t.Fatalf("git should not be wrapped after WithCommands replace")
	}
	got, changed := cfg.rewrite("only-this now")
	if !changed || got != "rtk only-this now" {
		t.Fatalf("replacement command not wrapped: %q", got)
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

func TestRewrite_WrapperFlagsWithValues(t *testing.T) {
	cfg := newCfg()
	cases := []struct{ in, want string }{
		// sudo -u <user>: the user argument is a flag value, not the program.
		{"sudo -u deploy git status", "sudo -u deploy rtk git status"},
		// nice -n <prio>: same pattern.
		{"nice -n 10 git status", "nice -n 10 rtk git status"},
		// env with its own flags and an env-assignment before the command.
		{"env -i PATH=/bin git status", "env -i PATH=/bin rtk git status"},
		// Multiple flags, some with values.
		{"sudo -E -u deploy -H git status", "sudo -E -u deploy -H rtk git status"},
		// Wrapper + flag landing on a non-allow-listed program -> passthrough.
		{"sudo -u deploy rm -rf /", "sudo -u deploy rm -rf /"},
	}
	for _, tc := range cases {
		got, _ := cfg.rewrite(tc.in)
		if got != tc.want {
			t.Errorf("rewrite(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRewrite_BareCommandArgumentsNotSwallowed(t *testing.T) {
	// Without a wrapper, a flag followed by a non-flag MUST NOT be
	// treated as a flag value; otherwise "echo git" would be rewritten
	// as "echo rtk git" (wrong — git is echo's argument, not a program).
	cfg := newCfg()
	if _, changed := cfg.rewrite("echo git"); changed {
		t.Fatalf("bare 'echo git' must not be rewritten")
	}
	if _, changed := cfg.rewrite("echo -v git"); changed {
		t.Fatalf("'echo -v git' must not be rewritten")
	}
}

func TestRewrite_SubshellsAndParens(t *testing.T) {
	cfg := newCfg()
	cases := []struct{ in, want string }{
		// A subshell groups commands; operators inside must not split the
		// outer command, but the subshell body is left alone (we cannot
		// safely rewrite inside parens without shell semantics).
		{"(cd foo && git status)", "(cd foo && git status)"},
		// Pipe OUTSIDE the subshell is a top-level separator.
		{"(cd foo && git status) | head", "(cd foo && git status) | rtk head"},
		// Multiple top-level commands with a subshell in one of them.
		{"(cd a && npm test) && git status", "(cd a && npm test) && rtk git status"},
	}
	for _, tc := range cases {
		got, _ := cfg.rewrite(tc.in)
		if got != tc.want {
			t.Errorf("rewrite(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRewrite_CommandSubstitution(t *testing.T) {
	cfg := newCfg()
	cases := []struct{ in, want string }{
		// $(...) — pipes inside the substitution must NOT split the
		// outer command, and the substitution body itself is left alone.
		{"git log --format=$(date +%s)", "rtk git log --format=$(date +%s)"},
		{"echo $(git log | head) && git status", "echo $(git log | head) && rtk git status"},
		// Backtick substitution gets the same treatment.
		{"echo `git log | head` && git status", "echo `git log | head` && rtk git status"},
		// Nested substitution.
		{"echo $(git log $(date +%s)) && npm test", "echo $(git log $(date +%s)) && rtk npm test"},
	}
	for _, tc := range cases {
		got, _ := cfg.rewrite(tc.in)
		if got != tc.want {
			t.Errorf("rewrite(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRewrite_RedirectionsAndBackground(t *testing.T) {
	cfg := newCfg()
	cases := []struct{ in, want string }{
		{"git status 2>&1", "rtk git status 2>&1"},
		{"git status > /tmp/out", "rtk git status > /tmp/out"},
		{"git status &", "rtk git status &"},
	}
	for _, tc := range cases {
		got, _ := cfg.rewrite(tc.in)
		if got != tc.want {
			t.Errorf("rewrite(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSplitSegments_RoundTripFidelity(t *testing.T) {
	// Whatever we parse should round-trip to the exact original string
	// when no rewriting happens — this guards against subtle byte-loss
	// bugs in splitSegments / joinSegments.
	cases := []string{
		"git status",
		"git status | head",
		"git status | head -n 20",
		"cd foo && git status; npm test || false",
		`git commit -m "fix && tidy"`,
		`git commit -m 'fix || tidy'`,
		"(cd foo && git status) | head",
		"echo $(git log | head) && git status",
		"echo `git log` && git status",
		"git log --format=$(date +%s)",
		"git status 2>&1",
		"git status &",
	}
	for _, in := range cases {
		segs, seps := splitSegments(in)
		got := joinSegments(segs, seps)
		if got != in {
			t.Errorf("round-trip failed:\n in:  %q\n out: %q", in, got)
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
