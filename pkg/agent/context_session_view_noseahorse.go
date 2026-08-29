//go:build mipsle || netbsd || (freebsd && arm)

package agent

// ResolveSessionSummary returns the fallback on platforms where the seahorse
// short-term memory engine is unavailable.
func ResolveSessionSummary(workspace, sessionKey string, budget int, fallback string) string {
	return fallback
}
