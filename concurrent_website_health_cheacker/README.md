# Concurrent Website Health Checker

## Overview

A Go REST API that concurrently checks the health of 15 public websites using a **worker pool** of goroutines. Each run measures response times and HTTP status codes, persists results to a local JSON file, and compares them against the previous run to show whether each endpoint got **faster**, **slower**, or remained **stable**. The API is served with [Gin](https://github.com/gin-gonic/gin).

---

## Project Structure

```
concurrent_website_health_cheacker/
├── main.go                   # Gin setup, route registration, server start
├── README.md                 # This file
├── go.mod                    # Module: health_checker
├── go.sum                    # Auto-generated dependency checksums
│
├── models/
│   └── models.go             # All shared struct types (CheckResult, RunResponse, etc.)
│
├── worker/
│   └── pool.go               # Worker pool, HTTP check logic, hardcoded site list
│
├── store/
│   └── store.go              # JSON file persistence, historical comparison logic
│
├── handler/
│   └── handler.go            # Gin handler functions (one per endpoint)
│
├── routes/
│   └── routes.go             # Route registration under /api/v1
│
└── data/
    └── stats.json            # Runtime-generated; stores up to 50 historical runs
```

---

## Architecture

### Worker Pool Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                         POST /api/v1/check                      │
└────────────────────────────┬────────────────────────────────────┘
                             │
                    handler.RunCheck()
                             │
                    worker.Run(Sites)
                             │
          ┌──────────────────▼──────────────────┐
          │          Dispatcher goroutine         │
          │  for url in Sites: jobChan <- url     │
          │  close(jobChan)                       │
          └──────────────────┬──────────────────-┘
                             │ buffered jobChan (cap 20)
          ┌──────────────────▼──────────────────┐
          │         Worker Pool (5 goroutines)    │
          │                                       │
          │  Worker 1 ──► checkSite(url) ──►┐     │
          │  Worker 2 ──► checkSite(url) ──►│     │
          │  Worker 3 ──► checkSite(url) ──►├──► resultChan
          │  Worker 4 ──► checkSite(url) ──►│     │
          │  Worker 5 ──► checkSite(url) ──►┘     │
          │                                       │
          │  sync.WaitGroup tracks completion      │
          └──────────────────┬──────────────────-┘
                             │ buffered resultChan (cap 20)
          ┌──────────────────▼──────────────────┐
          │             Collector                 │
          │  wg.Wait() → close(resultChan)        │
          │  drain into []CheckResult             │
          └──────────────────┬──────────────────-┘
                             │
                    store.Compare(lastRun, current)
                             │
                    store.Save(newRun)
                             │
                    JSON response → caller
```

**Why a buffered channel?** The job channel has capacity 20 (larger than the 15-URL list), so the dispatcher goroutine loads all URLs instantly without blocking. Workers consume at their own pace. The result channel has the same capacity so workers never block writing results — they finish and call `wg.Done()` immediately.

**Why `sync.WaitGroup` over a done channel?** The worker count is fixed and known at start. `WaitGroup` is the idiomatic Go primitive for "wait for N goroutines to finish" — cleaner than a secondary signalling channel.

### Historical Comparison

Each run is saved as a `HistoricalRun` in `data/stats.json`. On every new run:

1. The previous run's results are loaded from disk.
2. Each URL's current response time is compared to its previous response time.
3. A `RunComparison` is generated per URL with:
   - `delta_ms` — absolute change (positive = slower, negative = faster)
   - `delta_percent` — percentage change rounded to 2 decimal places
   - `trend` — one of:
     - `"faster"` — improved by more than 5%
     - `"slower"` — degraded by more than 5%
     - `"stable"` — within ±5% of previous
     - `"new"` — no prior data for this URL

The 5% threshold prevents noise (e.g. ±20ms on a 400ms baseline) from being reported as a meaningful change.

The file stores at most **50 runs** (`MaxStoredRuns`). Older runs are trimmed on write.

---

## API Reference

Base URL: `http://localhost:8080/api/v1`

### POST `/check`
Triggers a full concurrent health check of all 15 sites. Blocks until all workers finish (~3–10s depending on network). Persists the run and returns a comparison against the previous run.

**Response:**
```json
{
  "run_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "run_at": "2026-04-01T10:00:00Z",
  "total_ms": 3241,
  "healthy_sites": 12,
  "unhealthy_sites": 3,
  "results": [
    {
      "url": "https://www.google.com",
      "status_code": 200,
      "response_time_ms": 87,
      "healthy": true,
      "checked_at": "2026-04-01T10:00:00Z"
    },
    {
      "url": "https://httpbin.org/status/500",
      "status_code": 500,
      "response_time_ms": 341,
      "healthy": false,
      "checked_at": "2026-04-01T10:00:01Z"
    },
    {
      "url": "https://this-domain-does-not-exist-xyz.com",
      "status_code": 0,
      "response_time_ms": 10001,
      "healthy": false,
      "error": "dial tcp: no such host",
      "checked_at": "2026-04-01T10:00:01Z"
    }
  ],
  "comparisons": [
    {
      "url": "https://www.google.com",
      "previous_ms": 102,
      "current_ms": 87,
      "delta_ms": -15,
      "delta_percent": -14.71,
      "trend": "faster"
    },
    {
      "url": "https://httpbin.org/delay/3",
      "previous_ms": 3050,
      "current_ms": 3241,
      "delta_ms": 191,
      "delta_percent": 6.26,
      "trend": "slower"
    }
  ]
}
```

---

### GET `/history`
Returns historical runs, newest first. Accepts optional `?limit=N` query param (default 10).

```bash
curl http://localhost:8080/api/v1/history?limit=5
```

---

### GET `/history/:run_id`
Returns a single historical run by its UUID. Returns `404` if not found.

```bash
curl http://localhost:8080/api/v1/history/f47ac10b-58cc-4372-a567-0e02b2c3d479
```

---

### GET `/sites`
Returns the hardcoded list of URLs the checker monitors.

```json
{
  "sites": [
    "https://www.google.com",
    "https://www.github.com",
    "..."
  ],
  "count": 15
}
```

---

### GET `/health`
Simple liveness check for the API itself.

```json
{ "status": "ok" }
```

---

## Getting Started

### Prerequisites
- Go 1.25.5 or later
- Internet access (the checker makes outbound HTTP requests)

### Run

```bash
cd concurrent_website_health_cheacker
go mod tidy
go run .
```

The server starts on `http://localhost:8080`.

### Example Workflow

```bash
# 1. See what sites will be checked
curl http://localhost:8080/api/v1/sites | jq

# 2. Run your first health check
curl -X POST http://localhost:8080/api/v1/check | jq

# 3. Run again — comparisons will now show trends
curl -X POST http://localhost:8080/api/v1/check | jq '.comparisons'

# 4. Browse history
curl http://localhost:8080/api/v1/history | jq

# 5. Fetch a specific run
curl http://localhost:8080/api/v1/history/<run_id> | jq
```

---

## Website List

| URL | Category | Why included |
|-----|----------|--------------|
| `https://www.google.com` | Always-up, fast | Baseline fast responder; failure = network issue |
| `https://www.github.com` | Always-up, fast | Developer-familiar; slightly more latency than Google |
| `https://www.cloudflare.com` | Always-up, fast | CDN-backed; tests a different routing path |
| `https://www.wikipedia.org` | Always-up, medium | High-traffic public site; stable 200 |
| `https://www.reddit.com` | Always-up, variable | Known for TTFB variation; interesting for history |
| `https://httpbin.org/get` | Controlled, fast | Immediate 200 with request echo; reliable fixture |
| `https://httpbin.org/status/200` | Controlled, fast | Explicit status-code endpoint |
| `https://httpbin.org/status/404` | Controlled, unhealthy | Deliberate 404 — verifies `healthy: false` |
| `https://httpbin.org/status/500` | Controlled, unhealthy | Deliberate 500 — tests error classification |
| `https://httpbin.org/delay/1` | Controlled, slow | 1s artificial delay — shows worker pool advantage |
| `https://httpbin.org/delay/3` | Controlled, slow | 3s artificial delay — magnifies time comparison |
| `https://www.amazon.com` | Real-world, variable | CDN response variation; interesting for delta |
| `https://api.github.com` | API endpoint | JSON response; always 200; non-HTML test |
| `https://www.example.com` | Minimal server | IANA reference; tiny HTML; very predictable |
| `https://this-domain-does-not-exist-xyz.com` | DNS fail | Forces error path; confirms `healthy: false` |

**Parallelism benefit:** The two `httpbin.org/delay` endpoints (1s + 3s) overlap with the other 13 checks when run through the 5-worker pool. Total wall time is ~3–4s instead of the ~20s it would take sequentially.

---

## Design Decisions

**Why a worker pool instead of one goroutine per URL?**
A pool with a fixed worker count (5) limits peak concurrency. Unbounded goroutines (one per URL) would work at 15 URLs but would not scale safely to hundreds. The pool pattern is the production-grade approach.

**Why a JSON file instead of SQLite?**
Zero external dependencies beyond `gin` and `uuid`. The file is human-readable, easy to inspect, and sufficient for the volume of data this project produces. It mirrors the simplicity of the sibling projects in this repo.

**Why Gin instead of `net/http`?**
Consistent with `todoGolang` and `contactList` in this repo. Gin reduces boilerplate for route grouping, path parameters, and JSON responses.

**Why a 5% threshold for trend detection?**
Network latency has natural jitter. A ±20ms fluctuation on a 400ms endpoint is ~5% — calling that "slower" would produce noise on every run. The threshold makes `trend` meaningful: it only fires when there is a real change.

---

## Future Improvements

- CLI flags for `WorkerCount`, `HTTPTimeoutSec`, and port
- Slack/webhook notification when a site flips from healthy to unhealthy
- HTML dashboard served by Gin (server-side rendered stats table)
- Configurable URL list via POST body instead of hardcoded slice
- Scheduled auto-checks using a ticker goroutine
- Per-URL average/min/max aggregated across all stored runs
