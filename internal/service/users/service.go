package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/google/uuid"
)


type Service struct {
	store db.Store
}


func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}


func interfaceToString(val interface{}, fallback string) string {
	if val == nil {
		return fallback
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fallback
}


func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	
	if err := ValidateEmail(req.Email); err != nil {
		return nil, fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}
	if err := ValidateUsername(req.Username); err != nil {
		return nil, fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}
	if err := ValidatePassword(req.Password); err != nil {
		return nil, fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}
	if err := ValidateFullName(req.FullName); err != nil {
		return nil, fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}

	
	existingUser, err := s.store.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser.ID != uuid.Nil {
		return nil, fmt.Errorf("%w: user with this email already exists", util.ErrConflict)
	}

	
	existingByUsername, err := s.store.GetUserByUsername(ctx, req.Username)
	if err == nil && existingByUsername.ID != uuid.Nil {
		return nil, fmt.Errorf("%w: username already taken", util.ErrConflict)
	}

	
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	
	roles := []string{"user"}
	var levelSQL sql.NullString
	if req.Level != nil {
		levelSQL = sql.NullString{String: *req.Level, Valid: true}
	}
	var deptSQL sql.NullString
	if req.Department != nil {
		deptSQL = sql.NullString{String: *req.Department, Valid: true}
	}
	var majorSQL sql.NullString
	if req.Major != nil {
		majorSQL = sql.NullString{String: *req.Major, Valid: true}
	}
	var yearSQL sql.NullInt32
	if req.Year != nil {
		yearSQL = sql.NullInt32{Int32: *req.Year, Valid: true}
	}

	interests := req.Interests
	if interests == nil {
		interests = []string{}
	}

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		SpaceID:    req.SpaceID,
		Username:   req.Username,
		Email:      strings.ToLower(req.Email),
		Password:   hashedPassword,
		FullName:   req.FullName,
		Roles:      roles,
		Level:      levelSQL,
		Department: deptSQL,
		Major:      majorSQL,
		Year:       yearSQL,
		Interests:  interests,
	})
	if err != nil {
		if util.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("%w: user already exists", util.ErrConflict)
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.toUserResponse(user), nil
}


func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found", util.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &UserResponse{
		ID:             user.ID,
		SpaceID:        user.SpaceID,
		Username:       user.Username,
		Email:          user.Email,
		FullName:       user.FullName,
		Avatar:         nullStringToPtr(user.Avatar),
		Bio:            nullStringToPtr(user.Bio),
		Verified:       nullBoolToPtr(user.Verified),
		Roles:          user.Roles,
		Level:          nullStringToPtr(user.Level),
		Department:     nullStringToPtr(user.Department),
		Major:          nullStringToPtr(user.Major),
		Year:           nullInt32ToPtr(user.Year),
		Interests:      user.Interests,
		FollowersCount: nullInt32ToPtr(user.FollowersCount),
		FollowingCount: nullInt32ToPtr(user.FollowingCount),
		MentorStatus:   nullStringToPtr(user.MentorStatus),
		TutorStatus:    nullStringToPtr(user.TutorStatus),
		Status:         nullStringToPtr(user.Status),
		CreatedAt:      nullTimeToPtr(user.CreatedAt),
		UpdatedAt:      nullTimeToPtr(user.UpdatedAt),
		SpaceName:      &user.SpaceName,
		SpaceSlug:      &user.SpaceSlug,
	}, nil
}


func (s *Service) GetUserByUsername(ctx context.Context, username string) (*UserResponse, error) {
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found", util.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return s.toUserResponse(user), nil
}


