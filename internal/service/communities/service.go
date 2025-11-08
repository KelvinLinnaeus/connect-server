package communities

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// Service handles community business logic
type Service struct {
	store db.Store
}

// NewService creates a new community service
func NewService(store db.Store) *Service {
	return &Service{store: store}
}

// CreateCommunity creates a new community
func (s *Service) CreateCommunity(ctx context.Context, req CreateCommunityRequest) (*CommunityResponse, error) {
	var description, coverImage sql.NullString
	var isPublic sql.NullBool
	var settings pqtype.NullRawMessage

	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.CoverImage != nil {
		coverImage = sql.NullString{String: *req.CoverImage, Valid: true}
	}
	if req.IsPublic != nil {
		isPublic = sql.NullBool{Bool: *req.IsPublic, Valid: true}
	}
	if req.Settings != nil {
		settings = *req.Settings
	}

	createdBy := uuid.NullUUID{UUID: req.CreatedBy, Valid: true}

	community, err := s.store.CreateCommunity(ctx, db.CreateCommunityParams{
		SpaceID:     req.SpaceID,
		Name:        req.Name,
		Description: description,
		Category:    req.Category,
		CoverImage:  coverImage,
		IsPublic:    isPublic,
		CreatedBy:   createdBy,
		Settings:    settings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}

	return s.toCommunityResponse(community), nil
}

// GetCommunityByID gets a community by ID
func (s *Service) GetCommunityByID(ctx context.Context, userID, communityID uuid.UUID) (*CommunityDetailResponse, error) {
	community, err := s.store.GetCommunityByID(ctx, db.GetCommunityByIDParams{
		UserID: userID,
		ID:     communityID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("community not found")
		}
		return nil, fmt.Errorf("failed to get community: %w", err)
	}

	return s.toCommunityDetailResponse(community), nil
}

// GetCommunityBySlug gets a community by slug (name)
func (s *Service) GetCommunityBySlug(ctx context.Context, userID, spaceID uuid.UUID, slug string) (*CommunityDetailResponse, error) {
	community, err := s.store.GetCommunityBySlug(ctx, db.GetCommunityBySlugParams{
		UserID:  userID,
		SpaceID: spaceID,
		Lower:   slug,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("community not found")
		}
		return nil, fmt.Errorf("failed to get community: %w", err)
	}

	return s.toCommunityDetailFromSlugResponse(community), nil
}

// ListCommunities lists communities with pagination and sorting
func (s *Service) ListCommunities(ctx context.Context, params ListCommunitiesParams) ([]CommunityListResponse, error) {
	offset := (params.Page - 1) * params.Limit

	communities, err := s.store.ListCommunities(ctx, db.ListCommunitiesParams{
		UserID:  params.UserID,
		SpaceID: params.SpaceID,
		Column3: params.SortBy,
		Limit:   params.Limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list communities: %w", err)
	}

	return s.toCommunityListResponses(communities), nil
}

// SearchCommunities searches communities by name, description, or category
func (s *Service) SearchCommunities(ctx context.Context, params SearchCommunitiesParams) ([]CommunityListResponse, error) {
	communities, err := s.store.SearchCommunities(ctx, db.SearchCommunitiesParams{
		UserID:  params.UserID,
		SpaceID: params.SpaceID,
		Name:    "%" + params.Query + "%",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search communities: %w", err)
	}

	return s.toSearchCommunityListResponses(communities), nil
}

// UpdateCommunity updates a community
func (s *Service) UpdateCommunity(ctx context.Context, communityID uuid.UUID, req UpdateCommunityRequest) (*CommunityResponse, error) {
	var description, coverImage sql.NullString
	var isPublic sql.NullBool
	var settings pqtype.NullRawMessage

	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.CoverImage != nil {
		coverImage = sql.NullString{String: *req.CoverImage, Valid: true}
	}
	if req.IsPublic != nil {
		isPublic = sql.NullBool{Bool: *req.IsPublic, Valid: true}
	}
	if req.Settings != nil {
		settings = *req.Settings
	}

	community, err := s.store.UpdateCommunity(ctx, db.UpdateCommunityParams{
		ID:          communityID,
		Name:        req.Name,
		Description: description,
		CoverImage:  coverImage,
		Category:    req.Category,
		IsPublic:    isPublic,
		Settings:    settings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update community: %w", err)
	}

	return s.toCommunityResponse(community), nil
}

// GetCommunityMembers gets all members of a community
func (s *Service) GetCommunityMembers(ctx context.Context, communityID uuid.UUID) ([]CommunityMemberResponse, error) {
	members, err := s.store.GetCommunityMembers(ctx, communityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community members: %w", err)
	}

	return s.toCommunityMemberResponses(members), nil
}

// GetCommunityModerators gets all moderators of a community
func (s *Service) GetCommunityModerators(ctx context.Context, communityID uuid.UUID) ([]CommunityModeratorResponse, error) {
	moderators, err := s.store.GetCommunityModerators(ctx, communityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community moderators: %w", err)
	}

	return s.toCommunityModeratorResponses(moderators), nil
}

// GetCommunityAdmins gets all admins of a community
func (s *Service) GetCommunityAdmins(ctx context.Context, communityID uuid.UUID) ([]CommunityAdminResponse, error) {
	admins, err := s.store.GetCommunityAdmins(ctx, communityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community admins: %w", err)
	}

	return s.toCommunityAdminResponses(admins), nil
}

// JoinCommunity allows a user to join a community
func (s *Service) JoinCommunity(ctx context.Context, communityID, userID uuid.UUID) (*CommunityMembershipResponse, error) {
	membership, err := s.store.JoinCommunity(ctx, db.JoinCommunityParams{
		CommunityID: communityID,
		UserID:      userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to join community: %w", err)
	}

	// Update community stats asynchronously
	go s.store.UpdateCommunityStats(context.Background(), communityID)

	return s.toCommunityMembershipResponse(membership), nil
}

// LeaveCommunity allows a user to leave a community
func (s *Service) LeaveCommunity(ctx context.Context, communityID, userID uuid.UUID) error {
	err := s.store.LeaveCommunity(ctx, db.LeaveCommunityParams{
		CommunityID: communityID,
		UserID:      userID,
	})
	if err != nil {
		return fmt.Errorf("failed to leave community: %w", err)
	}

	// Update community stats asynchronously
	go s.store.UpdateCommunityStats(context.Background(), communityID)

	return nil
}

// AddCommunityModerator adds a moderator to a community
func (s *Service) AddCommunityModerator(ctx context.Context, communityID uuid.UUID, req AddModeratorRequest) (*CommunityMembershipResponse, error) {
	permissions := req.Permissions
	if permissions == nil {
		permissions = []string{}
	}

	membership, err := s.store.AddCommunityModerator(ctx, db.AddCommunityModeratorParams{
		CommunityID: communityID,
		UserID:      req.UserID,
		Permissions: permissions,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add moderator: %w", err)
	}

	return s.toCommunityMembershipResponse(membership), nil
}

// RemoveCommunityModerator removes a moderator from a community
func (s *Service) RemoveCommunityModerator(ctx context.Context, communityID, userID uuid.UUID) error {
	err := s.store.RemoveCommunityModerator(ctx, db.RemoveCommunityModeratorParams{
		CommunityID: communityID,
		UserID:      userID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove moderator: %w", err)
	}

	return nil
}

// IsCommunityAdmin checks if a user is a community admin
func (s *Service) IsCommunityAdmin(ctx context.Context, communityID, userID uuid.UUID) (bool, error) {
	isAdmin, err := s.store.IsCommunityAdmin(ctx, db.IsCommunityAdminParams{
		CommunityID: communityID,
		UserID:      userID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}

	return isAdmin, nil
}

// IsCommunityModerator checks if a user is a community moderator
func (s *Service) IsCommunityModerator(ctx context.Context, communityID, userID uuid.UUID) (bool, error) {
	isModerator, err := s.store.IsCommunityModerator(ctx, db.IsCommunityModeratorParams{
		CommunityID: communityID,
		UserID:      userID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check moderator status: %w", err)
	}

	return isModerator, nil
}

// GetUserCommunities gets all communities a user is a member of
func (s *Service) GetUserCommunities(ctx context.Context, userID, spaceID uuid.UUID) ([]UserCommunityResponse, error) {
	communities, err := s.store.GetUserCommunities(ctx, db.GetUserCommunitiesParams{
		UserID:  userID,
		SpaceID: spaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user communities: %w", err)
	}

	return s.toUserCommunityResponses(communities), nil
}

// GetCommunityCategories gets all available community categories
func (s *Service) GetCommunityCategories(ctx context.Context, spaceID uuid.UUID) ([]string, error) {
	categories, err := s.store.GetCommunityCategories(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community categories: %w", err)
	}

	return categories, nil
}

// Helper conversion functions

func (s *Service) toCommunityResponse(c db.Community) *CommunityResponse {
	resp := &CommunityResponse{
		ID:       c.ID,
		SpaceID:  c.SpaceID,
		Name:     c.Name,
		Category: c.Category,
	}

	if c.Description.Valid {
		resp.Description = &c.Description.String
	}
	if c.CoverImage.Valid {
		resp.CoverImage = &c.CoverImage.String
	}
	if c.MemberCount.Valid {
		resp.MemberCount = c.MemberCount.Int32
	}
	if c.Status.Valid {
		resp.Status = c.Status.String
	}
	if c.PostCount.Valid {
		resp.PostCount = c.PostCount.Int32
	}
	if c.IsPublic.Valid {
		resp.IsPublic = c.IsPublic.Bool
	}
	if c.CreatedBy.Valid {
		resp.CreatedBy = &c.CreatedBy.UUID
	}
	if c.Settings.Valid {
		resp.Settings = &c.Settings
	}
	if c.CreatedAt.Valid {
		resp.CreatedAt = &c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		resp.UpdatedAt = &c.UpdatedAt.Time
	}

	return resp
}

func (s *Service) toCommunityDetailResponse(c db.GetCommunityByIDRow) *CommunityDetailResponse {
	resp := &CommunityDetailResponse{
		ID:       c.ID,
		SpaceID:  c.SpaceID,
		Name:     c.Name,
		Category: c.Category,
		IsMember: c.IsMember,
	}

	if c.Description.Valid {
		resp.Description = &c.Description.String
	}
	if c.CoverImage.Valid {
		resp.CoverImage = &c.CoverImage.String
	}
	if c.MemberCount.Valid {
		resp.MemberCount = c.MemberCount.Int32
	}
	if c.Status.Valid {
		resp.Status = c.Status.String
	}
	if c.PostCount.Valid {
		resp.PostCount = c.PostCount.Int32
	}
	if c.IsPublic.Valid {
		resp.IsPublic = c.IsPublic.Bool
	}
	if c.CreatedBy.Valid {
		resp.CreatedBy = &c.CreatedBy.UUID
	}
	if c.Settings.Valid {
		resp.Settings = &c.Settings
	}
	if c.CreatedAt.Valid {
		resp.CreatedAt = &c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		resp.UpdatedAt = &c.UpdatedAt.Time
	}
	if c.CreatedByUsername != "" {
		resp.CreatedByUsername = &c.CreatedByUsername
	}
	if c.CreatedByFullName != "" {
		resp.CreatedByFullName = &c.CreatedByFullName
	}
	if c.UserRole.Valid {
		resp.UserRole = &c.UserRole.String
	}
	if c.ActualMemberCount > 0 {
		resp.ActualMemberCount = &c.ActualMemberCount
	}
	if c.ActualPostCount > 0 {
		resp.ActualPostCount = &c.ActualPostCount
	}

	return resp
}

func (s *Service) toCommunityDetailFromSlugResponse(c db.GetCommunityBySlugRow) *CommunityDetailResponse {
	resp := &CommunityDetailResponse{
		ID:       c.ID,
		SpaceID:  c.SpaceID,
		Name:     c.Name,
		Category: c.Category,
		IsMember: c.IsMember,
	}

	if c.Description.Valid {
		resp.Description = &c.Description.String
	}
	if c.CoverImage.Valid {
		resp.CoverImage = &c.CoverImage.String
	}
	if c.MemberCount.Valid {
		resp.MemberCount = c.MemberCount.Int32
	}
	if c.Status.Valid {
		resp.Status = c.Status.String
	}
	if c.PostCount.Valid {
		resp.PostCount = c.PostCount.Int32
	}
	if c.IsPublic.Valid {
		resp.IsPublic = c.IsPublic.Bool
	}
	if c.CreatedBy.Valid {
		resp.CreatedBy = &c.CreatedBy.UUID
	}
	if c.Settings.Valid {
		resp.Settings = &c.Settings
	}
	if c.CreatedAt.Valid {
		resp.CreatedAt = &c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		resp.UpdatedAt = &c.UpdatedAt.Time
	}
	if c.CreatedByUsername != "" {
		resp.CreatedByUsername = &c.CreatedByUsername
	}
	if c.CreatedByFullName != "" {
		resp.CreatedByFullName = &c.CreatedByFullName
	}
	if c.UserRole.Valid {
		resp.UserRole = &c.UserRole.String
	}

	return resp
}

func (s *Service) toCommunityListResponses(communities []db.ListCommunitiesRow) []CommunityListResponse {
	responses := make([]CommunityListResponse, len(communities))
	for i, c := range communities {
		responses[i] = s.toCommunityListResponse(c)
	}
	return responses
}

func (s *Service) toCommunityListResponse(c db.ListCommunitiesRow) CommunityListResponse {
	resp := CommunityListResponse{
		ID:                c.ID,
		SpaceID:           c.SpaceID,
		Name:              c.Name,
		Category:          c.Category,
		IsMember:          c.IsMember,
		ActualMemberCount: c.ActualMemberCount,
	}

	if c.Description.Valid {
		resp.Description = &c.Description.String
	}
	if c.CoverImage.Valid {
		resp.CoverImage = &c.CoverImage.String
	}
	if c.MemberCount.Valid {
		resp.MemberCount = c.MemberCount.Int32
	}
	if c.Status.Valid {
		resp.Status = c.Status.String
	}
	if c.PostCount.Valid {
		resp.PostCount = c.PostCount.Int32
	}
	if c.IsPublic.Valid {
		resp.IsPublic = c.IsPublic.Bool
	}
	if c.CreatedBy.Valid {
		resp.CreatedBy = &c.CreatedBy.UUID
	}
	if c.Settings.Valid {
		resp.Settings = &c.Settings
	}
	if c.CreatedAt.Valid {
		resp.CreatedAt = &c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		resp.UpdatedAt = &c.UpdatedAt.Time
	}
	if c.UserRole.Valid {
		resp.UserRole = &c.UserRole.String
	}

	return resp
}

func (s *Service) toSearchCommunityListResponses(communities []db.SearchCommunitiesRow) []CommunityListResponse {
	responses := make([]CommunityListResponse, len(communities))
	for i, c := range communities {
		responses[i] = s.toSearchCommunityListResponse(c)
	}
	return responses
}

func (s *Service) toSearchCommunityListResponse(c db.SearchCommunitiesRow) CommunityListResponse {
	resp := CommunityListResponse{
		ID:                c.ID,
		SpaceID:           c.SpaceID,
		Name:              c.Name,
		Category:          c.Category,
		IsMember:          c.IsMember,
		ActualMemberCount: c.ActualMemberCount,
	}

	if c.Description.Valid {
		resp.Description = &c.Description.String
	}
	if c.CoverImage.Valid {
		resp.CoverImage = &c.CoverImage.String
	}
	if c.MemberCount.Valid {
		resp.MemberCount = c.MemberCount.Int32
	}
	if c.Status.Valid {
		resp.Status = c.Status.String
	}
	if c.PostCount.Valid {
		resp.PostCount = c.PostCount.Int32
	}
	if c.IsPublic.Valid {
		resp.IsPublic = c.IsPublic.Bool
	}
	if c.CreatedBy.Valid {
		resp.CreatedBy = &c.CreatedBy.UUID
	}
	if c.Settings.Valid {
		resp.Settings = &c.Settings
	}
	if c.CreatedAt.Valid {
		resp.CreatedAt = &c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		resp.UpdatedAt = &c.UpdatedAt.Time
	}

	return resp
}

func (s *Service) toUserCommunityResponses(communities []db.GetUserCommunitiesRow) []UserCommunityResponse {
	responses := make([]UserCommunityResponse, len(communities))
	for i, c := range communities {
		responses[i] = s.toUserCommunityResponse(c)
	}
	return responses
}

func (s *Service) toUserCommunityResponse(c db.GetUserCommunitiesRow) UserCommunityResponse {
	resp := UserCommunityResponse{
		ID:       c.ID,
		SpaceID:  c.SpaceID,
		Name:     c.Name,
		Category: c.Category,
		UserRole: c.UserRole,
	}

	if c.Description.Valid {
		resp.Description = &c.Description.String
	}
	if c.CoverImage.Valid {
		resp.CoverImage = &c.CoverImage.String
	}
	if c.MemberCount.Valid {
		resp.MemberCount = c.MemberCount.Int32
	}
	if c.Status.Valid {
		resp.Status = c.Status.String
	}
	if c.PostCount.Valid {
		resp.PostCount = c.PostCount.Int32
	}
	if c.IsPublic.Valid {
		resp.IsPublic = c.IsPublic.Bool
	}
	if c.JoinedAt.Valid {
		resp.JoinedAt = &c.JoinedAt.Time
	}
	if c.CreatedAt.Valid {
		resp.CreatedAt = &c.CreatedAt.Time
	}
	if c.UpdatedAt.Valid {
		resp.UpdatedAt = &c.UpdatedAt.Time
	}

	return resp
}

func (s *Service) toCommunityMemberResponses(members []db.GetCommunityMembersRow) []CommunityMemberResponse {
	responses := make([]CommunityMemberResponse, len(members))
	for i, m := range members {
		responses[i] = s.toCommunityMemberResponse(m)
	}
	return responses
}

func (s *Service) toCommunityMemberResponse(m db.GetCommunityMembersRow) CommunityMemberResponse {
	resp := CommunityMemberResponse{
		ID:       m.ID,
		Username: m.Username,
		FullName: m.FullName,
		Role:     m.Role,
	}

	if m.Avatar.Valid {
		resp.Avatar = &m.Avatar.String
	}
	if m.Level.Valid {
		resp.Level = &m.Level.String
	}
	if m.Department.Valid {
		resp.Department = &m.Department.String
	}
	if m.Verified.Valid {
		resp.Verified = m.Verified.Bool
	}
	if m.JoinedAt.Valid {
		resp.JoinedAt = &m.JoinedAt.Time
	}

	return resp
}

func (s *Service) toCommunityModeratorResponses(moderators []db.GetCommunityModeratorsRow) []CommunityModeratorResponse {
	responses := make([]CommunityModeratorResponse, len(moderators))
	for i, m := range moderators {
		responses[i] = s.toCommunityModeratorResponse(m)
	}
	return responses
}

func (s *Service) toCommunityModeratorResponse(m db.GetCommunityModeratorsRow) CommunityModeratorResponse {
	resp := CommunityModeratorResponse{
		ID:          m.ID,
		Username:    m.Username,
		FullName:    m.FullName,
		Permissions: m.Permissions,
	}

	if m.Avatar.Valid {
		resp.Avatar = &m.Avatar.String
	}

	return resp
}

func (s *Service) toCommunityAdminResponses(admins []db.GetCommunityAdminsRow) []CommunityAdminResponse {
	responses := make([]CommunityAdminResponse, len(admins))
	for i, a := range admins {
		responses[i] = s.toCommunityAdminResponse(a)
	}
	return responses
}

func (s *Service) toCommunityAdminResponse(a db.GetCommunityAdminsRow) CommunityAdminResponse {
	resp := CommunityAdminResponse{
		ID:          a.ID,
		Username:    a.Username,
		FullName:    a.FullName,
		Permissions: a.Permissions,
	}

	if a.Avatar.Valid {
		resp.Avatar = &a.Avatar.String
	}

	return resp
}

func (s *Service) toCommunityMembershipResponse(m db.CommunityMember) *CommunityMembershipResponse {
	resp := &CommunityMembershipResponse{
		ID:          m.ID,
		CommunityID: m.CommunityID,
		UserID:      m.UserID,
		Role:        m.Role,
		Permissions: m.Permissions,
	}

	if m.JoinedAt.Valid {
		resp.JoinedAt = &m.JoinedAt.Time
	}

	return resp
}
