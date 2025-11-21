package services

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Agent-Field/agentfield/control-plane/internal/core/domain"
	"github.com/Agent-Field/agentfield/control-plane/internal/core/interfaces"
	"github.com/Agent-Field/agentfield/control-plane/internal/storage"
	"github.com/Agent-Field/agentfield/control-plane/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock AgentClient for testing
type mockAgentClient struct {
	mu                   sync.RWMutex
	statusResponses      map[string]*interfaces.AgentStatusResponse
	statusErrors         map[string]error
	mcpHealthResponses   map[string]*interfaces.MCPHealthResponse
	mcpHealthErrors      map[string]error
	getStatusCallCount   map[string]int
	getMCPHealthCallCount map[string]int
}

func newMockAgentClient() *mockAgentClient {
	return &mockAgentClient{
		statusResponses:      make(map[string]*interfaces.AgentStatusResponse),
		statusErrors:         make(map[string]error),
		mcpHealthResponses:   make(map[string]*interfaces.MCPHealthResponse),
		mcpHealthErrors:      make(map[string]error),
		getStatusCallCount:   make(map[string]int),
		getMCPHealthCallCount: make(map[string]int),
	}
}

func (m *mockAgentClient) GetAgentStatus(ctx context.Context, nodeID string) (*interfaces.AgentStatusResponse, error) {
	m.mu.Lock()
	m.getStatusCallCount[nodeID]++
	m.mu.Unlock()

	m.mu.RLock()
	defer m.mu.RUnlock()

	if err, ok := m.statusErrors[nodeID]; ok {
		return nil, err
	}
	if resp, ok := m.statusResponses[nodeID]; ok {
		return resp, nil
	}
	return nil, errors.New("agent not found")
}

func (m *mockAgentClient) GetMCPHealth(ctx context.Context, nodeID string) (*interfaces.MCPHealthResponse, error) {
	m.mu.Lock()
	m.getMCPHealthCallCount[nodeID]++
	m.mu.Unlock()

	m.mu.RLock()
	defer m.mu.RUnlock()

	if err, ok := m.mcpHealthErrors[nodeID]; ok {
		return nil, err
	}
	if resp, ok := m.mcpHealthResponses[nodeID]; ok {
		return resp, nil
	}
	return nil, errors.New("MCP not available")
}

func (m *mockAgentClient) RestartMCPServer(ctx context.Context, nodeID, alias string) error {
	return nil
}

func (m *mockAgentClient) GetMCPTools(ctx context.Context, nodeID, alias string) (*interfaces.MCPToolsResponse, error) {
	return nil, nil
}

func (m *mockAgentClient) ShutdownAgent(ctx context.Context, nodeID string, graceful bool, timeoutSeconds int) (*interfaces.AgentShutdownResponse, error) {
	return nil, nil
}

func (m *mockAgentClient) setStatusResponse(nodeID string, status string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.statusResponses[nodeID] = &interfaces.AgentStatusResponse{
		Status: status,
	}
}

func (m *mockAgentClient) setStatusError(nodeID string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.statusErrors[nodeID] = err
}

func (m *mockAgentClient) setMCPHealthResponse(nodeID string, response *interfaces.MCPHealthResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mcpHealthResponses[nodeID] = response
}


func (m *mockAgentClient) getStatusCallCountFor(nodeID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getStatusCallCount[nodeID]
}

func (m *mockAgentClient) getMCPHealthCallCountFor(nodeID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getMCPHealthCallCount[nodeID]
}

func setupHealthMonitorTest(t *testing.T) (*HealthMonitor, storage.StorageProvider, *mockAgentClient, *StatusManager, *PresenceManager) {
	t.Helper()

	provider, ctx := setupTestStorage(t)

	// Create status manager
	statusConfig := StatusManagerConfig{
		ReconcileInterval: 30 * time.Second,
	}
	statusManager := NewStatusManager(provider, statusConfig, nil, nil)

	// Create presence manager
	presenceConfig := PresenceManagerConfig{
		HeartbeatTTL:  5 * time.Second,
		SweepInterval: 1 * time.Second,
		HardEvictTTL:  10 * time.Second,
	}
	presenceManager := NewPresenceManager(statusManager, presenceConfig)

	// Create mock agent client
	mockClient := newMockAgentClient()

	// Create health monitor
	config := HealthMonitorConfig{
		CheckInterval: 100 * time.Millisecond,
	}
	hm := NewHealthMonitor(provider, config, nil, mockClient, statusManager, presenceManager)

	t.Cleanup(func() {
		hm.Stop()
		presenceManager.Stop()
		_ = provider.Close(ctx)
	})

	return hm, provider, mockClient, statusManager, presenceManager
}

