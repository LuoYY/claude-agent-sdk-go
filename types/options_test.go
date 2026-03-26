package types

import (
	"encoding/json"
	"testing"
)

// TestWithMaxThinkingTokens tests the WithMaxThinkingTokens builder method.
func TestWithMaxThinkingTokens(t *testing.T) {
	opts := NewClaudeAgentOptions()

	// Test setting max thinking tokens
	result := opts.WithMaxThinkingTokens(5000)

	// Verify the method returns the same instance for chaining
	if result != opts {
		t.Error("WithMaxThinkingTokens should return the same instance for chaining")
	}

	// Verify the value is set correctly
	if opts.MaxThinkingTokens == nil {
		t.Fatal("MaxThinkingTokens should not be nil after setting")
	}

	if *opts.MaxThinkingTokens != 5000 {
		t.Errorf("Expected MaxThinkingTokens to be 5000, got %d", *opts.MaxThinkingTokens)
	}
}

// TestWithMaxBudgetUSD tests the WithMaxBudgetUSD builder method.
func TestWithMaxBudgetUSD(t *testing.T) {
	opts := NewClaudeAgentOptions()

	// Test setting max budget
	result := opts.WithMaxBudgetUSD(10.50)

	// Verify the method returns the same instance for chaining
	if result != opts {
		t.Error("WithMaxBudgetUSD should return the same instance for chaining")
	}

	// Verify the value is set correctly
	if opts.MaxBudgetUSD == nil {
		t.Fatal("MaxBudgetUSD should not be nil after setting")
	}

	if *opts.MaxBudgetUSD != 10.50 {
		t.Errorf("Expected MaxBudgetUSD to be 10.50, got %.2f", *opts.MaxBudgetUSD)
	}
}

// TestOptionsChaining tests that the builder methods can be chained together.
func TestOptionsChaining(t *testing.T) {
	opts := NewClaudeAgentOptions().
		WithMaxThinkingTokens(8000).
		WithMaxBudgetUSD(25.00).
		WithModel("claude-3-5-sonnet-20241022").
		WithMaxTurns(10)

	// Verify all values are set correctly
	if opts.MaxThinkingTokens == nil || *opts.MaxThinkingTokens != 8000 {
		t.Error("MaxThinkingTokens not set correctly in chain")
	}

	if opts.MaxBudgetUSD == nil || *opts.MaxBudgetUSD != 25.00 {
		t.Error("MaxBudgetUSD not set correctly in chain")
	}

	if opts.Model == nil || *opts.Model != "claude-3-5-sonnet-20241022" {
		t.Error("Model not set correctly in chain")
	}

	if opts.MaxTurns == nil || *opts.MaxTurns != 10 {
		t.Error("MaxTurns not set correctly in chain")
	}
}

// TestNewClaudeAgentOptions tests that the constructor creates a valid instance.
func TestNewClaudeAgentOptions(t *testing.T) {
	opts := NewClaudeAgentOptions()

	// Check that optional fields are nil by default
	if opts.MaxThinkingTokens != nil {
		t.Error("MaxThinkingTokens should be nil by default")
	}

	if opts.MaxBudgetUSD != nil {
		t.Error("MaxBudgetUSD should be nil by default")
	}

	// Check that maps are initialized
	if opts.Env == nil {
		t.Error("Env should be initialized")
	}

	if opts.ExtraArgs == nil {
		t.Error("ExtraArgs should be initialized")
	}
}

// TestWithMaxThinkingTokensZeroValue tests that zero values can be set.
func TestWithMaxThinkingTokensZeroValue(t *testing.T) {
	opts := NewClaudeAgentOptions().WithMaxThinkingTokens(0)

	if opts.MaxThinkingTokens == nil {
		t.Fatal("MaxThinkingTokens should not be nil")
	}

	if *opts.MaxThinkingTokens != 0 {
		t.Errorf("Expected MaxThinkingTokens to be 0, got %d", *opts.MaxThinkingTokens)
	}
}

