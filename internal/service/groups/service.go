package groups

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// Service handles group business logic
type Service struct {
	store db.Store
}

// NewService creates a new group service
func NewService(store db.Store) *Service {
	return &Service{store: store}
}

// CreateGroup creates a new group
func (s *Service) CreateGroup(ctx context.Context, req CreateGroupRequest) (*GroupResponse, error) {
	var communityID uuid.NullUUID
	var description, avatar, banner sql.NullString
	var allowInvites, allowMemberPosts sql.NullBool
	var settings pqtype.NullRawMessage
	
	if req.CommunityID != nil {
		communityID = uuid.NullUUID{UUID: *req.CommunityID, Valid: true}
	}
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.Avatar != nil {
		avatar = sql.NullString{String: *req.Avatar, Valid: true}
	}
	if req.Banner != nil {
		banner = sql.NullString{String: *req.Banner, Valid: true}
	}
	if req.AllowInvites != nil {
		allowInvites = sql.NullBool{Bool: *req.AllowInvites, Valid: true}
	}
	if req.AllowMemberPosts != nil {
		allowMemberPosts = sql.NullBool{Bool: *req.AllowMemberPosts, Valid: true}
	}
	if req.Settings != nil {
		settings = *req.Settings
	}
	
	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}
	
	createdBy := uuid.NullUUID{UUID: req.CreatedBy, Valid: true}
	
	group, err := s.store.CreateGroup(ctx, db.CreateGroupParams{
		SpaceID:          req.SpaceID,
		CommunityID:      communityID,
		Name:             req.Name,
		Description:      description,
		Category:         req.Category,
		GroupType:        req.GroupType,
		Avatar:           avatar,
		Banner:           banner,
		AllowInvites:     allowInvites,
		AllowMemberPosts: allowMemberPosts,
		CreatedBy:        createdBy,
		Tags:             tags,
		Settings:         settings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}
	
	return s.toGroupResponse(group), nil
}

