package agent

import (
	"strings"
	"testing"

	"github.com/kuzmichus/neoclaw/pkg/providers"
)

func TestBuildSessionPromptView(t *testing.T) {
	history := []providers.Message{
		{Role: "user", Content: "hello"},
		{Role: "assistant", Content: "hi there"},
	}

	msgs := BuildSessionPromptView(t.TempDir(), "session-test", 200000, history, "prior summary")

	if len(msgs) == 0 {
		t.Fatal("expected at least one message")
	}
	if msgs[0].Role != "system" {
		t.Fatalf("expected first message to be system, got %q", msgs[0].Role)
	}
	if !strings.Contains(msgs[0].Content, "CONTEXT_SUMMARY") {
		t.Fatal("expected CONTEXT_SUMMARY marker in system prompt")
	}
	if !strings.Contains(msgs[0].Content, "prior summary") {
		t.Fatal("expected summary text to be injected into system prompt")
	}
	if len(msgs) < 3 {
		t.Fatalf("expected system + history messages, got %d", len(msgs))
	}
	if msgs[len(msgs)-1].Role != "assistant" {
		t.Fatalf("expected last history message to be assistant, got %q", msgs[len(msgs)-1].Role)
	}
}