func TestHealthMonitor_NewHealthMonitor(t *testing.T) {
	provider, ctx := setupTestStorage(t)
	defer provider.Close(ctx)

	statusManager := NewStatusManager(provider, StatusManagerConfig{}, nil, nil)
	presenceManager := NewPresenceManager(statusManager, PresenceManagerConfig{})
	mockClient := newMockAgentClient()

	config := HealthMonitorConfig{
		CheckInterval: 10 * time.Second,
	}

	hm := NewHealthMonitor(provider, config, nil, mockClient, statusManager, presenceManager)

	require.NotNil(t, hm)
	assert.Equal(t, 10*time.Second, hm.config.CheckInterval)
	assert.NotNil(t, hm.activeAgents)
	assert.NotNil(t, hm.mcpHealthCache)
	assert.NotNil(t, hm.stopCh)
}

func TestHealthMonitor_NewHealthMonitor_DefaultConfig(t *testing.T) {
	provider, ctx := setupTestStorage(t)
	defer provider.Close(ctx)

	statusManager := NewStatusManager(provider, StatusManagerConfig{}, nil, nil)
	presenceManager := NewPresenceManager(statusManager, PresenceManagerConfig{})
	mockClient := newMockAgentClient()

	// Pass zero config to test defaults
	hm := NewHealthMonitor(provider, HealthMonitorConfig{}, nil, mockClient, statusManager, presenceManager)

	require.NotNil(t, hm)
	assert.Equal(t, 10*time.Second, hm.config.CheckInterval, "Should use default check interval")
}

func TestHealthMonitor_RegisterAgent(t *testing.T) {
	hm, _, _, _, presenceManager := setupHealthMonitorTest(t)

	nodeID := "test-agent-1"
	baseURL := "http://localhost:8001"

	hm.RegisterAgent(nodeID, baseURL)

	// Verify agent is in active registry
	hm.agentsMutex.RLock()
	agent, exists := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	require.True(t, exists, "Agent should be registered")
	assert.Equal(t, nodeID, agent.NodeID)
	assert.Equal(t, baseURL, agent.BaseURL)
	assert.Equal(t, types.HealthStatusUnknown, agent.LastStatus)

	// Verify presence manager was notified
	assert.True(t, presenceManager.HasLease(nodeID), "Presence manager should track agent")
}

func TestHealthMonitor_RegisterAgent_MultipleAgents(t *testing.T) {
	hm, _, _, _, _ := setupHealthMonitorTest(t)

	agents := map[string]string{
		"agent-1": "http://localhost:8001",
		"agent-2": "http://localhost:8002",
		"agent-3": "http://localhost:8003",
	}

	for nodeID, baseURL := range agents {
		hm.RegisterAgent(nodeID, baseURL)
	}

	hm.agentsMutex.RLock()
	defer hm.agentsMutex.RUnlock()

	assert.Equal(t, 3, len(hm.activeAgents), "Should have 3 registered agents")
	for nodeID := range agents {
		_, exists := hm.activeAgents[nodeID]
		assert.True(t, exists, "Agent %s should be registered", nodeID)
	}
}

func TestHealthMonitor_UnregisterAgent(t *testing.T) {
	hm, provider, _, _, presenceManager := setupHealthMonitorTest(t)
	ctx := context.Background()

	nodeID := "test-agent-unregister"
	baseURL := "http://localhost:8001"

	// First register the agent in storage
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: baseURL,
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	// Register in health monitor
	hm.RegisterAgent(nodeID, baseURL)

	// Verify agent is registered
	require.True(t, presenceManager.HasLease(nodeID))

	// Unregister
	hm.UnregisterAgent(nodeID)

	// Verify agent is removed from active registry
	hm.agentsMutex.RLock()
	_, exists := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	assert.False(t, exists, "Agent should be unregistered")

	// Verify presence manager was notified
	assert.False(t, presenceManager.HasLease(nodeID), "Presence manager should not track agent")
}

