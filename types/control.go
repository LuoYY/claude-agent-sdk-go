package types

import "encoding/json"

// PermissionMode represents the permission mode for Claude.
type PermissionMode string

const (
	PermissionModeDefault           PermissionMode = "default"
	PermissionModeAcceptEdits       PermissionMode = "acceptEdits"
	PermissionModePlan              PermissionMode = "plan"
	PermissionModeBypassPermissions PermissionMode = "bypassPermissions"
	PermissionModeDontAsk           PermissionMode = "dontAsk"
	PermissionModeAuto              PermissionMode = "auto"
)

// PermissionBehavior represents the behavior for a permission rule.
type PermissionBehavior string

const (
	PermissionBehaviorAllow PermissionBehavior = "allow"
	PermissionBehaviorDeny  PermissionBehavior = "deny"
	PermissionBehaviorAsk   PermissionBehavior = "ask"
)

// PermissionUpdateDestination represents where permission updates should be saved.
type PermissionUpdateDestination string

const (
	DestinationUserSettings    PermissionUpdateDestination = "userSettings"
	DestinationProjectSettings PermissionUpdateDestination = "projectSettings"
	DestinationLocalSettings   PermissionUpdateDestination = "localSettings"
	DestinationSession         PermissionUpdateDestination = "session"
)

// PermissionRuleValue represents a permission rule.
type PermissionRuleValue struct {
	ToolName    string  `json:"toolName"`
	RuleContent *string `json:"ruleContent,omitempty"`
}

// PermissionUpdate represents a permission update configuration.
type PermissionUpdate struct {
	Type        string                       `json:"type"` // addRules, replaceRules, removeRules, setMode, addDirectories, removeDirectories
	Rules       []PermissionRuleValue        `json:"rules,omitempty"`
	Behavior    *PermissionBehavior          `json:"behavior,omitempty"`
	Mode        *PermissionMode              `json:"mode,omitempty"`
	Directories []string                     `json:"directories,omitempty"`
	Destination *PermissionUpdateDestination `json:"destination,omitempty"`
}

// PermissionResultAllow represents an allow permission result.
type PermissionResultAllow struct {
	Behavior           string                  `json:"behavior"` // "allow"
	UpdatedInput       *map[string]interface{} `json:"updated_input,omitempty"`
	UpdatedPermissions []PermissionUpdate      `json:"updated_permissions,omitempty"`
}

// PermissionResultDeny represents a deny permission result.
type PermissionResultDeny struct {
	Behavior  string `json:"behavior"` // "deny"
	Message   string `json:"message,omitempty"`
	Interrupt bool   `json:"interrupt,omitempty"`
}

// ToolPermissionContext provides context for tool permission callbacks.
type ToolPermissionContext struct {
	Signal      interface{}        `json:"signal,omitempty"` // Future: abort signal support
	Suggestions []PermissionUpdate `json:"suggestions,omitempty"`
}

// HookEvent represents a hook event type.
type HookEvent string

const (
	HookEventPreToolUse         HookEvent = "PreToolUse"
	HookEventPostToolUse        HookEvent = "PostToolUse"
	HookEventPostToolUseFailure HookEvent = "PostToolUseFailure"
	HookEventUserPromptSubmit   HookEvent = "UserPromptSubmit"
	HookEventStop               HookEvent = "Stop"
	HookEventSubagentStop       HookEvent = "SubagentStop"
	HookEventSubagentStart      HookEvent = "SubagentStart"
	HookEventPreCompact         HookEvent = "PreCompact"
	HookEventNotification       HookEvent = "Notification"
	HookEventPermissionRequest  HookEvent = "PermissionRequest"
)

// BaseHookInput contains common fields for all hook inputs.
type BaseHookInput struct {
	SessionID      string  `json:"session_id"`
	TranscriptPath string  `json:"transcript_path"`
	CWD            string  `json:"cwd"`
	PermissionMode *string `json:"permission_mode,omitempty"`
}

// PreToolUseHookInput represents input for PreToolUse hook events.
type PreToolUseHookInput struct {
	BaseHookInput
	HookEventName string                 `json:"hook_event_name"` // "PreToolUse"
	ToolName      string                 `json:"tool_name"`
	ToolInput     map[string]interface{} `json:"tool_input"`
	ToolUseID     *string                `json:"tool_use_id,omitempty"`
	// Subagent context fields (present when hook fires from inside a sub-agent)
	AgentID   *string `json:"agent_id,omitempty"`
	AgentType *string `json:"agent_type,omitempty"`
}