func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*UserResponse, error) {
	
	if err := ValidateFullName(req.FullName); err != nil {
		return nil, fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}

	var bioSQL sql.NullString
	if req.Bio != nil {
		bioSQL = sql.NullString{String: *req.Bio, Valid: true}
	}
	var avatarSQL sql.NullString
	if req.Avatar != nil {
		avatarSQL = sql.NullString{String: *req.Avatar, Valid: true}
	}
	var levelSQL sql.NullString
	if req.Level != nil {
		levelSQL = sql.NullString{String: *req.Level, Valid: true}
	}
	var deptSQL sql.NullString
	if req.Department != nil {
		deptSQL = sql.NullString{String: *req.Department, Valid: true}
	}
	var majorSQL sql.NullString
	if req.Major != nil {
		majorSQL = sql.NullString{String: *req.Major, Valid: true}
	}
	var yearSQL sql.NullInt32
	if req.Year != nil {
		yearSQL = sql.NullInt32{Int32: *req.Year, Valid: true}
	}

	interests := req.Interests
	if interests == nil {
		interests = []string{}
	}

	user, err := s.store.UpdateUser(ctx, db.UpdateUserParams{
		FullName:   req.FullName,
		Bio:        bioSQL,
		Avatar:     avatarSQL,
		Level:      levelSQL,
		Department: deptSQL,
		Major:      majorSQL,
		Year:       yearSQL,
		Interests:  interests,
		ID:         userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found", util.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.toUserResponse(user), nil
}


func (s *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, req UpdatePasswordRequest) error {
	
	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user not found", util.ErrNotFound)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	
	if err := auth.CheckPassword(req.OldPassword, user.Password); err != nil {
		return fmt.Errorf("%w: incorrect password", util.ErrUnauthorized)
	}

	
	if err := ValidatePassword(req.NewPassword); err != nil {
		return fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}

	
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	
	err = s.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		Password: hashedPassword,
		ID:       userID,
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}


func (s *Service) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	err := s.store.DeactivateUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	return nil
}


func (s *Service) SearchUsers(ctx context.Context, query string, spaceID uuid.UUID) ([]UserResponse, error) {
	users, err := s.store.SearchUsers(ctx, db.SearchUsersParams{
		
		SpaceID: spaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	result := make([]UserResponse, len(users))
	for i, user := range users {
		result[i] = UserResponse{
			ID:             user.ID,
			Username:       user.Username,
			FullName:       user.FullName,
			Avatar:         nullStringToPtr(user.Avatar),
			Bio:            nullStringToPtr(user.Bio),
			Level:          nullStringToPtr(user.Level),
			Department:     nullStringToPtr(user.Department),
			Major:          nullStringToPtr(user.Major),
			Verified:       nullBoolToPtr(user.Verified),
			FollowersCount: nullInt32ToPtr(user.FollowersCount),
			FollowingCount: nullInt32ToPtr(user.FollowingCount),
		}
	}

	return result, nil
}


func (s *Service) GetSuggestedUsers(ctx context.Context, userID uuid.UUID, spaceID uuid.UUID, limit, offset int32) ([]UserResponse, error) {
	
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.store.GetSuggestedUsers(ctx, db.GetSuggestedUsersParams{
		FollowerID: userID,
		SpaceID:    spaceID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested users: %w", err)
	}

	result := make([]UserResponse, len(users))
	for i, user := range users {
		result[i] = UserResponse{
			ID:             user.ID,
			Username:       interfaceToString(user.Username, "unknown_user"),
			FullName:       interfaceToString(user.FullName, "Unknown User"),
			Avatar:         nullStringToPtr(user.Avatar),
			Bio:            nullStringToPtr(user.Bio),
			Level:          nullStringToPtr(user.Level),
			Department:     nullStringToPtr(user.Department),
			Verified:       nullBoolToPtr(user.Verified),
			FollowersCount: nullInt32ToPtr(user.FollowersCount),
			FollowingCount: nullInt32ToPtr(user.FollowingCount),
			
			
		}
	}

	return result, nil
}


func (s *Service) Authenticate(ctx context.Context, email, password string) (*UserResponse, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: invalid credentials", util.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	
	if err := auth.CheckPassword(password, user.Password); err != nil {
		return nil, fmt.Errorf("%w: invalid credentials", util.ErrInvalidCredentials)
	}

	return s.toUserResponse(user), nil
}


func (s *Service) FollowUser(ctx context.Context, followerID, followingID, spaceID uuid.UUID) error {
	
	if followerID == followingID {
		return fmt.Errorf("%w: cannot follow yourself", util.ErrBadRequest)
	}

	
	_, err := s.store.GetUserByID(ctx, followingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: user not found", util.ErrNotFound)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	
	_, err = s.store.FollowUser(ctx, db.FollowUserParams{
		FollowerID:  followerID,
		FollowingID: followingID,
		SpaceID:     spaceID,
	})
	if err != nil {
		
		if strings.Contains(err.Error(), "no rows") {
			return nil 
		}
		return fmt.Errorf("failed to follow user: %w", err)
	}

	
	if err := s.store.IncrementFollowingCount(ctx, followerID); err != nil {
		return fmt.Errorf("failed to update following count: %w", err)
	}
	if err := s.store.IncrementFollowersCount(ctx, followingID); err != nil {
		return fmt.Errorf("failed to update followers count: %w", err)
	}

	return nil
}


func (s *Service) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	
	if followerID == followingID {
		return fmt.Errorf("%w: cannot unfollow yourself", util.ErrBadRequest)
	}

	
	err := s.store.UnfollowUser(ctx, db.UnfollowUserParams{
		FollowerID:  followerID,
		FollowingID: followingID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}

	
	if err := s.store.DecrementFollowingCount(ctx, followerID); err != nil {
		return fmt.Errorf("failed to update following count: %w", err)
	}
	if err := s.store.DecrementFollowersCount(ctx, followingID); err != nil {
		return fmt.Errorf("failed to update followers count: %w", err)
	}

	return nil
}


func (s *Service) CheckIfFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	isFollowing, err := s.store.CheckIfFollowing(ctx, db.CheckIfFollowingParams{
		FollowerID:  followerID,
		FollowingID: followingID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check follow status: %w", err)
	}
	return isFollowing, nil
}


func (s *Service) GetFollowers(ctx context.Context, userID uuid.UUID, page, limit int32) ([]UserFollowResponse, error) {
	offset := (page - 1) * limit

	followers, err := s.store.GetUserFollowers(ctx, db.GetUserFollowersParams{
		FollowingID: userID,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	response := make([]UserFollowResponse, len(followers))
	for i, follower := range followers {
		response[i] = UserFollowResponse{
			ID:             follower.ID,
			Username:       follower.Username,
			FullName:       follower.FullName,
			Avatar:         nullStringToPtr(follower.Avatar),
			Bio:            nullStringToPtr(follower.Bio),
			Verified:       nullBoolToPtr(follower.Verified),
			FollowersCount: nullInt32ToPtr(follower.FollowersCount),
			FollowingCount: nullInt32ToPtr(follower.FollowingCount),
			FollowedAt:     nullTimeToPtr(follower.FollowedAt),
		}
	}

	return response, nil
}


func (s *Service) GetFollowing(ctx context.Context, userID uuid.UUID, page, limit int32) ([]UserFollowResponse, error) {
	offset := (page - 1) * limit

	following, err := s.store.GetUserFollowing(ctx, db.GetUserFollowingParams{
		FollowerID: userID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	response := make([]UserFollowResponse, len(following))
	for i, user := range following {
		response[i] = UserFollowResponse{
			ID:             user.ID,
			Username:       user.Username,
			FullName:       user.FullName,
			Avatar:         nullStringToPtr(user.Avatar),
			Bio:            nullStringToPtr(user.Bio),
			Verified:       nullBoolToPtr(user.Verified),
			FollowersCount: nullInt32ToPtr(user.FollowersCount),
			FollowingCount: nullInt32ToPtr(user.FollowingCount),
			FollowedAt:     nullTimeToPtr(user.FollowedAt),
		}
	}

	return response, nil
}


func (s *Service) toUserResponse(user db.User) *UserResponse {
	return &UserResponse{
		ID:             user.ID,
		SpaceID:        user.SpaceID,
		Username:       user.Username,
		Email:          user.Email,
		FullName:       user.FullName,
		Avatar:         nullStringToPtr(user.Avatar),
		Bio:            nullStringToPtr(user.Bio),
		Verified:       nullBoolToPtr(user.Verified),
		Roles:          user.Roles,
		Level:          nullStringToPtr(user.Level),
		Department:     nullStringToPtr(user.Department),
		Major:          nullStringToPtr(user.Major),
		Year:           nullInt32ToPtr(user.Year),
		Interests:      user.Interests,
		FollowersCount: nullInt32ToPtr(user.FollowersCount),
		FollowingCount: nullInt32ToPtr(user.FollowingCount),
		MentorStatus:   nullStringToPtr(user.MentorStatus),
		TutorStatus:    nullStringToPtr(user.TutorStatus),
		Status:         nullStringToPtr(user.Status),
		CreatedAt:      nullTimeToPtr(user.CreatedAt),
		UpdatedAt:      nullTimeToPtr(user.UpdatedAt),
	}
}

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullBoolToPtr(nb sql.NullBool) *bool {
	if nb.Valid {
		return &nb.Bool
	}
	return nil
}

func nullInt32ToPtr(ni sql.NullInt32) *int32 {
	if ni.Valid {
		return &ni.Int32
	}
	return nil
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
