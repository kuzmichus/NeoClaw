package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/tools"
)

const (
	mermaidRenderPolicy = `# MERMAID DIAGRAMS\nThe chat interface can render Mermaid diagrams. Whenever a diagram would help the user understand information more clearly (architecture, flows, sequences, state machines, relationships, timelines, comparisons), wrap valid Mermaid source in a fenced code block whose language tag is "mermaid". The diagram is shown to the user visually, with a toggle to view and copy the source code. Do not hesitate to use diagrams for clearer, more illustrative communication.`

	svgRenderPolicy = `# SVG IMAGES\nThe chat interface can render SVG graphics. You can display an SVG image in two ways: (1) wrap valid SVG markup in a fenced code block whose language tag is "svg" — it will be shown to the user both visually and as source code, with a toggle and a "Copy code" button; (2) embed an inline image with a data URI, e.g. ` + "`![alt](data:image/svg+xml;base64,...)`" + `. Whenever an illustration, schematic, chart, icon, or any visual would make the answer clearer or more engaging, produce an SVG and do not hesitate to use it for more illustrative communication with the user.`
)

type toolDiscoveryPromptContributor struct {
	useBM25  bool
	useRegex bool
}

func (c toolDiscoveryPromptContributor) PromptSource() PromptSourceDescriptor {
	return PromptSourceDescriptor{
		ID:              PromptSourceToolDiscovery,
		Owner:           "tools",
		Description:     "Tool discovery instructions",
		Allowed:         []PromptPlacement{{Layer: PromptLayerCapability, Slot: PromptSlotTooling}},
		StableByDefault: true,
	}
}

func (c toolDiscoveryPromptContributor) ContributePrompt(
	_ context.Context,
	req PromptBuildRequest,
) ([]PromptPart, error) {
	if req.SuppressToolUseRule {
		return nil, nil
	}
	useBM25 := c.useBM25 && promptAllowsTool(req, tools.BM25SearchToolName)
	useRegex := c.useRegex && promptAllowsTool(req, tools.RegexSearchToolName)
	if !useBM25 && !useRegex {
		return nil, nil
	}
	content := formatToolDiscoveryRule(useBM25, useRegex)
	if strings.TrimSpace(content) == "" {
		return nil, nil
	}

	return []PromptPart{
		{
			ID:      "capability.tool_discovery",
			Layer:   PromptLayerCapability,
			Slot:    PromptSlotTooling,
			Source:  PromptSource{ID: PromptSourceToolDiscovery, Name: "tool_registry:discovery"},
			Title:   "tool discovery",
			Content: content,
			Stable:  true,
			Cache:   PromptCacheEphemeral,
		},
	}, nil
}

type mcpServerPromptContributor struct {
	serverName string
	toolCount  int
	deferred   bool
}

func (c mcpServerPromptContributor) PromptSource() PromptSourceDescriptor {
	return PromptSourceDescriptor{
		ID:              mcpPromptSourceID(c.serverName),
		Owner:           "mcp",
		Description:     fmt.Sprintf("MCP server %q capability prompt", c.serverName),
		Allowed:         []PromptPlacement{{Layer: PromptLayerCapability, Slot: PromptSlotMCP}},
		StableByDefault: true,
	}
}

func (c mcpServerPromptContributor) ContributePrompt(
	_ context.Context,
	req PromptBuildRequest,
) ([]PromptPart, error) {
	if req.SuppressToolUseRule {
		return nil, nil
	}
	serverName := strings.TrimSpace(c.serverName)
	if serverName == "" || c.toolCount <= 0 {
		return nil, nil
	}
	if len(req.AllowedTools) > 0 &&
		!promptAllowsToolPrefix(req, "mcp_"+promptSourceComponent(serverName)+"_") {
		return nil, nil
	}

	availability := "available as native tools"
	if c.deferred {
		availability = "hidden behind tool discovery until unlocked"
	}

	return []PromptPart{
		{
			ID:     "capability.mcp." + promptSourceComponent(serverName),
			Layer:  PromptLayerCapability,
			Slot:   PromptSlotMCP,
			Source: PromptSource{ID: mcpPromptSourceID(serverName), Name: "mcp:" + serverName},
			Title:  "MCP server capability",
			Content: fmt.Sprintf(
				"MCP server `%s` is connected. It contributes %d tool(s), currently %s.",
				serverName,
				c.toolCount,
				availability,
			),
			Stable: true,
			Cache:  PromptCacheEphemeral,
		},
	}, nil
}

type agentDiscoveryPromptContributor struct {
	agentID  string
	discover func(agentID string) []AgentDescriptor
}