// PostToolUseHookInput represents input for PostToolUse hook events.
type PostToolUseHookInput struct {
	BaseHookInput
	HookEventName string                 `json:"hook_event_name"` // "PostToolUse"
	ToolName      string                 `json:"tool_name"`
	ToolInput     map[string]interface{} `json:"tool_input"`
	ToolResponse  interface{}            `json:"tool_response"`
	ToolUseID     *string                `json:"tool_use_id,omitempty"`
	// Subagent context fields (present when hook fires from inside a sub-agent)
	AgentID   *string `json:"agent_id,omitempty"`
	AgentType *string `json:"agent_type,omitempty"`
}

// UserPromptSubmitHookInput represents input for UserPromptSubmit hook events.
type UserPromptSubmitHookInput struct {
	BaseHookInput
	HookEventName string `json:"hook_event_name"` // "UserPromptSubmit"
	Prompt        string `json:"prompt"`
}

// StopHookInput represents input for Stop hook events.
type StopHookInput struct {
	BaseHookInput
	HookEventName  string `json:"hook_event_name"` // "Stop"
	StopHookActive bool   `json:"stop_hook_active"`
}

// SubagentStopHookInput represents input for SubagentStop hook events.
type SubagentStopHookInput struct {
	BaseHookInput
	HookEventName       string  `json:"hook_event_name"` // "SubagentStop"
	StopHookActive      bool    `json:"stop_hook_active"`
	AgentID             string  `json:"agent_id"`
	AgentType           *string `json:"agent_type,omitempty"`
	AgentTranscriptPath *string `json:"agent_transcript_path,omitempty"`
}

// PreCompactHookInput represents input for PreCompact hook events.
type PreCompactHookInput struct {
	BaseHookInput
	HookEventName      string  `json:"hook_event_name"` // "PreCompact"
	Trigger            string  `json:"trigger"`         // "manual" or "auto"
	CustomInstructions *string `json:"custom_instructions,omitempty"`
}

// PostToolUseFailureHookInput represents input for PostToolUseFailure hook events.
type PostToolUseFailureHookInput struct {
	BaseHookInput
	HookEventName string                 `json:"hook_event_name"` // "PostToolUseFailure"
	ToolName      string                 `json:"tool_name"`
	ToolInput     map[string]interface{} `json:"tool_input"`
	ToolUseID     *string                `json:"tool_use_id,omitempty"`
	Error         string                 `json:"error"`
	IsInterrupt   *bool                  `json:"is_interrupt,omitempty"`
	// Subagent context fields (present when hook fires from inside a sub-agent)
	AgentID   *string `json:"agent_id,omitempty"`
	AgentType *string `json:"agent_type,omitempty"`
}

// NotificationHookInput represents input for Notification hook events.
type NotificationHookInput struct {
	BaseHookInput
	HookEventName    string  `json:"hook_event_name"` // "Notification"
	Message          string  `json:"message"`
	Title            *string `json:"title,omitempty"`
	NotificationType *string `json:"notification_type,omitempty"`
}

// SubagentStartHookInput represents input for SubagentStart hook events.
type SubagentStartHookInput struct {
	BaseHookInput
	HookEventName string  `json:"hook_event_name"` // "SubagentStart"
	AgentID       string  `json:"agent_id"`
	AgentType     *string `json:"agent_type,omitempty"`
}

// PermissionRequestHookInput represents input for PermissionRequest hook events.
type PermissionRequestHookInput struct {
	BaseHookInput
	HookEventName         string                 `json:"hook_event_name"` // "PermissionRequest"
	ToolName              string                 `json:"tool_name"`
	ToolInput             map[string]interface{} `json:"tool_input"`
	PermissionSuggestions []PermissionUpdate     `json:"permission_suggestions,omitempty"`
	// Subagent context fields (present when hook fires from inside a sub-agent)
	AgentID   *string `json:"agent_id,omitempty"`
	AgentType *string `json:"agent_type,omitempty"`
}

// HookSpecificOutput is an interface for all hook-specific outputs.
type HookSpecificOutput interface {
	GetHookEventName() string
}

// PreToolUseHookSpecificOutput represents hook-specific output for PreToolUse events.
type PreToolUseHookSpecificOutput struct {
	HookEventName            string                  `json:"hookEventName"`                // "PreToolUse"
	PermissionDecision       *string                 `json:"permissionDecision,omitempty"` // "allow", "deny", "ask"
	PermissionDecisionReason *string                 `json:"permissionDecisionReason,omitempty"`
	UpdatedInput             *map[string]interface{} `json:"updatedInput,omitempty"`
}

// GetHookEventName returns the hook event name.
func (h *PreToolUseHookSpecificOutput) GetHookEventName() string {
	return h.HookEventName
}

// PostToolUseHookSpecificOutput represents hook-specific output for PostToolUse events.
type PostToolUseHookSpecificOutput struct {
	HookEventName     string  `json:"hookEventName"` // "PostToolUse"
	AdditionalContext *string `json:"additionalContext,omitempty"`
}

