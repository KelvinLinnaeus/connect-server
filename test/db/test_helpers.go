package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// TestDB wraps database connection for testing
type TestDB struct {
	DB    *sql.DB
	Store db.Store
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *TestDB {
	// SECURITY: Test database credentials must be set via environment variable
	// Set TEST_DATABASE_URL environment variable before running tests
	// Example: export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/connect?sslmode=disable"
	databaseURL := "postgres://postgres:Kelvin@localhost:5432/connect_test?sslmode=disable"

	conn, err := sql.Open("postgres", databaseURL)
	require.NoError(t, err, "Failed to open database connection")
	require.NotNil(t, conn)

	// Configure connection pool to prevent exhaustion during tests
	conn.SetMaxOpenConns(10)                  // Lower limit for tests
	conn.SetMaxIdleConns(2)                   // Lower idle connections for tests
	conn.SetConnMaxLifetime(5 * time.Minute)  // Maximum connection lifetime

	err = conn.Ping()
	require.NoError(t, err, "Failed to ping database")

	store := db.NewStore(conn)

	return &TestDB{
		DB:    conn,
		Store: store,
	}
}

// TeardownTestDB closes the database connection
func (testDB *TestDB) TeardownTestDB() {
	testDB.DB.Close()
}

// CreateRandomUser creates a random user for testing
func CreateRandomUser(t *testing.T, store db.Store, spaceID uuid.UUID) db.User {
	hashedPassword, err := auth.HashPassword("Test123!@#")
	require.NoError(t, err)

	randomEmail := fmt.Sprintf("test_%s@example.com", uuid.New().String()[:8])
	randomUsername := fmt.Sprintf("user_%s", uuid.New().String()[:8])
	randomPhone := fmt.Sprintf("555%07d", uuid.New().ID()%10000000) // Generate random 10-digit phone starting with 555

	user, err := store.CreateUser(context.Background(), db.CreateUserParams{
		SpaceID:     spaceID,
		Username:    randomUsername,
		Email:       randomEmail,
		Password:    hashedPassword,
		FullName:    "Test User",
		Roles:       []string{"user"},
		PhoneNumber: randomPhone,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

// CreateTestSpace creates a space for testing
func CreateTestSpace(t *testing.T, db *sql.DB) uuid.UUID {
	// Create a unique space for each test run to avoid conflicts
	uniqueSlug := fmt.Sprintf("test-space-%s", uuid.New().String()[:8])
	uniqueName := fmt.Sprintf("Test Space %s", uuid.New().String()[:8])

	var spaceID uuid.UUID
	err := db.QueryRow(`
		INSERT INTO spaces (name, slug, type, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, uniqueName, uniqueSlug, "university", "active").Scan(&spaceID)
	require.NoError(t, err, "Failed to create test space")
	return spaceID
}

// CleanupTestData removes test data
func CleanupTestData(t *testing.T, db *sql.DB) {
	// Clean up test data in reverse order of dependencies
	// This list should include all tables to ensure complete cleanup between tests
	tables := []string{
		// Sessions and auth
		"user_sessions",
		"login_attempts",

		// Social features
		"likes",
		"comments",
		"posts",
		"follows",

		// Messaging
		"message_reads",
		"messages",
		"conversations",
		"conversation_participants",

		// Notifications
		"notifications",

		// Communities and Groups
		"group_applications",
		"group_roles",
		"group_members",
		"groups",
		"community_members",
		"communities",

		// Events and Announcements
		"event_registrations",
		"events",
		"announcements",

		// Mentorship
		"mentorship_sessions",
		"mentorship_requests",

		// Admin
		"space_activities",
		"audit_logs",
		"user_suspensions",

		// Users and Spaces (last)
		"users",
		"spaces",
	}

	for _, table := range tables {
		// Use CASCADE delete where appropriate, but explicit cleanup is safer
		_, _ = db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		// Don't require.NoError here as some tables might not exist in all test scenarios
	}
}

// CreateTestFollow creates a follow relationship for testing
func CreateTestFollow(t *testing.T, store db.Store, followerID, followingID, spaceID uuid.UUID) db.Follow {
	follow, err := store.FollowUser(context.Background(), db.FollowUserParams{
		FollowerID:  followerID,
		FollowingID: followingID,
		SpaceID:     spaceID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, follow)

	return follow
}