func TestHealthMonitor_UnregisterAgent_NonExistent(t *testing.T) {
	hm, _, _, _, _ := setupHealthMonitorTest(t)

	// Unregister non-existent agent should not panic
	hm.UnregisterAgent("non-existent-agent")

	hm.agentsMutex.RLock()
	count := len(hm.activeAgents)
	hm.agentsMutex.RUnlock()

	assert.Equal(t, 0, count, "Should have no registered agents")
}

func TestHealthMonitor_CheckAgentHealth_Healthy(t *testing.T) {
	hm, provider, mockClient, _, presenceManager := setupHealthMonitorTest(t)
	ctx := context.Background()

	nodeID := "test-agent-healthy"
	baseURL := "http://localhost:8001"

	// Register agent in storage
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: baseURL,
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	// Register in health monitor
	hm.RegisterAgent(nodeID, baseURL)

	// Set mock to return healthy status
	mockClient.setStatusResponse(nodeID, "running")

	// Perform health check
	hm.agentsMutex.RLock()
	activeAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	hm.checkAgentHealth(activeAgent)

	// Wait a bit for async updates
	time.Sleep(100 * time.Millisecond)

	// Verify status was updated to active
	hm.agentsMutex.RLock()
	updatedAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	assert.Equal(t, types.HealthStatusActive, updatedAgent.LastStatus, "Agent should be marked as active")
	assert.True(t, presenceManager.HasLease(nodeID), "Presence should be updated")
}

func TestHealthMonitor_CheckAgentHealth_Inactive(t *testing.T) {
	hm, provider, mockClient, _, _ := setupHealthMonitorTest(t)
	ctx := context.Background()

	nodeID := "test-agent-inactive"
	baseURL := "http://localhost:8001"

	// Register agent in storage
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: baseURL,
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	// Register in health monitor
	hm.RegisterAgent(nodeID, baseURL)

	// Set mock to return error (simulating agent offline)
	mockClient.setStatusError(nodeID, errors.New("connection refused"))

	// Perform health check
	hm.agentsMutex.RLock()
	activeAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	hm.checkAgentHealth(activeAgent)

	// Wait a bit for async updates
	time.Sleep(100 * time.Millisecond)

	// Verify status was updated to inactive
	hm.agentsMutex.RLock()
	updatedAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	assert.Equal(t, types.HealthStatusInactive, updatedAgent.LastStatus, "Agent should be marked as inactive")
}

func TestHealthMonitor_CheckAgentHealth_NotRunning(t *testing.T) {
	hm, provider, mockClient, _, _ := setupHealthMonitorTest(t)
	ctx := context.Background()

	nodeID := "test-agent-not-running"
	baseURL := "http://localhost:8001"

	// Register agent in storage
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: baseURL,
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	// Register in health monitor
	hm.RegisterAgent(nodeID, baseURL)

	// Set mock to return non-running status
	mockClient.setStatusResponse(nodeID, "stopped")

	// Perform health check
	hm.agentsMutex.RLock()
	activeAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	hm.checkAgentHealth(activeAgent)

	// Wait a bit for async updates
	time.Sleep(100 * time.Millisecond)

	// Verify status was updated to inactive
	hm.agentsMutex.RLock()
	updatedAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	assert.Equal(t, types.HealthStatusInactive, updatedAgent.LastStatus, "Agent should be marked as inactive when not running")
}

