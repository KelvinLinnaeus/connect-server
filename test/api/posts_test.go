package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	testhelpers "github.com/connect-univyn/connect_server/test/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// POST /api/posts - Create Post
// =============================================================================

func TestCreatePost(t *testing.T) {
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
			name: "ValidPost",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
				"content":  "Test post content",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "ValidPostWithMedia",
			body: map[string]interface{}{
				"space_id":  spaceID.String(),
				"content":   "Post with media",
				"media_url": "https://example.com/image.jpg",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingContent",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "MissingSpaceID",
			body: map[string]interface{}{
				"content": "Missing space ID",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "NoAuth",
			body: map[string]interface{}{
				"space_id": spaceID.String(),
				"content":  "Should fail",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "InvalidSpaceID",
			body: map[string]interface{}{
				"space_id": "invalid-uuid",
				"content":  "Invalid space",
			},
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := ts.MakeRequest(t, http.MethodPost, "/api/posts", tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusCreated {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "content")
				RequireFieldExists(t, data, "created_at")
			}
		})
	}
}

// =============================================================================
// GET /api/posts/:id - Get Post by ID
// =============================================================================

func TestGetPost(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	// Create a test post
	post, err := ts.TestDB.Store.CreatePost(context.Background(), db.CreatePostParams{
		SpaceID:  spaceID,
		AuthorID: user.ID,
		Content:  "Test post",
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		postID       string
		expectedCode int
	}{
		{
			name:         "ValidPost",
			postID:       post.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "NonexistentPost",
			postID:       uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "InvalidPostID",
			postID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/%s", tc.postID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)

			if tc.expectedCode == http.StatusOK {
				data := ParseSuccessResponse(t, recorder)
				RequireFieldExists(t, data, "id")
				RequireFieldExists(t, data, "content")
			}
		})
	}
}

// =============================================================================
// DELETE /api/posts/:id - Delete Post
// =============================================================================

