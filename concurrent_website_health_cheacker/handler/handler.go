package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"health_checker/models"
	"health_checker/store"
	"health_checker/worker"
)

// RunCheck triggers a full concurrent health check, persists the run, and returns
// results alongside a comparison against the previous run.
func RunCheck(ctx *gin.Context) {
	sf, err := store.Load()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Capture the last run's results for comparison before the new run.
	var previousResults []models.CheckResult
	if len(sf.Runs) > 0 {
		previousResults = sf.Runs[len(sf.Runs)-1].Results
	}

	start := time.Now()
	results := worker.Run(worker.Sites)
	totalMs := time.Since(start).Milliseconds()

	healthy, unhealthy := 0, 0
	for _, r := range results {
		if r.Healthy {
			healthy++
		} else {
			unhealthy++
		}
	}

	runID := uuid.New().String()
	runAt := time.Now()

	newRun := models.HistoricalRun{
		RunID:     runID,
		RunAt:     runAt,
		Results:   results,
		TotalMs:   totalMs,
		SiteCount: len(results),
	}

	sf.Runs = append(sf.Runs, newRun)
	sf.LastUpdated = runAt

	if err := store.Save(sf); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	comparisons := store.Compare(previousResults, results)

	ctx.JSON(http.StatusOK, models.RunResponse{
		RunID:          runID,
		RunAt:          runAt,
		Results:        results,
		Comparisons:    comparisons,
		TotalMs:        totalMs,
		HealthySites:   healthy,
		UnhealthySites: unhealthy,
	})
}

// GetHistory returns historical runs, newest first. Accepts ?limit=N (default 10).
func GetHistory(ctx *gin.Context) {
	sf, err := store.Load()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	limit := 10
	if q := ctx.Query("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			limit = n
		}
	}

	runs := sf.Runs
	// Reverse so newest is first.
	for i, j := 0, len(runs)-1; i < j; i, j = i+1, j-1 {
		runs[i], runs[j] = runs[j], runs[i]
	}

	if limit < len(runs) {
		runs = runs[:limit]
	}

	ctx.JSON(http.StatusOK, gin.H{
		"runs":  runs,
		"count": len(runs),
	})
}

// GetRunByID returns a single historical run by its UUID.
func GetRunByID(ctx *gin.Context) {
	runID := ctx.Param("run_id")

	sf, err := store.Load()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for _, run := range sf.Runs {
		if run.RunID == runID {
			ctx.JSON(http.StatusOK, run)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"message": "run not found"})
}

// GetSites returns the hardcoded list of monitored URLs.
func GetSites(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"sites": worker.Sites,
		"count": len(worker.Sites),
	})
}

// HealthCheck is a simple liveness endpoint for the API itself.
func HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
