# Go SDK vs Python SDK Feature Parity Guide

## At a Glance

| Metric | Value |
|--------|-------|
| **Go SDK Version** | v0.6.0 |
| **Python SDK Version** | v0.1.36 |
| **Feature Parity** | ~99% |
| **Status** | Production-Ready |
| **Last Updated** | 2026-02-15 |

## Overview

This document provides a comprehensive comparison of features between the Go Agent SDK and the official Python Agent SDK. Both SDKs are fully functional and production-ready for all use cases.

### Key Differences in Approach

- **Language Differences**: Python (async/await) vs Go (goroutines/channels)
- **Configuration**: Python (Pydantic dataclasses) vs Go (builder pattern)
- **Threading**: Python (GIL) vs Go (explicit concurrency)
- **Error Handling**: Python (exceptions) vs Go (error types)

---

## Complete Feature Comparison Matrix

### Core API & Session Management

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Query (one-shot)** | ✅ | ✅ | Both SDKs support simple one-off requests |
| **Client (interactive)** | ✅ | ✅ | Full bidirectional communication |
| **Message streaming** | ✅ | ✅ | Python: async generators, Go: channels |
| **Session resumption** | ✅ | ✅ | Continue previous conversations with session ID |
| **Session forking** | ✅ | ✅ | Create new branch from existing session |
| **Message history** | ✅ | ✅ | Access and replay message sequences |
| **Model selection** | ✅ | ✅ | Switch between available models |
| **Fallback model** | ✅ | ✅ | Use alternate model if primary unavailable |
| **Max turns control** | ✅ | ✅ | Limit conversation length |
| **Max budget control** | ✅ | ✅ | Cost limits in USD |

### Client Runtime Control

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Interrupt** | ✅ | ✅ | Stop current operation mid-stream |
| **Set model** | ✅ | ✅ | Change AI model mid-conversation |
| **Set permission mode** | ✅ | ✅ | Change permission mode at runtime |
| **Get MCP status** | ✅ | ✅ | Check MCP server connection status |
| **Get server info** | ✅ | ✅ | Retrieve server initialization info |
| **Rewind files** | ✅ | ✅ | Undo file changes to checkpoint |
| **File checkpointing** | ✅ | ✅ | Track file changes for rewinding |

### Message Types & Content Blocks

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **User messages** | ✅ | ✅ | With UUID and tool_use_result fields |
| **Assistant messages** | ✅ | ✅ | With error field support |
| **System messages** | ✅ | ✅ | System notifications |
| **Result messages** | ✅ | ✅ | With structured_output field |
| **Stream events** | ✅ | ✅ | Partial message updates |
| **Text content blocks** | ✅ | ✅ | Plain text responses |
| **Tool use blocks** | ✅ | ✅ | Tool invocation requests |
| **Tool result blocks** | ✅ | ✅ | Tool execution results |
| **Thinking blocks** | ✅ | ✅ | With signature field |

### Tool Integration & Permissions

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Tool permissions** | ✅ | ✅ | Permission callbacks for tool use |
| **Permission modes** | ✅ | ✅ | default, acceptEdits, plan, bypassPermissions |
| **Tool filtering** | ✅ | ✅ | AllowedTools, DisallowedTools |
| **Tool use callbacks** | ✅ | ✅ | React to tool execution |
| **Permission storage** | ✅ | ✅ | Save permissions to user/project/local settings |
| **Updated input support** | ✅ | ✅ | Callbacks can modify tool inputs |

### MCP (Model Context Protocol) Support

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **SDK MCP servers** | ✅ | ✅ | Create tools in-process |
| **External MCP servers** | ✅ | ✅ | stdio, SSE, HTTP connections |
| **Custom MCP servers** | ✅ | ✅ | Implement MCPServer interface |
| **Tool schema validation** | ✅ | ✅ | JSON schema input validation |
| **Tool listing** | ✅ | ✅ | Discover available tools dynamically |
| **Tool annotations** | ✅ | ✅ | readOnlyHint, destructiveHint, etc. |
| **MCP factory function** | ❌ | ✅ | Go SDK has convenient factory |

### Hook System

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **PreToolUse** | ✅ | ✅ | Before tool execution (with tool_use_id) |
| **PostToolUse** | ✅ | ✅ | After tool execution (with tool_use_id) |
| **PostToolUseFailure** | ✅ | ✅ | On tool execution failure |
| **UserPromptSubmit** | ✅ | ✅ | Before user input processing |
| **Stop** | ✅ | ✅ | Session stop event |
| **SubagentStop** | ✅ | ✅ | Subagent stops (with agent_id, type, transcript) |
| **SubagentStart** | ✅ | ✅ | Subagent starts (with agent_id, type) |
| **PreCompact** | ✅ | ✅ | Before context compaction |
| **Notification** | ✅ | ✅ | Notification events (message, title, type) |
| **PermissionRequest** | ✅ | ✅ | Permission request events |
| **Regex matching** | ✅ | ✅ | Filter hooks by tool name pattern |
| **Hook continuation** | ✅ | ✅ | Control whether to continue execution |
| **Async hook output** | ✅ | ✅ | Defer hook execution |

