package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
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
			name: "ValidGroup",
			body: map[string]interface{}{
				"space_id":   spaceID.String(),
				"name":       "Test Group",
				"slug":       "test-group",
				"category":   "general",
				"group_type": "public",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"space_id":   spaceID.String(),
				"name":       "Should Fail",
				"slug":       "should-fail",
				"category":   "general",
				"group_type": "public",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/groups", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestListGroups(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)

	url := fmt.Sprintf("/api/groups?space_id=%s&page=1&limit=20", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestJoinGroup(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	// Create a group
	group, err := ts.TestDB.Store.CreateGroup(context.Background(), db.CreateGroupParams{
		SpaceID:   spaceID,
		Name:      "Test Group",
		Category:  "general",
		GroupType: "study",
	})
	require.NoError(t, err)

	url := fmt.Sprintf("/api/groups/%s/join", group.ID.String())
	recorder := ts.MakeRequest(t, http.MethodPost, url, nil, token)
	CheckResponseCode(t, recorder, http.StatusOK)
}
