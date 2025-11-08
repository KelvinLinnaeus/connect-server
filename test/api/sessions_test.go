package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetSession(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	// Create test space and user
	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	// Create a session
	session, err := ts.TestDB.Store.CreateSession(context.Background(), db.CreateSessionParams{
		ID:           uuid.New(),
		UserID:       user.ID,
		Username:     user.Username,
		RefreshToken: "test-refresh-token",
		UserAgent:    "test-agent",
		IsBlocked:    false,
		SpaceID:      spaceID,
		ExpiresAt:    time.Now().Add(24 * time.Hour), // Set expiry to 24 hours from now
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		sessionID    string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidSession",
			sessionID:    session.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			sessionID:    session.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidSessionID",
			sessionID:    "invalid-uuid",
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/sessions/%s", tc.sessionID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}