// TestWithMaxBudgetUSDZeroValue tests that zero budget can be set.
func TestWithMaxBudgetUSDZeroValue(t *testing.T) {
	opts := NewClaudeAgentOptions().WithMaxBudgetUSD(0.0)

	if opts.MaxBudgetUSD == nil {
		t.Fatal("MaxBudgetUSD should not be nil")
	}

	if *opts.MaxBudgetUSD != 0.0 {
		t.Errorf("Expected MaxBudgetUSD to be 0.0, got %.2f", *opts.MaxBudgetUSD)
	}
}

// TestPluginConfig tests PluginConfig type and validation.
func TestPluginConfig(t *testing.T) {
	t.Run("NewLocalPluginConfig", func(t *testing.T) {
		plugin := NewLocalPluginConfig("/path/to/plugin")
		if plugin.Type != "local" {
			t.Errorf("expected Type 'local', got %s", plugin.Type)
		}
		if plugin.Path != "/path/to/plugin" {
			t.Errorf("expected Path '/path/to/plugin', got %s", plugin.Path)
		}
	})

	t.Run("NewPluginConfig with valid type", func(t *testing.T) {
		plugin, err := NewPluginConfig("local", "/path/to/plugin")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if plugin.Type != "local" {
			t.Errorf("expected Type 'local', got %s", plugin.Type)
		}
		if plugin.Path != "/path/to/plugin" {
			t.Errorf("expected Path '/path/to/plugin', got %s", plugin.Path)
		}
	})

	t.Run("NewPluginConfig with invalid type", func(t *testing.T) {
		_, err := NewPluginConfig("remote", "/path/to/plugin")
		if err == nil {
			t.Error("expected error for unsupported plugin type")
		}
	})

	t.Run("NewPluginConfig with empty path", func(t *testing.T) {
		_, err := NewPluginConfig("local", "")
		if err == nil {
			t.Error("expected error for empty path")
		}
	})
}

// TestClaudeAgentOptions_Plugins tests plugin builder methods.
func TestClaudeAgentOptions_Plugins(t *testing.T) {
	t.Run("WithPlugins", func(t *testing.T) {
		opts := NewClaudeAgentOptions()
		plugins := []PluginConfig{
			*NewLocalPluginConfig("/path/to/plugin1"),
			*NewLocalPluginConfig("/path/to/plugin2"),
		}
		opts.WithPlugins(plugins)

		if len(opts.Plugins) != 2 {
			t.Errorf("expected 2 plugins, got %d", len(opts.Plugins))
		}
	})

	t.Run("WithPlugin", func(t *testing.T) {
		opts := NewClaudeAgentOptions()
		plugin := *NewLocalPluginConfig("/path/to/plugin")
		opts.WithPlugin(plugin)

		if len(opts.Plugins) != 1 {
			t.Errorf("expected 1 plugin, got %d", len(opts.Plugins))
		}
		if opts.Plugins[0].Path != "/path/to/plugin" {
			t.Errorf("expected Path '/path/to/plugin', got %s", opts.Plugins[0].Path)
		}
	})

	t.Run("WithLocalPlugin", func(t *testing.T) {
		opts := NewClaudeAgentOptions()
		opts.WithLocalPlugin("/path/to/plugin")

		if len(opts.Plugins) != 1 {
			t.Errorf("expected 1 plugin, got %d", len(opts.Plugins))
		}
		if opts.Plugins[0].Type != "local" {
			t.Errorf("expected Type 'local', got %s", opts.Plugins[0].Type)
		}
		if opts.Plugins[0].Path != "/path/to/plugin" {
			t.Errorf("expected Path '/path/to/plugin', got %s", opts.Plugins[0].Path)
		}
	})

	t.Run("multiple plugins via WithPlugin", func(t *testing.T) {
		opts := NewClaudeAgentOptions()
		opts.WithPlugin(*NewLocalPluginConfig("/path/1")).
			WithPlugin(*NewLocalPluginConfig("/path/2")).
			WithPlugin(*NewLocalPluginConfig("/path/3"))

		if len(opts.Plugins) != 3 {
			t.Errorf("expected 3 plugins, got %d", len(opts.Plugins))
		}
	})

	t.Run("multiple plugins via WithLocalPlugin chaining", func(t *testing.T) {
		opts := NewClaudeAgentOptions()
		opts.WithLocalPlugin("/path/1").
			WithLocalPlugin("/path/2").
			WithLocalPlugin("/path/3")

		if len(opts.Plugins) != 3 {
			t.Errorf("expected 3 plugins, got %d", len(opts.Plugins))
		}

		// Verify paths
		expectedPaths := []string{"/path/1", "/path/2", "/path/3"}
		for i, plugin := range opts.Plugins {
			if plugin.Path != expectedPaths[i] {
				t.Errorf("plugin[%d].Path = %s, want %s", i, plugin.Path, expectedPaths[i])
			}
		}
	})

	t.Run("empty plugins by default", func(t *testing.T) {
		opts := NewClaudeAgentOptions()
		if opts.Plugins == nil {
			t.Error("Plugins should not be nil")
		}
		if len(opts.Plugins) != 0 {
			t.Errorf("expected 0 plugins by default, got %d", len(opts.Plugins))
		}
	})
}

