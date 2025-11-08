package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	testhelpers "github.com/connect-univyn/connect-server/test/db"
)

func TestCreateSpace(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	
	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	
	uniqueSlug := fmt.Sprintf("new-university-%s", uuid.New().String()[:8])

	testCases := []struct {
		name         string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name: "ValidSpace",
			body: map[string]interface{}{
				"name": "New University",
				"slug": uniqueSlug,
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingName",
			body: map[string]interface{}{
				"slug": fmt.Sprintf("test-slug-%s", uuid.New().String()[:8]),
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "MissingSlug",
			body: map[string]interface{}{
				"name": "Test Name",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"name": "Test Space",
				"slug": fmt.Sprintf("test-space-noauth-%s", uuid.New().String()[:8]),
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/spaces", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusCreated {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "name")
				RequireFieldExists(t, data, "slug")
			}
		})
	}
}

func TestListSpaces(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	
	testhelpers.CreateTestSpace(t, ts.TestDB.DB)

	recorder := ts.MakeRequest(t, http.MethodGet, "/api/spaces?page=1&limit=20", nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)

	data := ParseSuccessResponse(t, recorder)
	RequireFieldExists(t, data, "spaces")
}

func TestGetSpace(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	
	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)

	testCases := []struct {
		name         string
		spaceID      string
		expectedCode int
	}{
		{
			name:         "ValidID",
			spaceID:      spaceID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidID",
			spaceID:      "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "NonexistentID",
			spaceID:      "00000000-0000-0000-0000-000000000000",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/spaces/%s", tc.spaceID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestUpdateSpace(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	
	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		spaceID      string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name:    "ValidUpdate",
			spaceID: spaceID.String(),
			body: map[string]interface{}{
				"name":        "Updated University Name",
				"description": "This is an updated description",
				"location":    "New Location",
			},
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:    "NoAuth",
			spaceID: spaceID.String(),
			body: map[string]interface{}{
				"name": "Should Fail",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidID",
			spaceID:      "invalid-uuid",
			body:         map[string]interface{}{"name": "Test"},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/spaces/%s", tc.spaceID)
			recorder := ts.MakeRequest(t, http.MethodPut, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "name")
			}
		})
	}
}

func TestDeleteSpace(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	
	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		spaceID      string
		token        string
		expectedCode int
	}{
		{
			name:         "NoAuth",
			spaceID:      spaceID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ValidDelete",
			spaceID:      spaceID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/spaces/%s", tc.spaceID)
			recorder := ts.MakeRequest(t, http.MethodDelete, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}
