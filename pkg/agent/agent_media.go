// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package agent

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/h2non/filetype"

	"github.com/kuzmichus/neoclaw/pkg/logger"
	"github.com/kuzmichus/neoclaw/pkg/media"
	"github.com/kuzmichus/neoclaw/pkg/providers"
)

// genericPlaceholderRegex matches generic media placeholders emitted by various
// channels: [image], [image: photo], [image: filename.jpg] — but NOT path tags
// like [image:/path/to/file] (path tags have no space after the colon).
var (
	imagePlaceholderRegex = regexp.MustCompile(`\[image(:\s+[^\]]*)?\]`)
	audioPlaceholderRegex = regexp.MustCompile(`\[audio(:\s+[^\]]*)?\]`)
	videoPlaceholderRegex = regexp.MustCompile(`\[video(:\s+[^\]]*)?\]`)
	filePlaceholderRegex  = regexp.MustCompile(`\[file(:\s+[^\]]*)?\]`)
)

func normalizeCurrentTurnStart(messages []providers.Message, currentTurnStart int) int {
	if currentTurnStart < 0 {
		return 0
	}
	if currentTurnStart > len(messages) {
		return len(messages)
	}
	return currentTurnStart
}

func currentTurnMessages(messages []providers.Message, currentTurnStart int) []providers.Message {
	currentTurnStart = normalizeCurrentTurnStart(messages, currentTurnStart)
	return messages[currentTurnStart:]
}

// maxInlineTextAttachmentBytes caps how large a text-like attachment may be
// for its contents to be inlined into the prompt instead of only exposing a
// [file:/path] tag for tool-based reading.
const maxInlineTextAttachmentBytes = 64 * 1024

// isInlineTextMimeType reports whether files of this MIME type can be read
// as plain text and inlined into the prompt.
func isInlineTextMimeType(mimeType string) bool {
	if strings.HasPrefix(mimeType, "text/") {
		return true
	}
	switch mimeType {
	case "application/json",
		"application/xml",
		"application/x-yaml",
		"application/yaml":
		return true
	}
	return false
}

// readInlineTextFile reads a file as UTF-8 text for prompt inlining.
// Returns ok=false when the file cannot be read, exceeds the inline size
// cap, or is not valid UTF-8.
func readInlineTextFile(localPath string) (string, bool) {
	data, err := os.ReadFile(localPath)
	if err != nil || len(data) == 0 || len(data) > maxInlineTextAttachmentBytes {
		return "", false
	}
	if !utf8.Valid(data) {
		return "", false
	}
	return string(data), true
}

// formatInlineTextBlock renders an inlined text attachment with explicit
// delimiters so the model can tell attached content apart from user text.
func formatInlineTextBlock(filename, localPath, content string) string {
	name := strings.TrimSpace(filename)
	if name == "" {
		name = filepath.Base(localPath)
	}

	var b strings.Builder
	b.WriteString("--- attached file: ")
	b.WriteString(name)
	b.WriteString(" ---\n")
	b.WriteString(content)
	if !strings.HasSuffix(content, "\n") {
		b.WriteString("\n")
	}
	b.WriteString("--- end of attached file ---")
	return b.String()
}

