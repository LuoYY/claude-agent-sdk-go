package types

import (
	"encoding/json"
	"testing"
)

// TestPermissionModeConstants tests that permission mode constants are defined correctly.
func TestPermissionModeConstants(t *testing.T) {
	modes := []PermissionMode{
		PermissionModeDefault,
		PermissionModeAcceptEdits,
		PermissionModePlan,
		PermissionModeBypassPermissions,
	}

	for _, mode := range modes {
		if mode == "" {
			t.Error("permission mode should not be empty")
		}
	}
}

// TestPermissionUpdateMarshaling tests JSON marshaling of PermissionUpdate.
func TestPermissionUpdateMarshaling(t *testing.T) {
	behavior := PermissionBehaviorAllow
	update := &PermissionUpdate{
		Type: "addRules",
		Rules: []PermissionRuleValue{
			{
				ToolName:    "Bash",
				RuleContent: stringPtr("allow ls command"),
			},
		},
		Behavior: &behavior,
	}

	data, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("failed to marshal PermissionUpdate: %v", err)
	}

	var decoded PermissionUpdate
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal PermissionUpdate: %v", err)
	}

	if decoded.Type != update.Type {
		t.Errorf("type doesn't match")
	}
	if len(decoded.Rules) != len(update.Rules) {
		t.Errorf("rules length doesn't match")
	}
}

// TestSDKControlPermissionRequest tests JSON marshaling of SDKControlPermissionRequest.
func TestSDKControlPermissionRequest(t *testing.T) {
	req := &SDKControlPermissionRequest{
		Subtype:  "can_use_tool",
		ToolName: "Bash",
		Input: map[string]interface{}{
			"command": "ls -la",
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal SDKControlPermissionRequest: %v", err)
	}

	var decoded SDKControlPermissionRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal SDKControlPermissionRequest: %v", err)
	}

	if decoded.ToolName != req.ToolName {
		t.Errorf("tool name doesn't match")
	}
}

// TestHookEventConstants tests that hook event constants are defined correctly.
func TestHookEventConstants(t *testing.T) {
	events := []HookEvent{
		HookEventPreToolUse,
		HookEventPostToolUse,
		HookEventPostToolUseFailure,
		HookEventUserPromptSubmit,
		HookEventStop,
		HookEventSubagentStop,
		HookEventSubagentStart,
		HookEventPreCompact,
		HookEventNotification,
		HookEventPermissionRequest,
	}

	for _, event := range events {
		if event == "" {
			t.Error("hook event should not be empty")
		}
	}

	// Verify specific string values match Python SDK
	expectedValues := map[HookEvent]string{
		HookEventPreToolUse:         "PreToolUse",
		HookEventPostToolUse:        "PostToolUse",
		HookEventPostToolUseFailure: "PostToolUseFailure",
		HookEventUserPromptSubmit:   "UserPromptSubmit",
		HookEventStop:               "Stop",
		HookEventSubagentStop:       "SubagentStop",
		HookEventSubagentStart:      "SubagentStart",
		HookEventPreCompact:         "PreCompact",
		HookEventNotification:       "Notification",
		HookEventPermissionRequest:  "PermissionRequest",
	}

	for event, expected := range expectedValues {
		if string(event) != expected {
			t.Errorf("HookEvent %q should be %q", event, expected)
		}
	}
}

// TestPreToolUseHookInput tests JSON marshaling of PreToolUseHookInput.
func TestPreToolUseHookInput(t *testing.T) {
	input := &PreToolUseHookInput{
		BaseHookInput: BaseHookInput{
			SessionID:      "session-123",
			TranscriptPath: "/path/to/transcript",
			CWD:            "/home/user",
		},
		HookEventName: "PreToolUse",
		ToolName:      "Bash",
		ToolInput: map[string]interface{}{
			"command": "echo hello",
		},
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal PreToolUseHookInput: %v", err)
	}

	var decoded PreToolUseHookInput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal PreToolUseHookInput: %v", err)
	}

	if decoded.ToolName != input.ToolName {
		t.Errorf("tool name doesn't match")
	}
}

// TestPostToolUseFailureHookInput tests JSON marshaling of PostToolUseFailureHookInput.
func TestPostToolUseFailureHookInput(t *testing.T) {
	toolUseID := "tool-456"
	input := &PostToolUseFailureHookInput{
		BaseHookInput: BaseHookInput{
			SessionID:      "session-123",
			TranscriptPath: "/path/to/transcript",
			CWD:            "/home/user",
		},
		HookEventName: "PostToolUseFailure",
		ToolName:      "Bash",
		ToolInput: map[string]interface{}{
			"command": "rm -rf /",
		},
		ToolUseID: &toolUseID,
		Error:     "permission denied",
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal PostToolUseFailureHookInput: %v", err)
	}

	var decoded PostToolUseFailureHookInput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal PostToolUseFailureHookInput: %v", err)
	}

	if decoded.ToolName != input.ToolName {
		t.Errorf("tool name doesn't match")
	}
	if decoded.Error != input.Error {
		t.Errorf("error doesn't match: got %q, want %q", decoded.Error, input.Error)
	}
	if decoded.ToolUseID == nil || *decoded.ToolUseID != toolUseID {
		t.Errorf("tool_use_id doesn't match")
	}
}

