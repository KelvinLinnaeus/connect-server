package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	testhelpers "github.com/connect-univyn/connect-server/test/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)





func TestCreateCommunity(t *testing.T) {
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
			name: "ValidCommunity",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
				"name":     "Test Community",
				"category": "general",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "ValidCommunityWithDescription",
			body: map[string]interface{}{
				"space_id":    spaceID.String(),
				"name":        "Community With Desc",
				"category":    "academic",
				"description": "A test community description",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingName",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
				"category": "general",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "MissingCategory",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
				"name":     "Test Community",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
				"name":     "Should Fail",
				"category": "general",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "InvalidSpaceID",
			body: map[string]interface{}{
				"space_id": "invalid-uuid",
				"name":     "Test Community",
				"category": "general",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/communities", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusCreated {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "name")
				RequireFieldExists(t, data, "category")
			}
		})
	}
}





func TestListCommunities(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)

	testCases := []struct {
		name         string
		query        string
		expectedCode int
	}{
		{
			name:         "ValidListWithPagination",
			query:        fmt.Sprintf("?space_id=%s&page=1&limit=20", spaceID.String()),
			expectedCode: http.StatusOK,
		},
		{
			name:         "ListWithoutPagination",
			query:        fmt.Sprintf("?space_id=%s", spaceID.String()),
			expectedCode: http.StatusOK,
		},
		{
			name:         "EmptyQuery",
			query:        "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities%s", tc.query)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestGetCommunity(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	community, err := ts.TestDB.Store.CreateCommunity(context.Background(), db.CreateCommunityParams{
		SpaceID:  spaceID,
		Name:     "Test Community",
		Category: "general",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		communityID  string
		expectedCode int
	}{
		{
			name:         "ValidCommunity",
			communityID:  community.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "NonexistentCommunity",
			communityID:  uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "InvalidCommunityID",
			communityID:  "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/%s", tc.communityID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "name")
			}
		})
	}
}





func TestGetCommunityBySlug(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	
	_, err := ts.TestDB.Store.CreateCommunity(context.Background(), db.CreateCommunityParams{
		SpaceID:  spaceID,
		Name:     "test-community-slug",
		Category: "general",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		slug         string
		expectedCode int
	}{
		{
			name:         "ValidSlug",
			slug:         "test-community-slug",
			expectedCode: http.StatusOK,
		},
		{
			name:         "NonexistentSlug",
			slug:         "nonexistent-slug-12345",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/slug/%s?space_id=%s", tc.slug, spaceID.String())
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestUpdateCommunity(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	community, err := ts.TestDB.Store.CreateCommunity(context.Background(), db.CreateCommunityParams{
		SpaceID:  spaceID,
		Name:     "Test Community",
		Category: "general",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		communityID  string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name:        "ValidUpdate",
			communityID: community.ID.String(),
			body: map[string]interface{}{
				"name":        "Updated Community Name",
				"description": "Updated description",
				"category":    "general",
			},
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			communityID:  community.ID.String(),
			body:         map[string]interface{}{"name": "Should Fail"},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidCommunityID",
			communityID:  "invalid-uuid",
			body:         map[string]interface{}{"name": "Test"},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/%s", tc.communityID)
			recorder := ts.MakeRequest(t, http.MethodPut, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestJoinCommunity(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	community, err := ts.TestDB.Store.CreateCommunity(context.Background(), db.CreateCommunityParams{
		SpaceID:  spaceID,
		Name:     "Test Community",
		Category: "general",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		communityID  string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidJoin",
			communityID:  community.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			communityID:  community.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidCommunityID",
			communityID:  "invalid-uuid",
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/%s/join", tc.communityID)
			recorder := ts.MakeRequest(t, http.MethodPost, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestLeaveCommunity(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	community, err := ts.TestDB.Store.CreateCommunity(context.Background(), db.CreateCommunityParams{
		SpaceID:  spaceID,
		Name:     "Test Community",
		Category: "general",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		communityID  string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidLeave",
			communityID:  community.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			communityID:  community.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/%s/leave", tc.communityID)
			recorder := ts.MakeRequest(t, http.MethodPost, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestSearchCommunities(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)

	testCases := []struct {
		name         string
		query        string
		expectedCode int
	}{
		{
			name:         "ValidSearch",
			query:        fmt.Sprintf("?q=test&space_id=%s", spaceID.String()),
			expectedCode: http.StatusOK,
		},
		{
			name:         "SearchWithCategory",
			query:        fmt.Sprintf("?q=test&category=general&space_id=%s", spaceID.String()),
			expectedCode: http.StatusOK,
		},
		{
			name:         "EmptySearch",
			query:        "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/search%s", tc.query)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestGetCommunityCategories(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/communities/categories?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}





func TestGetCommunityMembers(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	community, err := ts.TestDB.Store.CreateCommunity(context.Background(), db.CreateCommunityParams{
		SpaceID:  spaceID,
		Name:     "Test Community",
		Category: "general",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		communityID  string
		expectedCode int
	}{
		{
			name:         "ValidMembers",
			communityID:  community.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidCommunityID",
			communityID:  "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/%s/members", tc.communityID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestGetCommunityModerators(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	community, err := ts.TestDB.Store.CreateCommunity(context.Background(), db.CreateCommunityParams{
		SpaceID:  spaceID,
		Name:     "Test Community",
		Category: "general",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		communityID  string
		expectedCode int
	}{
		{
			name:         "ValidModerators",
			communityID:  community.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidCommunityID",
			communityID:  "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/%s/moderators", tc.communityID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestGetCommunityAdmins(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	community, err := ts.TestDB.Store.CreateCommunity(context.Background(), db.CreateCommunityParams{
		SpaceID:  spaceID,
		Name:     "Test Community",
		Category: "general",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		communityID  string
		expectedCode int
	}{
		{
			name:         "ValidAdmins",
			communityID:  community.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidCommunityID",
			communityID:  "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/%s/admins", tc.communityID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestAddCommunityModerator(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	moderator := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	community, err := ts.TestDB.Store.CreateCommunity(context.Background(), db.CreateCommunityParams{
		SpaceID:  spaceID,
		Name:     "Test Community",
		Category: "general",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		communityID  string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name:        "ValidAddModerator",
			communityID: community.ID.String(),
			body: map[string]interface{}{
				"user_id": moderator.ID.String(),
			},
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			communityID:  community.ID.String(),
			body:         map[string]interface{}{"user_id": moderator.ID.String()},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/communities/%s/moderators", tc.communityID)
			recorder := ts.MakeRequest(t, http.MethodPost, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}





func TestGetUserCommunities(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidRequest",
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/communities?space_id=%s", spaceID.String())
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}
