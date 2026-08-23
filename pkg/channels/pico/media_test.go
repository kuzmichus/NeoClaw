package pico

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kuzmichus/neoclaw/pkg/bus"
	"github.com/kuzmichus/neoclaw/pkg/channels"
	"github.com/kuzmichus/neoclaw/pkg/config"
	"github.com/kuzmichus/neoclaw/pkg/media"
)

func newTestPicoChannelWithMediaStore(t *testing.T) (*PicoChannel, *bus.MessageBus, *media.FileMediaStore) {
	t.Helper()

	msgBus := bus.NewMessageBus()
	bc := &config.Channel{Type: config.ChannelPico, Enabled: true}
	cfg := &config.PicoSettings{}
	cfg.SetToken("test-token")
	ch, err := NewPicoChannel(bc, cfg, msgBus)
	if err != nil {
		t.Fatalf("NewPicoChannel: %v", err)
	}
	ch.ctx = context.Background()

	store := media.NewFileMediaStore()
	ch.SetMediaStore(store)
	t.Cleanup(func() {
		_ = os.RemoveAll(media.TempDir())
	})

	return ch, msgBus, store
}

const testPNGDataURL = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+X2ioAAAAASUVORK5CYII="

func TestValidateInlineDataURL_AcceptsDocuments(t *testing.T) {
	valid := []string{
		"data:text/plain;base64,aGVsbG8=",
		"data:text/markdown;base64,IyBUaXRsZQ==",
		"data:text/html;base64,PGh0bWw+",
		"data:application/pdf;base64,JVBERi0xLjQ=",
		"data:application/json;base64,e30=",
	}
	for _, dataURL := range valid {
		if err := validateInlineDataURL(dataURL); err != nil {
			t.Errorf("validateInlineDataURL(%q) error = %v, want nil", dataURL, err)
		}
	}
}

func TestValidateInlineDataURL_RejectsUnsupported(t *testing.T) {
	invalid := map[string]string{
		"":                                    "empty",
		"not-a-data-url":                      "no prefix",
		"data:application/zip;base64,UEsDBA=": "zip not allowed",
		"data:image/svg+xml;base64,PHN2Zz4=":  "svg not allowed",
		"data:application/pdf,JVBERi0xLjQ=":   "not base64",
		"data:application/pdf;base64,!!!":     "invalid base64",
	}
	for dataURL, name := range invalid {
		if err := validateInlineDataURL(dataURL); err == nil {
			t.Errorf("validateInlineDataURL(%q) = nil, want error (%s)", dataURL, name)
		}
	}
}