// GetHookEventName returns the hook event name.
func (h *PostToolUseHookSpecificOutput) GetHookEventName() string {
	return h.HookEventName
}

// UserPromptSubmitHookSpecificOutput represents hook-specific output for UserPromptSubmit events.
type UserPromptSubmitHookSpecificOutput struct {
	HookEventName     string  `json:"hookEventName"` // "UserPromptSubmit"
	AdditionalContext *string `json:"additionalContext,omitempty"`
}

// GetHookEventName returns the hook event name.
func (h *UserPromptSubmitHookSpecificOutput) GetHookEventName() string {
	return h.HookEventName
}

// PostToolUseFailureHookSpecificOutput represents hook-specific output for PostToolUseFailure events.
type PostToolUseFailureHookSpecificOutput struct {
	HookEventName     string  `json:"hookEventName"` // "PostToolUseFailure"
	AdditionalContext *string `json:"additionalContext,omitempty"`
}

// GetHookEventName returns the hook event name.
func (h *PostToolUseFailureHookSpecificOutput) GetHookEventName() string {
	return h.HookEventName
}

// NotificationHookSpecificOutput represents hook-specific output for Notification events.
type NotificationHookSpecificOutput struct {
	HookEventName string `json:"hookEventName"` // "Notification"
}

// GetHookEventName returns the hook event name.
func (h *NotificationHookSpecificOutput) GetHookEventName() string {
	return h.HookEventName
}

// SubagentStartHookSpecificOutput represents hook-specific output for SubagentStart events.
type SubagentStartHookSpecificOutput struct {
	HookEventName     string  `json:"hookEventName"` // "SubagentStart"
	AdditionalContext *string `json:"additionalContext,omitempty"`
}

// GetHookEventName returns the hook event name.
func (h *SubagentStartHookSpecificOutput) GetHookEventName() string {
	return h.HookEventName
}

// PermissionRequestHookSpecificOutput represents hook-specific output for PermissionRequest events.
type PermissionRequestHookSpecificOutput struct {
	HookEventName            string  `json:"hookEventName"`                // "PermissionRequest"
	PermissionDecision       *string `json:"permissionDecision,omitempty"` // "allow", "deny", "ask"
	PermissionDecisionReason *string `json:"permissionDecisionReason,omitempty"`
}

// GetHookEventName returns the hook event name.
func (h *PermissionRequestHookSpecificOutput) GetHookEventName() string {
	return h.HookEventName
}

// AsyncHookJSONOutput represents async hook output that defers hook execution.
type AsyncHookJSONOutput struct {
	Async        bool `json:"async"`
	AsyncTimeout *int `json:"asyncTimeout,omitempty"`
}

// SyncHookJSONOutput represents synchronous hook output with control and decision fields.
type SyncHookJSONOutput struct {
	// Common control fields
	Continue       *bool   `json:"continue,omitempty"`
	SuppressOutput *bool   `json:"suppressOutput,omitempty"`
	StopReason     *string `json:"stopReason,omitempty"`

	// Decision fields
	Decision      *string `json:"decision,omitempty"` // "block"
	SystemMessage *string `json:"systemMessage,omitempty"`
	Reason        *string `json:"reason,omitempty"`

	// Hook-specific outputs
	HookSpecificOutput interface{} `json:"hookSpecificOutput,omitempty"`
}

// HookContext provides context information for hook callbacks.
type HookContext struct {
	Signal interface{} `json:"signal,omitempty"` // Future: abort signal support
}

// SDKControlInterruptRequest represents an interrupt request.
type SDKControlInterruptRequest struct {
	Subtype string `json:"subtype"` // "interrupt"
}

// SDKControlPermissionRequest represents a permission request for tool use.
type SDKControlPermissionRequest struct {
	Subtype               string                 `json:"subtype"` // "can_use_tool"
	ToolName              string                 `json:"tool_name"`
	Input                 map[string]interface{} `json:"input"`
	PermissionSuggestions []PermissionUpdate     `json:"permission_suggestions,omitempty"`
	BlockedPath           *string                `json:"blocked_path,omitempty"`
}

// SDKControlInitializeRequest represents an initialization request.
type SDKControlInitializeRequest struct {
	Subtype string                 `json:"subtype"` // "initialize"
	Hooks   map[string]interface{} `json:"hooks,omitempty"`
}

// SDKControlSetPermissionModeRequest represents a request to set permission mode.
type SDKControlSetPermissionModeRequest struct {
	Subtype string `json:"subtype"` // "set_permission_mode"
	Mode    string `json:"mode"`
}

// SDKHookCallbackRequest represents a hook callback request.
type SDKHookCallbackRequest struct {
	Subtype    string      `json:"subtype"` // "hook_callback"
	CallbackID string      `json:"callback_id"`
	Input      interface{} `json:"input"`
	ToolUseID  *string     `json:"tool_use_id,omitempty"`
}

