package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	recorder := ts.MakeRequest(t, http.MethodGet, "/health", nil, "")

	CheckResponseCode(t, recorder, http.StatusOK)

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Equal(t, "ok", response["status"])
	require.Contains(t, response, "db")
}
