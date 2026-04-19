package rtk

import (
	"context"
	"os/exec"
	"sync"

	"github.com/schlunsen/claude-agent-sdk-go/types"
)

// DefaultCommands is the default allow-list of shell programs whose output
// benefits from RTK compression. Callers can replace or extend it with
// WithCommands / WithAddedCommands.
var DefaultCommands = []string{
	// VCS
	"git", "gh",
	// Node / JS
	"npm", "pnpm", "yarn", "bun", "npx",
	"eslint", "tsc", "prettier", "next",
	"jest", "vitest", "playwright",
	// Python
	"pip", "pytest", "ruff", "mypy",
	// Go
	"go", "golangci-lint",
	// Rust
	"cargo",
	// Ruby
	"bundle", "rspec", "rubocop",
	// Databases / ORMs
	"prisma",
	// Files / search
	"ls", "cat", "grep", "rg", "find", "diff", "tail", "head",
	// Containers / cloud / infra
	"docker", "kubectl", "aws",
	// Networking
	"curl", "wget",
}

// config holds resolved options for the hook.
type config struct {
	binary          string
	commands        map[string]struct{}
	blocked         map[string]struct{}
	ultraCompact    bool
	onlyIfInstalled bool

	installedOnce sync.Once
	installed     bool
}

// Option customises the RTK hook.
type Option func(*config)

// WithBinary overrides the name or absolute path of the rtk executable
// (defaults to "rtk", resolved via PATH).
func WithBinary(path string) Option {
	return func(c *config) {
		if path != "" {
			c.binary = path
		}
	}
}

// WithCommands replaces the default allow-list with the supplied commands.
// Commands are matched on their basename (e.g. "git", not "/usr/bin/git").
//
// Calling WithCommands with no arguments is treated as a no-op rather
// than silently wiping the allow-list, to avoid the footgun of
// accidentally disabling rtk entirely. If you really do want to start
// from an empty allow-list and add commands with WithAddedCommands, call
// WithCommands with an explicit empty slice: WithCommands([]string{}...).
func WithCommands(cmds ...string) Option {
	return func(c *config) {
		if cmds == nil {
			return
		}
		c.commands = toSet(cmds)
	}
}

// WithAddedCommands extends the current allow-list with additional commands.
func WithAddedCommands(cmds ...string) Option {
	return func(c *config) {
		if c.commands == nil {
			c.commands = map[string]struct{}{}
		}
		for _, cmd := range cmds {
			if cmd == "" {
				continue
			}
			c.commands[cmd] = struct{}{}
		}
	}
}

// WithBlocked suppresses rewriting for specific commands, overriding the
// allow-list. Useful for carving out exceptions from DefaultCommands.
func WithBlocked(cmds ...string) Option {
	return func(c *config) {
		if c.blocked == nil {
			c.blocked = map[string]struct{}{}
		}
		for _, cmd := range cmds {
			if cmd == "" {
				continue
			}
			c.blocked[cmd] = struct{}{}
		}
	}
}

// WithUltraCompact passes rtk's "-u" flag for denser output.
func WithUltraCompact(enabled bool) Option {
	return func(c *config) { c.ultraCompact = enabled }
}

// OnlyIfInstalled makes the hook a no-op when the rtk binary is not on PATH.
// Recommended for shared or distributed agent configurations so that users
// who have not installed rtk see unchanged behaviour instead of failing
// tool calls.
func OnlyIfInstalled() Option {
	return func(c *config) { c.onlyIfInstalled = true }
}

// Hook returns a types.HookMatcher ready to register on a PreToolUse hook.
// It matches the "Bash" tool only.
//
//	opts := types.NewClaudeAgentOptions().
//	    WithHook(types.HookEventPreToolUse, rtk.Hook(rtk.OnlyIfInstalled()))
func Hook(opts ...Option) types.HookMatcher {
	cfg := &config{
		binary:   "rtk",
		commands: toSet(DefaultCommands),
	}
	for _, o := range opts {
		if o != nil {
			o(cfg)
		}
	}

	matcher := "Bash"
	return types.HookMatcher{
		Matcher: &matcher,
		Hooks:   []types.HookCallbackFunc{cfg.callback},
	}
}

// passthrough is the canonical "do nothing" hook response.
func passthrough() map[string]interface{} {
	return map[string]interface{}{"continue": true}
}

// callback implements the PreToolUse hook. It never errors on malformed
// input; instead it returns a pass-through response so a misbehaving hook
// cannot break tool execution.
func (c *config) callback(_ context.Context, input interface{}, _ *string, _ types.HookContext) (interface{}, error) {
	if c.onlyIfInstalled && !c.isInstalled() {
		return passthrough(), nil
	}

	toolName, toolInput, ok := extractBashInput(input)
	if !ok || toolName != "Bash" {
		return passthrough(), nil
	}

	cmd, _ := toolInput["command"].(string)
	if cmd == "" {
		return passthrough(), nil
	}

	rewritten, changed := c.rewrite(cmd)
	if !changed {
		return passthrough(), nil
	}

	newInput := cloneInput(toolInput)
	if newInput == nil {
		newInput = map[string]interface{}{}
	}
	newInput["command"] = rewritten

	decision := "allow"
	return map[string]interface{}{
		"hookSpecificOutput": map[string]interface{}{
			"hookEventName":      "PreToolUse",
			"permissionDecision": decision,
			"updatedInput":       newInput,
		},
	}, nil
}

// isInstalled reports whether c.binary resolves on PATH. The result is
// cached for the lifetime of the config so repeated hook invocations do
// not re-stat the filesystem.
func (c *config) isInstalled() bool {
	c.installedOnce.Do(func() {
		_, err := exec.LookPath(c.binary)
		c.installed = err == nil
	})
	return c.installed
}

// extractBashInput unwraps the PreToolUse payload from either a
// *types.PreToolUseHookInput, the corresponding value type, or the raw
// map[string]interface{} that the control protocol delivers.
func extractBashInput(input interface{}) (toolName string, toolInput map[string]interface{}, ok bool) {
	switch v := input.(type) {
	case *types.PreToolUseHookInput:
		if v == nil {
			return "", nil, false
		}
		return v.ToolName, v.ToolInput, true
	case types.PreToolUseHookInput:
		return v.ToolName, v.ToolInput, true
	case map[string]interface{}:
		name, _ := v["tool_name"].(string)
		in, _ := v["tool_input"].(map[string]interface{})
		return name, in, true
	default:
		return "", nil, false
	}
}

func toSet(xs []string) map[string]struct{} {
	out := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		if x == "" {
			continue
		}
		out[x] = struct{}{}
	}
	return out
}
