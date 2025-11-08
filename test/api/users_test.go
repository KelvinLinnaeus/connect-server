package api_test

import (
	"fmt"
	"net/http"
	"testing"

	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)

	testCases := []struct {
		name         string
		body         map[string]interface{}
		expectedCode int
	}{
		{
			name: "ValidUser",
			body: map[string]interface{}{
				"space_id":  spaceID.String(),
				"username":  "newuser",
				"email":     "newuser@example.com",
				"password":  "SecurePass123!",
				"full_name": "New User",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingEmail",
			body: map[string]interface{}{
				"space_id":  spaceID.String(),
				"username":  "testuser",
				"password":  "SecurePass123!",
				"full_name": "Test User",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "InvalidEmail",
			body: map[string]interface{}{
				"space_id":  spaceID.String(),
				"username":  "testuser2",
				"email":     "invalid-email",
				"password":  "SecurePass123!",
				"full_name": "Test User",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "WeakPassword",
			body: map[string]interface{}{
				"space_id":  spaceID.String(),
				"username":  "testuser3",
				"email":     "test3@example.com",
				"password":  "weak",
				"full_name": "Test User",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/users", tc.body, "")
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusCreated {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "username")
				RequireFieldExists(t, data, "email")
				RequireFieldNotExists(t, data, "password") // Password should not be returned
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		userID       string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidUser",
			userID:       user.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			userID:       user.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidUserID",
			userID:       "invalid-uuid",
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "NonexistentUser",
			userID:       uuid.New().String(),
			token:        token,
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/%s", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "username")
				RequireFieldExists(t, data, "email")
			}
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	testCases := []struct {
		name         string
		username     string
		expectedCode int
	}{
		{
			name:         "ValidUsername",
			username:     user.Username,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NonexistentUsername",
			username:     "nonexistent_user_12345",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/username/%s?space_id=%s", tc.username, spaceID.String())
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "username")
			}
		})
	}
}

func TestSearchUsers(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	testCases := []struct {
		name         string
		query        string
		expectedCode int
	}{
		{
			name:         "WithQuery",
			query:        fmt.Sprintf("?q=test&space_id=%s", spaceID.String()),
			expectedCode: http.StatusOK,
		},
		{
			name:         "EmptyQuery",
			query:        fmt.Sprintf("?space_id=%s", spaceID.String()),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/search%s", tc.query)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		userID       string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name:   "ValidUpdate",
			userID: user.ID.String(),
			body: map[string]interface{}{
				"full_name": "Updated Name",
				"bio":       "Updated bio",
			},
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:   "NoAuth",
			userID: user.ID.String(),
			body: map[string]interface{}{
				"full_name": "Unauthorized Update",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:   "InvalidUserID",
			userID: "invalid-uuid",
			body: map[string]interface{}{
				"full_name": "Updated Name",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/%s", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodPut, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "full_name")
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		userID       string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name:   "ValidPasswordUpdate",
			userID: user.ID.String(),
			body: map[string]interface{}{
				"old_password": "Test123!@#",
				"new_password": "NewSecure123!@#",
			},
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:   "MissingOldPassword",
			userID: user.ID.String(),
			body: map[string]interface{}{
				"new_password": "NewSecure123!@#",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "WeakNewPassword",
			userID: user.ID.String(),
			body: map[string]interface{}{
				"old_password": "Test123!@#",
				"new_password": "weak",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "NoAuth",
			userID: user.ID.String(),
			body: map[string]interface{}{
				"old_password": "Test123!@#",
				"new_password": "NewSecure123!@#",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/%s/password", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodPut, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestDeactivateUser(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		userID       string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidDeactivation",
			userID:       user.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			userID:       user.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/%s", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodDelete, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestFollowUser(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	follower := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	following := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, follower.ID)

	testCases := []struct {
		name         string
		followingID  string
		spaceID      string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidFollow",
			followingID:  following.ID.String(),
			spaceID:      spaceID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			followingID:  following.ID.String(),
			spaceID:      spaceID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidUserID",
			followingID:  "invalid-id",
			spaceID:      spaceID.String(),
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "MissingSpaceID",
			followingID:  following.ID.String(),
			spaceID:      "",
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/%s/follow?space_id=%s", tc.followingID, tc.spaceID)
			recorder := ts.MakeRequest(t, http.MethodPost, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "message")
			}
		})
	}
}

func TestUnfollowUser(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	follower := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	following := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, follower.ID)

	// First, follow the user
	testhelpers.CreateTestFollow(t, ts.TestDB.Store, follower.ID, following.ID, spaceID)

	testCases := []struct {
		name         string
		followingID  string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidUnfollow",
			followingID:  following.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			followingID:  following.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidUserID",
			followingID:  "invalid-id",
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/%s/follow", tc.followingID)
			recorder := ts.MakeRequest(t, http.MethodDelete, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "message")
			}
		})
	}
}

func TestCheckIfFollowing(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	follower := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	following := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	notFollowing := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, follower.ID)

	// Create a follow relationship
	testhelpers.CreateTestFollow(t, ts.TestDB.Store, follower.ID, following.ID, spaceID)

	testCases := []struct {
		name            string
		targetUserID    string
		token           string
		expectedCode    int
		expectedFollowing bool
	}{
		{
			name:            "IsFollowing",
			targetUserID:    following.ID.String(),
			token:           token,
			expectedCode:    http.StatusOK,
			expectedFollowing: true,
		},
		{
			name:            "NotFollowing",
			targetUserID:    notFollowing.ID.String(),
			token:           token,
			expectedCode:    http.StatusOK,
			expectedFollowing: false,
		},
		{
			name:            "NoAuth",
			targetUserID:    following.ID.String(),
			token:           "",
			expectedCode:    http.StatusUnauthorized,
			expectedFollowing: false,
		},
		{
			name:            "InvalidUserID",
			targetUserID:    "invalid-id",
			token:           token,
			expectedCode:    http.StatusBadRequest,
			expectedFollowing: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/%s/following/status", tc.targetUserID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "is_following")

				isFollowing, ok := data["is_following"].(bool)
				if !ok {
					t.Errorf("is_following should be a boolean")
				}

				if isFollowing != tc.expectedFollowing {
					t.Errorf("Expected is_following=%v, got %v", tc.expectedFollowing, isFollowing)
				}
			}
		})
	}
}

