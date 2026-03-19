// internal/api.go
package internal

import (
	"encoding/json"
	"net/http"
	"time"
)

// API provides HTTP endpoints for Evmon
type API struct {
	store Store
}

// NewAPI creates a new API instance
func NewAPI(store Store) *API {
	return &API{store: store}
}

// RegisterRoutes registers HTTP handlers on the given mux
func (api *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/status", api.handleStatus)
	mux.HandleFunc("/history", api.handleHistory)
}

// handleStatus returns the current status of all services
func (api *API) handleStatus(w http.ResponseWriter, r *http.Request) {
	// For simplicity, fetch all services and their current status
	// You might want to optimize this later with a dedicated query

	type ServiceStatus struct {
		ServiceID string `json:"service_id"`
		Status    Status `json:"status"`
	}

	// Placeholder: assuming store has a way to list all services
	// For MVP, you can hardcode or extend Store interface with ListServices
	services := []Service{} // TODO: replace with store.ListServices()

	var results []ServiceStatus
	for _, svc := range services {
		status, err := api.store.GetCurrentStatus(svc.ID)
		if err != nil {
			// If unknown, skip for now
			continue
		}
		results = append(results, ServiceStatus{
			ServiceID: svc.ID,
			Status:    status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// handleHistory returns events for a given service ID and optional time window
func (api *API) handleHistory(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Query().Get("service_id")
	if serviceID == "" {
		http.Error(w, "missing service_id", http.StatusBadRequest)
		return
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to time.Time
	var err error

	if fromStr == "" {
		from = time.Now().Add(-24 * time.Hour)
	} else {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			http.Error(w, "invalid from timestamp", http.StatusBadRequest)
			return
		}
	}

	if toStr == "" {
		to = time.Now()
	} else {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			http.Error(w, "invalid to timestamp", http.StatusBadRequest)
			return
		}
	}

	events, err := api.store.GetEventsInRange(serviceID, from, to)
	if err != nil {
		http.Error(w, "error fetching events: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}