// resolveMediaRefs resolves media:// refs in messages.
// For user messages: images get path tags only ([image:/path]) so the LLM
// can decide whether to view them via load_image or operate on the file.
// For tool messages: images are base64-encoded and appended as a synthetic
// user message only after the contiguous tool-message block ends, so we don't
// break the tool-results-must-immediately-follow-assistant constraint that
// LLM APIs enforce.
// Only tool messages from the current turn may emit the synthetic user
// follow-up; historical tool results stay as plain path-tagged history.
// Non-image files always get path tags regardless of role. Small text-like
// attachments on the current-turn user message are additionally inlined into
// the prompt so the model sees their content without a read_file round trip;
// historical messages only ever get path tags to avoid context bloat.
// Returns a new slice; original messages are not mutated.
func resolveMediaRefs(
	messages []providers.Message,
	store media.MediaStore,
	maxSize int,
	currentTurnStart int,
) []providers.Message {
	if store == nil {
		return messages
	}
	currentTurnStart = normalizeCurrentTurnStart(messages, currentTurnStart)

	result := make([]providers.Message, 0, len(messages))
	var pendingToolImages []string

	for idx, m := range messages {
		// When leaving a tool-message block, flush any accumulated images
		// as a synthetic user message.
		if m.Role != "tool" && len(pendingToolImages) > 0 {
			result = append(result, toolImageFollowUpPromptMessage(pendingToolImages))
			pendingToolImages = nil
		}

		if len(m.Media) == 0 {
			result = append(result, m)
			if idx == len(messages)-1 && len(pendingToolImages) > 0 {
				result = append(result, toolImageFollowUpPromptMessage(pendingToolImages))
				pendingToolImages = nil
			}
			continue
		}

		msg := m
		resolved := make([]string, 0, len(m.Media))
		var pathTags []string
		var inlineBlocks []string

		for _, ref := range m.Media {
			if !strings.HasPrefix(ref, "media://") {
				resolved = append(resolved, ref)
				continue
			}

			localPath, meta, err := store.ResolveWithMeta(ref)
			if err != nil {
				fields := map[string]any{
					"ref":   ref,
					"error": err.Error(),
				}
				if idx < currentTurnStart {
					logger.DebugCF("agent", "Skipped stale historical media ref", fields)
				} else {
					logger.WarnCF("agent", "Failed to resolve media ref", fields)
				}
				continue
			}

			info, err := os.Stat(localPath)
			if err != nil {
				logger.WarnCF("agent", "Failed to stat media file", map[string]any{
					"path":  localPath,
					"error": err.Error(),
				})
				continue
			}

			mime := detectMIME(localPath, meta)
			pathTags = append(pathTags, buildPathTag(mime, localPath))

			if m.Role == "user" &&
				idx >= currentTurnStart &&
				info.Size() <= maxInlineTextAttachmentBytes &&
				isInlineTextMimeType(mime) {
				if text, ok := readInlineTextFile(localPath); ok {
					inlineBlocks = append(
						inlineBlocks,
						formatInlineTextBlock(meta.Filename, localPath, text),
					)
				}
			}

			if m.Role == "tool" && idx >= currentTurnStart && strings.HasPrefix(mime, "image/") {
				dataURL := encodeImageToDataURL(localPath, mime, info, maxSize)
				if dataURL != "" {
					pendingToolImages = append(pendingToolImages, dataURL)
				}
			}
		}

		msg.Media = resolved
		if len(pathTags) > 0 {
			msg.Content = injectPathTags(msg.Content, pathTags)
		}
		if len(inlineBlocks) > 0 {
			msg.Content = strings.Join(append([]string{msg.Content}, inlineBlocks...), "\n\n")
		}
		result = append(result, msg)

		// If this is the last message and we have pending images, flush them.
		if idx == len(messages)-1 && len(pendingToolImages) > 0 {
			result = append(result, toolImageFollowUpPromptMessage(pendingToolImages))
			pendingToolImages = nil
		}
	}

	return result
}

// encodeImageToDataURL base64-encodes an image file into a data URL.
// Returns empty string if the file exceeds maxSize or encoding fails.
func encodeImageToDataURL(localPath, mime string, info os.FileInfo, maxSize int) string {
	if info.Size() > int64(maxSize) {
		logger.WarnCF("agent", "Media file too large, skipping", map[string]any{
			"path":     localPath,
			"size":     info.Size(),
			"max_size": maxSize,
		})
		return ""
	}

	f, err := os.Open(localPath)
	if err != nil {
		logger.WarnCF("agent", "Failed to open media file", map[string]any{
			"path":  localPath,
			"error": err.Error(),
		})
		return ""
	}
	defer f.Close()

	prefix := "data:" + mime + ";base64,"
	encodedLen := base64.StdEncoding.EncodedLen(int(info.Size()))
	var buf bytes.Buffer
	buf.Grow(len(prefix) + encodedLen)
	buf.WriteString(prefix)

	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	_, copyErr := io.Copy(encoder, f)
	closeErr := encoder.Close()
	if copyErr != nil {
		logger.WarnCF("agent", "Failed to encode media file", map[string]any{
			"path":  localPath,
			"error": copyErr.Error(),
		})
		return ""
	}
	if closeErr != nil {
		logger.WarnCF("agent", "Failed to close base64 encoder", map[string]any{
			"path":  localPath,
			"error": closeErr.Error(),
		})
		return ""
	}

	return buf.String()
}

