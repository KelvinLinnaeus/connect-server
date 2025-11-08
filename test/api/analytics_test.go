package api_test

import (
	"fmt"
	"net/http"
	"testing"

	testhelpers "github.com/connect-univyn/connect_server/test/db"
)

func TestCreateReport(t *testing.T) {
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
			name: "ValidReport",
			body: map[string]interface{}{
				"space_id":     spaceID.String(),
				"content_type": "post",
				"content_id":   "550e8400-e29b-41d4-a716-446655440000",
				"reason":       "spam content",
				"description":  "This is spam content",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingReason",
			body: map[string]interface{}{
				"content_type": "post",
				"content_id":   "550e8400-e29b-41d4-a716-446655440000",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"content_type": "post",
				"content_id":   "550e8400-e29b-41d4-a716-446655440000",
				"reason":       "spam content",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/analytics/reports", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestGetPendingReports(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/reports/pending?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetModerationQueue(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/moderation/queue?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetContentModerationStats(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/moderation/stats?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetSystemMetrics(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/metrics/system?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetSpaceStats(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/metrics/space?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetEngagementMetrics(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/engagement/metrics?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetUserActivityStats(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/activity/stats?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetUserGrowth(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/users/growth?period=30d&space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetUserEngagementRanking(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/users/ranking?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetTopPosts(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/top/posts?period=7d&limit=10&space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetTopCommunities(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/top/communities?period=7d&limit=10&space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetTopGroups(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/top/groups?period=7d&limit=10&space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetMentoringStats(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/mentorship/mentoring?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetTutoringStats(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/mentorship/tutoring?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetPopularIndustries(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/mentorship/industries?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetPopularSubjects(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/analytics/mentorship/subjects?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}
