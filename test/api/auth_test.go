package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/connect-univyn/connect_server/internal/util"
	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	// Create test space and user
	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	testCases := []struct {
		name          string
		body          map[string]interface{}
		expectedCode  int
		checkResponse func(t *testing.T, recorder *http.Response)
	}{
		{
			name: "ValidCredentials",
			body: map[string]interface{}{
				"email":    user.Email,
				"password": "Test123!@#",
			},
			expectedCode: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *http.Response) {
				var response util.SuccessResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				data := response.Data.(map[string]interface{})
				require.NotEmpty(t, data["access_token"])
				require.NotEmpty(t, data["refresh_token"])
				require.NotNil(t, data["user"])
			},
		},
		{
			name: "InvalidPassword",
			body: map[string]interface{}{
				"email":    user.Email,
				"password": "wrongpassword",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "NonexistentUser",
			body: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "Test123!@#",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "MissingEmail",
			body: map[string]interface{}{
				"password": "Test123!@#",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "MissingPassword",
			body: map[string]interface{}{
				"email": user.Email,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/users/login", tc.body, "")
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.checkResponse != nil {
				tc.checkResponse(t, recorder.Result())
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	// Create test space and user
	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	// Login to get tokens
	loginBody := map[string]interface{}{
		"email":    user.Email,
		"password": "Test123!@#",
	}
	loginRecorder := ts.MakeRequest(t, http.MethodPost, "/api/users/login", loginBody, "")
	require.Equal(t, http.StatusOK, loginRecorder.Code)

	var loginResponse util.SuccessResponse
	err := json.Unmarshal(loginRecorder.Body.Bytes(), &loginResponse)
	require.NoError(t, err)

	loginData := loginResponse.Data.(map[string]interface{})
	refreshToken := loginData["refresh_token"].(string)

	testCases := []struct {
		name         string
		body         map[string]interface{}
		expectedCode int
	}{
		{
			name: "ValidRefreshToken",
			body: map[string]interface{}{
				"refresh_token": refreshToken,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "InvalidRefreshToken",
			body: map[string]interface{}{
				"refresh_token": "invalid-token",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "MissingRefreshToken",
			body:         map[string]interface{}{},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/users/refresh", tc.body, "")
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				var response util.SuccessResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				data := response.Data.(map[string]interface{})
				require.NotEmpty(t, data["access_token"])
			}
		})
	}
}

func TestLogout(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	// Create test space and user
	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	// Create auth token
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidToken",
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoToken",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidToken",
			token:        "invalid-token",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/users/logout", nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}