func TestDeletePost(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	// Create a test post
	post, err := ts.TestDB.Store.CreatePost(context.Background(), db.CreatePostParams{
		SpaceID:  spaceID,
		AuthorID: user.ID,
		Content:  "Test post to delete",
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		postID       string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidDeletion",
			postID:       post.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			postID:       uuid.New().String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidPostID",
			postID:       "invalid-uuid",
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/%s", tc.postID)
			recorder := ts.MakeRequest(t, http.MethodDelete, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/feed - Get User Feed (Authenticated)
// =============================================================================

func TestGetUserFeed(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	testCases := []struct {
		name         string
		query        string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidFeedWithPagination",
			query:        "?page=1&limit=20",
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "FeedWithoutPagination",
			query:        "",
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			query:        "?page=1&limit=20",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/feed%s", tc.query)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/search - Search Posts
// =============================================================================

func TestSearchPosts(t *testing.T) {
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
			name:         "SearchWithoutQuery",
			query:        fmt.Sprintf("?space_id=%s", spaceID.String()),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "EmptySearch",
			query:        "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/search%s", tc.query)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/advanced-search - Advanced Search Posts
// =============================================================================

func TestAdvancedSearchPosts(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)

	testCases := []struct {
		name         string
		query        string
		expectedCode int
	}{
		{
			name:         "ValidAdvancedSearch",
			query:        fmt.Sprintf("?q=test&space_id=%s&sort_by=created_at", spaceID.String()),
			expectedCode: http.StatusOK,
		},
		{
			name:         "SearchWithFilters",
			query:        fmt.Sprintf("?space_id=%s&sort_by=likes", spaceID.String()),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/advanced-search%s", tc.query)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/trending - Get Trending Posts
// =============================================================================

func TestGetTrendingPosts(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)

	testCases := []struct {
		name         string
		query        string
		expectedCode int
	}{
		{
			name:         "ValidTrendingRequest",
			query:        fmt.Sprintf("?page=1&limit=10&space_id=%s", spaceID.String()),
			expectedCode: http.StatusOK,
		},
		{
			name:         "TrendingWithoutPagination",
			query:        fmt.Sprintf("?space_id=%s", spaceID.String()),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/trending%s", tc.query)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/:id/comments - Get Post Comments
// =============================================================================

func TestGetPostComments(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	post, err := ts.TestDB.Store.CreatePost(context.Background(), db.CreatePostParams{
		SpaceID:  spaceID,
		AuthorID: user.ID,
		Content:  "Test post",
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		postID       string
		expectedCode int
	}{
		{
			name:         "ValidPostComments",
			postID:       post.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidPostID",
			postID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "NonexistentPost",
			postID:       uuid.New().String(),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/%s/comments", tc.postID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// POST /api/posts/:id/comments - Create Comment
// =============================================================================

func TestCreateComment(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	post, err := ts.TestDB.Store.CreatePost(context.Background(), db.CreatePostParams{
		SpaceID:  spaceID,
		AuthorID: user.ID,
		Content:  "Test post",
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		postID       string
		body         map[string]interface{}
		token        string
		expectedCode int
	}{
		{
			name:   "ValidComment",
			postID: post.ID.String(),
			body: map[string]interface{}{
				"content": "Test comment",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name:   "MissingContent",
			postID: post.ID.String(),
			body:   map[string]interface{}{},
			token:  token,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "NoAuth",
			postID: post.ID.String(),
			body: map[string]interface{}{
				"content": "Should fail",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/%s/comments", tc.postID)
			recorder := ts.MakeRequest(t, http.MethodPost, url, tc.body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/:id/likes - Get Post Likes
// =============================================================================

func TestGetPostLikes(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	post, err := ts.TestDB.Store.CreatePost(context.Background(), db.CreatePostParams{
		SpaceID:  spaceID,
		AuthorID: user.ID,
		Content:  "Test post",
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		postID       string
		expectedCode int
	}{
		{
			name:         "ValidPostLikes",
			postID:       post.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidPostID",
			postID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/%s/likes", tc.postID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// POST /api/posts/:id/like - Toggle Post Like
// =============================================================================

func TestTogglePostLike(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	post, err := ts.TestDB.Store.CreatePost(context.Background(), db.CreatePostParams{
		SpaceID:  spaceID,
		AuthorID: user.ID,
		Content:  "Test post",
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		postID       string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidLike",
			postID:       post.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			postID:       post.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidPostID",
			postID:       "invalid-uuid",
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/%s/like", tc.postID)
			recorder := ts.MakeRequest(t, http.MethodPost, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/liked - Get User Liked Posts
// =============================================================================

func TestGetUserLikedPosts(t *testing.T) {
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
			recorder := ts.MakeRequest(t, http.MethodGet, "/api/posts/liked", nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// POST /api/posts/:id/repost - Create Repost
// =============================================================================

func TestCreateRepost(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	post, err := ts.TestDB.Store.CreatePost(context.Background(), db.CreatePostParams{
		SpaceID:  spaceID,
		AuthorID: user.ID,
		Content:  "Test post to repost",
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		postID       string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidRepost",
			postID:       post.ID.String(),
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "NoAuth",
			postID:       post.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidPostID",
			postID:       "invalid-uuid",
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/%s/repost", tc.postID)
			body := map[string]interface{}{}
			recorder := ts.MakeRequest(t, http.MethodPost, url, body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// PUT /api/posts/:id/pin - Pin Post
// =============================================================================

func TestPinPost(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	post, err := ts.TestDB.Store.CreatePost(context.Background(), db.CreatePostParams{
		SpaceID:  spaceID,
		AuthorID: user.ID,
		Content:  "Test post to pin",
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		postID       string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidPin",
			postID:       post.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			postID:       post.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/%s/pin", tc.postID)
			body := map[string]interface{}{}
			recorder := ts.MakeRequest(t, http.MethodPut, url, body, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/user/:user_id - Get User Posts
// =============================================================================

func TestGetUserPosts(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	testCases := []struct {
		name         string
		userID       string
		expectedCode int
	}{
		{
			name:         "ValidUserPosts",
			userID:       user.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidUserID",
			userID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "NonexistentUser",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/user/%s", tc.userID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/community/:community_id - Get Community Posts
// =============================================================================

func TestGetCommunityPosts(t *testing.T) {
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
			name:         "ValidCommunityPosts",
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
			url := fmt.Sprintf("/api/posts/community/%s", tc.communityID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// GET /api/posts/group/:group_id - Get Group Posts
// =============================================================================

func TestGetGroupPosts(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)

	group, err := ts.TestDB.Store.CreateGroup(context.Background(), db.CreateGroupParams{
		SpaceID:  spaceID,
		Name:     "Test Group",
		Category: "general",
		GroupType: "open",
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		groupID      string
		expectedCode int
	}{
		{
			name:         "ValidGroupPosts",
			groupID:      group.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "InvalidGroupID",
			groupID:      "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/posts/group/%s", tc.groupID)
			recorder := ts.MakeRequest(t, http.MethodGet, url, nil, "")
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}

// =============================================================================
// POST /api/comments/:id/like - Toggle Comment Like
// =============================================================================

func TestToggleCommentLike(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown()

	spaceID := testhelpers.CreateTestSpace(t, ts.TestDB.DB)
	user := testhelpers.CreateRandomUser(t, ts.TestDB.Store, spaceID)
	token := ts.CreateAuthToken(t, user.ID)

	// Create post and comment first
	post, err := ts.TestDB.Store.CreatePost(context.Background(), db.CreatePostParams{
		SpaceID:  spaceID,
		AuthorID: user.ID,
		Content:  "Test post",
	})
	require.NoError(t, err)

	comment, err := ts.TestDB.Store.CreateComment(context.Background(), db.CreateCommentParams{
		PostID:   post.ID,
		AuthorID: user.ID,
		Content:  "Test comment",
	})
	require.NoError(t, err)

	testCases := []struct {
		name         string
		commentID    string
		token        string
		expectedCode int
	}{
		{
			name:         "ValidCommentLike",
			commentID:    comment.ID.String(),
			token:        token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "NoAuth",
			commentID:    comment.ID.String(),
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "InvalidCommentID",
			commentID:    "invalid-uuid",
			token:        token,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/comments/%s/like", tc.commentID)
			recorder := ts.MakeRequest(t, http.MethodPost, url, nil, tc.token)
			CheckResponseCode(t, recorder, tc.expectedCode)
		})
	}
}
