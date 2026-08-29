//go:build mipsle || netbsd || (freebsd && arm)

package agent

import "github.com/kuzmichus/neoclaw/pkg/providers"

// resolveSessionContext falls back to the raw JSONL history and meta summary on
// platforms where the seahorse short-term memory engine is unavailable.
func resolveSessionContext(workspace, sessionKey string, budget int, fallbackHistory []providers.Message, fallbackSummary string) ([]providers.Message, string) {
	return fallbackHistory, fallbackSummary
}