// TestNotificationHookInput tests JSON marshaling of NotificationHookInput.
func TestNotificationHookInput(t *testing.T) {
	title := "Build Complete"
	notifType := "info"
	input := &NotificationHookInput{
		BaseHookInput: BaseHookInput{
			SessionID:      "session-123",
			TranscriptPath: "/path/to/transcript",
			CWD:            "/home/user",
		},
		HookEventName:    "Notification",
		Message:          "Build succeeded",
		Title:            &title,
		NotificationType: &notifType,
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal NotificationHookInput: %v", err)
	}

	var decoded NotificationHookInput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal NotificationHookInput: %v", err)
	}

	if decoded.Message != input.Message {
		t.Errorf("message doesn't match")
	}
	if decoded.Title == nil || *decoded.Title != title {
		t.Errorf("title doesn't match")
	}
}

// TestSubagentStartHookInput tests JSON marshaling of SubagentStartHookInput.
func TestSubagentStartHookInput(t *testing.T) {
	agentType := "custom"
	input := &SubagentStartHookInput{
		BaseHookInput: BaseHookInput{
			SessionID:      "session-123",
			TranscriptPath: "/path/to/transcript",
			CWD:            "/home/user",
		},
		HookEventName: "SubagentStart",
		AgentID:       "agent-789",
		AgentType:     &agentType,
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal SubagentStartHookInput: %v", err)
	}

	var decoded SubagentStartHookInput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal SubagentStartHookInput: %v", err)
	}

	if decoded.AgentID != input.AgentID {
		t.Errorf("agent_id doesn't match")
	}
	if decoded.AgentType == nil || *decoded.AgentType != agentType {
		t.Errorf("agent_type doesn't match")
	}
}

// TestPermissionRequestHookInput tests JSON marshaling of PermissionRequestHookInput.
func TestPermissionRequestHookInput(t *testing.T) {
	input := &PermissionRequestHookInput{
		BaseHookInput: BaseHookInput{
			SessionID:      "session-123",
			TranscriptPath: "/path/to/transcript",
			CWD:            "/home/user",
		},
		HookEventName: "PermissionRequest",
		ToolName:      "Write",
		ToolInput: map[string]interface{}{
			"file_path": "/etc/hosts",
		},
		PermissionSuggestions: []PermissionUpdate{
			{
				Type: "addRules",
				Rules: []PermissionRuleValue{
					{ToolName: "Write", RuleContent: stringPtr("allow /etc/hosts")},
				},
			},
		},
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal PermissionRequestHookInput: %v", err)
	}

	var decoded PermissionRequestHookInput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal PermissionRequestHookInput: %v", err)
	}

	if decoded.ToolName != input.ToolName {
		t.Errorf("tool name doesn't match")
	}
	if len(decoded.PermissionSuggestions) != 1 {
		t.Errorf("expected 1 permission suggestion, got %d", len(decoded.PermissionSuggestions))
	}
}

// TestNewHookSpecificOutputTypes tests the GetHookEventName method on new output types.
func TestNewHookSpecificOutputTypes(t *testing.T) {
	t.Run("PostToolUseFailureHookSpecificOutput", func(t *testing.T) {
		output := &PostToolUseFailureHookSpecificOutput{
			HookEventName:     "PostToolUseFailure",
			AdditionalContext: stringPtr("Consider retrying"),
		}
		if output.GetHookEventName() != "PostToolUseFailure" {
			t.Errorf("expected PostToolUseFailure, got %s", output.GetHookEventName())
		}
	})

	t.Run("NotificationHookSpecificOutput", func(t *testing.T) {
		output := &NotificationHookSpecificOutput{
			HookEventName: "Notification",
		}
		if output.GetHookEventName() != "Notification" {
			t.Errorf("expected Notification, got %s", output.GetHookEventName())
		}
	})

	t.Run("SubagentStartHookSpecificOutput", func(t *testing.T) {
		output := &SubagentStartHookSpecificOutput{
			HookEventName:     "SubagentStart",
			AdditionalContext: stringPtr("Starting agent"),
		}
		if output.GetHookEventName() != "SubagentStart" {
			t.Errorf("expected SubagentStart, got %s", output.GetHookEventName())
		}
	})

	t.Run("PermissionRequestHookSpecificOutput", func(t *testing.T) {
		decision := "allow"
		reason := "Trusted tool"
		output := &PermissionRequestHookSpecificOutput{
			HookEventName:            "PermissionRequest",
			PermissionDecision:       &decision,
			PermissionDecisionReason: &reason,
		}
		if output.GetHookEventName() != "PermissionRequest" {
			t.Errorf("expected PermissionRequest, got %s", output.GetHookEventName())
		}
	})
}

