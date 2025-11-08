package api_test

import (
	"database/sql"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	testhelpers "github.com/connect-univyn/connect-server/test/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
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
			name: "ValidEvent",
			body: map[string]interface{}{
				"title":       "Tech Meetup",
				"description": "Monthly tech meetup",
				"location":    "Building A",
				"start_date":  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"end_date":    time.Now().Add(26 * time.Hour).Format(time.RFC3339),
				"category":    "technology",
				"space_id":    spaceID.String(),
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingTitle",
			body: map[string]interface{}{
				"description": "Event without title",
				"start_time":  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"title":      "Unauthorized Event",
				"start_time": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/events", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestListEvents(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	
	_, err := ts.TestDB.Store.CreateEvent(context.Background(), db.CreateEventParams{
		SpaceID:     spaceID,
		Title:       "Test Event",
		Description: sql.NullString{String: "Test Description", Valid: true},
		Category:    "technology",
		StartDate:   time.Now().Add(24 * time.Hour),
		EndDate:     time.Now().Add(26 * time.Hour),
		Organizer:   uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	require.NoError(t, err)

	url := fmt.Sprintf("/api/events?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetEvent(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	event, err := ts.TestDB.Store.CreateEvent(context.Background(), db.CreateEventParams{
		SpaceID:     spaceID,
		Title:       "Test Event",
		Description: sql.NullString{String: "Test Description", Valid: true},
		Category:    "technology",
		StartDate:   time.Now().Add(24 * time.Hour),
		EndDate:     time.Now().Add(26 * time.Hour),
		Organizer:   uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		eventID      string
		expectedCode int
	}{
		{
			name:         "ValidEvent",
			eventID:      event.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidEventID",
			eventID:      "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "NonexistentEvent",
			eventID:      uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/events/%s", tc.eventID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestRegisterForEvent(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	organizer := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	attendee := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, attendee.ID)

	event, err := ts.TestDB.Store.CreateEvent(context.Background(), db.CreateEventParams{
		SpaceID:     spaceID,
		Title:       "Test Event",
		Description: sql.NullString{String: "Test Description", Valid: true},
		Category:    "technology",
		StartDate:   time.Now().Add(24 * time.Hour),
		EndDate:     time.Now().Add(26 * time.Hour),
		Organizer:   uuid.NullUUID{UUID: organizer.ID, Valid: true},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		eventID      string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidRegistration",
			eventID:      event.ID.String(),
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "NoAuth",
			eventID:      event.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/events/%s/register", tc.eventID)
			recorder := ts.MakeRequest(t, http.MethodPost, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestGetUpcomingEvents(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/events/upcoming?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestSearchEvents(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/events/search?q=tech&space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}
