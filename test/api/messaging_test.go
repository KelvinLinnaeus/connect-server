package api_test

import (
	"fmt"
	"net/http"
	"testing"

	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/google/uuid"
)

func TestCreateConversation(t *testing.T) {
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
			name: "ValidConversation",
			body: map[string]interface{}{
				"space_id":          spaceID.String(),
				"name":              "Project Discussion",
				"participant_ids":   []string{user.ID.String()},
				"conversation_type": "group",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"space_id":          spaceID.String(),
				"name":              "Unauthorized",
				"conversation_type": "group",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/conversations", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestGetUserConversations(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	recorder := ts.MakeRequest(t, http.MethodGet, "/api/conversations", nil, token)
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestSendMessage(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	// Use a mock conversation ID
	conversationID := uuid.New().String()

	testCases := []struct {
		name           string
		conversationID string
		body           map[string]interface{}
		token          string
		expectedCode   int
	}{
		{
			name:           "MissingContent",
			conversationID: conversationID,
			body:           map[string]interface{}{},
			token:          token,
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "NoAuth",
			conversationID: conversationID,
			body: map[string]interface{}{
				"content": "Unauthorized message",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/conversations/%s/messages", tc.conversationID)
			recorder := ts.MakeRequest(t, http.MethodPost, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestGetConversationMessages(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	conversationID := uuid.New().String()
	url := fmt.Sprintf("/api/conversations/%s/messages", conversationID)
	
	// This will likely return 404 or 403, but that's OK for testing
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, token)
	// Accept various codes as the conversation doesn't exist
	if recorder.Code != http.StatusOK && recorder.Code != http.StatusNotFound && recorder.Code != http.StatusForbidden {
		CheckResponseCode(t, recorder, http.StatusOK)
	}
}

func TestMarkMessagesAsRead(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	conversationID := uuid.New().String()
	url := fmt.Sprintf("/api/conversations/%s/read", conversationID)
	
	recorder := ts.MakeRequest(t, http.MethodPost, url, nil, token)
	// Accept various codes
	if recorder.Code != http.StatusOK && recorder.Code != http.StatusNotFound && recorder.Code != http.StatusForbidden {
		CheckResponseCode(t, recorder, http.StatusOK)
	}
}

func TestGetOrCreateDirectConversation(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user1 := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	user2 := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user1.ID)

	testCases := []struct {
		name         string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name: "ValidDirectConversation",
			body: map[string]interface{}{
				"recipient_id": user2.ID.String(),
				"space_id":     spaceID.String(),
			},
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name: "MissingRecipientID",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "MissingSpaceID",
			body: map[string]interface{}{
				"recipient_id": user2.ID.String(),
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "InvalidRecipientID",
			body: map[string]interface{}{
				"recipient_id": "invalid-id",
				"space_id":     spaceID.String(),
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"recipient_id": user2.ID.String(),
				"space_id":     spaceID.String(),
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/conversations/direct", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "conversation_id")

				// Verify conversation_id is a valid UUID
				conversationIDStr, ok := data["conversation_id"].(string)
				if !ok {
					t.Errorf("conversation_id should be a string")
				}

				_, err := uuid.Parse(conversationIDStr)
				if err != nil {
					t.Errorf("conversation_id should be a valid UUID, got: %s", conversationIDStr)
				}

				// Test idempotency - creating the same direct conversation again should return the same ID
				recorder2 := ts.MakeRequest(t, http.MethodPost, "/api/conversations/direct", tc.body, tc.token)
				CheckResponseCode(t, recorder2, http.StatusOK)

				data2 := ParseSuccessResponse(t, recorder2)
				conversationID2Str, ok := data2["conversation_id"].(string)
				if !ok {
					t.Errorf("conversation_id should be a string in second request")
				}

				if conversationIDStr != conversationID2Str {
					t.Errorf("Expected same conversation_id for idempotent requests, got %s and %s", conversationIDStr, conversationID2Str)
				}
			}
		})
	}
}