// TestSubagentContextFields tests that subagent context fields are properly serialized
// on tool-lifecycle hook inputs (matching Python SDK's _SubagentContextMixin).
func TestSubagentContextFields(t *testing.T) {
	agentID := "agent-123"
	agentType := "code-reviewer"

	t.Run("PreToolUseHookInput with subagent context", func(t *testing.T) {
		toolUseID := "tool-789"
		input := &PreToolUseHookInput{
			BaseHookInput: BaseHookInput{SessionID: "s1", TranscriptPath: "/t", CWD: "/c"},
			HookEventName: "PreToolUse",
			ToolName:      "Bash",
			ToolInput:     map[string]interface{}{"command": "ls"},
			ToolUseID:     &toolUseID,
			AgentID:       &agentID,
			AgentType:     &agentType,
		}
		data, _ := json.Marshal(input)
		var decoded PreToolUseHookInput
		_ = json.Unmarshal(data, &decoded)
		if decoded.AgentID == nil || *decoded.AgentID != agentID {
			t.Errorf("agent_id mismatch")
		}
		if decoded.ToolUseID == nil || *decoded.ToolUseID != toolUseID {
			t.Errorf("tool_use_id mismatch")
		}
	})

	t.Run("PostToolUseHookInput with subagent context", func(t *testing.T) {
		toolUseID := "tool-456"
		input := &PostToolUseHookInput{
			BaseHookInput: BaseHookInput{SessionID: "s1", TranscriptPath: "/t", CWD: "/c"},
			HookEventName: "PostToolUse",
			ToolName:      "Bash",
			ToolInput:     map[string]interface{}{"command": "ls"},
			ToolResponse:  "file1.txt",
			ToolUseID:     &toolUseID,
			AgentID:       &agentID,
			AgentType:     &agentType,
		}
		data, _ := json.Marshal(input)
		var decoded PostToolUseHookInput
		_ = json.Unmarshal(data, &decoded)
		if decoded.ToolUseID == nil || *decoded.ToolUseID != toolUseID {
			t.Errorf("tool_use_id mismatch")
		}
		if decoded.AgentType == nil || *decoded.AgentType != agentType {
			t.Errorf("agent_type mismatch")
		}
	})

	t.Run("PostToolUseFailureHookInput with is_interrupt", func(t *testing.T) {
		isInterrupt := true
		input := &PostToolUseFailureHookInput{
			BaseHookInput: BaseHookInput{SessionID: "s1", TranscriptPath: "/t", CWD: "/c"},
			HookEventName: "PostToolUseFailure",
			ToolName:      "Bash",
			ToolInput:     map[string]interface{}{"command": "rm"},
			Error:         "denied",
			IsInterrupt:   &isInterrupt,
			AgentID:       &agentID,
			AgentType:     &agentType,
		}
		data, _ := json.Marshal(input)
		var decoded PostToolUseFailureHookInput
		_ = json.Unmarshal(data, &decoded)
		if decoded.IsInterrupt == nil || !*decoded.IsInterrupt {
			t.Errorf("is_interrupt mismatch")
		}
		if decoded.AgentID == nil || *decoded.AgentID != agentID {
			t.Errorf("agent_id mismatch")
		}
	})

	t.Run("SubagentStopHookInput with agent fields", func(t *testing.T) {
		transcriptPath := "/path/to/agent/transcript"
		input := &SubagentStopHookInput{
			BaseHookInput:       BaseHookInput{SessionID: "s1", TranscriptPath: "/t", CWD: "/c"},
			HookEventName:       "SubagentStop",
			StopHookActive:      true,
			AgentID:             agentID,
			AgentType:           &agentType,
			AgentTranscriptPath: &transcriptPath,
		}
		data, _ := json.Marshal(input)
		var decoded SubagentStopHookInput
		_ = json.Unmarshal(data, &decoded)
		if decoded.AgentID != agentID {
			t.Errorf("agent_id mismatch: got %q", decoded.AgentID)
		}
		if decoded.AgentTranscriptPath == nil || *decoded.AgentTranscriptPath != transcriptPath {
			t.Errorf("agent_transcript_path mismatch")
		}
	})

	t.Run("PermissionRequestHookInput with subagent context", func(t *testing.T) {
		input := &PermissionRequestHookInput{
			BaseHookInput: BaseHookInput{SessionID: "s1", TranscriptPath: "/t", CWD: "/c"},
			HookEventName: "PermissionRequest",
			ToolName:      "Write",
			ToolInput:     map[string]interface{}{"path": "/tmp/x"},
			AgentID:       &agentID,
			AgentType:     &agentType,
		}
		data, _ := json.Marshal(input)
		var decoded PermissionRequestHookInput
		_ = json.Unmarshal(data, &decoded)
		if decoded.AgentID == nil || *decoded.AgentID != agentID {
			t.Errorf("agent_id mismatch")
		}
	})
}

// Helper function to create a string pointer.
func stringPtr(s string) *string {
	return &s
}
