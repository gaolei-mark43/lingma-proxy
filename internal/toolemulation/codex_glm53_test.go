package toolemulation

import (
	"strings"
	"testing"
)

func TestInjectToolingFramesToolsAsExternalExecutorProtocol(t *testing.T) {
	prompt := InjectTooling("", []ToolDef{{
		Name:        "exec_command",
		Description: "Run a shell command",
		InputSchema: map[string]any{
			"properties": map[string]any{
				"cmd": map[string]any{"type": "string"},
			},
			"required": []any{"cmd"},
		},
	}}, ToolChoice{Mode: "auto"}, nil)

	for _, want := range []string{
		"IMPORTANT EXTERNAL TOOL PROTOCOL",
		"NOT QoderCN/Lingma native tools",
		"external Codex/client executor",
		"Native QoderCN/Lingma tool availability is irrelevant",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q:\n%s", want, prompt)
		}
	}
}

func TestForceToolingPromptFramesExternalProtocol(t *testing.T) {
	prompt := ForceToolingPrompt(ToolChoice{Mode: "auto"})
	for _, want := range []string{
		"EXTERNAL tool protocol",
		"not a QoderCN/Lingma native tool call",
		"Do not ask the user to execute anything manually",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("force prompt missing %q:\n%s", want, prompt)
		}
	}
}

func TestInferToolCallsSupportsCodexExecCommandCmdArgument(t *testing.T) {
	calls := InferToolCallsFromText(
		"当前会话没有可用的命令执行工具。你可以运行 vm_stat 查看内存占用。",
		[]ToolDef{{
			Name: "exec_command",
			InputSchema: map[string]any{
				"properties": map[string]any{
					"cmd": map[string]any{"type": "string"},
				},
				"required": []any{"cmd"},
			},
		}},
	)
	if len(calls) != 1 {
		t.Fatalf("call count = %d, calls = %+v", len(calls), calls)
	}
	if calls[0].Name != "exec_command" {
		t.Fatalf("tool name = %q", calls[0].Name)
	}
	cmd, _ := calls[0].Arguments["cmd"].(string)
	if !strings.Contains(cmd, "vm_stat") {
		t.Fatalf("unexpected cmd = %q", cmd)
	}
	if _, exists := calls[0].Arguments["command"]; exists {
		t.Fatalf("unexpected legacy command arg: %+v", calls[0].Arguments)
	}
}

func TestLooksLikeRefusalRecognizesQoderNoToolMode(t *testing.T) {
	for _, text := range []string{
		"当前会话处于无工具可用模式。",
		"当前会话没有可用的命令执行工具。",
	} {
		if !LooksLikeRefusal(text) {
			t.Fatalf("LooksLikeRefusal(%q) = false", text)
		}
	}
}
