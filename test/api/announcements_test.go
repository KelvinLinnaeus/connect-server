package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateAnnouncement(t *testing.T) {
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
			name: "ValidAnnouncement",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
				"title":    "Important Update",
				"content":  "System maintenance scheduled",
				"type":     "info",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingTitle",
			body: map[string]interface{}{
				"content": "Content without title",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"title":   "Unauthorized",
				"content": "Test",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/announcements", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestListAnnouncements(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	// Create test announcement
	_, err := ts.TestDB.Store.CreateAnnouncement(context.Background(), db.CreateAnnouncementParams{
		SpaceID: spaceID,
		Title:   "Test Announcement",
		Content: "Test Content",
		Type:    "general",
		TargetAudience: []string{},
		AuthorID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	require.NoError(t, err)

	url := fmt.Sprintf("/api/announcements?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetAnnouncement(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	announcement, err := ts.TestDB.Store.CreateAnnouncement(context.Background(), db.CreateAnnouncementParams{
		SpaceID:        spaceID,
		Title:          "Test Announcement",
		Content:        "Test Content",
		Type:           "general",
		TargetAudience: []string{},
		AuthorID:       uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	require.NoError(t, err)

	testCases := []struct {
		name           string
		announcementID string
		expectedCode   int
	}{
		{
			name:           "ValidAnnouncement",
			announcementID: announcement.ID.String(),
			expectedCode:   http.StatusOK,
		},
		{
			name:           "InvalidAnnouncementID",
			announcementID: "invalid-uuid",
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "NonexistentAnnouncement",
			announcementID: uuid.New().String(),
			expectedCode:   http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/announcements/%s", tc.announcementID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestUpdateAnnouncement(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	announcement, err := ts.TestDB.Store.CreateAnnouncement(context.Background(), db.CreateAnnouncementParams{
		SpaceID:        spaceID,
		Title:          "Original Title",
		Content:        "Original Content",
		Type:           "general",
		TargetAudience: []string{},
		AuthorID:       uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	require.NoError(t, err)

	testCases := []struct {
		name           string
		announcementID string
		body           map[string]interface{}
		token          string
		expectedCode   int
	}{
		{
			name:           "ValidUpdate",
			announcementID: announcement.ID.String(),
			body: map[string]interface{}{
				"title":   "Updated Title",
				"content": "Updated Content",
				"type":    "info",
			},
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:           "NoAuth",
			announcementID: announcement.ID.String(),
			body: map[string]interface{}{
				"title": "Unauthorized Update",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/announcements/%s", tc.announcementID)
			recorder := ts.MakeRequest(t, http.MethodPut, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}