func TestHealthMonitor_CheckAgentHealth_StatusTransitions(t *testing.T) {
	hm, provider, mockClient, _, _ := setupHealthMonitorTest(t)
	ctx := context.Background()

	nodeID := "test-agent-transitions"
	baseURL := "http://localhost:8001"

	// Register agent in storage
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: baseURL,
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	// Register in health monitor
	hm.RegisterAgent(nodeID, baseURL)

	// Get active agent
	hm.agentsMutex.RLock()
	activeAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	// Test transition: Unknown -> Active
	mockClient.setStatusResponse(nodeID, "running")
	hm.checkAgentHealth(activeAgent)
	time.Sleep(100 * time.Millisecond)

	hm.agentsMutex.RLock()
	assert.Equal(t, types.HealthStatusActive, hm.activeAgents[nodeID].LastStatus)
	hm.agentsMutex.RUnlock()

	// Test transition: Active -> Inactive
	mockClient.setStatusError(nodeID, errors.New("connection refused"))
	hm.checkAgentHealth(activeAgent)
	time.Sleep(100 * time.Millisecond)

	hm.agentsMutex.RLock()
	assert.Equal(t, types.HealthStatusInactive, hm.activeAgents[nodeID].LastStatus)
	hm.agentsMutex.RUnlock()

	// Test transition: Inactive -> Active (should work after debouncing period)
	// Wait for debounce period
	time.Sleep(31 * time.Second)

	mockClient.setStatusResponse(nodeID, "running")
	mockClient.statusErrors = make(map[string]error) // Clear errors
	hm.checkAgentHealth(activeAgent)
	time.Sleep(100 * time.Millisecond)

	hm.agentsMutex.RLock()
	assert.Equal(t, types.HealthStatusActive, hm.activeAgents[nodeID].LastStatus)
	hm.agentsMutex.RUnlock()
}

func TestHealthMonitor_CheckAgentHealth_UnregisteredAgent(t *testing.T) {
	hm, _, mockClient, _, _ := setupHealthMonitorTest(t)

	nodeID := "test-agent-unregistered"
	baseURL := "http://localhost:8001"

	// Create agent but don't register it
	activeAgent := &ActiveAgent{
		NodeID:  nodeID,
		BaseURL: baseURL,
	}

	mockClient.setStatusResponse(nodeID, "running")

	// Should skip check for unregistered agent
	hm.checkAgentHealth(activeAgent)

	// Verify no status calls were made (agent was not in registry)
	assert.Equal(t, 0, mockClient.getStatusCallCountFor(nodeID))
}

func TestHealthMonitor_MCP_CheckMCPHealth(t *testing.T) {
	hm, provider, mockClient, _, _ := setupHealthMonitorTest(t)
	ctx := context.Background()

	nodeID := "test-agent-mcp"
	baseURL := "http://localhost:8001"

	// Register agent in storage
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: baseURL,
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	// Register in health monitor
	hm.RegisterAgent(nodeID, baseURL)

	// Set mock MCP health response
	mcpResponse := &interfaces.MCPHealthResponse{
		Summary: interfaces.MCPSummary{
			TotalServers:   3,
			RunningServers: 3,
			TotalTools:     15,
			OverallHealth:  0.95,
		},
	}
	mockClient.setMCPHealthResponse(nodeID, mcpResponse)

	// Set agent as healthy first
	mockClient.setStatusResponse(nodeID, "running")

	// Perform health check (should trigger MCP check for active agents)
	hm.agentsMutex.RLock()
	activeAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	hm.checkAgentHealth(activeAgent)
	time.Sleep(200 * time.Millisecond)

	// Verify MCP health was checked
	assert.Greater(t, mockClient.getMCPHealthCallCountFor(nodeID), 0, "MCP health should be checked for active agent")

	// Verify MCP health is cached
	cache := hm.GetMCPHealthCache()
	mcpData, exists := cache[nodeID]
	require.True(t, exists, "MCP health should be cached")
	assert.Equal(t, 3, mcpData.TotalServers)
	assert.Equal(t, 3, mcpData.RunningServers)
	assert.Equal(t, 15, mcpData.TotalTools)
	assert.Equal(t, 0.95, mcpData.OverallHealth)
}

