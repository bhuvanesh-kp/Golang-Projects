package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"health_checker/routes"
	"health_checker/worker"
)

func main() {
	fmt.Println("=== Concurrent Website Health Checker ===")
	fmt.Printf("Monitoring %d sites with %d workers\n", len(worker.Sites), worker.WorkerCount)
	fmt.Println("API listening on http://localhost:8080")
	fmt.Println()
	fmt.Println("Endpoints:")
	fmt.Println("  POST http://localhost:8080/api/v1/check       - run health check")
	fmt.Println("  GET  http://localhost:8080/api/v1/history     - browse history")
	fmt.Println("  GET  http://localhost:8080/api/v1/history/:id - single run")
	fmt.Println("  GET  http://localhost:8080/api/v1/sites       - site list")
	fmt.Println("  GET  http://localhost:8080/api/v1/health      - liveness")
	fmt.Println()

	r := gin.Default()
	routes.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
