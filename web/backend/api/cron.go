package api

import (
	"net/http"
	"net/http/httputil"

	"github.com/kuzmichus/neoclaw/pkg/logger"
)

// registerCronRoutes exposes cron management by proxying to the running
// gateway, which owns the CronService as the single source of truth.
func (h *Handler) registerCronRoutes(mux *http.ServeMux) {
	mux.Handle("/api/cron/", h.handleCronProxy())
}

// handleCronProxy forwards /api/cron/* to the gateway's cron REST API,
// injecting the gateway auth token so the request is authorized upstream.
func (h *Handler) handleCronProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.gatewayAvailableForProxy() {
			http.Error(w, "Gateway unavailable", http.StatusServiceUnavailable)
			return
		}

		gateway.mu.Lock()
		token := ""
		if gateway.pidData != nil {
			token = gateway.pidData.Token
		}
		gateway.mu.Unlock()

		target := h.gatewayProxyURL()
		proxy := &httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetURL(target)
				if token != "" {
					pr.Out.Header.Set("Authorization", "Bearer "+token)
				}
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				logger.Errorf("Failed to proxy cron request: %v", err)
				http.Error(w, "Gateway unavailable: "+err.Error(), http.StatusBadGateway)
			},
		}
		proxy.ServeHTTP(w, r)
	}
}
