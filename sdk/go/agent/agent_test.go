package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleReasonerAsyncPostsStatus(t *testing.T) {
	callbackCh := make(chan map[string]any, 1)
	callbackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		var payload map[string]any
		if err := dec.Decode(&payload); err == nil {
			callbackCh <- payload
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer callbackServer.Close()

	cfg := Config{
		NodeID:        "node-1",
		Version:       "1.0.0",
		TeamID:        "team",
		AgentFieldURL: callbackServer.URL,
		ListenAddress: ":0",
		PublicURL:     "http://localhost:0",
		Logger:        log.New(io.Discard, "[test] ", 0),
	}

	agent, err := New(cfg)
	require.NoError(t, err)

	agent.RegisterReasoner("demo", func(ctx context.Context, input map[string]any) (any, error) {
		time.Sleep(10 * time.Millisecond)
		return map[string]any{"ok": true}, nil
	})

	server := httptest.NewServer(agent.handler())
	defer server.Close()

	reqBody := []byte(`{"value":42}`)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/reasoners/demo", bytes.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Execution-ID", "exec-test")
	req.Header.Set("X-Run-ID", "run-1")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)
	resp.Body.Close()

	select {
	case payload := <-callbackCh:
		assert.Equal(t, "exec-test", payload["execution_id"])
		assert.Equal(t, "succeeded", payload["status"])
		result, ok := payload["result"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, true, result["ok"])
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for callback payload")
	}
}
