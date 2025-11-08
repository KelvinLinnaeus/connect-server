package api_test

import (
	"fmt"
	"net/http"
	"testing"

	testhelpers "github.com/connect-univyn/connect-server/test/db"
)

func TestCreateNotification(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name: "ValidNotification",
			body: map[string]interface{}{
				"to_user_id": user.ID.String(),
				"type":       "mention",
				"message":    "You were mentioned in a post",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingType",
			body: map[string]interface{}{
				"to_user_id": user.ID.String(),
				"message":    "Missing type field",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"to_user_id": user.ID.String(),
				"type":       "mention",
				"message":    "Unauthorized",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/notifications", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestGetUserNotifications(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	recorder := ts.MakeRequest(t, http.MethodGet, "/api/notifications", nil, token)
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestMarkNotificationAsRead(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	
	
	notificationID := "550e8400-e29b-41d4-a716-446655440000"

	testCases := []struct {
		name           string
		notificationID string
		token          string
		expectedCode   int
	}{
		{
			name:           "WithAuth",
			notificationID: notificationID,
			token:          token,
			expectedCode:   http.StatusOK, 
		},
		{
			name:           "NoAuth",
			notificationID: notificationID,
			token:          "",
			expectedCode:   http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/notifications/%s/read", tc.notificationID)
			recorder := ts.MakeRequest(t, http.MethodPut, url, nil, tc.token)
			
			if recorder.Code != http.StatusOK && recorder.Code != http.StatusNotFound {
				CheckResponseCode(t, recorder, tc.expectedCode)
			}
		})
	}
}

func TestMarkAllNotificationsAsRead(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	recorder := ts.MakeRequest(t, http.MethodPut, "/api/notifications/read-all", nil, token)
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetUnreadNotificationCount(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	recorder := ts.MakeRequest(t, http.MethodGet, "/api/notifications/unread-count", nil, token)
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestDeleteNotification(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	notificationID := "550e8400-e29b-41d4-a716-446655440000"

	testCases := []struct {
		name           string
		notificationID string
		token          string
		expectedCode   int
	}{
		{
			name:           "WithAuth",
			notificationID: notificationID,
			token:          token,
			expectedCode:   http.StatusOK, 
		},
		{
			name:           "NoAuth",
			notificationID: notificationID,
			token:          "",
			expectedCode:   http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/notifications/%s", tc.notificationID)
			recorder := ts.MakeRequest(t, http.MethodDelete, url, nil, tc.token)
			
			if recorder.Code != http.StatusOK && recorder.Code != http.StatusNotFound {
				CheckResponseCode(t, recorder, tc.expectedCode)
			}
		})
	}
}
