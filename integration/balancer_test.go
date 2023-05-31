package integration

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

func TestBalancer(t *testing.T) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		t.Skip("Integration test is not enabled")
	}

	resp1, err := client.Get(fmt.Sprintf("%s/api/v1/some-data2", baseAddress))
	require.NoError(t, err)
	assert.Equal(t, "server1:8080", resp1.Header.Get("lb-from"))

	resp2, err := client.Get(fmt.Sprintf("%s/api/v1/some-data5", baseAddress))
	require.NoError(t, err)
	assert.Equal(t, "server2:8080", resp2.Header.Get("lb-from"))

	resp3, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	require.NoError(t, err)
	assert.Equal(t, "server3:8080", resp3.Header.Get("lb-from"))

	respr, err := client.Get(fmt.Sprintf("%s/api/v1/some-data2", baseAddress))
	require.NoError(t, err)
	assert.Equal(t, "server1:8080", respr.Header.Get("lb-from"))
}

func BenchmarkBalancer(b *testing.B) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		b.Skip("Integration test is not enabled")
	}

	for i := 0; i < b.N; i++ {
		_, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		require.NoError(b, err)
	}
}
