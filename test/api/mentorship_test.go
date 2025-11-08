package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/stretchr/testify/require"
)

func TestCreateMentorProfile(t *testing.T) {
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
			name: "ValidMentorProfile",
			body: map[string]interface{}{
				"space_id":    spaceID.String(),
				"industry":    "technology",
				"experience":  5,
				"specialties": []string{"backend", "cloud"},
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingBio",
			body: map[string]interface{}{
				"industries": []string{"technology"},
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"bio":        "Unauthorized",
				"industries": []string{"technology"},
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/mentorship/mentors/profile", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestSearchMentors(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/mentorship/mentors/search?industry=technology&space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestGetMentorProfile(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	_, err := ts.TestDB.Store.CreateMentorProfile(context.Background(), db.CreateMentorProfileParams{
		UserID:      user.ID,
		SpaceID:     spaceID,
		Industry:    "technology",
		Experience:  5,
		Specialties: []string{},
	})
	require.NoError(t, err)

	url := fmt.Sprintf("/api/mentorship/mentors/profile/%s?space_id=%s", user.ID.String(), spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestCreateTutorProfile(t *testing.T) {
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
			name: "ValidTutorProfile",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
				"subjects": []string{"mathematics", "physics"},
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"bio":      "Unauthorized",
				"subjects": []string{"mathematics"},
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/mentorship/tutors/profile", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestSearchTutors(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/mentorship/tutors/search?subject=mathematics&space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestCreateMentoringSession(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	mentor := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	mentee := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, mentee.ID)

	_, err := ts.TestDB.Store.CreateMentorProfile(context.Background(), db.CreateMentorProfileParams{
		UserID:      mentor.ID,
		SpaceID:     spaceID,
		Industry:    "technology",
		Experience:  5,
		Specialties: []string{},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name: "ValidSession",
			body: map[string]interface{}{
				"mentor_id":    mentor.ID.String(),
				"space_id":     spaceID.String(),
				"topic":        "Career guidance",
				"scheduled_at": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"duration":     60,
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"mentor_id": mentor.ID.String(),
				"scheduled_at":      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/mentorship/mentoring/sessions", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestGetUserMentoringSessions(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	url := fmt.Sprintf("/api/mentorship/mentoring/sessions?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, token)
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestCreateMentorApplication(t *testing.T) {
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
			name: "ValidApplication",
			body: map[string]interface{}{
				"space_id":    spaceID.String(),
				"industry":    "technology",
				"experience":  5,
				"specialties": []string{"backend", "cloud"},
				"motivation":  "I want to help others grow in their careers and share my knowledge and experience in software engineering",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"bio":        "Unauthorized",
				"industries": []string{"technology"},
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/mentorship/mentors/applications", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestGetPendingMentorApplications(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/mentorship/mentors/applications/pending?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}

func TestCreateTutoringSession(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	tutor := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	student := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, student.ID)

	tutorProfile, err := ts.TestDB.Store.CreateTutorProfile(context.Background(), db.CreateTutorProfileParams{
		UserID:   tutor.ID,
		SpaceID:  spaceID,
		Subjects: []string{"mathematics"},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name: "ValidSession",
			body: map[string]interface{}{
				"tutor_id":     tutor.ID.String(),
				"space_id":     spaceID.String(),
				"subject":      "Calculus",
				"scheduled_at": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"duration":     60,
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"tutor_profile_id": tutorProfile.ID.String(),
				"scheduled_at":     time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/mentorship/tutoring/sessions", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

func TestCreateTutorApplication(t *testing.T) {
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
			name: "ValidApplication",
			body: map[string]interface{}{
				"space_id":       spaceID.String(),
				"subjects":       []string{"Mathematics", "Physics"},
				"experience":     "I have been tutoring for 3 years and have helped over 50 students improve their grades",
				"qualifications": "Bachelor's degree in Mathematics, Master's degree in Education",
				"motivation":     "I am passionate about teaching and want to help students understand complex concepts in a simple way",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingSubjects",
			body: map[string]interface{}{
				"space_id":       spaceID.String(),
				"experience":     "I have been tutoring for 3 years",
				"qualifications": "Bachelor's degree in Mathematics",
				"motivation":     "I am passionate about teaching",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "MissingExperience",
			body: map[string]interface{}{
				"space_id":       spaceID.String(),
				"subjects":       []string{"Mathematics"},
				"qualifications": "Bachelor's degree in Mathematics",
				"motivation":     "I am passionate about teaching",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "MissingQualifications",
			body: map[string]interface{}{
				"space_id":   spaceID.String(),
				"subjects":   []string{"Mathematics"},
				"experience": "I have been tutoring for 3 years",
				"motivation": "I am passionate about teaching",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "ShortMotivation",
			body: map[string]interface{}{
				"space_id":       spaceID.String(),
				"subjects":       []string{"Mathematics"},
				"experience":     "I have been tutoring for 3 years",
				"qualifications": "Bachelor's degree in Mathematics",
				"motivation":     "Too short",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"space_id":       spaceID.String(),
				"subjects":       []string{"Mathematics"},
				"experience":     "I have been tutoring for 3 years",
				"qualifications": "Bachelor's degree in Mathematics",
				"motivation":     "I am passionate about teaching and want to help students understand complex concepts",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/mentorship/tutors/applications", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusCreated {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "subjects")
				RequireFieldExists(t, data, "experience")
				RequireFieldExists(t, data, "qualifications")
				RequireFieldExists(t, data, "motivation")
				RequireFieldExists(t, data, "status")
			}
		})
	}
}

func TestGetPendingTutorApplications(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	url := fmt.Sprintf("/api/mentorship/tutors/applications/pending?space_id=%s", spaceID.String())
	recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
	CheckResponseCode(t, recorder, http.StatusOK)
}
