package gateway

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/kuzmichus/neoclaw/pkg/cron"
)

func int64Ptr(v int64) *int64 { return &v }

func newTestCronAPI(t *testing.T, authToken string) *cronAPI {
	t.Helper()
	storePath := filepath.Join(t.TempDir(), "jobs.json")
	cs := cron.NewCronService(storePath, nil)
	if err := cs.Start(); err != nil {
		t.Fatalf("failed to start cron service: %v", err)
	}
	t.Cleanup(cs.Stop)
	return &cronAPI{
		rs: &services{CronService: cs, authToken: authToken},
	}
}

func TestCronAPI_Lifecycle(t *testing.T) {
	api := newTestCronAPI(t, "")

	body := cronJobRequest{
		Name:     "daily",
		Enabled:  true,
		Schedule: cron.CronSchedule{Kind: "cron", Expr: "0 0 * * *"},
		Message:  "hello",
		Channel:  "cli",
		To:       "session-1",
	}
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/cron/jobs", bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	api.handleCronCollection(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body %s", rec.Code, rec.Body.String())
	}
	var created cron.CronJob
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected job id")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/cron/jobs", nil)
	rec = httptest.NewRecorder()
	api.handleCronCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d", rec.Code)
	}
	var list []cron.CronJob
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 job, got %d", len(list))
	}

	req = httptest.NewRequest(http.MethodGet, "/api/cron/jobs/"+created.ID, nil)
	rec = httptest.NewRecorder()
	api.handleCronJob(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d", rec.Code)
	}

	body.Message = "updated"
	raw, _ = json.Marshal(body)
	req = httptest.NewRequest(http.MethodPut, "/api/cron/jobs/"+created.ID, bytes.NewReader(raw))
	rec = httptest.NewRecorder()
	api.handleCronJob(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, body %s", rec.Code, rec.Body.String())
	}
	var updated cron.CronJob
	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode updated: %v", err)
	}
	if updated.Payload.Message != "updated" {
		t.Fatalf("expected updated message, got %q", updated.Payload.Message)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/cron/jobs/"+created.ID, nil)
	rec = httptest.NewRecorder()
	api.handleCronJob(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("delete status = %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/cron/jobs/"+created.ID, nil)
	rec = httptest.NewRecorder()
	api.handleCronJob(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", rec.Code)
	}
}

func TestCronAPI_RequiresAuth(t *testing.T) {
	api := newTestCronAPI(t, "secret")

	req := httptest.NewRequest(http.MethodGet, "/api/cron/jobs", nil)
	rec := httptest.NewRecorder()
	api.handleCronCollection(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/cron/jobs", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec = httptest.NewRecorder()
	api.handleCronCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with token, got %d", rec.Code)
	}
}

func TestCronAPI_ValidateSchedule(t *testing.T) {
	api := newTestCronAPI(t, "")

	body := cronJobRequest{
		Name:     "bad",
		Schedule: cron.CronSchedule{Kind: "cron"}, // missing expr
		Message:  "x",
	}
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/cron/jobs", bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	api.handleCronCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing expr, got %d", rec.Code)
	}
}

func TestCronAPI_EveryInterval(t *testing.T) {
	api := newTestCronAPI(t, "")
	body := cronJobRequest{
		Name:     "interval",
		Schedule: cron.CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)},
		Message:  "ping",
	}
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/cron/jobs", bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	api.handleCronCollection(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body %s", rec.Code, rec.Body.String())
	}
}