### Extended Thinking

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **ThinkingConfig adaptive** | ✅ | ✅ | Default adaptive thinking |
| **ThinkingConfig enabled** | ✅ | ✅ | With token budget control |
| **ThinkingConfig disabled** | ✅ | ✅ | Disable thinking |
| **Effort levels** | ✅ | ✅ | low, medium, high, max |
| **Max thinking tokens** | ✅ | ✅ | Deprecated (use ThinkingConfig) |
| **Thinking blocks** | ✅ | ✅ | Access reasoning process |

### Structured Output

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Output format** | ✅ | ✅ | JSON schema for structured responses |
| **Structured output result** | ✅ | ✅ | Validated output in ResultMessage |

### Sandbox Configuration

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Sandbox enabled** | ✅ | ✅ | Enable bash sandboxing |
| **Auto-allow bash** | ✅ | ✅ | Auto-approve when sandboxed |
| **Excluded commands** | ✅ | ✅ | Commands that bypass sandbox |
| **Network config** | ✅ | ✅ | Unix sockets, local binding, proxy |
| **Ignore violations** | ✅ | ✅ | File and network exemptions |

### System Configuration

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **System prompt** | ✅ | ✅ | Custom system instructions |
| **System prompt presets** | ✅ | ✅ | claude_code preset with append |
| **Setting sources** | ✅ | ✅ | user, project, local |
| **Model parameter** | ✅ | ✅ | Specify model version |

### Cost Management

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Token counting** | ✅ | ✅ | Access input/output token counts |
| **Cost summary** | ✅ | ✅ | Get usage statistics per message |
| **Budget limiting** | ✅ | ✅ | Max budget in USD |

### Plugin Support

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Local plugins** | ✅ | ✅ | Load from directory |
| **Plugin discovery** | ✅ | ✅ | Auto-detect plugin.json |

### Beta Features

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Beta registration** | ✅ | ✅ | Opt-in to experimental features |
| **Beta list** | ✅ | ✅ | Extended context, new models, etc. |

### Subagent Support

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Agent definitions** | ✅ | ✅ | Custom agents with tools and models |
| **Execution modes** | ✅ | ✅ | Sequential, parallel, auto |
| **Timeout/max turns** | ✅ | ✅ | Per-agent limits |
| **Execution config** | ✅ | ✅ | Global concurrency and error handling |

### Error Handling & Validation

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Typed errors** | ✅ | ✅ | Specific error types for different failures |
| **Error unwrapping** | ✅ | ✅ | Access root cause of errors |
| **CLI not found** | ✅ | ✅ | Detect missing Claude Code CLI |
| **Connection errors** | ✅ | ✅ | Network/subprocess failures |
| **Session not found** | ✅ | ✅ | Resume session missing |
| **JSON validation** | ✅ | ✅ | Parse and validate JSON responses |

### Transport & Infrastructure

| Feature | Python | Go | Notes |
|---------|--------|-----|-------|
| **Subprocess transport** | ✅ | ✅ | Connect via Claude Code CLI |
| **CLI discovery** | ✅ | ✅ | Auto-find Claude Code CLI |
| **CLI version check** | ✅ | ✅ | Validate compatible CLI version |
| **Environment variables** | ✅ | ✅ | Pass env to subprocess |
| **Extra CLI args** | ✅ | ✅ | Pass arbitrary flags |
| **Stderr callbacks** | ✅ | ✅ | Capture CLI debug output |
| **Stderr file logging** | ✅ | ✅ | SDK-managed log files |

---

## Version Compatibility

| Go SDK Version | Python SDK Version | Compatible |
|---|---|---|
| v0.6.0+ | v0.1.36+ | ✅ Yes |
| v0.2.0+ | v0.1.0+ | ✅ Yes |

Both SDKs use the same control protocol and are forward-compatible with different versions.

---

## When to Use Each SDK

### Choose Go SDK When:
- ✅ You need concurrent request handling
- ✅ You want strong type safety
- ✅ You prefer Go's concurrency model
- ✅ You need high performance
- ✅ You want a single binary deployment
- ✅ You want zero external dependencies

### Choose Python SDK When:
- ✅ You prefer Python's syntax
- ✅ You need rapid prototyping
- ✅ You're working with Jupyter notebooks
- ✅ You have an existing Python codebase
- ✅ You prefer async/await patterns

---

**Last Verified**: 2026-02-15

For the latest SDK versions and features, see:
- [Go SDK Releases](https://github.com/schlunsen/claude-agent-sdk-go/releases)
- [Python SDK Releases](https://github.com/anthropics/claude-agent-sdk-python/releases)