// TestWithBetas tests the WithBetas builder method.
func TestWithBetas(t *testing.T) {
	t.Run("WithBetas sets multiple beta flags", func(t *testing.T) {
		opts := NewClaudeAgentOptions()
		betas := []string{"context-1m-2025-08-07"}

		result := opts.WithBetas(betas)

		// Verify the method returns the same instance for chaining
		if result != opts {
			t.Error("WithBetas should return the same instance for chaining")
		}

		// Verify the values are set correctly
		if len(opts.Betas) != 1 {
			t.Errorf("expected 1 beta, got %d", len(opts.Betas))
		}

		if opts.Betas[0] != "context-1m-2025-08-07" {
			t.Errorf("expected beta 'context-1m-2025-08-07', got %s", opts.Betas[0])
		}
	})

	t.Run("WithBetas empty list", func(t *testing.T) {
		opts := NewClaudeAgentOptions().WithBetas([]string{})

		if len(opts.Betas) != 0 {
			t.Errorf("expected 0 betas, got %d", len(opts.Betas))
		}
	})

	t.Run("WithBetas replaces existing betas", func(t *testing.T) {
		opts := NewClaudeAgentOptions().
			WithBeta("beta-1").
			WithBeta("beta-2").
			WithBetas([]string{"beta-3", "beta-4"})

		if len(opts.Betas) != 2 {
			t.Errorf("expected 2 betas after WithBetas, got %d", len(opts.Betas))
		}

		if opts.Betas[0] != "beta-3" || opts.Betas[1] != "beta-4" {
			t.Errorf("expected betas [beta-3, beta-4], got %v", opts.Betas)
		}
	})
}

// TestWithBeta tests the WithBeta builder method.
func TestWithBeta(t *testing.T) {
	t.Run("WithBeta adds single beta flag", func(t *testing.T) {
		opts := NewClaudeAgentOptions()

		result := opts.WithBeta("context-1m-2025-08-07")

		// Verify the method returns the same instance for chaining
		if result != opts {
			t.Error("WithBeta should return the same instance for chaining")
		}

		// Verify the value is set correctly
		if len(opts.Betas) != 1 {
			t.Errorf("expected 1 beta, got %d", len(opts.Betas))
		}

		if opts.Betas[0] != "context-1m-2025-08-07" {
			t.Errorf("expected beta 'context-1m-2025-08-07', got %s", opts.Betas[0])
		}
	})

	t.Run("WithBeta multiple calls accumulate", func(t *testing.T) {
		opts := NewClaudeAgentOptions().
			WithBeta("beta-1").
			WithBeta("beta-2").
			WithBeta("beta-3")

		if len(opts.Betas) != 3 {
			t.Errorf("expected 3 betas, got %d", len(opts.Betas))
		}

		expectedBetas := []string{"beta-1", "beta-2", "beta-3"}
		for i, beta := range opts.Betas {
			if beta != expectedBetas[i] {
				t.Errorf("beta[%d] = %s, expected %s", i, beta, expectedBetas[i])
			}
		}
	})

	t.Run("empty betas by default", func(t *testing.T) {
		opts := NewClaudeAgentOptions()
		if opts.Betas == nil {
			t.Error("Betas should not be nil")
		}
		if len(opts.Betas) != 0 {
			t.Errorf("expected 0 betas by default, got %d", len(opts.Betas))
		}
	})
}

