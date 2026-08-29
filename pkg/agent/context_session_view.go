package agent

import "github.com/kuzmichus/neoclaw/pkg/providers"

// BuildSessionPromptView reconstructs the message list that the agent would send
// to the LLM for a stored session. It reuses the live ContextBuilder so the
// system prompt, dynamic context and CONTEXT_SUMMARY injection match production
// behavior.
//
// Tool-call schemas are intentionally omitted (they require the live agent's
// tool registry) and over-budget history trimming is not reproduced, since this
// view is meant for human inspection of the conversation context.
func BuildSessionPromptView(workspace string, history []providers.Message, summary string) []providers.Message {
	cb := NewContextBuilder(workspace)
	req := PromptBuildRequest{
		History: history,
		Summary: summary,
	}
	return cb.BuildMessagesFromPrompt(req)
}
