package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nikitakolesnik/pet-proj/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ProducerURL = "http://localhost:8080"
	ConsumerURL = "http://localhost:8081"
	MonitorURL  = "http://localhost:8082"
)

func TestEndToEndFlow(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Wait for services to be ready
	waitForService(t, ProducerURL+"/health", "producer")
	waitForService(t, ConsumerURL+"/health", "consumer")
	waitForService(t, MonitorURL+"/health", "monitor")

	t.Run("SendEvent", func(t *testing.T) {
		// Create test event
		event := models.Event{
			Type:   models.EventTypeUserAction,
			UserID: "test_user_123",
			Data: map[string]interface{}{
				"action": "test_action",
				"page":   "/test-page",
			},
			Source: models.SourceProducer,
		}

		// Send event to producer
		resp, err := sendEvent(t, event)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		eventID, ok := response["event_id"].(string)
		require.True(t, ok)
		assert.NotEmpty(t, eventID)

		// Wait for event to be processed
		time.Sleep(2 * time.Second)

		// Check if event was processed by consumer
		processedEvent, err := getProcessedEvent(t, eventID)
		require.NoError(t, err)
		assert.NotNil(t, processedEvent)
		assert.Equal(t, "processed", processedEvent["status"])
	})

	t.Run("GenerateRandomEvents", func(t *testing.T) {
		// Generate multiple random events
		resp, err := http.Post(ProducerURL+"/api/v1/generate?count=5", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		count, ok := response["count"].(float64)
		require.True(t, ok)
		assert.Equal(t, float64(5), count)
	})

	t.Run("GetTransactionStats", func(t *testing.T) {
		// Wait for some transactions to be recorded
		time.Sleep(5 * time.Second)

		// Get stats from monitor
		resp, err := http.Get(MonitorURL + "/api/v1/stats")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var stats map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&stats)
		require.NoError(t, err)
		assert.NotNil(t, stats)
	})

	t.Run("GetSystemHealth", func(t *testing.T) {
		// Get system health
		resp, err := http.Get(MonitorURL + "/api/v1/health")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var health map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&health)
		require.NoError(t, err)

		overallStatus, ok := health["overall_status"].(string)
		require.True(t, ok)
		assert.Equal(t, "healthy", overallStatus)

		services, ok := health["services"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "ok", services["kafka"])
		assert.Equal(t, "ok", services["redis"])
		assert.Equal(t, "ok", services["postgres"])
	})
}

func TestProducerService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	waitForService(t, ProducerURL+"/health", "producer")

	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get(ProducerURL + "/health")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var health map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&health)
		require.NoError(t, err)
		assert.Equal(t, "healthy", health["status"])
		assert.Equal(t, "producer", health["service"])
	})

	t.Run("GetStats", func(t *testing.T) {
		resp, err := http.Get(ProducerURL + "/api/v1/stats")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var stats map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&stats)
		require.NoError(t, err)
		assert.NotNil(t, stats)
	})

	t.Run("Metrics", func(t *testing.T) {
		resp, err := http.Get(ProducerURL + "/metrics")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "text/plain")
	})
}

func TestConsumerService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	waitForService(t, ConsumerURL+"/health", "consumer")

	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get(ConsumerURL + "/health")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var health map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&health)
		require.NoError(t, err)
		assert.Equal(t, "healthy", health["status"])
		assert.Equal(t, "consumer", health["service"])
	})

	t.Run("GetStats", func(t *testing.T) {
		resp, err := http.Get(ConsumerURL + "/api/v1/stats")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var stats map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&stats)
		require.NoError(t, err)
		assert.NotNil(t, stats)
	})
}

func TestMonitorService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	waitForService(t, MonitorURL+"/health", "monitor")

	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get(MonitorURL + "/health")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var health map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&health)
		require.NoError(t, err)
		assert.Equal(t, "healthy", health["status"])
		assert.Equal(t, "monitor", health["service"])
	})

	t.Run("GetTransactions", func(t *testing.T) {
		resp, err := http.Get(MonitorURL + "/api/v1/transactions?limit=10")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Contains(t, response, "transactions")
		assert.Contains(t, response, "count")
	})

	t.Run("GetDashboard", func(t *testing.T) {
		resp, err := http.Get(MonitorURL + "/api/v1/dashboard")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var dashboard map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&dashboard)
		require.NoError(t, err)
		assert.Contains(t, dashboard, "health")
		assert.Contains(t, dashboard, "recent_transactions")
		assert.Contains(t, dashboard, "stats")
	})
}

// Helper functions

func waitForService(t *testing.T, url, serviceName string) {
	t.Helper()

	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}

	t.Fatalf("Service %s at %s is not ready after %d seconds", serviceName, url, maxRetries)
}

func sendEvent(t *testing.T, event models.Event) (*http.Response, error) {
	t.Helper()

	jsonData, err := json.Marshal(event)
	require.NoError(t, err)

	return http.Post(ProducerURL+"/api/v1/events", "application/json",
		strings.NewReader(string(jsonData)))
}

func getProcessedEvent(t *testing.T, eventID string) (map[string]interface{}, error) {
	t.Helper()

	resp, err := http.Get(ConsumerURL + "/api/v1/events/" + eventID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get processed event: status %d", resp.StatusCode)
	}

	var processedEvent map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&processedEvent)
	return processedEvent, err
}