// TestAgentDefinitionSkills tests the Skills field on AgentDefinition.
func TestAgentDefinitionSkills(t *testing.T) {
	t.Run("agent with skills", func(t *testing.T) {
		agent := AgentDefinition{
			Description: "Test agent",
			Prompt:      "Test prompt",
			Skills:      []string{"code-review", "testing", "documentation"},
		}

		if len(agent.Skills) != 3 {
			t.Errorf("expected 3 skills, got %d", len(agent.Skills))
		}
		if agent.Skills[0] != "code-review" {
			t.Errorf("expected first skill to be 'code-review', got %s", agent.Skills[0])
		}
	})

	t.Run("agent without skills omits field in JSON", func(t *testing.T) {
		agent := AgentDefinition{
			Description: "Test agent",
			Prompt:      "Test prompt",
		}

		data, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		// Skills should be omitted from JSON when nil
		var decoded map[string]interface{}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if _, exists := decoded["skills"]; exists {
			t.Error("skills should be omitted from JSON when nil")
		}
	})

	t.Run("agent with skills serializes correctly", func(t *testing.T) {
		agent := AgentDefinition{
			Description: "Test agent",
			Prompt:      "Test prompt",
			Skills:      []string{"coding", "review"},
			Tools:       []string{"Read", "Write"},
		}

		data, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		var decoded AgentDefinition
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if len(decoded.Skills) != 2 {
			t.Errorf("expected 2 skills after roundtrip, got %d", len(decoded.Skills))
		}
		if len(decoded.Tools) != 2 {
			t.Errorf("expected 2 tools after roundtrip, got %d", len(decoded.Tools))
		}
	})

	t.Run("WithAgent includes skills", func(t *testing.T) {
		opts := NewClaudeAgentOptions().
			WithAgent("reviewer", AgentDefinition{
				Description: "Code reviewer",
				Prompt:      "Review code",
				Skills:      []string{"code-review"},
			})

		agent, ok := opts.Agents["reviewer"]
		if !ok {
			t.Fatal("agent 'reviewer' not found")
		}
		if len(agent.Skills) != 1 || agent.Skills[0] != "code-review" {
			t.Errorf("expected skills [code-review], got %v", agent.Skills)
		}
	})
}

// TestThinkingConfig tests the ThinkingConfig constructors and serialization.
func TestThinkingConfig(t *testing.T) {
	t.Run("adaptive", func(t *testing.T) {
		config := NewThinkingAdaptive()
		if config.Type != "adaptive" {
			t.Errorf("expected type 'adaptive', got %s", config.Type)
		}
		if config.BudgetTokens != nil {
			t.Error("BudgetTokens should be nil for adaptive")
		}
	})

	t.Run("enabled", func(t *testing.T) {
		config := NewThinkingEnabled(10000)
		if config.Type != "enabled" {
			t.Errorf("expected type 'enabled', got %s", config.Type)
		}
		if config.BudgetTokens == nil || *config.BudgetTokens != 10000 {
			t.Errorf("BudgetTokens should be 10000")
		}
	})

	t.Run("disabled", func(t *testing.T) {
		config := NewThinkingDisabled()
		if config.Type != "disabled" {
			t.Errorf("expected type 'disabled', got %s", config.Type)
		}
	})

	t.Run("JSON roundtrip", func(t *testing.T) {
		config := NewThinkingEnabled(5000)
		data, err := json.Marshal(config)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}
		var decoded ThinkingConfig
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if decoded.Type != "enabled" || decoded.BudgetTokens == nil || *decoded.BudgetTokens != 5000 {
			t.Error("roundtrip mismatch")
		}
	})
}

// TestEffortLevel tests effort level constants.
func TestEffortLevel(t *testing.T) {
	levels := map[EffortLevel]string{
		EffortLow:    "low",
		EffortMedium: "medium",
		EffortHigh:   "high",
		EffortMax:    "max",
	}
	for level, expected := range levels {
		if string(level) != expected {
			t.Errorf("EffortLevel %q should be %q", level, expected)
		}
	}
}

