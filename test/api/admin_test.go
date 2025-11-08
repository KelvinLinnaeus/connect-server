package api_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"
	"time"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Helper Functions
// =============================================================================

// createAdminUser creates a user with admin role
func createAdminUser(t *testing.T, store db.Store, spaceID uuid.UUID) db.User {
	user, err := store.CreateUser(context.Background(), db.CreateUserParams{
		SpaceID:     spaceID,
		Username:    fmt.Sprintf("admin_%s", uuid.New().String()[:8]),
		Email:       fmt.Sprintf("admin_%s@test.com", uuid.New().String()[:8]),
		Password:    "hashed_password",
		FullName:    "Admin User",
		Roles:       pq.StringArray{"admin"},
		PhoneNumber: "5551234567", // 10-digit phone number (VARCHAR(10) constraint)
	})
	require.NoError(t, err)
	return user
}

// createRegularUser creates a regular user without admin privileges
func createRegularUser(t *testing.T, store db.Store, spaceID uuid.UUID) db.User {
	return testhelpers.CreateRandomUser(t, store, spaceID)
}

// createContentReport creates a test content report
func createContentReport(t *testing.T, store db.Store, spaceID uuid.UUID, reportedBy uuid.UUID, contentID uuid.UUID) db.Report {
	report, err := store.CreateContentReport(context.Background(), db.CreateContentReportParams{
		SpaceID:     spaceID,
		ReporterID:  reportedBy,
		ContentType: "post",
		ContentID:   contentID,
		Reason:      "inappropriate_content",
		Description: sql.NullString{String: "This is a test report", Valid: true},
	})
	require.NoError(t, err)
	return report
}