func TestHealthMonitor_MCP_HealthChange(t *testing.T) {
	hm, provider, mockClient, _, _ := setupHealthMonitorTest(t)
	ctx := context.Background()

	nodeID := "test-agent-mcp-change"
	baseURL := "http://localhost:8001"

	// Register agent in storage
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: baseURL,
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	// Register in health monitor
	hm.RegisterAgent(nodeID, baseURL)

	// Set agent as healthy
	mockClient.setStatusResponse(nodeID, "running")

	// First MCP health check
	mcpResponse1 := &interfaces.MCPHealthResponse{
		Summary: interfaces.MCPSummary{
			TotalServers:   3,
			RunningServers: 3,
			TotalTools:     15,
			OverallHealth:  0.95,
		},
	}
	mockClient.setMCPHealthResponse(nodeID, mcpResponse1)

	hm.agentsMutex.RLock()
	activeAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	hm.checkAgentHealth(activeAgent)
	time.Sleep(200 * time.Millisecond)

	// Change MCP health
	mcpResponse2 := &interfaces.MCPHealthResponse{
		Summary: interfaces.MCPSummary{
			TotalServers:   3,
			RunningServers: 2, // One server failed
			TotalTools:     10,
			OverallHealth:  0.67,
		},
	}
	mockClient.setMCPHealthResponse(nodeID, mcpResponse2)

	// Second health check
	hm.checkAgentHealth(activeAgent)
	time.Sleep(200 * time.Millisecond)

	// Verify MCP health was updated
	cache := hm.GetMCPHealthCache()
	mcpData, exists := cache[nodeID]
	require.True(t, exists)
	assert.Equal(t, 2, mcpData.RunningServers, "MCP health should be updated")
	assert.Equal(t, 0.67, mcpData.OverallHealth)
}

func TestHealthMonitor_MCP_NoChange(t *testing.T) {
	hm, provider, mockClient, _, _ := setupHealthMonitorTest(t)
	ctx := context.Background()

	nodeID := "test-agent-mcp-no-change"
	baseURL := "http://localhost:8001"

	// Register agent in storage
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: baseURL,
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	// Register in health monitor
	hm.RegisterAgent(nodeID, baseURL)

	// Set agent as healthy
	mockClient.setStatusResponse(nodeID, "running")

	// Set MCP health response
	mcpResponse := &interfaces.MCPHealthResponse{
		Summary: interfaces.MCPSummary{
			TotalServers:   3,
			RunningServers: 3,
			TotalTools:     15,
			OverallHealth:  0.95,
		},
	}
	mockClient.setMCPHealthResponse(nodeID, mcpResponse)

	hm.agentsMutex.RLock()
	activeAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	// First check
	hm.checkAgentHealth(activeAgent)
	time.Sleep(200 * time.Millisecond)

	// Verify hasMCPHealthChanged returns false for same data
	newSummary := &domain.MCPSummaryData{
		TotalServers:   3,
		RunningServers: 3,
		TotalTools:     15,
		OverallHealth:  0.95,
	}
	hasChanged := hm.hasMCPHealthChanged(nodeID, newSummary)
	assert.False(t, hasChanged, "Should detect no change in MCP health")
}

func TestHealthMonitor_MCP_InactiveAgent(t *testing.T) {
	hm, provider, mockClient, _, _ := setupHealthMonitorTest(t)
	ctx := context.Background()

	nodeID := "test-agent-mcp-inactive"
	baseURL := "http://localhost:8001"

	// Register agent in storage
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: baseURL,
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	// Register in health monitor
	hm.RegisterAgent(nodeID, baseURL)

	// Set agent as inactive
	mockClient.setStatusError(nodeID, errors.New("connection refused"))

	hm.agentsMutex.RLock()
	activeAgent := hm.activeAgents[nodeID]
	hm.agentsMutex.RUnlock()

	hm.checkAgentHealth(activeAgent)
	time.Sleep(200 * time.Millisecond)

	// MCP health should NOT be checked for inactive agents
	assert.Equal(t, 0, mockClient.getMCPHealthCallCountFor(nodeID), "MCP health should not be checked for inactive agent")
}

func TestHealthMonitor_GetMCPHealthCache(t *testing.T) {
	hm, _, _, _, _ := setupHealthMonitorTest(t)

	// Add some test data to cache
	hm.mcpCacheMutex.Lock()
	hm.mcpHealthCache["agent-1"] = &domain.MCPSummaryData{
		TotalServers:   3,
		RunningServers: 3,
		TotalTools:     15,
		OverallHealth:  0.95,
	}
	hm.mcpHealthCache["agent-2"] = &domain.MCPSummaryData{
		TotalServers:   2,
		RunningServers: 1,
		TotalTools:     8,
		OverallHealth:  0.50,
	}
	hm.mcpCacheMutex.Unlock()

	// Get cache
	cache := hm.GetMCPHealthCache()

	// Verify cache contents
	assert.Equal(t, 2, len(cache))
	assert.Contains(t, cache, "agent-1")
	assert.Contains(t, cache, "agent-2")
	assert.Equal(t, 3, cache["agent-1"].TotalServers)
	assert.Equal(t, 1, cache["agent-2"].RunningServers)
}

