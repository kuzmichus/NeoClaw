//go:build !mipsle && !netbsd && !(freebsd && arm)

package agent

import (
	"context"
	"os"
	"path/filepath"

	"github.com/kuzmichus/neoclaw/pkg/providers"
	"github.com/kuzmichus/neoclaw/pkg/seahorse"
)

// resolveSessionContext reconstructs the LLM conversation context for a session
// from the seahorse store: the assembled message history plus the stored
// compression summary. Assemble is a read-only operation (it does not invoke the
// LLM or mutate stored data), so a no-op completion function is sufficient.
//
// When seahorse is unavailable for the workspace, or holds no data for the
// session, it falls back to the raw JSONL history and meta summary so the view
// still reproduces the uncompressed portion of the conversation.
func resolveSessionContext(workspace, sessionKey string, budget int, fallbackHistory []providers.Message, fallbackSummary string) ([]providers.Message, string) {
	if budget <= 0 {
		budget = 200000
	}

	dbPath := filepath.Join(workspace, "sessions", "seahorse.db")
	if _, err := os.Stat(dbPath); err != nil {
		return fallbackHistory, fallbackSummary
	}

	noopComplete := func(ctx context.Context, prompt string, opts seahorse.CompleteOptions) (string, error) {
		return "", nil
	}
	engine, err := seahorse.NewEngine(seahorse.Config{DBPath: dbPath}, noopComplete)
	if err != nil {
		return fallbackHistory, fallbackSummary
	}
	defer engine.Close()

	result, err := engine.Assemble(context.Background(), sessionKey, seahorse.AssembleInput{Budget: budget})
	if err != nil || result == nil {
		return fallbackHistory, fallbackSummary
	}
	if len(result.Messages) == 0 && result.Summary == "" {
		return fallbackHistory, fallbackSummary
	}

	history := seahorseToProviderMessages(result)
	summary := result.Summary
	if summary == "" {
		summary = fallbackSummary
	}
	return history, summary
}
