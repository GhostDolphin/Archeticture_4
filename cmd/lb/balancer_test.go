package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadBalancer(t *testing.T) {
	healthChecker := &LoadBalancerHealthChecker{
		serverHealthStatus: map[string]bool{
			"server1:8080": true,
			"server2:8080": true,
			"server3:8080": true,
		},
	}

	balancer := &LoadBalancer{
		healthChecker: healthChecker,
	}

	server1 := balancer.balance("/check")
	server1SecondTime := balancer.balance("/check")
	server2 := balancer.balance("/check2")
	server3 := balancer.balance("/check5")

	assert.Equal(t, "server1:8080", server1)
	assert.Equal(t, server1, server1SecondTime)
	assert.Equal(t, "server2:8080", server2)
	assert.Equal(t, "server3:8080", server3)
}

func TestLoadBalancerHealthChecker(t *testing.T) {
	healthChecker := &LoadBalancerHealthChecker{
		serverHealthStatus: map[string]bool{},
		health:             mockHealth,
	}

	healthChecker.CheckAllServers()
	assert.Equal(t, map[string]bool{"server1:8080": true, "server2:8080": false, "server3:8080": false}, healthChecker.serverHealthStatus)

	healthyServers := healthChecker.GetHealthyServers()
	assert.Equal(t, []string{"server1:8080"}, healthyServers)

	healthChecker.health = mockHealthAllTrue
	healthChecker.CheckAllServers()
	healthyServers = healthChecker.GetHealthyServers()
	assert.Equal(t, []string{"server1:8080", "server2:8080", "server3:8080"}, healthyServers)
}

func TestLoadBalancerNoAvailableServers(t *testing.T) {
	healthChecker := &LoadBalancerHealthChecker{
		serverHealthStatus: map[string]bool{
			"server1:8080": false,
			"server2:8080": false,
			"server3:8080": false,
		},
	}

	balancer := &LoadBalancer{
		healthChecker: healthChecker,
	}

	server := balancer.balance("/check")

	assert.Equal(t, "", server)
}

func mockHealth(dst string) bool {
	if dst == "server1:8080" {
		return true
	} else if dst == "server2:8080" {
		return false
	}
	return false
}

func mockHealthAllTrue(dst string) bool {
	return true
}