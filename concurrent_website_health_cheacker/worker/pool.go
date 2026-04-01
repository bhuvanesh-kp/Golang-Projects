package worker

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"health_checker/models"
)

const (
	WorkerCount    = 5
	JobBufferSize  = 20
	ResultBuffer   = 20
	HTTPTimeoutSec = 10
)

// Sites is the hardcoded list of public URLs to health-check.
var Sites = []string{
	"https://www.google.com",
	"https://www.github.com",
	"https://www.cloudflare.com",
	"https://www.wikipedia.org",
	"https://www.reddit.com",
	"https://httpbin.org/get",
	"https://httpbin.org/status/200",
	"https://httpbin.org/status/404",
	"https://httpbin.org/status/500",
	"https://httpbin.org/delay/1",
	"https://httpbin.org/delay/3",
	"https://www.amazon.com",
	"https://api.github.com",
	"https://www.example.com",
	"https://this-domain-does-not-exist-xyz.com",
}

var httpClient = &http.Client{
	Timeout: HTTPTimeoutSec * time.Second,
	// Do not follow redirects so we capture the real status code.
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// Run dispatches all URLs through the worker pool and returns aggregated results.
// It blocks until every URL has been checked.
func Run(urls []string) []models.CheckResult {
	jobChan := make(chan string, JobBufferSize)
	resultChan := make(chan models.CheckResult, ResultBuffer)

	var wg sync.WaitGroup

	// Start workers.
	for i := 0; i < WorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobChan {
				resultChan <- checkSite(url)
			}
		}()
	}

	// Dispatcher: load all URLs into the job channel, then close it.
	go func() {
		for _, url := range urls {
			jobChan <- url
		}
		close(jobChan)
	}()

	// Close the result channel once all workers are done.
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results.
	var results []models.CheckResult
	for r := range resultChan {
		results = append(results, r)
	}
	return results
}

// checkSite performs a single HTTP GET against url and returns a CheckResult.
func checkSite(url string) models.CheckResult {
	start := time.Now()
	result := models.CheckResult{
		URL:       url,
		CheckedAt: start,
	}

	resp, err := httpClient.Get(url)
	elapsed := time.Since(start).Milliseconds()
	result.ResponseTimeMs = elapsed

	if err != nil {
		result.Healthy = false
		result.Error = fmt.Sprintf("%v", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Healthy = resp.StatusCode >= 200 && resp.StatusCode < 400
	return result
}