func (c agentDiscoveryPromptContributor) PromptSource() PromptSourceDescriptor {
	return PromptSourceDescriptor{
		ID:              PromptSourceAgentDiscovery,
		Owner:           "agent",
		Description:     "Structured multi-agent discovery registry",
		Allowed:         []PromptPlacement{{Layer: PromptLayerCapability, Slot: PromptSlotTooling}},
		StableByDefault: false,
	}
}

func (c agentDiscoveryPromptContributor) ContributePrompt(
	_ context.Context,
	req PromptBuildRequest,
) ([]PromptPart, error) {
	if req.SuppressToolUseRule {
		return nil, nil
	}
	if !promptAllowsTool(req, "spawn") {
		return nil, nil
	}
	if c.discover == nil {
		return nil, nil
	}
	content := formatAgentDiscoverySection(c.discover(c.agentID))
	if strings.TrimSpace(content) == "" {
		return nil, nil
	}

	return []PromptPart{
		{
			ID:      "capability.agent_discovery",
			Layer:   PromptLayerCapability,
			Slot:    PromptSlotTooling,
			Source:  PromptSource{ID: PromptSourceAgentDiscovery, Name: "agent:discovery"},
			Title:   "agent discovery",
			Content: content,
			Stable:  false,
			Cache:   PromptCacheNone,
		},
	}, nil
}

// channelRichRenderPromptContributor adds the Mermaid/SVG rich-rendering output
// policy parts, but only for channels that can actually render them (the pico
// dashboard and its client). Other channels (telegram, slack, cli, ...) must not
// receive these instructions because their UIs cannot render the diagrams.
type channelRichRenderPromptContributor struct{}

func (c channelRichRenderPromptContributor) PromptSource() PromptSourceDescriptor {
	return PromptSourceDescriptor{
		ID:              PromptSourceOutputPolicy,
		Owner:           "agent",
		Description:     "Channel-aware rich rendering (Mermaid/SVG) output policy",
		Allowed:         []PromptPlacement{{Layer: PromptLayerContext, Slot: PromptSlotOutput}},
		StableByDefault: true,
	}
}

func (c channelRichRenderPromptContributor) ContributePrompt(
	_ context.Context,
	req PromptBuildRequest,
) ([]PromptPart, error) {
	if req.Channel != config.ChannelPico && req.Channel != config.ChannelPicoClient {
		return nil, nil
	}

	return []PromptPart{
		{
			ID:      "context.output_policy.mermaid",
			Layer:   PromptLayerContext,
			Slot:    PromptSlotOutput,
			Source:  PromptSource{ID: PromptSourceOutputPolicy, Name: "mermaid_render"},
			Title:   "mermaid diagram rendering",
			Content: mermaidRenderPolicy,
			Stable:  true,
			Cache:   PromptCacheEphemeral,
		},
		{
			ID:      "context.output_policy.svg",
			Layer:   PromptLayerContext,
			Slot:    PromptSlotOutput,
			Source:  PromptSource{ID: PromptSourceOutputPolicy, Name: "svg_render"},
			Title:   "svg image rendering",
			Content: svgRenderPolicy,
			Stable:  true,
			Cache:   PromptCacheEphemeral,
		},
	}, nil
}

func mcpPromptSourceID(serverName string) PromptSourceID {
	return PromptSourceID("mcp:" + promptSourceComponent(serverName))
}

func promptSourceComponent(value string) string {
	const maxLen = 64

	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "unnamed"
	}

	var b strings.Builder
	lastWasSep := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastWasSep = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastWasSep = false
		case r == '-' || r == '_':
			if !lastWasSep && b.Len() > 0 {
				b.WriteRune(r)
				lastWasSep = true
			}
		default:
			if !lastWasSep && b.Len() > 0 {
				b.WriteRune('_')
				lastWasSep = true
			}
		}
	}

	result := strings.Trim(b.String(), "_")
	if result == "" {
		return "unnamed"
	}
	if len(result) > maxLen {
		return result[:maxLen]
	}
	return result
}

func promptAllowsTool(req PromptBuildRequest, name string) bool {
	if len(req.AllowedTools) == 0 {
		return true
	}
	allowed := cleanAllowedSet(req.AllowedTools)
	_, ok := allowed[strings.ToLower(strings.TrimSpace(name))]
	return ok
}

func promptAllowsToolPrefix(req PromptBuildRequest, prefix string) bool {
	if len(req.AllowedTools) == 0 {
		return true
	}
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	if prefix == "" {
		return false
	}
	for _, name := range req.AllowedTools {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(name)), prefix) {
			return true
		}
	}
	return false
}
