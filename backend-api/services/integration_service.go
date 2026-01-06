package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"idash-backend-api/database"
)

type integrationService struct {
	mutex      sync.RWMutex
	running    bool
	lastStatus IntegrationStatus
}

var integrationOnce sync.Once
var integrationInstance *integrationService

// Integration returns the singleton integration service
func Integration() *integrationService {
	integrationOnce.Do(func() {
		integrationInstance = &integrationService{
			lastStatus: IntegrationStatus{Status: "idle", LastUpdated: time.Now()},
		}
	})
	return integrationInstance
}

// IntegrationStatus describes the state of the integration pipeline
type IntegrationStatus struct {
	Status      string      `json:"status"`
	StartedBy   interface{} `json:"started_by"`
	StartedAt   *time.Time  `json:"started_at"`
	CompletedAt *time.Time  `json:"completed_at"`
	Duration    string      `json:"duration"`
	Message     string      `json:"message"`
	Errors      []string    `json:"errors"`
	LastUpdated time.Time   `json:"last_updated"`
}

// Status returns a copy of the latest status snapshot
func (s *integrationService) Status() IntegrationStatus {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastStatus
}

// IsRunning returns true when a sync is already in progress
func (s *integrationService) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// RunFullSync executes the full data refresh + AI regeneration pipeline
func (s *integrationService) RunFullSync(startedBy interface{}) {
	s.mutex.Lock()
	if s.running {
		s.mutex.Unlock()
		return
	}
	s.running = true
	startedAt := time.Now()
	s.lastStatus = IntegrationStatus{
		Status:      "running",
		StartedBy:   startedBy,
		StartedAt:   &startedAt,
		LastUpdated: time.Now(),
	}
	s.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	errs := make([]string, 0)

	syncSteps := []struct {
		label  string
		action func(context.Context) error
	}{
		{"Refresh materialized views", refreshMaterializedViews},
		{"Regenerate forecasts", regenerateForecasts},
		{"Refresh cached metrics", refreshCachedMetrics},
	}

	for _, step := range syncSteps {
		if err := step.action(ctx); err != nil {
			errs = append(errs, step.label+": "+err.Error())
			log.Printf("[integration] step %s failed: %v", step.label, err)
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.running = false
	completed := time.Now()
	s.lastStatus.CompletedAt = &completed
	s.lastStatus.Duration = completed.Sub(*s.lastStatus.StartedAt).Round(time.Second).String()
	s.lastStatus.LastUpdated = time.Now()

	if len(errs) > 0 {
		s.lastStatus.Status = "completed_with_errors"
		s.lastStatus.Errors = errs
		s.lastStatus.Message = "Integration completed with warnings"
	} else {
		s.lastStatus.Status = "completed"
		s.lastStatus.Message = "Integration pipeline completed successfully"
	}
}

func refreshMaterializedViews(ctx context.Context) error {
	err := database.DB.WithContext(ctx).Exec("CALL refresh_materialized_views()").Error
	if err != nil {
		// tolerate different MySQL messages when the stored procedure is absent
		le := strings.ToLower(err.Error())
		if strings.Contains(le, "doesn't exist") || strings.Contains(le, "does not exist") || strings.Contains(le, "procedure") || strings.Contains(le, "error 1305") {
			// missing procedure is non-fatal for environments without materialized view support
			return nil
		}
		return err
	}
	return nil
}

func regenerateForecasts(ctx context.Context) error {
	types := []string{"sales", "production", "finance", "inventory", "hr", "supplychain"}
	for _, t := range types {
		if err := triggerForecast(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

func triggerForecast(ctx context.Context, forecastType string) error {
	url := os.Getenv("AI_SERVICE_URL")
	if url == "" {
		url = "http://localhost:8000"
	}
	// Map forecast types to ai-service endpoints
	supported := map[string]string{
		"sales":      "/api/forecast/sales",
		"production": "/api/forecast/production",
		"finance":    "/api/forecast/finance",
		"inventory":  "/api/forecast/inventory",
	}

	path, ok := supported[forecastType]
	if !ok {
		// If the AI service doesn't support this forecast type, skip with a logged warning
		log.Printf("[integration] skipping unsupported forecast type: %s", forecastType)
		return nil
	}

	// Provide minimal payload; ai-service may require start/end dates, so include defaults
	start := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	end := time.Now().Format("2006-01-02")
	reqBody, _ := json.Marshal(map[string]interface{}{"forecast_type": forecastType, "start_date": start, "end_date": end, "days": 30})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return errors.New("forecast service returned status " + res.Status)
	}

	return nil
}

func refreshCachedMetrics(ctx context.Context) error {
	return database.DB.WithContext(ctx).Exec("DELETE FROM api_cache").Error
}
