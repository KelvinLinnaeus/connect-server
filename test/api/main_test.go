package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/connect-univyn/connect_server/internal/api"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// TestServer wraps the API server for testing
type TestServer struct {
	Server *api.Server
	Config util.Config
	TestDB *testhelpers.TestDB
	t      *testing.T
}

// SetupTestServer creates a test server with database
func SetupTestServer(t *testing.T) *TestServer {
	gin.SetMode(gin.ReleaseMode)
	testDB := testhelpers.SetupTestDB(t)

	config := util.Config{
		TokenSymmetricKey:    "12345678901234567890123456789012",
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 24 * time.Hour,
		Environment:          "test",
		RateLimitDefault:     100,
		CORSAllowedOrigins:   "http://localhost:3000,http://localhost:5173", // Test frontend URLs
		LiveEnabled:          true,                                          // Disable live features in tests
	}

	server, err := api.NewServer(config, testDB.Store)
	require.NoError(t, err)

	return &TestServer{
		Server: server,
		Config: config,
		TestDB: testDB,
		t:      t,
	}
}

// Teardown cleans up the test server
func (ts *TestServer) Teardown() {
	testhelpers.CleanupTestData(ts.t, ts.TestDB.DB)
	ts.TestDB.TeardownTestDB()
}

// MakeRequest makes an HTTP request to the test server
func (ts *TestServer) MakeRequest(t *testing.T, method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(data)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	request, err := http.NewRequest(method, url, reqBody)
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	if token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	recorder := httptest.NewRecorder()
	ts.Server.GetRouter().ServeHTTP(recorder, request)

	return recorder
}

// CreateAuthToken creates a JWT token for testing
func (ts *TestServer) CreateAuthToken(t *testing.T, userID uuid.UUID) string {
	tokenMaker, err := auth.NewPasetoMaker(ts.Config.TokenSymmetricKey)
	require.NoError(t, err)

	// Get user to extract username and spaceID
	user, err := ts.TestDB.Store.GetUserByID(context.Background(), userID)
	require.NoError(t, err)

	token, _, err := tokenMaker.CreateToken(
		userID.String(),
		user.Username,
		user.SpaceID.String(),
		ts.Config.AccessTokenDuration,
	)
	require.NoError(t, err)

	return token
}

// ParseSuccessResponse parses a successful API response
func ParseSuccessResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]interface{} {
	var response util.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data to be map[string]interface{}, got %T", response.Data)
	}

	return data
}

// ParseErrorResponse parses an error API response
func ParseErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	errorData, ok := response["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected error to be map[string]interface{}, got %T", response["error"])
	}

	return errorData
}

// CheckResponseCode checks if response code matches expected
func CheckResponseCode(t *testing.T, recorder *httptest.ResponseRecorder, expectedCode int) {
	if recorder.Code != expectedCode {
		t.Logf("Response body: %s", recorder.Body.String())
	}
	require.Equal(t, expectedCode, recorder.Code)
}

// RequireFieldExists checks if a field exists in the data
func RequireFieldExists(t *testing.T, data map[string]interface{}, field string) {
	_, exists := data[field]
	require.True(t, exists, "Field %s should exist", field)
}

// RequireFieldNotExists checks if a field does not exist in the data
func RequireFieldNotExists(t *testing.T, data map[string]interface{}, field string) {
	_, exists := data[field]
	require.False(t, exists, "Field %s should not exist", field)
}