func TestHandleMessageSend_MaterializesDocumentAttachment(t *testing.T) {
	ch, msgBus, store := newTestPicoChannelWithMediaStore(t)

	payload := "attachment text content"
	ch.handleMessageSend(&picoConn{id: "conn-1", sessionID: "sess-doc"}, PicoMessage{
		Type:      TypeMessageSend,
		ID:        "msg-doc-1",
		SessionID: "sess-doc",
		Payload: map[string]any{
			PayloadKeyContent: "please review",
			"attachments": []any{
				map[string]any{
					"type":     "file",
					"url":      "data:text/plain;base64," + base64.StdEncoding.EncodeToString([]byte(payload)),
					"filename": "notes.txt",
				},
			},
		},
	})

	select {
	case inbound := <-msgBus.InboundChan():
		if inbound.Content != "please review" {
			t.Fatalf("content = %q, want 'please review'", inbound.Content)
		}
		if len(inbound.Media) != 1 {
			t.Fatalf("len(media) = %d, want 1", len(inbound.Media))
		}
		ref := inbound.Media[0]
		if !strings.HasPrefix(ref, "media://") {
			t.Fatalf("media[0] = %q, want media:// ref", ref)
		}
		localPath, meta, err := store.ResolveWithMeta(ref)
		if err != nil {
			t.Fatalf("ResolveWithMeta() error = %v", err)
		}
		if meta.Filename != "notes.txt" || meta.ContentType != "text/plain" {
			t.Fatalf("meta = %#v, want notes.txt/text/plain", meta)
		}
		data, err := os.ReadFile(localPath)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", localPath, err)
		}
		if string(data) != payload {
			t.Fatalf("stored file = %q, want %q", data, payload)
		}
		if filepath.Base(localPath) != "notes.txt" {
			t.Fatalf("stored basename = %q, want notes.txt", filepath.Base(localPath))
		}
		if _, err := os.Stat(localPath); err != nil {
			t.Fatalf("stored file missing: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("expected inbound pico message")
	}
}

func TestHandleMessageSend_KeepsImagesInline(t *testing.T) {
	ch, msgBus, _ := newTestPicoChannelWithMediaStore(t)

	ch.handleMessageSend(&picoConn{id: "conn-1", sessionID: "sess-img"}, PicoMessage{
		Type:      TypeMessageSend,
		ID:        "msg-img-1",
		SessionID: "sess-img",
		Payload: map[string]any{
			"media": []any{testPNGDataURL},
		},
	})

	select {
	case inbound := <-msgBus.InboundChan():
		if len(inbound.Media) != 1 || inbound.Media[0] != testPNGDataURL {
			t.Fatalf("media = %#v, want inline image data URL", inbound.Media)
		}
	case <-time.After(time.Second):
		t.Fatal("expected inbound pico message")
	}
}

func TestHandleMessageSend_RejectsDocumentWithoutMediaStore(t *testing.T) {
	msgBus := bus.NewMessageBus()
	bc := &config.Channel{Type: config.ChannelPico, Enabled: true}
	cfg := &config.PicoSettings{}
	cfg.SetToken("test-token")
	ch, err := NewPicoChannel(bc, cfg, msgBus)
	if err != nil {
		t.Fatalf("NewPicoChannel: %v", err)
	}
	ch.ctx = context.Background()

	pc := &picoConn{id: "conn-1", sessionID: "sess-nostore"}
	pc.closed.Store(true)

	ch.handleMessageSend(pc, PicoMessage{
		Type:      TypeMessageSend,
		ID:        "msg-nostore",
		SessionID: "sess-nostore",
		Payload: map[string]any{
			"attachments": []any{
				map[string]any{
					"type": "file",
					"url":  "data:text/plain;base64,aGVsbG8=",
				},
			},
		},
	})

	select {
	case inbound := <-msgBus.InboundChan():
		t.Fatalf("unexpected inbound message %#v, want none", inbound)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestSanitizeDocumentFilename(t *testing.T) {
	cases := []struct {
		in       string
		mimeType string
		want     string
	}{
		{"report.pdf", "application/pdf", "report.pdf"},
		{"../../etc/passwd", "text/plain", "passwd.txt"},
		{"", "application/pdf", "document.pdf"},
		{"no-extension", "application/pdf", "no-extension.pdf"},
		{"my notes (final).md", "text/markdown", "my notes _final_.md"},
		{strings.Repeat("x", 200) + ".txt", "text/plain", strings.Repeat("x", 116) + ".txt"},
	}
	for _, tc := range cases {
		got := sanitizeDocumentFilename(tc.in, tc.mimeType)
		if got != tc.want {
			t.Errorf("sanitizeDocumentFilename(%q, %q) = %q, want %q", tc.in, tc.mimeType, got, tc.want)
		}
		if got != filepath.Base(got) {
			t.Errorf("sanitizeDocumentFilename(%q) = %q, must be a bare filename", tc.in, got)
		}
	}
}

func TestResolveInboundMedia_ScopeRegistration(t *testing.T) {
	ch, _, store := newTestPicoChannelWithMediaStore(t)

	items := []inlineMedia{
		{dataURL: testPNGDataURL},
		{dataURL: "data:application/pdf;base64,JVBERi0xLjQ=", filename: "doc.pdf"},
	}
	scope := channels.BuildMediaScope("pico", "pico:sess-x", "msg-x")

	resolved, err := ch.resolveInboundMedia(items, scope)
	if err != nil {
		t.Fatalf("resolveInboundMedia() error = %v", err)
	}
	if len(resolved) != 2 {
		t.Fatalf("len(resolved) = %d, want 2", len(resolved))
	}
	if resolved[0] != testPNGDataURL {
		t.Fatalf("resolved[0] = %q, want inline image", resolved[0])
	}
	if !strings.HasPrefix(resolved[1], "media://") {
		t.Fatalf("resolved[1] = %q, want media:// ref", resolved[1])
	}

	if err := store.ReleaseAll(scope); err != nil {
		t.Fatalf("ReleaseAll() error = %v", err)
	}
	if _, _, err := store.ResolveWithMeta(resolved[1]); err == nil {
		t.Fatal("ref should be released with scope")
	}
}