// SDKControlMcpMessageRequest represents an MCP message request.
type SDKControlMcpMessageRequest struct {
	Subtype    string      `json:"subtype"` // "mcp_message"
	ServerName string      `json:"server_name"`
	Message    interface{} `json:"message"`
}

// SDKControlRequest represents a control request from the CLI.
type SDKControlRequest struct {
	Type      string          `json:"type"` // "control_request"
	RequestID string          `json:"request_id"`
	Request   json.RawMessage `json:"request"` // Union type - needs custom unmarshaling
}

// ControlResponse represents a successful control response.
type ControlResponse struct {
	Subtype   string                 `json:"subtype"` // "success"
	RequestID string                 `json:"request_id"`
	Response  map[string]interface{} `json:"response,omitempty"`
}

// ControlErrorResponse represents an error control response.
type ControlErrorResponse struct {
	Subtype   string `json:"subtype"` // "error"
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
}

// SDKControlResponse represents a control response to the CLI.
type SDKControlResponse struct {
	Type     string          `json:"type"`     // "control_response"
	Response json.RawMessage `json:"response"` // Union type - needs custom unmarshaling
}

// MCPServer represents an MCP server interface for handling MCP messages.
// This is a minimal interface for routing MCP JSONRPC messages.
// Concrete implementations can use the MCP SDK or custom logic.
type MCPServer interface {
	// HandleMessage handles an incoming JSONRPC message and returns the response.
	HandleMessage(message map[string]interface{}) (map[string]interface{}, error)

	// Name returns the server name.
	Name() string

	// Version returns the server version.
	Version() string
}

// ContextUsageCategory represents a single context usage category
// (system prompt, tools, messages, etc.).
type ContextUsageCategory struct {
	Name       string `json:"name"`
	Tokens     int    `json:"tokens"`
	Color      string `json:"color"`
	IsDeferred *bool  `json:"isDeferred,omitempty"`
}

// ContextUsageResponse provides a breakdown of current context window usage by category,
// matching the data shown by the `/context` command in the CLI.
type ContextUsageResponse struct {
	// Categories is token usage broken down by category (system prompt, tools, messages, etc.)
	Categories []ContextUsageCategory `json:"categories"`

	// TotalTokens is total tokens currently in the context window
	TotalTokens int `json:"totalTokens"`

	// MaxTokens is effective maximum tokens (may be reduced by autocompact buffer)
	MaxTokens int `json:"maxTokens"`

	// RawMaxTokens is the raw model context window size
	RawMaxTokens int `json:"rawMaxTokens"`

	// Percentage is percentage of context window used (0-100)
	Percentage float64 `json:"percentage"`

	// Model is model name the context usage is calculated for
	Model string `json:"model"`

	// IsAutoCompactEnabled is whether autocompact is enabled for this session
	IsAutoCompactEnabled bool `json:"isAutoCompactEnabled"`

	// MemoryFiles are CLAUDE.md and memory files loaded, with path, type, and token counts
	MemoryFiles []map[string]interface{} `json:"memoryFiles,omitempty"`

	// MCPTools are MCP tools with name, serverName, tokens, and isLoaded status
	MCPTools []map[string]interface{} `json:"mcpTools,omitempty"`

	// Agents are agent definitions with agentType, source, and token counts
	Agents []map[string]interface{} `json:"agents,omitempty"`

	// GridRows is visual grid representation used by the CLI context display
	GridRows [][]map[string]interface{} `json:"gridRows,omitempty"`

	// AutoCompactThreshold is the token threshold at which autocompact triggers
	AutoCompactThreshold *int `json:"autoCompactThreshold,omitempty"`

	// DeferredBuiltinTools are built-in tools deferred from the initial tool list
	DeferredBuiltinTools []map[string]interface{} `json:"deferredBuiltinTools,omitempty"`

	// SystemTools are system (built-in) tools with name and token counts
	SystemTools []map[string]interface{} `json:"systemTools,omitempty"`

	// SystemPromptSections are system prompt sections with name and token counts
	SystemPromptSections []map[string]interface{} `json:"systemPromptSections,omitempty"`

	// SlashCommands is slash command usage summary
	SlashCommands map[string]interface{} `json:"slashCommands,omitempty"`

	// Skills is skill usage summary with frontmatter breakdown
	Skills map[string]interface{} `json:"skills,omitempty"`

	// MessageBreakdown is detailed breakdown of message tokens by type (tool calls, results, etc.)
	MessageBreakdown map[string]interface{} `json:"messageBreakdown,omitempty"`

	// APIUsage is cumulative API usage for the session
	APIUsage map[string]interface{} `json:"apiUsage,omitempty"`
}