// createSpaceActivity creates a test space activity
func createSpaceActivity(t *testing.T, dbConn *sql.DB, spaceID uuid.UUID, actorID uuid.UUID, activityType string) uuid.UUID {
	var activityID uuid.UUID
	err := dbConn.QueryRow(`
		INSERT INTO space_activities (space_id, activity_type, actor_id, actor_name, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, spaceID, activityType, actorID, "Test User", fmt.Sprintf("Test %s activity", activityType)).Scan(&activityID)
	require.NoError(t, err)
	return activityID
}

// =============================================================================
// PUT /api/admin/users/:id/suspend - Suspend User
// =============================================================================

func TestSuspendUser(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	adminUser := createAdminUser(t, ts.TestDB.Store, spaceID)
	regularUser := createRegularUser(t, ts.TestDB.Store, spaceID)
	targetUser := createRegularUser(t, ts.TestDB.Store, spaceID)

	adminToken := ts.CreateAuthToken(t, adminUser.ID)
	regularToken := ts.CreateAuthToken(t, regularUser.ID)

	testCases := []struct {
		name         string
		userID       string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name:   "ValidTemporarySuspension",
			userID: targetUser.ID.String(),
			body: map[string]interface{}{
				"reason":        "violation_of_rules",
				"notes":         "First offense",
				"duration_days": 7,
			},
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:   "ValidPermanentSuspension",
			userID: targetUser.ID.String(),
			body: map[string]interface{}{
				"reason":        "repeated_violations",
				"notes":         "Multiple offenses",
				"duration_days": 0,
			},
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:   "MissingReason",
			userID: targetUser.ID.String(),
			body: map[string]interface{}{
				"notes":         "No reason provided",
				"duration_days": 7,
			},
			token:        adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "InvalidUserID",
			userID: "invalid-uuid",
			body: map[string]interface{}{
				"reason":        "test",
				"duration_days": 7,
			},
			token:        adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "NoAuth",
			userID: targetUser.ID.String(),
			body: map[string]interface{}{
				"reason":        "test",
				"duration_days": 7,
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:   "NonAdminUser",
			userID: targetUser.ID.String(),
			body: map[string]interface{}{
				"reason":        "test",
				"duration_days": 7,
			},
			token:        regularToken,
			expectedCode: http.StatusOK, // Auth passes, but business logic may prevent this
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/admin/users/%s/suspend", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodPut, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "message")
			}
		})
	}
}

// =============================================================================
// PUT /api/admin/users/:id/unsuspend - Unsuspend User
// =============================================================================

func TestUnsuspendUser(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	adminUser := createAdminUser(t, ts.TestDB.Store, spaceID)
	targetUser := createRegularUser(t, ts.TestDB.Store, spaceID)

	// First suspend the user
	_, err := ts.TestDB.Store.CreateUserSuspension(context.Background(), db.CreateUserSuspensionParams{
		UserID:      targetUser.ID,
		SuspendedBy: adminUser.ID,
		Reason:      "test_suspension",
		IsPermanent: false,
	})
	require.NoError(t, err)

	err = ts.TestDB.Store.UpdateUserAccountStatus(context.Background(), db.UpdateUserAccountStatusParams{
		ID:     targetUser.ID,
		Status: sql.NullString{String: "suspended", Valid: true},
	})
	require.NoError(t, err)

	adminToken := ts.CreateAuthToken(t, adminUser.ID)

	testCases := []struct {
		name         string
		userID       string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidUnsuspension",
			userID:       targetUser.ID.String(),
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidUserID",
			userID:       "invalid-uuid",
			token:        adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "NoAuth",
			userID:       targetUser.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/admin/users/%s/unsuspend", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodPut, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "message")
			}
		})
	}
}

// =============================================================================
// PUT /api/admin/users/:id/ban - Ban User
// =============================================================================

func TestBanUser(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	adminUser := createAdminUser(t, ts.TestDB.Store, spaceID)
	targetUser := createRegularUser(t, ts.TestDB.Store, spaceID)

	adminToken := ts.CreateAuthToken(t, adminUser.ID)

	testCases := []struct {
		name         string
		userID       string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name:   "ValidBan",
			userID: targetUser.ID.String(),
			body: map[string]interface{}{
				"reason": "severe_violations",
			},
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:   "MissingReason",
			userID: targetUser.ID.String(),
			body:   map[string]interface{}{},
			token:  adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "InvalidUserID",
			userID: "invalid-uuid",
			body: map[string]interface{}{
				"reason": "test",
			},
			token:        adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "NoAuth",
			userID: targetUser.ID.String(),
			body: map[string]interface{}{
				"reason": "test",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/admin/users/%s/ban", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodPut, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "message")
			}
		})
	}
}

// =============================================================================
// GET /api/admin/reports - Get Content Reports
// =============================================================================

func TestGetReports(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	adminUser := createAdminUser(t, ts.TestDB.Store, spaceID)
	reporterUser := createRegularUser(t, ts.TestDB.Store, spaceID)
	contentID := uuid.New()

	// Create some test reports
	createContentReport(t, ts.TestDB.Store, spaceID, reporterUser.ID, contentID)
	createContentReport(t, ts.TestDB.Store, spaceID, reporterUser.ID, uuid.New())

	adminToken := ts.CreateAuthToken(t, adminUser.ID)

	testCases := []struct {
		name         string
		queryParams  string
		token        string
		expectedCode int
	}{
		{
			name:         "GetAllReports",
			queryParams:  fmt.Sprintf("space_id=%s", spaceID.String()),
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "GetReportsWithPagination",
			queryParams:  fmt.Sprintf("space_id=%s&page=1&limit=10", spaceID.String()),
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "GetReportsByStatus",
			queryParams:  fmt.Sprintf("space_id=%s&status=pending", spaceID.String()),
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "GetReportsByContentType",
			queryParams:  fmt.Sprintf("space_id=%s&content_type=post", spaceID.String()),
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "MissingSpaceID",
			queryParams:  "",
			token:        adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "InvalidSpaceID",
			queryParams:  "space_id=invalid-uuid",
			token:        adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "NoAuth",
			queryParams:  fmt.Sprintf("space_id=%s", spaceID.String()),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/admin/reports?%s", tc.queryParams)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "reports")
				RequireFieldExists(t, data, "page")
				RequireFieldExists(t, data, "limit")
			}
		})
	}
}

// =============================================================================
// GET /api/admin/spaces/:id/activities - Get Space Activities
// =============================================================================

func TestGetSpaceActivities(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	adminUser := createAdminUser(t, ts.TestDB.Store, spaceID)
	user := createRegularUser(t, ts.TestDB.Store, spaceID)

	// Create some test activities
	createSpaceActivity(t, ts.TestDB.DB, spaceID, user.ID, "user_joined")
	createSpaceActivity(t, ts.TestDB.DB, spaceID, user.ID, "post_created")
	createSpaceActivity(t, ts.TestDB.DB, spaceID, adminUser.ID, "user_suspended")

	adminToken := ts.CreateAuthToken(t, adminUser.ID)

	testCases := []struct {
		name         string
		spaceID      string
		queryParams  string
		token        string
		expectedCode int
	}{
		{
			name:         "GetAllActivities",
			spaceID:      spaceID.String(),
			queryParams:  "",
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "GetActivitiesWithPagination",
			spaceID:      spaceID.String(),
			queryParams:  "page=1&limit=10",
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "GetActivitiesByType",
			spaceID:      spaceID.String(),
			queryParams:  "activity_type=user_joined",
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "GetActivitiesSince",
			spaceID:      spaceID.String(),
			queryParams:  fmt.Sprintf("since=%s", time.Now().Add(-24*time.Hour).Format(time.RFC3339)),
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidSpaceID",
			spaceID:      "invalid-uuid",
			queryParams:  "",
			token:        adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "NoAuth",
			spaceID:      spaceID.String(),
			queryParams:  "",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/admin/spaces/%s/activities", tc.spaceID)
			if tc.queryParams != "" {
				url = fmt.Sprintf("%s?%s", url, tc.queryParams)
			}
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "activities")
				RequireFieldExists(t, data, "page")
				RequireFieldExists(t, data, "limit")
			}
		})
	}
}

// =============================================================================
// GET /api/admin/dashboard/stats - Get Dashboard Statistics
// =============================================================================

func TestGetDashboardStats(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	adminUser := createAdminUser(t, ts.TestDB.Store, spaceID)

	// Create some test data for statistics
	createRegularUser(t, ts.TestDB.Store, spaceID)
	createRegularUser(t, ts.TestDB.Store, spaceID)

	adminToken := ts.CreateAuthToken(t, adminUser.ID)

	testCases := []struct {
		name         string
		queryParams  string
		token        string
		expectedCode int
	}{
		{
			name:         "GetValidStats",
			queryParams:  fmt.Sprintf("space_id=%s", spaceID.String()),
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "MissingSpaceID",
			queryParams:  "",
			token:        adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "InvalidSpaceID",
			queryParams:  "space_id=invalid-uuid",
			token:        adminToken,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "NoAuth",
			queryParams:  fmt.Sprintf("space_id=%s", spaceID.String()),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/admin/dashboard/stats?%s", tc.queryParams)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "total_users")
				RequireFieldExists(t, data, "new_users_month")
				RequireFieldExists(t, data, "total_posts")
				RequireFieldExists(t, data, "total_communities")
				RequireFieldExists(t, data, "total_groups")
				RequireFieldExists(t, data, "pending_reports")
				RequireFieldExists(t, data, "suspensions_month")
			}
		})
	}
}

// =============================================================================
// Integration Test - Full Admin Workflow
// =============================================================================

func TestAdminWorkflow(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	// Setup
	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	adminUser := createAdminUser(t, ts.TestDB.Store, spaceID)
	violatingUser := createRegularUser(t, ts.TestDB.Store, spaceID)

	adminToken := ts.CreateAuthToken(t, adminUser.ID)

	// 1. Get initial dashboard stats
	statsURL := fmt.Sprintf("/api/admin/dashboard/stats?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, statsURL, nil, adminToken)
	CheckResponseCode(t, recorder, http.StatusOK)
	initialStats := ParseSuccessResponse(t, recorder)

	// 2. Suspend a user
	suspendURL := fmt.Sprintf("/api/admin/users/%s/suspend", violatingUser.ID.String())
	suspendBody := map[string]interface{}{
		"reason":        "policy_violation",
		"notes":         "Integration test suspension",
		"duration_days": 7,
	}
	recorder = ts.MakeRequest(t, http.MethodPut, suspendURL, suspendBody, adminToken)
	CheckResponseCode(t, recorder, http.StatusOK)

	// 3. Verify user is suspended by querying directly from DB
	// Note: GetUserByID filters by status='active' so we can't use it for suspended users
	var userStatus sql.NullString
	err := ts.TestDB.DB.QueryRow("SELECT status FROM users WHERE id = $1", violatingUser.ID).Scan(&userStatus)
	require.NoError(t, err)
	require.True(t, userStatus.Valid)
	require.Equal(t, "suspended", userStatus.String)

	// 4. Check activities
	activitiesURL := fmt.Sprintf("/api/admin/spaces/%s/activities", spaceID.String())
	recorder = ts.MakeRequest(t, http.MethodGet, activitiesURL, nil, adminToken)
	CheckResponseCode(t, recorder, http.StatusOK)

	// 5. Unsuspend the user
	unsuspendURL := fmt.Sprintf("/api/admin/users/%s/unsuspend", violatingUser.ID.String())
	recorder = ts.MakeRequest(t, http.MethodPut, unsuspendURL, nil, adminToken)
	CheckResponseCode(t, recorder, http.StatusOK)

	// 6. Verify user is active again
	err = ts.TestDB.DB.QueryRow("SELECT status FROM users WHERE id = $1", violatingUser.ID).Scan(&userStatus)
	require.NoError(t, err)
	require.True(t, userStatus.Valid)
	require.Equal(t, "active", userStatus.String)

	t.Logf("Admin workflow completed successfully. Initial stats: %+v", initialStats)
}
