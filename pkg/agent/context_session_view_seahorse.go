//go:build !mipsle && !netbsd && !(freebsd && arm)

package agent

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/kuzmichus/neoclaw/pkg/seahorse"
)

// ResolveSessionSummary returns the seahorse-assembled summary for a session,
// falling back to the provided value when seahorse is not the active manager or
// no seahorse database exists for the workspace. Assemble does not invoke the
// LLM, so a no-op completion function is sufficient.
func ResolveSessionSummary(workspace, sessionKey string, budget int, fallback string) string {
	if strings.TrimSpace(fallback) != "" {
		return fallback
	}
	if budget <= 0 {
		budget = 200000
	}

	dbPath := filepath.Join(workspace, "sessions", "seahorse.db")
	if _, err := os.Stat(dbPath); err != nil {
		return fallback
	}

	noopComplete := func(ctx context.Context, prompt string, opts seahorse.CompleteOptions) (string, error) {
		return "", nil
	}
	engine, err := seahorse.NewEngine(seahorse.Config{DBPath: dbPath}, noopComplete)
	if err != nil {
		return fallback
	}
	defer engine.Close()

	result, err := engine.Assemble(context.Background(), sessionKey, seahorse.AssembleInput{Budget: budget})
	if err != nil || result == nil {
		return fallback
	}
	return result.Summary
}
