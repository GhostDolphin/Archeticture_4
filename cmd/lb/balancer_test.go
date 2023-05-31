package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalancer(t *testing.T) {
	isHealthy:= &IsHealthy{}
	isHealthy.status = map[string]bool{
		"server1:8080": true,
		"server2:8080": true,
		"server3:8080": true,
	}

	balance := &Balancer{}
	balance.isHealthy = isHealthy

	server1 := balance.doBalancer("/check")
	server1again := balance.doBalancer("/check")
	server2 := balance.doBalancer("/check2")
	server3 := balance.doBalancer("/check5")

	assert.Equal(t, "server1:8080", server1)
	assert.Equal(t, server1, server1again)
	assert.Equal(t, "server2:8080", server2)
	assert.Equal(t, "server3:8080", server3)
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

func TestHealthChecker(t *testing.T) {
	isHealthy := &IsHealthy{}
	isHealthy.status = map[string]bool{}
	isHealthy.health = mockHealth

	isHealthy.CheckAll()
	assert.Equal(t, map[string]bool{"server1:8080": true, "server2:8080": false, "server3:8080": false}, isHealthy.status)

	allHealthy := isHealthy.AllHealthy()
	assert.Equal(t, []string{"server1:8080"}, allHealthy)

	isHealthy.health = mockHealthAllTrue
	isHealthy.CheckAll()
	allHealthy = isHealthy.AllHealthy()
	assert.Equal(t, []string{"server1:8080", "server2:8080", "server3:8080"}, allHealthy)
}

func TestNoAvailableServers(t *testing.T) {
	healthchecker := &IsHealthy{}
	healthchecker.status = map[string]bool{
		"server1:8080": false,
		"server2:8080": false,
		"server3:8080": false,
	}

	balance := &Balancer{}
	balance.isHealthy = healthchecker

	server := balance.doBalancer("/check")

	assert.Equal(t, "", server)
}
