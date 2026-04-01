package models

import "time"

// CheckResult holds the outcome of a single HTTP health probe.
type CheckResult struct {
	URL            string    `json:"url"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int64     `json:"response_time_ms"`
	Healthy        bool      `json:"healthy"`
	Error          string    `json:"error,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}

// RunComparison shows how a single URL's response time changed between runs.
// Trend values: "faster" (<-5%), "slower" (>+5%), "stable" (within ±5%), "new" (no prior data).
type RunComparison struct {
	URL          string  `json:"url"`
	PreviousMs   int64   `json:"previous_ms"`
	CurrentMs    int64   `json:"current_ms"`
	DeltaMs      int64   `json:"delta_ms"`      // positive = slower, negative = faster
	DeltaPercent float64 `json:"delta_percent"` // rounded to 2 decimal places
	Trend        string  `json:"trend"`
}

// HistoricalRun is one complete check run persisted to stats.json.
type HistoricalRun struct {
	RunID     string        `json:"run_id"`
	RunAt     time.Time     `json:"run_at"`
	Results   []CheckResult `json:"results"`
	TotalMs   int64         `json:"total_ms"`
	SiteCount int           `json:"site_count"`
}

// RunResponse is returned by POST /api/v1/check.
type RunResponse struct {
	RunID          string          `json:"run_id"`
	RunAt          time.Time       `json:"run_at"`
	Results        []CheckResult   `json:"results"`
	Comparisons    []RunComparison `json:"comparisons"`
	TotalMs        int64           `json:"total_ms"`
	HealthySites   int             `json:"healthy_sites"`
	UnhealthySites int             `json:"unhealthy_sites"`
}

// StatsFile is the root object in data/stats.json.
type StatsFile struct {
	LastUpdated time.Time       `json:"last_updated"`
	Runs        []HistoricalRun `json:"runs"`
}