// GetGroupByID gets a group by ID
func (s *Service) GetGroupByID(ctx context.Context, userID, groupID uuid.UUID) (*GroupDetailResponse, error) {
	group, err := s.store.GetGroupByID(ctx, db.GetGroupByIDParams{
		UserID: userID,
		ID:     groupID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	
	return s.toGroupDetailResponse(group), nil
}

// ListGroups lists groups with pagination and sorting
func (s *Service) ListGroups(ctx context.Context, params ListGroupsParams) ([]GroupListResponse, error) {
	offset := (params.Page - 1) * params.Limit
	
	groups, err := s.store.ListGroups(ctx, db.ListGroupsParams{
		UserID:  params.UserID,
		SpaceID: params.SpaceID,
		Column3: params.SortBy,
		Limit:   params.Limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	
	return s.toGroupListResponses(groups), nil
}

// SearchGroups searches groups by name, description, or tags
func (s *Service) SearchGroups(ctx context.Context, params SearchGroupsParams) ([]GroupListResponse, error) {
	groups, err := s.store.SearchGroups(ctx, db.SearchGroupsParams{
		UserID:  params.UserID,
		SpaceID: params.SpaceID,
		Name:    "%" + params.Query + "%",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search groups: %w", err)
	}
	
	return s.toSearchGroupListResponses(groups), nil
}

// UpdateGroup updates a group
func (s *Service) UpdateGroup(ctx context.Context, groupID uuid.UUID, req UpdateGroupRequest) (*GroupResponse, error) {
	var description, avatar, banner sql.NullString
	var allowInvites, allowMemberPosts sql.NullBool
	var settings pqtype.NullRawMessage
	
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.Avatar != nil {
		avatar = sql.NullString{String: *req.Avatar, Valid: true}
	}
	if req.Banner != nil {
		banner = sql.NullString{String: *req.Banner, Valid: true}
	}
	if req.AllowInvites != nil {
		allowInvites = sql.NullBool{Bool: *req.AllowInvites, Valid: true}
	}
	if req.AllowMemberPosts != nil {
		allowMemberPosts = sql.NullBool{Bool: *req.AllowMemberPosts, Valid: true}
	}
	if req.Settings != nil {
		settings = *req.Settings
	}
	
	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}
	
	group, err := s.store.UpdateGroup(ctx, db.UpdateGroupParams{
		ID:               groupID,
		Name:             req.Name,
		Description:      description,
		Category:         req.Category,
		Avatar:           avatar,
		Banner:           banner,
		AllowInvites:     allowInvites,
		AllowMemberPosts: allowMemberPosts,
		Tags:             tags,
		Settings:         settings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}
	
	return s.toGroupResponse(group), nil
}

// GetUserGroups gets all groups a user is a member of
func (s *Service) GetUserGroups(ctx context.Context, userID, spaceID uuid.UUID) ([]UserGroupResponse, error) {
	groups, err := s.store.GetUserGroups(ctx, db.GetUserGroupsParams{
		UserID:  userID,
		SpaceID: spaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	
	return s.toUserGroupResponses(groups), nil
}

// JoinGroup allows a user to join a group
func (s *Service) JoinGroup(ctx context.Context, groupID, userID uuid.UUID, invitedBy *uuid.UUID) (*GroupMembershipResponse, error) {
	var inviter uuid.NullUUID
	if invitedBy != nil {
		inviter = uuid.NullUUID{UUID: *invitedBy, Valid: true}
	}
	
	membership, err := s.store.JoinGroup(ctx, db.JoinGroupParams{
		GroupID:   groupID,
		UserID:    userID,
		InvitedBy: inviter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}
	
	// Update group stats asynchronously
	go s.store.UpdateGroupStats(context.Background(), groupID)
	
	return s.toGroupMembershipResponse(membership), nil
}

// LeaveGroup allows a user to leave a group
func (s *Service) LeaveGroup(ctx context.Context, groupID, userID uuid.UUID) error {
	err := s.store.LeaveGroup(ctx, db.LeaveGroupParams{
		GroupID: groupID,
		UserID:  userID,
	})
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}
	
	// Update group stats asynchronously
	go s.store.UpdateGroupStats(context.Background(), groupID)
	
	return nil
}

// GetGroupJoinRequests gets all join requests for a group
func (s *Service) GetGroupJoinRequests(ctx context.Context, groupID uuid.UUID) ([]GroupMemberResponse, error) {
	requests, err := s.store.GetGroupJoinRequests(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get join requests: %w", err)
	}
	
	return s.toGroupMemberResponses(requests), nil
}

// AddGroupAdmin adds an admin to a group
func (s *Service) AddGroupAdmin(ctx context.Context, groupID uuid.UUID, req AddGroupAdminRequest) (*GroupMembershipResponse, error) {
	permissions := req.Permissions
	if permissions == nil {
		permissions = []string{}
	}
	
	membership, err := s.store.AddGroupAdmin(ctx, db.AddGroupAdminParams{
		GroupID:       groupID,
		UserID:        req.UserID,
		Permissions:   permissions,
		Permissions_2: permissions,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add admin: %w", err)
	}
	
	return s.toGroupMembershipResponse(membership), nil
}

// RemoveGroupAdmin removes an admin from a group
func (s *Service) RemoveGroupAdmin(ctx context.Context, groupID, userID uuid.UUID) error {
	err := s.store.RemoveGroupAdmin(ctx, db.RemoveGroupAdminParams{
		GroupID: groupID,
		UserID:  userID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove admin: %w", err)
	}
	
	return nil
}

// AddGroupModerator adds a moderator to a group
func (s *Service) AddGroupModerator(ctx context.Context, groupID uuid.UUID, req AddGroupModeratorRequest) (*GroupMembershipResponse, error) {
	permissions := req.Permissions
	if permissions == nil {
		permissions = []string{}
	}
	
	membership, err := s.store.AddGroupModerator(ctx, db.AddGroupModeratorParams{
		GroupID:       groupID,
		UserID:        req.UserID,
		Permissions:   permissions,
		Permissions_2: permissions,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add moderator: %w", err)
	}
	
	return s.toGroupMembershipResponse(membership), nil
}

// RemoveGroupModerator removes a moderator from a group
func (s *Service) RemoveGroupModerator(ctx context.Context, groupID, userID uuid.UUID) error {
	err := s.store.RemoveGroupModerator(ctx, db.RemoveGroupModeratorParams{
		GroupID: groupID,
		UserID:  userID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove moderator: %w", err)
	}
	
	return nil
}

// IsGroupAdmin checks if a user is a group admin
func (s *Service) IsGroupAdmin(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	isAdmin, err := s.store.IsGroupAdmin(ctx, db.IsGroupAdminParams{
		GroupID: groupID,
		UserID:  userID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}
	
	return isAdmin, nil
}

// IsGroupModerator checks if a user is a group moderator
func (s *Service) IsGroupModerator(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	isModerator, err := s.store.IsGroupModerator(ctx, db.IsGroupModeratorParams{
		GroupID: groupID,
		UserID:  userID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check moderator status: %w", err)
	}
	
	return isModerator, nil
}

// UpdateGroupMemberRole updates a member's role
func (s *Service) UpdateGroupMemberRole(ctx context.Context, groupID, userID uuid.UUID, req UpdateMemberRoleRequest) error {
	permissions := req.Permissions
	if permissions == nil {
		permissions = []string{}
	}
	
	err := s.store.UpdateGroupMemberRole(ctx, db.UpdateGroupMemberRoleParams{
		GroupID:     groupID,
		UserID:      userID,
		Role:        req.Role,
		Permissions: permissions,
	})
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}
	
	return nil
}

// CreateProjectRole creates a new project role
func (s *Service) CreateProjectRole(ctx context.Context, groupID uuid.UUID, req CreateProjectRoleRequest) (*ProjectRoleResponse, error) {
	var description, requirements sql.NullString
	
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.Requirements != nil {
		requirements = sql.NullString{String: *req.Requirements, Valid: true}
	}
	
	skillsRequired := req.SkillsRequired
	if skillsRequired == nil {
		skillsRequired = []string{}
	}
	
	role, err := s.store.CreateProjectRole(ctx, db.CreateProjectRoleParams{
		GroupID:        groupID,
		Name:           req.Name,
		Description:    description,
		SlotsTotal:     req.SlotsTotal,
		Requirements:   requirements,
		SkillsRequired: skillsRequired,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create project role: %w", err)
	}
	
	return s.toProjectRoleResponse(role), nil
}

// GetProjectRoles gets all project roles for a group
func (s *Service) GetProjectRoles(ctx context.Context, groupID uuid.UUID) ([]ProjectRoleResponse, error) {
	roles, err := s.store.GetProjectRoles(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project roles: %w", err)
	}
	
	return s.toProjectRoleResponses(roles), nil
}

// ApplyForProjectRole applies for a project role
func (s *Service) ApplyForProjectRole(ctx context.Context, roleID, userID uuid.UUID, req ApplyForRoleRequest) (*RoleApplicationResponse, error) {
	var message sql.NullString
	if req.Message != nil {
		message = sql.NullString{String: *req.Message, Valid: true}
	}
	
	application, err := s.store.ApplyForProjectRole(ctx, db.ApplyForProjectRoleParams{
		RoleID:  roleID,
		UserID:  userID,
		Message: message,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to apply for role: %w", err)
	}
	
	return s.toSimpleRoleApplicationResponse(application), nil
}

// GetRoleApplications gets all role applications for a group
func (s *Service) GetRoleApplications(ctx context.Context, groupID uuid.UUID) ([]RoleApplicationResponse, error) {
	applications, err := s.store.GetRoleApplications(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role applications: %w", err)
	}
	
	return s.toRoleApplicationResponses(applications), nil
}

// Helper conversion functions

func (s *Service) toGroupResponse(g db.Group) *GroupResponse {
	resp := &GroupResponse{
		ID:        g.ID,
		SpaceID:   g.SpaceID,
		Name:      g.Name,
		Category:  g.Category,
		GroupType: g.GroupType,
		Tags:      g.Tags,
	}
	
	if g.CommunityID.Valid {
		resp.CommunityID = &g.CommunityID.UUID
	}
	if g.Description.Valid {
		resp.Description = &g.Description.String
	}
	if g.Avatar.Valid {
		resp.Avatar = &g.Avatar.String
	}
	if g.Banner.Valid {
		resp.Banner = &g.Banner.String
	}
	if g.MemberCount.Valid {
		resp.MemberCount = g.MemberCount.Int32
	}
	if g.PostCount.Valid {
		resp.PostCount = g.PostCount.Int32
	}
	if g.Status.Valid {
		resp.Status = g.Status.String
	}
	if g.Visibility.Valid {
		resp.Visibility = g.Visibility.String
	}
	if g.AllowInvites.Valid {
		resp.AllowInvites = g.AllowInvites.Bool
	}
	if g.AllowMemberPosts.Valid {
		resp.AllowMemberPosts = g.AllowMemberPosts.Bool
	}
	if g.CreatedBy.Valid {
		resp.CreatedBy = &g.CreatedBy.UUID
	}
	if g.Settings.Valid {
		resp.Settings = &g.Settings
	}
	if g.CreatedAt.Valid {
		resp.CreatedAt = &g.CreatedAt.Time
	}
	if g.UpdatedAt.Valid {
		resp.UpdatedAt = &g.UpdatedAt.Time
	}
	
	return resp
}

func (s *Service) toGroupDetailResponse(g db.GetGroupByIDRow) *GroupDetailResponse {
	resp := &GroupDetailResponse{
		ID:                g.ID,
		SpaceID:           g.SpaceID,
		Name:              g.Name,
		Category:          g.Category,
		GroupType:         g.GroupType,
		Tags:              g.Tags,
		CreatedByUsername: g.CreatedByUsername,
		CreatedByFullName: g.CreatedByFullName,
		IsMember:          g.IsMember,
		ActualMemberCount: g.ActualMemberCount,
		ActualPostCount:   g.ActualPostCount,
	}
	
	if g.CommunityID.Valid {
		resp.CommunityID = &g.CommunityID.UUID
	}
	if g.Description.Valid {
		resp.Description = &g.Description.String
	}
	if g.Avatar.Valid {
		resp.Avatar = &g.Avatar.String
	}
	if g.Banner.Valid {
		resp.Banner = &g.Banner.String
	}
	if g.MemberCount.Valid {
		resp.MemberCount = g.MemberCount.Int32
	}
	if g.PostCount.Valid {
		resp.PostCount = g.PostCount.Int32
	}
	if g.Status.Valid {
		resp.Status = g.Status.String
	}
	if g.Visibility.Valid {
		resp.Visibility = g.Visibility.String
	}
	if g.AllowInvites.Valid {
		resp.AllowInvites = g.AllowInvites.Bool
	}
	if g.AllowMemberPosts.Valid {
		resp.AllowMemberPosts = g.AllowMemberPosts.Bool
	}
	if g.CreatedBy.Valid {
		resp.CreatedBy = &g.CreatedBy.UUID
	}
	if g.Settings.Valid {
		resp.Settings = &g.Settings
	}
	if g.CreatedAt.Valid {
		resp.CreatedAt = &g.CreatedAt.Time
	}
	if g.UpdatedAt.Valid {
		resp.UpdatedAt = &g.UpdatedAt.Time
	}
	if g.CommunityName.Valid {
		resp.CommunityName = &g.CommunityName.String
	}
	if g.UserRole.Valid {
		resp.UserRole = &g.UserRole.String
	}
	
	return resp
}

func (s *Service) toGroupListResponses(groups []db.ListGroupsRow) []GroupListResponse {
	responses := make([]GroupListResponse, len(groups))
	for i, g := range groups {
		responses[i] = s.toGroupListResponse(g)
	}
	return responses
}

func (s *Service) toGroupListResponse(g db.ListGroupsRow) GroupListResponse {
	resp := GroupListResponse{
		ID:                g.ID,
		SpaceID:           g.SpaceID,
		Name:              g.Name,
		Category:          g.Category,
		GroupType:         g.GroupType,
		Tags:              g.Tags,
		IsMember:          g.IsMember,
		ActualMemberCount: g.ActualMemberCount,
	}
	
	if g.CommunityID.Valid {
		resp.CommunityID = &g.CommunityID.UUID
	}
	if g.Description.Valid {
		resp.Description = &g.Description.String
	}
	if g.Avatar.Valid {
		resp.Avatar = &g.Avatar.String
	}
	if g.Banner.Valid {
		resp.Banner = &g.Banner.String
	}
	if g.MemberCount.Valid {
		resp.MemberCount = g.MemberCount.Int32
	}
	if g.PostCount.Valid {
		resp.PostCount = g.PostCount.Int32
	}
	if g.Status.Valid {
		resp.Status = g.Status.String
	}
	if g.Visibility.Valid {
		resp.Visibility = g.Visibility.String
	}
	if g.AllowInvites.Valid {
		resp.AllowInvites = g.AllowInvites.Bool
	}
	if g.AllowMemberPosts.Valid {
		resp.AllowMemberPosts = g.AllowMemberPosts.Bool
	}
	if g.CreatedBy.Valid {
		resp.CreatedBy = &g.CreatedBy.UUID
	}
	if g.Settings.Valid {
		resp.Settings = &g.Settings
	}
	if g.CreatedAt.Valid {
		resp.CreatedAt = &g.CreatedAt.Time
	}
	if g.UpdatedAt.Valid {
		resp.UpdatedAt = &g.UpdatedAt.Time
	}
	if g.CommunityName.Valid {
		resp.CommunityName = &g.CommunityName.String
	}
	if g.UserRole.Valid {
		resp.UserRole = &g.UserRole.String
	}
	
	return resp
}

func (s *Service) toSearchGroupListResponses(groups []db.SearchGroupsRow) []GroupListResponse {
	responses := make([]GroupListResponse, len(groups))
	for i, g := range groups {
		responses[i] = s.toSearchGroupListResponse(g)
	}
	return responses
}

func (s *Service) toSearchGroupListResponse(g db.SearchGroupsRow) GroupListResponse {
	resp := GroupListResponse{
		ID:                g.ID,
		SpaceID:           g.SpaceID,
		Name:              g.Name,
		Category:          g.Category,
		GroupType:         g.GroupType,
		Tags:              g.Tags,
		IsMember:          g.IsMember,
		ActualMemberCount: g.ActualMemberCount,
	}
	
	if g.CommunityID.Valid {
		resp.CommunityID = &g.CommunityID.UUID
	}
	if g.Description.Valid {
		resp.Description = &g.Description.String
	}
	if g.Avatar.Valid {
		resp.Avatar = &g.Avatar.String
	}
	if g.Banner.Valid {
		resp.Banner = &g.Banner.String
	}
	if g.MemberCount.Valid {
		resp.MemberCount = g.MemberCount.Int32
	}
	if g.PostCount.Valid {
		resp.PostCount = g.PostCount.Int32
	}
	if g.Status.Valid {
		resp.Status = g.Status.String
	}
	if g.Visibility.Valid {
		resp.Visibility = g.Visibility.String
	}
	if g.AllowInvites.Valid {
		resp.AllowInvites = g.AllowInvites.Bool
	}
	if g.AllowMemberPosts.Valid {
		resp.AllowMemberPosts = g.AllowMemberPosts.Bool
	}
	if g.CreatedBy.Valid {
		resp.CreatedBy = &g.CreatedBy.UUID
	}
	if g.Settings.Valid {
		resp.Settings = &g.Settings
	}
	if g.CreatedAt.Valid {
		resp.CreatedAt = &g.CreatedAt.Time
	}
	if g.UpdatedAt.Valid {
		resp.UpdatedAt = &g.UpdatedAt.Time
	}
	if g.CommunityName.Valid {
		resp.CommunityName = &g.CommunityName.String
	}
	
	return resp
}

func (s *Service) toUserGroupResponses(groups []db.GetUserGroupsRow) []UserGroupResponse {
	responses := make([]UserGroupResponse, len(groups))
	for i, g := range groups {
		responses[i] = s.toUserGroupResponse(g)
	}
	return responses
}

func (s *Service) toUserGroupResponse(g db.GetUserGroupsRow) UserGroupResponse {
	resp := UserGroupResponse{
		ID:        g.ID,
		SpaceID:   g.SpaceID,
		Name:      g.Name,
		Category:  g.Category,
		GroupType: g.GroupType,
		Tags:      g.Tags,
		UserRole:  g.UserRole,
	}
	
	if g.CommunityID.Valid {
		resp.CommunityID = &g.CommunityID.UUID
	}
	if g.Description.Valid {
		resp.Description = &g.Description.String
	}
	if g.Avatar.Valid {
		resp.Avatar = &g.Avatar.String
	}
	if g.Banner.Valid {
		resp.Banner = &g.Banner.String
	}
	if g.MemberCount.Valid {
		resp.MemberCount = g.MemberCount.Int32
	}
	if g.PostCount.Valid {
		resp.PostCount = g.PostCount.Int32
	}
	if g.Status.Valid {
		resp.Status = g.Status.String
	}
	if g.Visibility.Valid {
		resp.Visibility = g.Visibility.String
	}
	if g.AllowInvites.Valid {
		resp.AllowInvites = g.AllowInvites.Bool
	}
	if g.AllowMemberPosts.Valid {
		resp.AllowMemberPosts = g.AllowMemberPosts.Bool
	}
	if g.CreatedBy.Valid {
		resp.CreatedBy = &g.CreatedBy.UUID
	}
	if g.Settings.Valid {
		resp.Settings = &g.Settings
	}
	if g.CreatedAt.Valid {
		resp.CreatedAt = &g.CreatedAt.Time
	}
	if g.UpdatedAt.Valid {
		resp.UpdatedAt = &g.UpdatedAt.Time
	}
	if g.CommunityName.Valid {
		resp.CommunityName = &g.CommunityName.String
	}
	if g.JoinedAt.Valid {
		resp.JoinedAt = &g.JoinedAt.Time
	}
	
	return resp
}

func (s *Service) toGroupMemberResponses(members []db.GetGroupJoinRequestsRow) []GroupMemberResponse {
	responses := make([]GroupMemberResponse, len(members))
	for i, m := range members {
		responses[i] = s.toGroupMemberResponse(m)
	}
	return responses
}

func (s *Service) toGroupMemberResponse(m db.GetGroupJoinRequestsRow) GroupMemberResponse {
	resp := GroupMemberResponse{
		ID:          m.ID,
		Username:    m.Username,
		FullName:    m.FullName,
		Role:        m.Role,
		Permissions: m.Permissions,
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

func (s *Service) toGroupMembershipResponse(m db.GroupMember) *GroupMembershipResponse {
	resp := &GroupMembershipResponse{
		ID:          m.ID,
		GroupID:     m.GroupID,
		UserID:      m.UserID,
		Role:        m.Role,
		Permissions: m.Permissions,
	}
	
	if m.JoinedAt.Valid {
		resp.JoinedAt = &m.JoinedAt.Time
	}
	if m.InvitedBy.Valid {
		resp.InvitedBy = &m.InvitedBy.UUID
	}
	
	return resp
}

func (s *Service) toProjectRoleResponses(roles []db.GroupRole) []ProjectRoleResponse {
	responses := make([]ProjectRoleResponse, len(roles))
	for i, r := range roles {
		responses[i] = *s.toProjectRoleResponse(r)
	}
	return responses
}

func (s *Service) toProjectRoleResponse(r db.GroupRole) *ProjectRoleResponse {
	slotsFilled := int32(0)
	if r.SlotsFilled.Valid {
		slotsFilled = r.SlotsFilled.Int32
	}

	resp := &ProjectRoleResponse{
		ID:             r.ID,
		GroupID:        r.GroupID,
		Name:           r.Name,
		SlotsTotal:     r.SlotsTotal,
		SlotsFilled:    slotsFilled,
		SkillsRequired: r.SkillsRequired,
	}
	
	if r.Description.Valid {
		resp.Description = &r.Description.String
	}
	if r.Requirements.Valid {
		resp.Requirements = &r.Requirements.String
	}
	if r.CreatedAt.Valid {
		resp.CreatedAt = &r.CreatedAt.Time
	}
	if r.UpdatedAt.Valid {
		resp.UpdatedAt = &r.UpdatedAt.Time
	}
	
	return resp
}

func (s *Service) toSimpleRoleApplicationResponse(a db.GroupApplication) *RoleApplicationResponse {
	resp := &RoleApplicationResponse{
		ID:     a.ID,
		RoleID: a.RoleID,
		UserID: a.UserID,
	}
	
	if a.Message.Valid {
		resp.Message = &a.Message.String
	}
	if a.Status.Valid {
		resp.Status = a.Status.String
	}
	if a.AppliedAt.Valid {
		resp.AppliedAt = &a.AppliedAt.Time
	}
	if a.ReviewedAt.Valid {
		resp.ReviewedAt = &a.ReviewedAt.Time
	}
	if a.ReviewedBy.Valid {
		resp.ReviewedBy = &a.ReviewedBy.UUID
	}
	if a.ReviewNotes.Valid {
		resp.ReviewNotes = &a.ReviewNotes.String
	}
	
	return resp
}

func (s *Service) toRoleApplicationResponses(applications []db.GetRoleApplicationsRow) []RoleApplicationResponse {
	responses := make([]RoleApplicationResponse, len(applications))
	for i, a := range applications {
		responses[i] = s.toRoleApplicationResponse(a)
	}
	return responses
}

func (s *Service) toRoleApplicationResponse(a db.GetRoleApplicationsRow) RoleApplicationResponse {
	resp := RoleApplicationResponse{
		ID:       a.ID,
		RoleID:   a.RoleID,
		UserID:   a.UserID,
		Username: a.Username,
		FullName: a.FullName,
		RoleName: a.RoleName,
	}
	
	if a.Message.Valid {
		resp.Message = &a.Message.String
	}
	if a.Status.Valid {
		resp.Status = a.Status.String
	}
	if a.AppliedAt.Valid {
		resp.AppliedAt = &a.AppliedAt.Time
	}
	if a.ReviewedAt.Valid {
		resp.ReviewedAt = &a.ReviewedAt.Time
	}
	if a.ReviewedBy.Valid {
		resp.ReviewedBy = &a.ReviewedBy.UUID
	}
	if a.ReviewNotes.Valid {
		resp.ReviewNotes = &a.ReviewNotes.String
	}
	if a.Avatar.Valid {
		resp.Avatar = &a.Avatar.String
	}
	
	return resp
}
