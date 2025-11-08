package posts

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/connect-univyn/connect-server/internal/live"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/sqlc-dev/pqtype"
)


type Service struct {
	store       db.Store
	liveService *live.Service
}


func NewService(store db.Store, liveService *live.Service) *Service {
	return &Service{
		store:       store,
		liveService: liveService,
	}
}


func interfaceToStringPtr(val interface{}) *string {
	if val == nil {
		return nil
	}
	if str, ok := val.(string); ok {
		return &str
	}
	return nil
}


func (s *Service) CreatePost(ctx context.Context, req CreatePostRequest) (*PostResponse, error) {
	var communityID, groupID, parentPostID, quotedPostID uuid.NullUUID
	var visibility sql.NullString

	if req.CommunityID != nil {
		communityID = uuid.NullUUID{UUID: *req.CommunityID, Valid: true}
	}
	if req.GroupID != nil {
		groupID = uuid.NullUUID{UUID: *req.GroupID, Valid: true}
	}
	if req.ParentPostID != nil {
		parentPostID = uuid.NullUUID{UUID: *req.ParentPostID, Valid: true}
	}
	if req.QuotedPostID != nil {
		quotedPostID = uuid.NullUUID{UUID: *req.QuotedPostID, Valid: true}
	}
	if req.Visibility != "" {
		visibility = sql.NullString{String: req.Visibility, Valid: true}
	} else {
		visibility = sql.NullString{String: "public", Valid: true}
	}

	var media pqtype.NullRawMessage
	if req.Media != nil {
		media = *req.Media
	}

	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}

	post, err := s.store.CreatePost(ctx, db.CreatePostParams{
		AuthorID:     req.AuthorID,
		SpaceID:      req.SpaceID,
		CommunityID:  communityID,
		GroupID:      groupID,
		ParentPostID: parentPostID,
		QuotedPostID: quotedPostID,
		Content:      req.Content,
		Media:        media,
		Tags:         tags,
		Visibility:   visibility,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	response := s.toPostResponse(post)

	
	if s.liveService != nil {
		postPayload := map[string]interface{}{
			"id":         post.ID.String(),
			"author_id":  post.AuthorID.String(),
			"space_id":   post.SpaceID.String(),
			"content":    post.Content,
			"tags":       post.Tags,
			"visibility": response.Visibility,
			"created_at": post.CreatedAt.Time.Unix(),
		}
		if response.Media != nil {
			postPayload["media"] = response.Media
		}
		if req.CommunityID != nil {
			postPayload["community_id"] = req.CommunityID.String()
		}
		if req.GroupID != nil {
			postPayload["group_id"] = req.GroupID.String()
		}

		if err := s.liveService.PublishPostCreated(ctx, req.SpaceID, req.AuthorID, postPayload); err != nil {
			log.Error().Err(err).Msg("Failed to publish post.created event")
		}
	}

	return response, nil
}


func (s *Service) GetPostByID(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (*PostResponse, error) {
	post, err := s.store.GetPostByID(ctx, db.GetPostByIDParams{
		UserID: userID,
		ID:     postID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return s.toDetailedPostResponse(post), nil
}


func (s *Service) DeletePost(ctx context.Context, postID uuid.UUID, authorID uuid.UUID) error {
	err := s.store.DeletePost(ctx, db.DeletePostParams{
		ID:       postID,
		AuthorID: authorID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}


func (s *Service) GetUserPosts(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*PostResponse, error) {
	posts, err := s.store.GetUserPosts(ctx, db.GetUserPostsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}

	return s.toUserPostResponses(posts), nil
}


func (s *Service) GetUserFeed(ctx context.Context, userID uuid.UUID, spaceID uuid.UUID, limit, offset int32) ([]*PostResponse, error) {
	posts, err := s.store.GetUserFeed(ctx, db.GetUserFeedParams{
		UserID:  userID,
		SpaceID: spaceID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user feed: %w", err)
	}

	return s.toUserFeedResponses(posts), nil
}


func (s *Service) GetCommunityPosts(ctx context.Context, userID uuid.UUID, communityID uuid.UUID, limit, offset int32) ([]*PostResponse, error) {
	posts, err := s.store.GetCommunityPosts(ctx, db.GetCommunityPostsParams{
		UserID:      userID,
		CommunityID: uuid.NullUUID{UUID: communityID, Valid: true},
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get community posts: %w", err)
	}

	return s.toCommunityPostResponses(posts), nil
}


func (s *Service) GetGroupPosts(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, limit, offset int32) ([]*PostResponse, error) {
	posts, err := s.store.GetGroupPosts(ctx, db.GetGroupPostsParams{
		UserID:  userID,
		GroupID: uuid.NullUUID{UUID: groupID, Valid: true},
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get group posts: %w", err)
	}

	return s.toGroupPostResponses(posts), nil
}


func (s *Service) GetTrendingPosts(ctx context.Context, spaceID uuid.UUID) ([]*PostResponse, error) {
	posts, err := s.store.GetTrendingPosts(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending posts: %w", err)
	}

	return s.toTrendingPostResponses(posts), nil
}


func (s *Service) GetTrendingTopics(ctx context.Context, spaceID uuid.UUID, limit, offset int32) ([]TrendingTopicResponse, error) {
	
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	topics, err := s.store.GetTrendingTopics(ctx, db.GetTrendingTopicsParams{
		SpaceID: spaceID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get trending topics: %w", err)
	}

	responses := make([]TrendingTopicResponse, len(topics))
	for i, topic := range topics {
		
		name, ok := topic.Name.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert topic name to string")
		}

		responses[i] = TrendingTopicResponse{
			ID:         topic.ID,
			Name:       name,
			Category:   topic.Category,
			PostsCount: topic.PostsCount,
			TrendScore: topic.TrendScore,
		}
	}

	return responses, nil
}


func (s *Service) SearchPosts(ctx context.Context, query string, spaceID uuid.UUID, limit, offset int32) ([]*PostResponse, error) {
	posts, err := s.store.SearchPosts(ctx, db.SearchPostsParams{
		SpaceID:        spaceID,
		PlaintoTsquery: query,
		Content:        "%" + query + "%",
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}

	return s.toSearchPostResponses(posts), nil
}


func (s *Service) AdvancedSearchPosts(ctx context.Context, query string, spaceID uuid.UUID, limit, offset int32) ([]*PostResponse, error) {
	posts, err := s.store.AdvancedSearchPosts(ctx, db.AdvancedSearchPostsParams{
		PlaintoTsquery: query,
		SpaceID:        spaceID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to perform advanced search: %w", err)
	}

	return s.toAdvancedSearchPostResponses(posts), nil
}


func (s *Service) GetUserLikedPosts(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*PostResponse, error) {
	posts, err := s.store.GetUserLikedPosts(ctx, db.GetUserLikedPostsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get liked posts: %w", err)
	}

	return s.toUserLikedPostResponses(posts), nil
}


func (s *Service) GetPostComments(ctx context.Context, postID uuid.UUID) ([]*CommentResponse, error) {
	comments, err := s.store.GetPostComments(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post comments: %w", err)
	}

	return s.toCommentResponses(comments), nil
}


func (s *Service) GetPostLikes(ctx context.Context, postID uuid.UUID) ([]*UserLikeResponse, error) {
	likes, err := s.store.GetPostLikes(ctx, uuid.NullUUID{UUID: postID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get post likes: %w", err)
	}

	return s.toUserLikeResponses(likes), nil
}


func (s *Service) CreateComment(ctx context.Context, req CreateCommentRequest) (*CommentResponse, error) {
	var parentCommentID uuid.NullUUID
	if req.ParentCommentID != nil {
		parentCommentID = uuid.NullUUID{UUID: *req.ParentCommentID, Valid: true}
	}

	comment, err := s.store.CreateComment(ctx, db.CreateCommentParams{
		PostID:          req.PostID,
		AuthorID:        req.AuthorID,
		ParentCommentID: parentCommentID,
		Content:         req.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	response := s.toSimpleCommentResponse(comment)

	
	if s.liveService != nil {
		commentPayload := map[string]interface{}{
			"id":         comment.ID.String(),
			"post_id":    comment.PostID.String(),
			"author_id":  comment.AuthorID.String(),
			"content":    comment.Content,
			"created_at": comment.CreatedAt.Time.Unix(),
		}
		if req.ParentCommentID != nil {
			commentPayload["parent_comment_id"] = req.ParentCommentID.String()
		}

		if err := s.liveService.PublishCommentCreated(ctx, req.PostID, req.AuthorID, commentPayload); err != nil {
			log.Error().Err(err).Msg("Failed to publish comment.created event")
		}
	}

	return response, nil
}


func (s *Service) CreateRepost(ctx context.Context, req CreateRepostParams) (*PostResponse, error) {
	var quotedPostID uuid.NullUUID
	var visibility sql.NullString

	if req.QuotedPostID != nil {
		quotedPostID = uuid.NullUUID{UUID: *req.QuotedPostID, Valid: true}
	}

	if req.Visibility != "" {
		visibility = sql.NullString{String: req.Visibility, Valid: true}
	} else {
		visibility = sql.NullString{String: "public", Valid: true}
	}

	post, err := s.store.CreateRepost(ctx, db.CreateRepostParams{
		AuthorID:     req.AuthorID,
		SpaceID:      req.SpaceID,
		QuotedPostID: quotedPostID,
		Content:      req.Content,
		Visibility:   visibility,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create repost: %w", err)
	}

	return s.toPostResponse(post), nil
}


func (s *Service) TogglePostLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (int32, error) {
	
	postDetail, err := s.store.GetPostByID(ctx, db.GetPostByIDParams{
		UserID: userID,
		ID:     postID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get post details: %w", err)
	}

	likesCount, err := s.store.TogglePostLike(ctx, db.TogglePostLikeParams{
		UserID: userID,
		ID:     postID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to toggle post like: %w", err)
	}

	var count int32
	if likesCount.Valid {
		count = likesCount.Int32
	}

	
	if s.liveService != nil {
		if err := s.liveService.PublishPostLiked(ctx, postID, userID, postDetail.SpaceID, int(count)); err != nil {
			log.Error().Err(err).Msg("Failed to publish post.liked event")
		}
	}

	return count, nil
}


func (s *Service) ToggleCommentLike(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (bool, error) {
	liked, err := s.store.ToggleCommentLike(ctx, db.ToggleCommentLikeParams{
		UserID:    userID,
		CommentID: uuid.NullUUID{UUID: commentID, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("failed to toggle comment like: %w", err)
	}

	return liked, nil
}


func (s *Service) PinPost(ctx context.Context, postID uuid.UUID, isPinned bool) error {
	err := s.store.PinPost(ctx, db.PinPostParams{
		IsPinned: sql.NullBool{Bool: isPinned, Valid: true},
		ID:       postID,
	})
	if err != nil {
		return fmt.Errorf("failed to pin/unpin post: %w", err)
	}
	return nil
}


func (s *Service) IncrementPostViews(ctx context.Context, postID uuid.UUID) error {
	err := s.store.IncrementPostViews(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to increment post views: %w", err)
	}
	return nil
}



func (s *Service) toPostResponse(post db.Post) *PostResponse {
	resp := &PostResponse{
		ID:            post.ID,
		AuthorID:      post.AuthorID,
		SpaceID:       post.SpaceID,
		Content:       post.Content,
		Tags:          post.Tags,
		LikesCount:    post.LikesCount.Int32,
		CommentsCount: post.CommentsCount.Int32,
		RepostsCount:  post.RepostsCount.Int32,
		QuotesCount:   post.QuotesCount.Int32,
		ViewsCount:    post.ViewsCount.Int32,
		IsPinned:      post.IsPinned.Bool,
		Visibility:    post.Visibility.String,
		Status:        post.Status.String,
	}

	if post.CommunityID.Valid {
		resp.CommunityID = &post.CommunityID.UUID
	}
	if post.GroupID.Valid {
		resp.GroupID = &post.GroupID.UUID
	}
	if post.ParentPostID.Valid {
		resp.ParentPostID = &post.ParentPostID.UUID
	}
	if post.QuotedPostID.Valid {
		resp.QuotedPostID = &post.QuotedPostID.UUID
	}
	if post.CreatedAt.Valid {
		resp.CreatedAt = &post.CreatedAt.Time
	}
	if post.UpdatedAt.Valid {
		resp.UpdatedAt = &post.UpdatedAt.Time
	}

	return resp
}

func (s *Service) toDetailedPostResponse(post db.GetPostByIDRow) *PostResponse {
	resp := &PostResponse{
		ID:            post.ID,
		AuthorID:      post.AuthorID,
		SpaceID:       post.SpaceID,
		Content:       post.Content,
		Tags:          post.Tags,
		LikesCount:    post.LikesCount.Int32,
		CommentsCount: post.CommentsCount.Int32,
		RepostsCount:  post.RepostsCount.Int32,
		QuotesCount:   post.QuotesCount.Int32,
		ViewsCount:    post.ViewsCount.Int32,
		IsPinned:      post.IsPinned.Bool,
		Visibility:    post.Visibility.String,
		Status:        post.Status.String,
	}

	if post.CommunityID.Valid {
		resp.CommunityID = &post.CommunityID.UUID
	}
	if post.GroupID.Valid {
		resp.GroupID = &post.GroupID.UUID
	}
	if post.ParentPostID.Valid {
		resp.ParentPostID = &post.ParentPostID.UUID
	}
	if post.QuotedPostID.Valid {
		resp.QuotedPostID = &post.QuotedPostID.UUID
	}
	if post.CreatedAt.Valid {
		resp.CreatedAt = &post.CreatedAt.Time
	}
	if post.UpdatedAt.Valid {
		resp.UpdatedAt = &post.UpdatedAt.Time
	}

	resp.Username = interfaceToStringPtr(post.Username)
	resp.FullName = interfaceToStringPtr(post.FullName)
	if post.AuthorAvatar.Valid {
		resp.AuthorAvatar = &post.AuthorAvatar.String
	}
	if post.AuthorVerified.Valid {
		resp.AuthorVerified = &post.AuthorVerified.Bool
	}
	if post.CommunityName.Valid {
		resp.CommunityName = &post.CommunityName.String
	}
	if post.GroupName.Valid {
		resp.GroupName = &post.GroupName.String
	}
	if post.QuotedContent.Valid {
		resp.QuotedContent = &post.QuotedContent.String
	}
	if post.QuotedAuthorID.Valid {
		resp.QuotedAuthorID = &post.QuotedAuthorID.UUID
	}
	resp.QuotedUsername = interfaceToStringPtr(post.QuotedUsername)
	resp.QuotedFullName = interfaceToStringPtr(post.QuotedFullName)

	resp.IsLiked = &post.IsLiked

	return resp
}

func (s *Service) toUserPostResponses(posts []db.GetUserPostsRow) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		resp := &PostResponse{
			ID:            post.ID,
			AuthorID:      post.AuthorID,
			SpaceID:       post.SpaceID,
			Content:       post.Content,
			Tags:          post.Tags,
			LikesCount:    post.LikesCount.Int32,
			CommentsCount: post.CommentsCount.Int32,
			RepostsCount:  post.RepostsCount.Int32,
			QuotesCount:   post.QuotesCount.Int32,
			ViewsCount:    post.ViewsCount.Int32,
			IsPinned:      post.IsPinned.Bool,
			Visibility:    post.Visibility.String,
			Status:        post.Status.String,
			Username:      interfaceToStringPtr(post.Username),
			FullName:      interfaceToStringPtr(post.FullName),
			IsLiked:       &post.IsLiked,
		}

		if post.CommunityID.Valid {
			resp.CommunityID = &post.CommunityID.UUID
		}
		if post.GroupID.Valid {
			resp.GroupID = &post.GroupID.UUID
		}
		if post.ParentPostID.Valid {
			resp.ParentPostID = &post.ParentPostID.UUID
		}
		if post.QuotedPostID.Valid {
			resp.QuotedPostID = &post.QuotedPostID.UUID
		}
		if post.CreatedAt.Valid {
			resp.CreatedAt = &post.CreatedAt.Time
		}
		if post.UpdatedAt.Valid {
			resp.UpdatedAt = &post.UpdatedAt.Time
		}
		if post.AuthorAvatar.Valid {
			resp.AuthorAvatar = &post.AuthorAvatar.String
		}
		if post.CommunityName.Valid {
			resp.CommunityName = &post.CommunityName.String
		}
		if post.GroupName.Valid {
			resp.GroupName = &post.GroupName.String
		}

		responses[i] = resp
	}
	return responses
}

func (s *Service) toUserFeedResponses(posts []db.GetUserFeedRow) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		resp := &PostResponse{
			ID:             post.ID,
			AuthorID:       post.AuthorID,
			SpaceID:        post.SpaceID,
			Content:        post.Content,
			Tags:           post.Tags,
			LikesCount:     post.LikesCount.Int32,
			CommentsCount:  post.CommentsCount.Int32,
			RepostsCount:   post.RepostsCount.Int32,
			QuotesCount:    post.QuotesCount.Int32,
			ViewsCount:     post.ViewsCount.Int32,
			IsPinned:       post.IsPinned.Bool,
			Visibility:     post.Visibility.String,
			Status:         post.Status.String,
			Username:       interfaceToStringPtr(post.Username),
			FullName:       interfaceToStringPtr(post.FullName),
			IsLiked:        &post.IsLiked,
			AuthorVerified: &post.AuthorVerified.Bool,
		}

		if post.CommunityID.Valid {
			resp.CommunityID = &post.CommunityID.UUID
		}
		if post.GroupID.Valid {
			resp.GroupID = &post.GroupID.UUID
		}
		if post.ParentPostID.Valid {
			resp.ParentPostID = &post.ParentPostID.UUID
		}
		if post.QuotedPostID.Valid {
			resp.QuotedPostID = &post.QuotedPostID.UUID
		}
		if post.CreatedAt.Valid {
			resp.CreatedAt = &post.CreatedAt.Time
		}
		if post.UpdatedAt.Valid {
			resp.UpdatedAt = &post.UpdatedAt.Time
		}
		if post.AuthorAvatar.Valid {
			resp.AuthorAvatar = &post.AuthorAvatar.String
		}
		if post.CommunityName.Valid {
			resp.CommunityName = &post.CommunityName.String
		}
		if post.GroupName.Valid {
			resp.GroupName = &post.GroupName.String
		}

		responses[i] = resp
	}
	return responses
}

func (s *Service) toCommunityPostResponses(posts []db.GetCommunityPostsRow) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		resp := &PostResponse{
			ID:            post.ID,
			AuthorID:      post.AuthorID,
			SpaceID:       post.SpaceID,
			Content:       post.Content,
			Tags:          post.Tags,
			LikesCount:    post.LikesCount.Int32,
			CommentsCount: post.CommentsCount.Int32,
			RepostsCount:  post.RepostsCount.Int32,
			QuotesCount:   post.QuotesCount.Int32,
			ViewsCount:    post.ViewsCount.Int32,
			IsPinned:      post.IsPinned.Bool,
			Visibility:    post.Visibility.String,
			Status:        post.Status.String,
			Username:      interfaceToStringPtr(post.Username),
			FullName:      interfaceToStringPtr(post.FullName),
			IsLiked:       &post.IsLiked,
		}

		if post.CommunityID.Valid {
			resp.CommunityID = &post.CommunityID.UUID
		}
		if post.GroupID.Valid {
			resp.GroupID = &post.GroupID.UUID
		}
		if post.ParentPostID.Valid {
			resp.ParentPostID = &post.ParentPostID.UUID
		}
		if post.QuotedPostID.Valid {
			resp.QuotedPostID = &post.QuotedPostID.UUID
		}
		if post.CreatedAt.Valid {
			resp.CreatedAt = &post.CreatedAt.Time
		}
		if post.UpdatedAt.Valid {
			resp.UpdatedAt = &post.UpdatedAt.Time
		}
		if post.AuthorAvatar.Valid {
			resp.AuthorAvatar = &post.AuthorAvatar.String
		}

		responses[i] = resp
	}
	return responses
}

func (s *Service) toGroupPostResponses(posts []db.GetGroupPostsRow) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		resp := &PostResponse{
			ID:            post.ID,
			AuthorID:      post.AuthorID,
			SpaceID:       post.SpaceID,
			Content:       post.Content,
			Tags:          post.Tags,
			LikesCount:    post.LikesCount.Int32,
			CommentsCount: post.CommentsCount.Int32,
			RepostsCount:  post.RepostsCount.Int32,
			QuotesCount:   post.QuotesCount.Int32,
			ViewsCount:    post.ViewsCount.Int32,
			IsPinned:      post.IsPinned.Bool,
			Visibility:    post.Visibility.String,
			Status:        post.Status.String,
			Username:      interfaceToStringPtr(post.Username),
			FullName:      interfaceToStringPtr(post.FullName),
			IsLiked:       &post.IsLiked,
		}

		if post.CommunityID.Valid {
			resp.CommunityID = &post.CommunityID.UUID
		}
		if post.GroupID.Valid {
			resp.GroupID = &post.GroupID.UUID
		}
		if post.ParentPostID.Valid {
			resp.ParentPostID = &post.ParentPostID.UUID
		}
		if post.QuotedPostID.Valid {
			resp.QuotedPostID = &post.QuotedPostID.UUID
		}
		if post.CreatedAt.Valid {
			resp.CreatedAt = &post.CreatedAt.Time
		}
		if post.UpdatedAt.Valid {
			resp.UpdatedAt = &post.UpdatedAt.Time
		}
		if post.AuthorAvatar.Valid {
			resp.AuthorAvatar = &post.AuthorAvatar.String
		}

		responses[i] = resp
	}
	return responses
}

func (s *Service) toTrendingPostResponses(posts []db.GetTrendingPostsRow) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		resp := &PostResponse{
			ID:              post.ID,
			AuthorID:        post.AuthorID,
			SpaceID:         post.SpaceID,
			Content:         post.Content,
			Tags:            post.Tags,
			LikesCount:      post.LikesCount.Int32,
			CommentsCount:   post.CommentsCount.Int32,
			RepostsCount:    post.RepostsCount.Int32,
			QuotesCount:     post.QuotesCount.Int32,
			ViewsCount:      post.ViewsCount.Int32,
			IsPinned:        post.IsPinned.Bool,
			Visibility:      post.Visibility.String,
			Status:          post.Status.String,
			Username:        interfaceToStringPtr(post.Username),
			FullName:        interfaceToStringPtr(post.FullName),
			EngagementScore: &post.EngagementScore,
		}

		if post.CommunityID.Valid {
			resp.CommunityID = &post.CommunityID.UUID
		}
		if post.GroupID.Valid {
			resp.GroupID = &post.GroupID.UUID
		}
		if post.ParentPostID.Valid {
			resp.ParentPostID = &post.ParentPostID.UUID
		}
		if post.QuotedPostID.Valid {
			resp.QuotedPostID = &post.QuotedPostID.UUID
		}
		if post.CreatedAt.Valid {
			resp.CreatedAt = &post.CreatedAt.Time
		}
		if post.UpdatedAt.Valid {
			resp.UpdatedAt = &post.UpdatedAt.Time
		}
		if post.AuthorAvatar.Valid {
			resp.AuthorAvatar = &post.AuthorAvatar.String
		}

		responses[i] = resp
	}
	return responses
}

func (s *Service) toSearchPostResponses(posts []db.SearchPostsRow) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		resp := &PostResponse{
			ID:            post.ID,
			AuthorID:      post.AuthorID,
			SpaceID:       post.SpaceID,
			Content:       post.Content,
			Tags:          post.Tags,
			LikesCount:    post.LikesCount.Int32,
			CommentsCount: post.CommentsCount.Int32,
			RepostsCount:  post.RepostsCount.Int32,
			QuotesCount:   post.QuotesCount.Int32,
			ViewsCount:    post.ViewsCount.Int32,
			IsPinned:      post.IsPinned.Bool,
			Visibility:    post.Visibility.String,
			Status:        post.Status.String,
			Username:      interfaceToStringPtr(post.Username),
			FullName:      interfaceToStringPtr(post.FullName),
		}

		rank := post.Rank
		resp.RelevanceScore = &rank

		if post.CommunityID.Valid {
			resp.CommunityID = &post.CommunityID.UUID
		}
		if post.GroupID.Valid {
			resp.GroupID = &post.GroupID.UUID
		}
		if post.ParentPostID.Valid {
			resp.ParentPostID = &post.ParentPostID.UUID
		}
		if post.QuotedPostID.Valid {
			resp.QuotedPostID = &post.QuotedPostID.UUID
		}
		if post.CreatedAt.Valid {
			resp.CreatedAt = &post.CreatedAt.Time
		}
		if post.UpdatedAt.Valid {
			resp.UpdatedAt = &post.UpdatedAt.Time
		}
		if post.AuthorAvatar.Valid {
			resp.AuthorAvatar = &post.AuthorAvatar.String
		}
		if post.CommunityName.Valid {
			resp.CommunityName = &post.CommunityName.String
		}
		if post.GroupName.Valid {
			resp.GroupName = &post.GroupName.String
		}

		responses[i] = resp
	}
	return responses
}

func (s *Service) toAdvancedSearchPostResponses(posts []db.AdvancedSearchPostsRow) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		resp := &PostResponse{
			ID:             post.ID,
			AuthorID:       post.AuthorID,
			SpaceID:        post.SpaceID,
			Content:        post.Content,
			Tags:           post.Tags,
			LikesCount:     post.LikesCount.Int32,
			CommentsCount:  post.CommentsCount.Int32,
			RepostsCount:   post.RepostsCount.Int32,
			QuotesCount:    post.QuotesCount.Int32,
			ViewsCount:     post.ViewsCount.Int32,
			IsPinned:       post.IsPinned.Bool,
			Visibility:     post.Visibility.String,
			Status:         post.Status.String,
			Username:       interfaceToStringPtr(post.Username),
			FullName:       interfaceToStringPtr(post.FullName),
			RelevanceScore: &post.RelevanceScore,
		}

		if post.CommunityID.Valid {
			resp.CommunityID = &post.CommunityID.UUID
		}
		if post.GroupID.Valid {
			resp.GroupID = &post.GroupID.UUID
		}
		if post.ParentPostID.Valid {
			resp.ParentPostID = &post.ParentPostID.UUID
		}
		if post.QuotedPostID.Valid {
			resp.QuotedPostID = &post.QuotedPostID.UUID
		}
		if post.CreatedAt.Valid {
			resp.CreatedAt = &post.CreatedAt.Time
		}
		if post.UpdatedAt.Valid {
			resp.UpdatedAt = &post.UpdatedAt.Time
		}
		if post.AuthorAvatar.Valid {
			resp.AuthorAvatar = &post.AuthorAvatar.String
		}
		if post.CommunityName.Valid {
			resp.CommunityName = &post.CommunityName.String
		}
		if post.GroupName.Valid {
			resp.GroupName = &post.GroupName.String
		}

		responses[i] = resp
	}
	return responses
}

func (s *Service) toUserLikedPostResponses(posts []db.GetUserLikedPostsRow) []*PostResponse {
	responses := make([]*PostResponse, len(posts))
	for i, post := range posts {
		resp := &PostResponse{
			ID:            post.ID,
			AuthorID:      post.AuthorID,
			SpaceID:       post.SpaceID,
			Content:       post.Content,
			Tags:          post.Tags,
			LikesCount:    post.LikesCount.Int32,
			CommentsCount: post.CommentsCount.Int32,
			RepostsCount:  post.RepostsCount.Int32,
			QuotesCount:   post.QuotesCount.Int32,
			ViewsCount:    post.ViewsCount.Int32,
			IsPinned:      post.IsPinned.Bool,
			Visibility:    post.Visibility.String,
			Status:        post.Status.String,
			Username:      interfaceToStringPtr(post.Username),
			FullName:      interfaceToStringPtr(post.FullName),
		}

		if post.CommunityID.Valid {
			resp.CommunityID = &post.CommunityID.UUID
		}
		if post.GroupID.Valid {
			resp.GroupID = &post.GroupID.UUID
		}
		if post.ParentPostID.Valid {
			resp.ParentPostID = &post.ParentPostID.UUID
		}
		if post.QuotedPostID.Valid {
			resp.QuotedPostID = &post.QuotedPostID.UUID
		}
		if post.CreatedAt.Valid {
			resp.CreatedAt = &post.CreatedAt.Time
		}
		if post.UpdatedAt.Valid {
			resp.UpdatedAt = &post.UpdatedAt.Time
		}
		if post.AuthorAvatar.Valid {
			resp.AuthorAvatar = &post.AuthorAvatar.String
		}

		responses[i] = resp
	}
	return responses
}

func (s *Service) toCommentResponses(comments []db.GetPostCommentsRow) []*CommentResponse {
	responses := make([]*CommentResponse, len(comments))
	for i, comment := range comments {
		resp := &CommentResponse{
			ID:         comment.ID,
			PostID:     comment.PostID,
			AuthorID:   comment.AuthorID,
			Content:    comment.Content,
			LikesCount: comment.LikesCount.Int32,
			Status:     comment.Status.String,
			Username:   &comment.Username,
			FullName:   &comment.FullName,
		}

		if comment.ParentCommentID.Valid {
			resp.ParentCommentID = &comment.ParentCommentID.UUID
		}
		if comment.CreatedAt.Valid {
			resp.CreatedAt = &comment.CreatedAt.Time
		}
		if comment.UpdatedAt.Valid {
			resp.UpdatedAt = &comment.UpdatedAt.Time
		}
		if comment.Avatar.Valid {
			resp.Avatar = &comment.Avatar.String
		}
		resp.Depth = &comment.Depth

		responses[i] = resp
	}
	return responses
}

func (s *Service) toSimpleCommentResponse(comment db.Comment) *CommentResponse {
	resp := &CommentResponse{
		ID:         comment.ID,
		PostID:     comment.PostID,
		AuthorID:   comment.AuthorID,
		Content:    comment.Content,
		LikesCount: comment.LikesCount.Int32,
		Status:     comment.Status.String,
	}

	if comment.ParentCommentID.Valid {
		resp.ParentCommentID = &comment.ParentCommentID.UUID
	}
	if comment.CreatedAt.Valid {
		resp.CreatedAt = &comment.CreatedAt.Time
	}
	if comment.UpdatedAt.Valid {
		resp.UpdatedAt = &comment.UpdatedAt.Time
	}

	return resp
}

func (s *Service) toUserLikeResponses(likes []db.GetPostLikesRow) []*UserLikeResponse {
	responses := make([]*UserLikeResponse, len(likes))
	for i, like := range likes {
		resp := &UserLikeResponse{
			ID:       like.ID,
			Username: like.Username,
			FullName: like.FullName,
		}

		if like.Avatar.Valid {
			resp.Avatar = &like.Avatar.String
		}
		if like.CreatedAt.Valid {
			resp.CreatedAt = &like.CreatedAt.Time
		}

		responses[i] = resp
	}
	return responses
}
