// Package rtk provides a drop-in PreToolUse hook that wraps noisy shell
// commands issued by Claude through the Bash tool with the RTK CLI proxy
// (https://github.com/rtk-ai/rtk), compressing their output before it is
// returned to the model. RTK reports 60-90% token reductions on commands
// like "git status", "npm test", "cargo build", etc.
//
// # Basic usage
//
//	import (
//	    "github.com/schlunsen/claude-agent-sdk-go/types"
//	    "github.com/schlunsen/claude-agent-sdk-go/middleware/rtk"
//	)
//
//	opts := types.NewClaudeAgentOptions().
//	    WithHook(types.HookEventPreToolUse, rtk.Hook(
//	        rtk.WithUltraCompact(true),
//	        rtk.OnlyIfInstalled(),
//	    ))
//
// # How it works
//
// The hook matches only the "Bash" tool. For each shell invocation Claude
// wants to run, it parses the command line, detects the first program in
// each pipeline segment (e.g. "cd dir && git status | head") and, if that
// program is in the configured allow-list, rewrites the segment so the
// first token is prefixed with "rtk " (e.g. "cd dir && rtk git status | head").
//
// Commands that are already prefixed with rtk, not in the allow-list, or
// explicitly blocked are passed through unchanged. If OnlyIfInstalled is
// set and the rtk binary is not on PATH, the hook is a no-op.
//
// # Scope
//
// The hook only affects the built-in Bash tool. Built-in Read, Grep, and
// Glob tools bypass the hook entirely; steer the agent toward shell
// equivalents (cat, rg, find) if you want maximum savings.
//
// # Extending
//
// Use WithCommands to replace the default set, WithAddedCommands to extend
// it, WithBlocked to suppress rewriting for specific programs, WithBinary
// to change the rtk executable path, and WithUltraCompact to pass rtk's
// "-u" flag for even denser output.
package rtk