// TestNewOptionsFields tests builder methods for new ClaudeAgentOptions fields.
func TestNewOptionsFields(t *testing.T) {
	t.Run("WithThinking", func(t *testing.T) {
		opts := NewClaudeAgentOptions().WithThinking(NewThinkingAdaptive())
		if opts.Thinking == nil || opts.Thinking.Type != "adaptive" {
			t.Error("Thinking should be adaptive")
		}
	})

	t.Run("WithEffort", func(t *testing.T) {
		opts := NewClaudeAgentOptions().WithEffort(EffortHigh)
		if opts.Effort == nil || *opts.Effort != EffortHigh {
			t.Error("Effort should be high")
		}
	})

	t.Run("WithFallbackModel", func(t *testing.T) {
		opts := NewClaudeAgentOptions().WithFallbackModel("claude-3-5-haiku-latest")
		if opts.FallbackModel == nil || *opts.FallbackModel != "claude-3-5-haiku-latest" {
			t.Error("FallbackModel mismatch")
		}
	})

	t.Run("WithOutputFormat", func(t *testing.T) {
		schema := map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"answer": map[string]interface{}{"type": "string"},
			},
		}
		opts := NewClaudeAgentOptions().WithOutputFormat(schema)
		if opts.OutputFormat == nil {
			t.Error("OutputFormat should not be nil")
		}
	})

	t.Run("WithSandbox", func(t *testing.T) {
		opts := NewClaudeAgentOptions().WithSandbox(&SandboxSettings{Type: "docker"})
		if opts.Sandbox == nil || opts.Sandbox.Type != "docker" {
			t.Error("Sandbox should be docker")
		}
	})

	t.Run("WithEnableFileCheckpointing", func(t *testing.T) {
		opts := NewClaudeAgentOptions().WithEnableFileCheckpointing(true)
		if !opts.EnableFileCheckpointing {
			t.Error("EnableFileCheckpointing should be true")
		}
	})

	t.Run("WithSystemPromptFile", func(t *testing.T) {
		opts := NewClaudeAgentOptions().WithSystemPromptFile("/path/to/prompt.md")
		spf, ok := opts.SystemPrompt.(SystemPromptFile)
		if !ok {
			t.Fatal("SystemPrompt should be SystemPromptFile")
		}
		if spf.Type != "file" || spf.Path != "/path/to/prompt.md" {
			t.Errorf("SystemPromptFile mismatch: %+v", spf)
		}
	})
}

// TestAgentDefinitionNewFields tests memory and mcpServers fields on AgentDefinition.
func TestAgentDefinitionNewFields(t *testing.T) {
	t.Run("with memory", func(t *testing.T) {
		memory := "project"
		agent := AgentDefinition{
			Description: "Test",
			Prompt:      "Test",
			Memory:      &memory,
		}
		data, _ := json.Marshal(agent)
		var decoded map[string]interface{}
		_ = json.Unmarshal(data, &decoded)
		if decoded["memory"] != "project" {
			t.Errorf("memory should be 'project' in JSON")
		}
	})

	t.Run("with mcpServers", func(t *testing.T) {
		agent := AgentDefinition{
			Description: "Test",
			Prompt:      "Test",
			McpServers:  []interface{}{"server1", map[string]interface{}{"url": "http://localhost"}},
		}
		data, _ := json.Marshal(agent)
		var decoded AgentDefinition
		_ = json.Unmarshal(data, &decoded)
		if len(decoded.McpServers) != 2 {
			t.Errorf("expected 2 mcpServers, got %d", len(decoded.McpServers))
		}
	})

	t.Run("omits empty fields", func(t *testing.T) {
		agent := AgentDefinition{Description: "Test", Prompt: "Test"}
		data, _ := json.Marshal(agent)
		var decoded map[string]interface{}
		_ = json.Unmarshal(data, &decoded)
		for _, field := range []string{"memory", "mcpServers", "skills"} {
			if _, exists := decoded[field]; exists {
				t.Errorf("%s should be omitted from JSON when empty", field)
			}
		}
	})
}

// TestHookMatcherTimeout tests the timeout field on HookMatcher.
func TestHookMatcherTimeout(t *testing.T) {
	timeout := 30.0
	matcher := HookMatcher{
		Matcher: stringPtr("Bash"),
		Timeout: &timeout,
	}

	data, err := json.Marshal(matcher)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded["timeout"] != 30.0 {
		t.Errorf("expected timeout 30.0, got %v", decoded["timeout"])
	}
	if decoded["matcher"] != "Bash" {
		t.Errorf("expected matcher 'Bash', got %v", decoded["matcher"])
	}
}