func TestHealthMonitor_ConcurrentAccess(t *testing.T) {
	hm, provider, mockClient, _, _ := setupHealthMonitorTest(t)
	ctx := context.Background()

	// Register multiple agents
	agents := []string{"agent-1", "agent-2", "agent-3", "agent-4", "agent-5"}
	for _, nodeID := range agents {
		agent := &types.AgentNode{
			ID:      nodeID,
			BaseURL: "http://localhost:800" + nodeID[len(nodeID)-1:],
		}
		err := provider.RegisterAgent(ctx, agent)
		require.NoError(t, err)

		hm.RegisterAgent(nodeID, agent.BaseURL)
		mockClient.setStatusResponse(nodeID, "running")
	}

	var wg sync.WaitGroup

	// Concurrent health checks
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hm.checkActiveAgents()
		}()
	}

	// Concurrent register/unregister
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			nodeID := "temp-agent-" + string(rune('0'+idx))
			hm.RegisterAgent(nodeID, "http://localhost:9000")
			time.Sleep(10 * time.Millisecond)
			hm.UnregisterAgent(nodeID)
		}(i)
	}

	// Concurrent MCP cache access
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = hm.GetMCPHealthCache()
		}()
	}

	wg.Wait()

	// Verify no race conditions
	hm.agentsMutex.RLock()
	activeCount := len(hm.activeAgents)
	hm.agentsMutex.RUnlock()

	assert.Equal(t, 5, activeCount, "Should have 5 active agents after concurrent operations")
}

func TestHealthMonitor_StartStop(t *testing.T) {
	hm, _, _, _, _ := setupHealthMonitorTest(t)

	// Start in goroutine
	go hm.Start()

	// Let it run for a bit
	time.Sleep(300 * time.Millisecond)

	// Stop should not block
	hm.Stop()

	// Verify stop worked (stop channel closed)
	select {
	case <-hm.stopCh:
		// Expected: channel is closed
	default:
		t.Fatal("Stop channel should be closed")
	}
}

func TestHealthMonitor_PeriodicChecks(t *testing.T) {
	provider, ctx := setupTestStorage(t)
	defer provider.Close(ctx)

	mockClient := newMockAgentClient()
	statusManager := NewStatusManager(provider, StatusManagerConfig{}, nil, nil)
	presenceManager := NewPresenceManager(statusManager, PresenceManagerConfig{})

	// Use very short interval for testing
	config := HealthMonitorConfig{
		CheckInterval: 50 * time.Millisecond,
	}
	hm := NewHealthMonitor(provider, config, nil, mockClient, statusManager, presenceManager)

	// Register agent
	nodeID := "test-periodic"
	agent := &types.AgentNode{
		ID:      nodeID,
		BaseURL: "http://localhost:8001",
	}
	err := provider.RegisterAgent(ctx, agent)
	require.NoError(t, err)

	hm.RegisterAgent(nodeID, agent.BaseURL)
	mockClient.setStatusResponse(nodeID, "running")

	// Start monitoring
	go hm.Start()

	// Let it run for multiple check intervals
	time.Sleep(250 * time.Millisecond)

	// Stop
	hm.Stop()

	// Verify multiple checks occurred
	callCount := mockClient.getStatusCallCountFor(nodeID)
	assert.GreaterOrEqual(t, callCount, 3, "Should have performed multiple periodic checks")
}

func TestHealthMonitor_CheckActiveAgents_NoAgents(t *testing.T) {
	hm, _, _, _, _ := setupHealthMonitorTest(t)

	// Should not panic with no agents
	hm.checkActiveAgents()

	hm.agentsMutex.RLock()
	count := len(hm.activeAgents)
	hm.agentsMutex.RUnlock()

	assert.Equal(t, 0, count)
}
