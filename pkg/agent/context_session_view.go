package agent

import "github.com/kuzmichus/neoclaw/pkg/providers"

// BuildSessionPromptView reconstructs the message list that the agent would send
// to the LLM for a stored session. It reuses the live ContextBuilder so the
// system prompt, dynamic context and CONTEXT_SUMMARY injection match production
// behavior.
//
// The conversation history and compression summary are resolved from the seahorse
// store when available (via resolveSessionContext), which correctly reproduces
// sessions whose older messages were compacted out of the JSONL. When seahorse is
// unavailable or holds no data for the session, the raw JSONL history and meta
// summary are used as the fallback.
//
// Tool-call schemas are intentionally omitted (they require the live agent's
// tool registry) and over-budget history trimming is not reproduced, since this
// view is meant for human inspection of the conversation context.
func BuildSessionPromptView(workspace, sessionKey string, budget int, fallbackHistory []providers.Message, fallbackSummary string) []providers.Message {
	history, summary := resolveSessionContext(workspace, sessionKey, budget, fallbackHistory, fallbackSummary)
	cb := NewContextBuilder(workspace)
	req := PromptBuildRequest{
		History: history,
		Summary: summary,
	}
	return cb.BuildMessagesFromPrompt(req)
}