// TestSubagentExecutionModeConstants tests the SubagentExecutionMode enum values.
func TestSubagentExecutionModeConstants(t *testing.T) {
	tests := []struct {
		mode     SubagentExecutionMode
		expected string
	}{
		{SubagentExecutionModeSequential, "sequential"},
		{SubagentExecutionModeParallel, "parallel"},
		{SubagentExecutionModeAuto, "auto"},
	}

	for _, tt := range tests {
		if string(tt.mode) != tt.expected {
			t.Errorf("SubagentExecutionMode = %q, expected %q", string(tt.mode), tt.expected)
		}
	}
}

// TestMultiInvocationModeConstants tests the MultiInvocationMode enum values.
func TestMultiInvocationModeConstants(t *testing.T) {
	tests := []struct {
		mode     MultiInvocationMode
		expected string
	}{
		{MultiInvocationModeSequential, "sequential"},
		{MultiInvocationModeParallel, "parallel"},
		{MultiInvocationModeError, "error"},
	}

	for _, tt := range tests {
		if string(tt.mode) != tt.expected {
			t.Errorf("MultiInvocationMode = %q, expected %q", string(tt.mode), tt.expected)
		}
	}
}

// TestSubagentErrorHandlingConstants tests the SubagentErrorHandling enum values.
func TestSubagentErrorHandlingConstants(t *testing.T) {
	tests := []struct {
		mode     SubagentErrorHandling
		expected string
	}{
		{SubagentErrorHandlingFailFast, "fail_fast"},
		{SubagentErrorHandlingContinue, "continue"},
	}

	for _, tt := range tests {
		if string(tt.mode) != tt.expected {
			t.Errorf("SubagentErrorHandling = %q, expected %q", string(tt.mode), tt.expected)
		}
	}
}

// TestNewSubagentExecutionConfig tests the NewSubagentExecutionConfig constructor.
func TestNewSubagentExecutionConfig(t *testing.T) {
	t.Run("creates config with sensible defaults", func(t *testing.T) {
		config := NewSubagentExecutionConfig()

		if config.MultiInvocation != MultiInvocationModeSequential {
			t.Errorf("expected MultiInvocation to be sequential, got %s", config.MultiInvocation)
		}

		if config.MaxConcurrent != 3 {
			t.Errorf("expected MaxConcurrent to be 3, got %d", config.MaxConcurrent)
		}

		if config.ErrorHandling != SubagentErrorHandlingContinue {
			t.Errorf("expected ErrorHandling to be continue, got %s", config.ErrorHandling)
		}
	})

	t.Run("can be customized after creation", func(t *testing.T) {
		config := NewSubagentExecutionConfig()
		config.MultiInvocation = MultiInvocationModeParallel
		config.MaxConcurrent = 5
		config.ErrorHandling = SubagentErrorHandlingFailFast

		if config.MultiInvocation != MultiInvocationModeParallel {
			t.Errorf("expected MultiInvocation to be parallel, got %s", config.MultiInvocation)
		}

		if config.MaxConcurrent != 5 {
			t.Errorf("expected MaxConcurrent to be 5, got %d", config.MaxConcurrent)
		}

		if config.ErrorHandling != SubagentErrorHandlingFailFast {
			t.Errorf("expected ErrorHandling to be fail_fast, got %s", config.ErrorHandling)
		}
	})
}

