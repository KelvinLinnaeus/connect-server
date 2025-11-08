package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)


type TestDB struct {
	DB    *sql.DB
	Store db.Store
}


func SetupTestDB(t *testing.T) *TestDB {
	
	
	
	databaseURL := "postgres://postgres:Kelvin@localhost:5432/connect_test?sslmode=disable"

	conn, err := sql.Open("postgres", databaseURL)
	require.NoError(t, err, "Failed to open database connection")
	require.NotNil(t, conn)

	
	conn.SetMaxOpenConns(10)                  
	conn.SetMaxIdleConns(2)                   
	conn.SetConnMaxLifetime(5 * time.Minute)  

	err = conn.Ping()
	require.NoError(t, err, "Failed to ping database")

	store := db.NewStore(conn)

	return &TestDB{
		DB:    conn,
		Store: store,
	}
}


func (testDB *TestDB) TeardownTestDB() {
	testDB.DB.Close()
}


func CreateRandomUser(t *testing.T, store db.Store, spaceID uuid.UUID) db.User {
	hashedPassword, err := auth.HashPassword("Test123!@#")
	require.NoError(t, err)

	randomEmail := fmt.Sprintf("test_%s@example.com", uuid.New().String()[:8])
	randomUsername := fmt.Sprintf("user_%s", uuid.New().String()[:8])
	randomPhone := fmt.Sprintf("555%07d", uuid.New().ID()%10000000) 

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


func CreateTestSpace(t *testing.T, db *sql.DB) uuid.UUID {
	
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


func CleanupTestData(t *testing.T, db *sql.DB) {
	
	
	tables := []string{
		
		"user_sessions",
		"login_attempts",

		
		"likes",
		"comments",
		"posts",
		"follows",

		
		"message_reads",
		"messages",
		"conversations",
		"conversation_participants",

		
		"notifications",

		
		"group_applications",
		"group_roles",
		"group_members",
		"groups",
		"community_members",
		"communities",

		
		"event_registrations",
		"events",
		"announcements",

		
		"mentorship_sessions",
		"mentorship_requests",

		
		"space_activities",
		"audit_logs",
		"user_suspensions",

		
		"users",
		"spaces",
	}

	for _, table := range tables {
		
		_, _ = db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		
	}
}


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
