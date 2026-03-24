package lib

import (
	roundrobin "LoadBalancer/Algorithms/RoundRobin"
	"fmt"
	"net/http"
	"time"
)			

func HealthCheck(node *roundrobin.Node) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(node.Url + "/health")
	fmt.Printf("Health check for server %s: %d", node.Server,resp.StatusCode)
	if err != nil || resp.StatusCode != http.StatusOK {
		node.FailCount++
	} else {
		node.FailCount = 0
		node.Healthy = true
	}

	if node.FailCount >= 3 {
		node.Healthy = false
	}

}