package gateway

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/kuzmichus/neoclaw/pkg/agent"
	"github.com/kuzmichus/neoclaw/pkg/constants"
	"github.com/kuzmichus/neoclaw/pkg/cron"
)

// cronAPI exposes cron job management over the shared gateway HTTP mux. It
// operates on the in-memory CronService so that scheduling and execution
// remain the single source of truth owned by the running gateway.
type cronAPI struct {
	rs *services
	al *agent.AgentLoop
}

// registerCronRoutes wires the cron management endpoints onto the channel
// manager's shared mux. It must be called after SetupHTTPServerListeners.
func (rs *services) registerCronRoutes(al *agent.AgentLoop) {
	api := &cronAPI{rs: rs, al: al}
	rs.ChannelManager.RegisterRoute("/api/cron/jobs", api.handleCronCollection)
	rs.ChannelManager.RegisterRoute("/api/cron/jobs/", api.handleCronJob)
}

func extractBearerToken(header string) string {
	const prefix = "Bearer "
	if len(header) < len(prefix) {
		return ""
	}
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return header[len(prefix):]
}

func (a *cronAPI) authorized(w http.ResponseWriter, r *http.Request) bool {
	token := a.rs.authToken
	if token == "" {
		return true
	}
	given := extractBearerToken(r.Header.Get("Authorization"))
	if subtle.ConstantTimeCompare([]byte(given), []byte(token)) != 1 {
		writeCronError(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	return true
}

func writeCronError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeCronJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// cronJobRequest is the payload for creating or updating a cron job.
type cronJobRequest struct {
	Name     string            `json:"name"`
	Enabled  bool              `json:"enabled"`
	Schedule cron.CronSchedule `json:"schedule"`
	Message  string            `json:"message"`
	Command  string            `json:"command"`
	Channel  string            `json:"channel"`
	To       string            `json:"to"`
}

func validateSchedule(s cron.CronSchedule) error {
	switch s.Kind {
	case "at":
		if s.AtMS == nil || *s.AtMS <= 0 {
			return errors.New("schedule.atMs is required for kind 'at'")
		}
	case "every":
		if s.EveryMS == nil || *s.EveryMS <= 0 {
			return errors.New("schedule.everyMs is required for kind 'every'")
		}
	case "cron":
		if s.Expr == "" {
			return errors.New("schedule.expr is required for kind 'cron'")
		}
	default:
		return fmt.Errorf("invalid schedule kind %q", s.Kind)
	}
	return nil
}

// validateCommand mirrors the security rule enforced by the cron tool: command
// execution is allowed only when enabled in config and restricted to internal
// or explicitly allow-listed remote channels.
func (a *cronAPI) validateCommand(channel, to, command string) error {
	if command == "" {
		return nil
	}
	cfg := a.al.GetConfig()
	if !cfg.Tools.Cron.AllowCommand {
		return errors.New("command execution is disabled (tools.cron.allow_command)")
	}
	allowedRemotes := cfg.Tools.Cron.CommandAllowedRemotes
	if !constants.IsInternalChannel(channel) && !isCronCommandAllowedRemote(channel, to, allowedRemotes) {
		return errors.New("command execution is restricted to internal channels or configured remote channels")
	}
	return nil
}

func isCronCommandAllowedRemote(channel, chatID string, allowed []string) bool {
	if channel == "" {
		return false
	}
	target := channel
	if chatID != "" {
		target = channel + ":" + chatID
	}
	for _, entry := range allowed {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if entry == "*" || entry == channel || entry == target {
			return true
		}
	}
	return false
}

// handleCronCollection handles GET (list) and POST (create) on /api/cron/jobs.
func (a *cronAPI) handleCronCollection(w http.ResponseWriter, r *http.Request) {
	if !a.authorized(w, r) {
		return
	}
	cs := a.rs.CronService
	if cs == nil {
		writeCronError(w, http.StatusServiceUnavailable, "cron service unavailable")
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeCronJSON(w, http.StatusOK, cs.ListJobs(true))
	case http.MethodPost:
		var req cronJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeCronError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
			return
		}
		if req.Name == "" {
			writeCronError(w, http.StatusBadRequest, "name is required")
			return
		}
		if err := validateSchedule(req.Schedule); err != nil {
			writeCronError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := a.validateCommand(req.Channel, req.To, req.Command); err != nil {
			writeCronError(w, http.StatusForbidden, err.Error())
			return
		}
		job, err := cs.AddJob(req.Name, req.Schedule, req.Message, req.Channel, req.To)
		if err != nil {
			writeCronError(w, http.StatusBadRequest, err.Error())
			return
		}
		if req.Command != "" {
			job.Payload.Command = req.Command
			if err := cs.UpdateJob(job); err != nil {
				writeCronError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		writeCronJSON(w, http.StatusCreated, job)
	default:
		writeCronError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleCronJob handles per-job operations on /api/cron/jobs/{id}.
func (a *cronAPI) handleCronJob(w http.ResponseWriter, r *http.Request) {
	if !a.authorized(w, r) {
		return
	}
	cs := a.rs.CronService
	if cs == nil {
		writeCronError(w, http.StatusServiceUnavailable, "cron service unavailable")
		return
	}

	rest := strings.TrimPrefix(r.URL.Path, "/api/cron/jobs/")
	rest = strings.Trim(rest, "/")
	id := rest
	if id == "" {
		writeCronError(w, http.StatusNotFound, "job not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		job, ok := cs.GetJob(id)
		if !ok {
			writeCronError(w, http.StatusNotFound, "job not found")
			return
		}
		writeCronJSON(w, http.StatusOK, job)
	case http.MethodPut:
		var req cronJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeCronError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
			return
		}
		existing, ok := cs.GetJob(id)
		if !ok {
			writeCronError(w, http.StatusNotFound, "job not found")
			return
		}
		if err := validateSchedule(req.Schedule); err != nil {
			writeCronError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := a.validateCommand(req.Channel, req.To, req.Command); err != nil {
			writeCronError(w, http.StatusForbidden, err.Error())
			return
		}
		existing.Name = req.Name
		existing.Enabled = req.Enabled
		existing.Schedule = req.Schedule
		existing.Payload.Message = req.Message
		existing.Payload.Command = req.Command
		existing.Payload.Channel = req.Channel
		existing.Payload.To = req.To
		if err := cs.UpdateJob(existing); err != nil {
			writeCronError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeCronJSON(w, http.StatusOK, existing)
	case http.MethodDelete:
		if !cs.RemoveJob(id) {
			writeCronError(w, http.StatusNotFound, "job not found")
			return
		}
		writeCronJSON(w, http.StatusOK, map[string]string{"status": "removed"})
	default:
		writeCronError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