func buildArtifactTags(store media.MediaStore, refs []string) []string {
	if store == nil || len(refs) == 0 {
		return nil
	}

	tags := make([]string, 0, len(refs))
	for _, ref := range refs {
		localPath, meta, err := store.ResolveWithMeta(ref)
		if err != nil {
			continue
		}
		mime := detectMIME(localPath, meta)
		tags = append(tags, buildPathTag(mime, localPath))
	}

	return tags
}

func buildProviderAttachments(store media.MediaStore, refs []string) []providers.Attachment {
	if store == nil || len(refs) == 0 {
		return nil
	}

	attachments := make([]providers.Attachment, 0, len(refs))
	for _, ref := range refs {
		attachment := providers.Attachment{Ref: ref}
		if _, meta, err := store.ResolveWithMeta(ref); err == nil {
			attachment.Filename = meta.Filename
			attachment.ContentType = meta.ContentType
			attachment.Type = inferMediaType(meta.Filename, meta.ContentType)
		}
		attachments = append(attachments, attachment)
	}

	return attachments
}

// detectMIME determines the MIME type from metadata or magic-bytes detection.
// Returns empty string if detection fails.
func detectMIME(localPath string, meta media.MediaMeta) string {
	if meta.ContentType != "" {
		return meta.ContentType
	}
	kind, err := filetype.MatchFile(localPath)
	if err != nil || kind == filetype.Unknown {
		return ""
	}
	return kind.MIME.Value
}

// buildPathTag creates a structured tag exposing the local file path.
// Tag type is derived from MIME: [image:/path], [audio:/path], [video:/path], or [file:/path].
func buildPathTag(mime, localPath string) string {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return "[image:" + localPath + "]"
	case strings.HasPrefix(mime, "audio/"):
		return "[audio:" + localPath + "]"
	case strings.HasPrefix(mime, "video/"):
		return "[video:" + localPath + "]"
	default:
		return "[file:" + localPath + "]"
	}
}

// injectPathTags replaces generic media tags in content with path-bearing versions,
// or appends if no matching generic tag is found. Channels emit a few different
// placeholder formats — [image], [image: photo], [image: filename.jpg] — so we
// match all of them via regex while leaving path tags ([image:/path]) untouched.
//
// When content is structured data (e.g., JSON from Feishu interactive cards or
// post messages), tags are only injected via placeholder replacement — never
// appended — to avoid corrupting the payload.
func injectPathTags(content string, tags []string) string {
	isStructured := looksLikeJSON(content)
	for _, tag := range tags {
		var pattern *regexp.Regexp
		switch {
		case strings.HasPrefix(tag, "[image:"):
			pattern = imagePlaceholderRegex
		case strings.HasPrefix(tag, "[audio:"):
			pattern = audioPlaceholderRegex
		case strings.HasPrefix(tag, "[video:"):
			pattern = videoPlaceholderRegex
		case strings.HasPrefix(tag, "[file:"):
			pattern = filePlaceholderRegex
		}

		if pattern != nil {
			if loc := pattern.FindStringIndex(content); loc != nil {
				content = content[:loc[0]] + tag + content[loc[1]:]
				continue
			}
		}

		if isStructured {
			content = tag + "\n" + content
			continue
		}

		if content == "" {
			content = tag
		} else {
			content += " " + tag
		}
	}
	return content
}

func looksLikeJSON(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 1 && s[0] == '{'
}
