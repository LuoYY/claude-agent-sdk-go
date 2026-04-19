// Example: reducing token usage with RTK.
//
// This example registers the drop-in RTK PreToolUse hook so every Bash
// command Claude issues is transparently wrapped with "rtk <cmd>",
// compressing output before it returns to the model.
//
// Prerequisites:
//
//	brew install rtk            # or: curl https://rtk.ai/install.sh | sh
//
// If rtk is not installed, OnlyIfInstalled() makes the hook a no-op so
// this example still runs correctly; commands just won't be compressed.
package main

import (
	"context"
	"fmt"
	"log"

	claude "github.com/schlunsen/claude-agent-sdk-go"
	"github.com/schlunsen/claude-agent-sdk-go/middleware/rtk"
	"github.com/schlunsen/claude-agent-sdk-go/types"
)

func main() {
	ctx := context.Background()

	opts := types.NewClaudeAgentOptions().
		WithModel("claude-sonnet-4-5-20250929").
		WithAllowedTools("Bash").
		WithPermissionMode(types.PermissionModeBypassPermissions).
		WithHook(types.HookEventPreToolUse, rtk.Hook(
			rtk.WithUltraCompact(true), // pass "-u" for denser output
			rtk.OnlyIfInstalled(),      // no-op if rtk binary is missing
		))

	messages, err := claude.Query(ctx,
		"Run 'git status' and tell me what branch I'm on and whether the tree is clean.",
		opts,
	)
	if err != nil {
		log.Fatalf("query failed: %v", err)
	}

	for msg := range messages {
		switch m := msg.(type) {
		case *types.AssistantMessage:
			for _, block := range m.Content {
				switch b := block.(type) {
				case *types.TextBlock:
					fmt.Println(b.Text)
				case *types.ToolUseBlock:
					// Note: if rtk is installed, the command here will
					// already be the rewritten "rtk git status" form.
					if cmd, ok := b.Input["command"].(string); ok {
						fmt.Printf("[tool] %s -> %s\n", b.Name, cmd)
					}
				}
			}
		case *types.ResultMessage:
			fmt.Println("--- done ---")
		}
	}
}
