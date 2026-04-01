package store

import (
	"encoding/json"
	"errors"
	"math"
	"os"

	"health_checker/models"
)

const (
	StatsFilePath = "data/stats.json"
	MaxStoredRuns = 50
)

// Load reads stats.json and returns its contents.
// If the file does not exist, an empty StatsFile is returned without error.
func Load() (models.StatsFile, error) {
	var sf models.StatsFile

	data, err := os.ReadFile(StatsFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return sf, nil
		}
		return sf, err
	}

	if err := json.Unmarshal(data, &sf); err != nil {
		return sf, err
	}
	return sf, nil
}

// Save writes sf to stats.json, trimming to MaxStoredRuns if needed.
// It creates the data/ directory if it does not exist.
func Save(sf models.StatsFile) error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	if len(sf.Runs) > MaxStoredRuns {
		sf.Runs = sf.Runs[len(sf.Runs)-MaxStoredRuns:]
	}

	data, err := json.MarshalIndent(sf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(StatsFilePath, data, 0644)
}

// Compare produces a RunComparison for each URL in current against previous.
// URLs with no prior data get trend "new".
func Compare(previous, current []models.CheckResult) []models.RunComparison {
	prevMap := make(map[string]int64, len(previous))
	for _, r := range previous {
		prevMap[r.URL] = r.ResponseTimeMs
	}

	comparisons := make([]models.RunComparison, 0, len(current))
	for _, r := range current {
		prevMs, exists := prevMap[r.URL]
		if !exists {
			comparisons = append(comparisons, models.RunComparison{
				URL:       r.URL,
				CurrentMs: r.ResponseTimeMs,
				Trend:     "new",
			})
			continue
		}

		deltaMs := r.ResponseTimeMs - prevMs
		var deltaPercent float64
		if prevMs > 0 {
			deltaPercent = math.Round((float64(deltaMs)/float64(prevMs))*10000) / 100
		}

		trend := "stable"
		switch {
		case deltaPercent < -5:
			trend = "faster"
		case deltaPercent > 5:
			trend = "slower"
		}

		comparisons = append(comparisons, models.RunComparison{
			URL:          r.URL,
			PreviousMs:   prevMs,
			CurrentMs:    r.ResponseTimeMs,
			DeltaMs:      deltaMs,
			DeltaPercent: deltaPercent,
			Trend:        trend,
		})
	}
	return comparisons
}