// TestAgentDefinitionWithExecutionControl tests AgentDefinition with new execution control fields.
func TestAgentDefinitionWithExecutionControl(t *testing.T) {
	t.Run("agent with execution mode", func(t *testing.T) {
		mode := SubagentExecutionModeParallel
		agent := AgentDefinition{
			Description:   "Test agent",
			Prompt:        "Test prompt",
			ExecutionMode: &mode,
		}

		if agent.ExecutionMode == nil {
			t.Fatal("ExecutionMode should not be nil")
		}

		if *agent.ExecutionMode != SubagentExecutionModeParallel {
			t.Errorf("expected ExecutionMode to be parallel, got %s", *agent.ExecutionMode)
		}
	})

	t.Run("agent with timeout", func(t *testing.T) {
		timeout := 30.5
		agent := AgentDefinition{
			Description: "Test agent",
			Prompt:      "Test prompt",
			Timeout:     &timeout,
		}

		if agent.Timeout == nil {
			t.Fatal("Timeout should not be nil")
		}

		if *agent.Timeout != 30.5 {
			t.Errorf("expected Timeout to be 30.5, got %f", *agent.Timeout)
		}
	})

	t.Run("agent with max turns", func(t *testing.T) {
		maxTurns := 5
		agent := AgentDefinition{
			Description: "Test agent",
			Prompt:      "Test prompt",
			MaxTurns:    &maxTurns,
		}

		if agent.MaxTurns == nil {
			t.Fatal("MaxTurns should not be nil")
		}

		if *agent.MaxTurns != 5 {
			t.Errorf("expected MaxTurns to be 5, got %d", *agent.MaxTurns)
		}
	})

	t.Run("agent with all execution control fields", func(t *testing.T) {
		mode := SubagentExecutionModeSequential
		timeout := 60.0
		maxTurns := 10
		agent := AgentDefinition{
			Description:   "Full agent",
			Prompt:        "Full prompt",
			Tools:         []string{"Read", "Write"},
			ExecutionMode: &mode,
			Timeout:       &timeout,
			MaxTurns:      &maxTurns,
		}

		if agent.ExecutionMode == nil || *agent.ExecutionMode != SubagentExecutionModeSequential {
			t.Errorf("ExecutionMode mismatch")
		}

		if agent.Timeout == nil || *agent.Timeout != 60.0 {
			t.Errorf("Timeout mismatch")
		}

		if agent.MaxTurns == nil || *agent.MaxTurns != 10 {
			t.Errorf("MaxTurns mismatch")
		}
	})
}

// TestWithSubagentExecution tests the WithSubagentExecution builder method.
func TestWithSubagentExecution(t *testing.T) {
	t.Run("sets subagent execution config", func(t *testing.T) {
		opts := NewClaudeAgentOptions()
		config := NewSubagentExecutionConfig()
		config.MaxConcurrent = 5

		result := opts.WithSubagentExecution(config)

		// Verify the method returns the same instance for chaining
		if result != opts {
			t.Error("WithSubagentExecution should return the same instance for chaining")
		}

		// Verify the value is set
		if opts.SubagentExecution == nil {
			t.Fatal("SubagentExecution should not be nil after setting")
		}

		if opts.SubagentExecution.MaxConcurrent != 5 {
			t.Errorf("expected MaxConcurrent to be 5, got %d", opts.SubagentExecution.MaxConcurrent)
		}
	})

	t.Run("replaces existing config", func(t *testing.T) {
		opts := NewClaudeAgentOptions()

		config1 := NewSubagentExecutionConfig()
		config1.MaxConcurrent = 2
		opts.WithSubagentExecution(config1)

		config2 := NewSubagentExecutionConfig()
		config2.MaxConcurrent = 8
		opts.WithSubagentExecution(config2)

		if opts.SubagentExecution.MaxConcurrent != 8 {
			t.Errorf("expected MaxConcurrent to be 8 after replacement, got %d", opts.SubagentExecution.MaxConcurrent)
		}
	})

	t.Run("method chaining works", func(t *testing.T) {
		config := NewSubagentExecutionConfig()
		config.MultiInvocation = MultiInvocationModeParallel

		opts := NewClaudeAgentOptions().
			WithModel("claude-opus-4-5-latest").
			WithSubagentExecution(config).
			WithAgent("test", AgentDefinition{
				Description: "Test",
				Prompt:      "Test",
			})

		if opts.SubagentExecution == nil {
			t.Fatal("SubagentExecution should be set")
		}

		if opts.SubagentExecution.MultiInvocation != MultiInvocationModeParallel {
			t.Errorf("expected MultiInvocation to be parallel")
		}

		if opts.Model == nil || *opts.Model != "claude-opus-4-5-latest" {
			t.Errorf("Model should be set to claude-opus-4-5-latest")
		}

		if _, ok := opts.Agents["test"]; !ok {
			t.Errorf("Agent 'test' should be set")
		}
	})
}