func TestGetFollowers(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	follower1 := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	follower2 := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	// Create followers
	testhelpers.CreateTestFollow(t, ts.TestDB.Store, follower1.ID, user.ID, spaceID)
	testhelpers.CreateTestFollow(t, ts.TestDB.Store, follower2.ID, user.ID, spaceID)

	testCases := []struct {
		name         string
		userID       string
		expectedCode int
		minFollowers int
	}{
		{
			name:         "ValidGetFollowers",
			userID:       user.ID.String(),
			expectedCode: http.StatusOK,
			minFollowers: 2,
		},
		{
			name:         "InvalidUserID",
			userID:       "invalid-id",
			expectedCode: http.StatusBadRequest,
			minFollowers: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/%s/followers", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				followers, ok := data["data"].([]interface{})
				if !ok {
					// If data is not directly an array, it might be wrapped
					dataMap, ok := data["data"].(map[string]interface{})
					if ok {
						followers, _ = dataMap["followers"].([]interface{})
					}
				}

				if followers != nil && len(followers) < tc.minFollowers {
					t.Errorf("Expected at least %d followers, got %d", tc.minFollowers, len(followers))
				}
			}
		})
	}
}

func TestGetFollowing(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	following1 := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	following2 := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	// Create following relationships
	testhelpers.CreateTestFollow(t, ts.TestDB.Store, user.ID, following1.ID, spaceID)
	testhelpers.CreateTestFollow(t, ts.TestDB.Store, user.ID, following2.ID, spaceID)

	testCases := []struct {
		name          string
		userID        string
		expectedCode  int
		minFollowing  int
	}{
		{
			name:          "ValidGetFollowing",
			userID:        user.ID.String(),
			expectedCode:  http.StatusOK,
			minFollowing:  2,
		},
		{
			name:          "InvalidUserID",
			userID:        "invalid-id",
			expectedCode:  http.StatusBadRequest,
			minFollowing:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/users/%s/following", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				following, ok := data["data"].([]interface{})
				if !ok {
					// If data is not directly an array, it might be wrapped
					dataMap, ok := data["data"].(map[string]interface{})
					if ok {
						following, _ = dataMap["following"].([]interface{})
					}
				}

				if following != nil && len(following) < tc.minFollowing {
					t.Errorf("Expected at least %d following, got %d", tc.minFollowing, len(following))
				}
			}
		})
	}
}
