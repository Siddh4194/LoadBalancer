package lib

import (
	roundrobin "LoadBalancer/Algorithms/RoundRobin"
	"LoadBalancer/ai"
	"LoadBalancer/logger"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func HealthCheck(node *roundrobin.Node) {

	start := time.Now()

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	logger.Log.Info("health_check_started",
		zap.String("server", node.Server),
	)

	resp, err := client.Get(node.Url + "/health")

	duration := time.Since(start)

	if err != nil {

		node.FailCount++

		logger.Log.Error("health_check_failed",
			zap.String("server", node.Server),
			zap.String("url", node.Url),
			zap.Int("fail_count", node.FailCount),
			zap.Duration("latency", duration),
			zap.Error(err),
		)

	} else {

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {

			node.FailCount++
			logger.Log.Warn("health_check_unhealthy_status",
				zap.String("server", node.Server),
				zap.Int("status_code", resp.StatusCode),
				zap.Int("fail_count", node.FailCount),
				zap.Duration("latency", duration),
			)

		} else {

			node.FailCount = 0
			node.Healthy = true

			logger.Log.Info("health_check_success",
				zap.String("server", node.Server),
				zap.Int("status_code", resp.StatusCode),
				zap.Duration("latency", duration),
			)
		}
	}

	if node.FailCount >= 3 {

		node.Healthy = false

		fmt.Printf("inside the health checks")
		ai.AnalyzeInfrastructure("Health check returned unhealthy status")
		logger.Log.Error("server_marked_unhealthy",
			zap.String("server", node.Server),
			zap.Int("fail_count", node.FailCount),
		)
	}
}
