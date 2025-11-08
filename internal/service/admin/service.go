package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Service struct {
	store db.Store
}

func NewService(store db.Store) *Service {
	return &Service{store: store}
}



func (s *Service) SuspendUser(ctx context.Context, req SuspendUserRequest) error {
	
	var suspendedUntil sql.NullTime
	if req.DurationDays > 0 {
		suspendedUntil = sql.NullTime{
			Time:  time.Now().AddDate(0, 0, req.DurationDays),
			Valid: true,
		}
	}

	_, err := s.store.CreateUserSuspension(ctx, db.CreateUserSuspensionParams{
		UserID:         req.UserID,
		SuspendedBy:    req.SuspendedBy,
		Reason:         req.Reason,
		Notes:          sql.NullString{String: req.Notes, Valid: req.Notes != ""},
		SuspendedUntil: suspendedUntil,
		IsPermanent:    req.IsPermanent,
	})
	if err != nil {
		return fmt.Errorf("failed to create suspension: %w", err)
	}

	
	err = s.store.UpdateUserAccountStatus(ctx, db.UpdateUserAccountStatusParams{
		ID:     req.UserID,
		Status: sql.NullString{String: "suspended", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	
	details := fmt.Sprintf(`{"reason": "%s", "duration_days": %d}`, req.Reason, req.DurationDays)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  req.SuspendedBy,
		Action:       "suspend_user",
		ResourceType: "user",
		ResourceID:   uuid.NullUUID{UUID: req.UserID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return nil
}

func (s *Service) UnsuspendUser(ctx context.Context, userID, adminID uuid.UUID) error {
	err := s.store.LiftSuspension(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to unsuspend user: %w", err)
	}

	
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "unsuspend_user",
		ResourceType: "user",
		ResourceID:   uuid.NullUUID{UUID: userID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
	})

	return nil
}

func (s *Service) BanUser(ctx context.Context, userID, adminID uuid.UUID, reason string) error {
	
	_, err := s.store.CreateUserSuspension(ctx, db.CreateUserSuspensionParams{
		UserID:      userID,
		SuspendedBy: adminID,
		Reason:      reason,
		Notes:       sql.NullString{String: "Permanent ban", Valid: true},
		IsPermanent: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create ban: %w", err)
	}

	
	err = s.store.UpdateUserAccountStatus(ctx, db.UpdateUserAccountStatusParams{
		ID:     userID,
		Status: sql.NullString{String: "banned", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	
	details := fmt.Sprintf(`{"reason": "%s"}`, reason)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "ban_user",
		ResourceType: "user",
		ResourceID:   uuid.NullUUID{UUID: userID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return nil
}



func (s *Service) GetContentReports(ctx context.Context, req GetReportsRequest) ([]ContentReportResponse, error) {
	
	reports, err := s.store.GetContentReports(ctx, db.GetContentReportsParams{
		SpaceID:     req.SpaceID,
		Status:      sql.NullString{String: req.Status, Valid: req.Status != ""},
		ContentType: req.ContentType,
		Limit:       req.Limit,
		Offset:      req.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	responses := make([]ContentReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = ContentReportResponse{
			ID:          report.ID,
			SpaceID:     report.SpaceID,
			ContentType: report.ContentType,
			ContentID:   report.ContentID,
			Reason:      report.Reason,
			Description: report.Description.String,
			Status:      report.Status.String,
			CreatedAt:   report.CreatedAt.Time,
		}
		if report.ReviewedAt.Valid {
			responses[i].ResolvedAt = &report.ReviewedAt.Time
		}
	}

	return responses, nil
}



func (s *Service) GetSpaceActivities(ctx context.Context, spaceID uuid.UUID, activityType string, since time.Time, limit, offset int32) ([]SpaceActivityResponse, error) {
	
	activities, err := s.store.GetSpaceActivities(ctx, db.GetSpaceActivitiesParams{
		SpaceID:      spaceID,
		ActivityType: activityType,
		CreatedAt:    since,
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}

	responses := make([]SpaceActivityResponse, len(activities))
	for i, activity := range activities {
		responses[i] = SpaceActivityResponse{
			ID:           activity.ID,
			ActivityType: activity.ActivityType,
			Description:  activity.Description,
			CreatedAt:    activity.CreatedAt,
		}
		if activity.ActorID.Valid {
			responses[i].ActorID = &activity.ActorID.UUID
		}
		if activity.ActorName.Valid {
			responses[i].ActorName = activity.ActorName.String
		}
	}

	return responses, nil
}



func (s *Service) GetDashboardStats(ctx context.Context, spaceID uuid.UUID) (DashboardStatsResponse, error) {
	stats, err := s.store.GetAdminDashboardStats(ctx, spaceID)
	if err != nil {
		return DashboardStatsResponse{}, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	return DashboardStatsResponse{
		TotalUsers:       stats.TotalUsers,
		NewUsersMonth:    stats.NewUsersMonth,
		TotalPosts:       stats.TotalPosts,
		TotalCommunities: stats.TotalCommunities,
		TotalGroups:      stats.TotalGroups,
		PendingReports:   stats.PendingReports,
		SuspensionsMonth: stats.SuspensionsMonth,
	}, nil
}



func (s *Service) GetUsers(ctx context.Context, spaceID uuid.UUID, limit, offset int32) ([]UserResponse, int64, error) {
	
	users, err := s.store.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	
	count := int64(len(users))

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			FullName:   user.FullName,
			Email:      user.Email,
			Avatar:     user.Avatar.String,
			Status:     user.Status.String,
			CreatedAt:  user.CreatedAt.Time,
			Roles:      user.Roles,
			Department: user.Department.String,
		}
	}

	return responses, count, nil
}

func (s *Service) DeleteUser(ctx context.Context, userID, adminID uuid.UUID) error {
	
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "delete_user",
		ResourceType: "user",
		ResourceID:   uuid.NullUUID{UUID: userID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
	})

	
	err := s.store.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}



func (s *Service) GetTutorApplications(ctx context.Context, spaceID uuid.UUID, limit, offset int32) ([]GetAllTutorApplicationsResponse, error) {
	var applications []db.GetAllTutorApplicationsRow
	var err error

	applications, err = s.store.GetAllTutorApplications(ctx, db.GetAllTutorApplicationsParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get tutor applications: %w", err)
	}

	responses := make([]GetAllTutorApplicationsResponse, len(applications))
	for i, app := range applications {
		responses[i] = GetAllTutorApplicationsResponse{
			ID:          app.ID,
			UserID:      app.UserID,
			Status:      app.Status.String,
			SubmittedAt: app.SubmittedAt.Time,
			ApplicantID: app.ApplicantID,
			Subjects:    app.Subjects,
			HourlyRate:  app.HourlyRate.String,
			FullName:    app.FullName,
		}
	}

	return responses, nil
}

func (s *Service) ApproveTutorApplication(ctx context.Context, appID, adminID uuid.UUID, notes string) error {
	application, err := s.store.GetUserTutorApplicationStatusById(ctx, appID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no mentor application found")
		}
		return fmt.Errorf("failed to get mentor application: %w", err)
	}

	if application.String == "approved" {
		return fmt.Errorf("mentor already approved")
	}

	tutor, err := s.store.UpdateTutorApplicationStatus(ctx, db.UpdateTutorApplicationStatusParams{
		ID:            appID,
		Status:        sql.NullString{String: "approved", Valid: true},
		ReviewedBy:    uuid.NullUUID{UUID: adminID, Valid: true},
		ReviewerNotes: sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to approve tutor application: %w", err)
	}

	
	details := fmt.Sprintf(`{"notes": "%s"}`, notes)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "approve_tutor_application",
		ResourceType: "tutor_application",
		ResourceID:   uuid.NullUUID{UUID: appID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	_, err = s.store.CreateTutorProfile(ctx, db.CreateTutorProfileParams{
		UserID:         tutor.ApplicantID,
		SpaceID:        tutor.SpaceID,
		Subjects:       tutor.Subjects,
		HourlyRate:     tutor.HourlyRate,
		Description:    tutor.Qualifications,
		Availability:   pqtype.NullRawMessage{RawMessage: tutor.Availability, Valid: true},
		Experience:     tutor.Experience,
		Verified:       sql.NullBool{Bool: true, Valid: true},
		Qualifications: tutor.Qualifications,
	})

	if err != nil {
		return fmt.Errorf("failed to create tutor profile: %w", err)
	}

	return nil
}

func (s *Service) RejectTutorApplication(ctx context.Context, appID, adminID uuid.UUID, notes string) error {
	_, err := s.store.UpdateTutorApplicationStatus(ctx, db.UpdateTutorApplicationStatusParams{
		ID:            appID,
		Status:        sql.NullString{String: "rejected", Valid: true},
		ReviewedBy:    uuid.NullUUID{UUID: adminID, Valid: true},
		ReviewerNotes: sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to reject tutor application: %w", err)
	}

	
	details := fmt.Sprintf(`{"notes": "%s"}`, notes)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "reject_tutor_application",
		ResourceType: "tutor_application",
		ResourceID:   uuid.NullUUID{UUID: appID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return nil
}

func (s *Service) GetMentorApplications(ctx context.Context, spaceID uuid.UUID, limit, offset int32) ([]GetAllMentorApplicationsResponse, error) {
	var applications []db.GetAllMentorApplicationsRow

	applications, err := s.store.GetAllMentorApplications(ctx, db.GetAllMentorApplicationsParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get mentor applications: %w", err)
	}

	responses := make([]GetAllMentorApplicationsResponse, len(applications))
	for i, app := range applications {
		responses[i] = GetAllMentorApplicationsResponse{
			ID:          app.ID,
			ApplicantID: app.ApplicantID,
			Status:      app.Status.String,
			SubmittedAt: app.SubmittedAt.Time,
			Industry:    app.Industry,
			Experience:  app.Experience,
			Specialties: app.Specialties,
			FullName:    app.FullName,
			UserID:      app.UserID,
			Position:    app.Position.String,
			Company:     app.Company.String,
		}

	}

	return responses, nil
}

func (s *Service) ApproveMentorApplication(ctx context.Context, appID, adminID uuid.UUID, notes string) error {
	application, err := s.store.GetUserMentorApplicationStatusById(ctx, appID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no mentor application found")
		}
		return fmt.Errorf("failed to get mentor application: %w", err)
	}

	if application.String == "approved" {
		return fmt.Errorf("mentor already approved")
	}

	mentor, err := s.store.UpdateMentorApplicationStatus(ctx, db.UpdateMentorApplicationStatusParams{
		ID:            appID,
		Status:        sql.NullString{String: "approved", Valid: true},
		ReviewedBy:    uuid.NullUUID{UUID: adminID, Valid: true},
		ReviewerNotes: sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to approve mentor application: %w", err)
	}

	_, err = s.store.CreateMentorProfile(ctx, db.CreateMentorProfileParams{
		UserID:       mentor.ApplicantID,
		SpaceID:      mentor.SpaceID,
		Industry:     mentor.Industry,
		Company:      mentor.Company,
		Position:     mentor.Position,
		Experience:   mentor.Experience,
		Specialties:  mentor.Specialties,
		Description:  mentor.ApproachDescription,
		Availability: pqtype.NullRawMessage{RawMessage: mentor.Availability, Valid: true},
		Verified:     sql.NullBool{Bool: true, Valid: true},
	})

	if err != nil {
		return fmt.Errorf("failed to create mentor profile: %w", err)
	}

	
	details := fmt.Sprintf(`{"notes": "%s"}`, notes)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "approve_mentor_application",
		ResourceType: "mentor_application",
		ResourceID:   uuid.NullUUID{UUID: appID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return nil
}

func (s *Service) RejectMentorApplication(ctx context.Context, appID, adminID uuid.UUID, notes string) error {
	_, err := s.store.UpdateMentorApplicationStatus(ctx, db.UpdateMentorApplicationStatusParams{
		ID:            appID,
		Status:        sql.NullString{String: "rejected", Valid: true},
		ReviewedBy:    uuid.NullUUID{UUID: adminID, Valid: true},
		ReviewerNotes: sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to reject mentor application: %w", err)
	}

	
	details := fmt.Sprintf(`{"notes": "%s"}`, notes)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "reject_mentor_application",
		ResourceType: "mentor_application",
		ResourceID:   uuid.NullUUID{UUID: appID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return nil
}



func (s *Service) ResolveReport(ctx context.Context, reportID, adminID uuid.UUID, action, notes string) error {
	
	actionJSON := []byte(fmt.Sprintf(`{"action": "%s"}`, action))

	_, err := s.store.UpdateContentReportWithAction(ctx, db.UpdateContentReportWithActionParams{
		ID:              reportID,
		Status:          sql.NullString{String: "resolved", Valid: true},
		ReviewedBy:      uuid.NullUUID{UUID: adminID, Valid: true},
		ModerationNotes: sql.NullString{String: notes, Valid: notes != ""},
		ActionsTaken:    pqtype.NullRawMessage{RawMessage: actionJSON, Valid: action != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to resolve report: %w", err)
	}

	
	details := fmt.Sprintf(`{"action": "%s", "notes": "%s"}`, action, notes)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "resolve_report",
		ResourceType: "content_report",
		ResourceID:   uuid.NullUUID{UUID: reportID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return nil
}

func (s *Service) EscalateReport(ctx context.Context, reportID, adminID uuid.UUID) error {
	_, err := s.store.UpdateContentReportPriority(ctx, db.UpdateContentReportPriorityParams{
		ID:       reportID,
		Priority: sql.NullString{String: "urgent", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to escalate report: %w", err)
	}

	
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "escalate_report",
		ResourceType: "content_report",
		ResourceID:   uuid.NullUUID{UUID: reportID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
	})

	return nil
}



func (s *Service) GetGroups(ctx context.Context, spaceID uuid.UUID, status string, limit, offset int32) ([]GroupResponse, int64, error) {
	var groups []db.Group
	var err error

	if status != "" {
		groups, err = s.store.GetGroupsByStatus(ctx, db.GetGroupsByStatusParams{
			SpaceID: spaceID,
			Status:  sql.NullString{String: status, Valid: true},
			Limit:   limit,
			Offset:  offset,
		})
	} else {
		groups, err = s.store.GetGroupsBySpaceID(ctx, db.GetGroupsBySpaceIDParams{
			SpaceID: spaceID,
			Limit:   limit,
			Offset:  offset,
		})
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get groups: %w", err)
	}

	
	count := int64(len(groups))

	responses := make([]GroupResponse, len(groups))
	for i, group := range groups {
		responses[i] = GroupResponse{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description.String,
			Status:      group.Status.String,
			CreatedAt:   group.CreatedAt.Time,
		}
	}

	return responses, count, nil
}

func (s *Service) ApproveGroup(ctx context.Context, groupID, adminID uuid.UUID) error {
	err := s.store.UpdateGroupStatus(ctx, db.UpdateGroupStatusParams{
		ID:     groupID,
		Status: sql.NullString{String: "active", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to approve group: %w", err)
	}

	
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "approve_group",
		ResourceType: "group",
		ResourceID:   uuid.NullUUID{UUID: groupID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
	})

	return nil
}

func (s *Service) RejectGroup(ctx context.Context, groupID, adminID uuid.UUID, reason string) error {
	err := s.store.UpdateGroupStatus(ctx, db.UpdateGroupStatusParams{
		ID:     groupID,
		Status: sql.NullString{String: "rejected", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to reject group: %w", err)
	}

	
	details := fmt.Sprintf(`{"reason": "%s"}`, reason)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "reject_group",
		ResourceType: "group",
		ResourceID:   uuid.NullUUID{UUID: groupID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return nil
}

func (s *Service) DeleteGroup(ctx context.Context, groupID, adminID uuid.UUID) error {
	
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "delete_group",
		ResourceType: "group",
		ResourceID:   uuid.NullUUID{UUID: groupID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
	})

	
	err := s.store.DeleteGroup(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}



func (s *Service) CheckAdminPermission(ctx context.Context, userID uuid.UUID) (bool, error) {
	return s.store.CheckAdminPermission(ctx, userID)
}



func (s *Service) GetAllSettings(ctx context.Context) ([]db.SystemSetting, error) {
	settings, err := s.store.GetAllSystemSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}
	return settings, nil
}

func (s *Service) UpdateSetting(ctx context.Context, key string, value []byte, description string, updatedBy uuid.UUID) (*db.SystemSetting, error) {
	setting, err := s.store.UpsertSystemSetting(ctx, db.UpsertSystemSettingParams{
		Key:         key,
		Value:       value,
		Description: sql.NullString{String: description, Valid: description != ""},
		UpdatedBy:   uuid.NullUUID{UUID: updatedBy, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update setting: %w", err)
	}

	
	details := fmt.Sprintf(`{"key": "%s"}`, key)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  updatedBy,
		Action:       "update_system_setting",
		ResourceType: "system_setting",
		ResourceID:   uuid.NullUUID{UUID: setting.ID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &setting, nil
}



func (s *Service) GetUserGrowth(ctx context.Context, spaceID uuid.UUID, since time.Time) ([]db.GetUserGrowthDataRow, error) {
	data, err := s.store.GetUserGrowthData(ctx, db.GetUserGrowthDataParams{
		SpaceID:   spaceID,
		CreatedAt: sql.NullTime{Time: since, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user growth data: %w", err)
	}
	return data, nil
}

func (s *Service) GetContentGrowth(ctx context.Context, spaceID uuid.UUID, since time.Time) ([]db.GetContentGrowthDataRow, error) {
	data, err := s.store.GetContentGrowthData(ctx, db.GetContentGrowthDataParams{
		SpaceID:   spaceID,
		CreatedAt: sql.NullTime{Time: since, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get content growth data: %w", err)
	}
	return data, nil
}

func (s *Service) GetActivityStats(ctx context.Context, spaceID uuid.UUID, since time.Time) (*db.GetActivityStatsRow, error) {
	stats, err := s.store.GetActivityStats(ctx, db.GetActivityStatsParams{
		SpaceID:   spaceID,
		CreatedAt: since,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get activity stats: %w", err)
	}
	return &stats, nil
}



func (s *Service) GetAllAdmins(ctx context.Context, status string, limit, offset int32) ([]db.GetAllAdminUsersRow, error) {
	admins, err := s.store.GetAllAdminUsers(ctx, db.GetAllAdminUsersParams{
		Status: sql.NullString{String: status, Valid: status != ""},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get admins: %w", err)
	}
	return admins, nil
}

func (s *Service) UpdateAdminRole(ctx context.Context, userID, adminID uuid.UUID, roles []string) error {
	_, err := s.store.UpdateUserRole(ctx, db.UpdateUserRoleParams{
		ID:    userID,
		Roles: roles,
	})
	if err != nil {
		return fmt.Errorf("failed to update admin role: %w", err)
	}

	
	rolesJSON := fmt.Sprintf(`{"roles": %v}`, roles)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "update_admin_role",
		ResourceType: "user",
		ResourceID:   uuid.NullUUID{UUID: userID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(rolesJSON), Valid: true},
	})

	return nil
}

func (s *Service) UpdateAdminStatus(ctx context.Context, userID, adminID uuid.UUID, status string) error {
	err := s.store.UpdateUserAccountStatus(ctx, db.UpdateUserAccountStatusParams{
		ID:     userID,
		Status: sql.NullString{String: status, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update admin status: %w", err)
	}

	
	details := fmt.Sprintf(`{"status": "%s"}`, status)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "update_admin_status",
		ResourceType: "user",
		ResourceID:   uuid.NullUUID{UUID: userID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return nil
}



func (s *Service) GetNotifications(ctx context.Context, userID uuid.UUID, typeFilter, priority string, isRead *bool, limit, offset int32) ([]db.GetUserNotificationsRow, error) {
	
	
	notifications, err := s.store.GetUserNotifications(ctx, db.GetUserNotificationsParams{
		ToUserID: userID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	
	var filtered []db.GetUserNotificationsRow
	for _, n := range notifications {
		if typeFilter != "" && n.Type != typeFilter {
			continue
		}
		if priority != "" && (!n.Priority.Valid || n.Priority.String != priority) {
			continue
		}
		if isRead != nil && (!n.IsRead.Valid || n.IsRead.Bool != *isRead) {
			continue
		}
		filtered = append(filtered, n)
	}

	return filtered, nil
}

func (s *Service) MarkNotificationAsRead(ctx context.Context, notificationID uuid.UUID) error {
	err := s.store.MarkAsRead(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

func (s *Service) DeleteNotification(ctx context.Context, notificationID uuid.UUID) error {
	err := s.store.DeleteNotification(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

func (s *Service) MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error {
	err := s.store.MarkAllAsRead(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}



func (s *Service) GetAllCommunities(ctx context.Context, spaceID uuid.UUID, category, status string, limit, offset int32) ([]GetAllCommunitiesResponse, error) {
	communities, err := s.store.ListAllCommunitiesAdmin(ctx, db.ListAllCommunitiesAdminParams{
		SpaceID:  spaceID,
		Category: category,
		Status:   sql.NullString{String: status, Valid: status != ""},
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get communities: %w", err)
	}

	result := make([]GetAllCommunitiesResponse, len(communities))
	for i, community := range communities {
		result[i] = GetAllCommunitiesResponse{
			ID:                community.ID,
			SpaceID:           community.SpaceID,
			Name:              community.Name,
			Description:       community.Description.String,
			Category:          community.Category,
			CoverImage:        community.CoverImage.String,
			MemberCount:       community.MemberCount.Int32,
			Status:            community.Status.String,
			PostCount:         community.PostCount.Int32,
			IsPublic:          community.IsPublic.Bool,
			CreatedBy:         community.CreatedBy.UUID,
			Settings:          community.Settings.RawMessage,
			CreatedAt:         community.CreatedAt.Time,
			UpdatedAt:         community.UpdatedAt.Time,
			CreatedByUsername: community.CreatedByUsername.String,
			CreatedByFullName: community.CreatedByFullName.String,
			ActualMemberCount: community.ActualMemberCount,
			ActualPostCount:   community.ActualPostCount,
		}
	}
	return result, nil
}

func (s *Service) CreateCommunity(ctx context.Context, adminID uuid.UUID, req CreateCommunityRequest) (*db.Community, error) {
	
	
	
	
	
	
	
	

	community, err := s.store.CreateCommunity(ctx, db.CreateCommunityParams{
		SpaceID:     req.SpaceID,
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Category:    req.Category,
		CoverImage:  sql.NullString{String: req.CoverImage, Valid: req.CoverImage != ""},
		IsPublic:    sql.NullBool{Bool: req.IsPublic, Valid: true},
		CreatedBy:   uuid.NullUUID{UUID: adminID, Valid: true},
		Settings:    pqtype.NullRawMessage{RawMessage: req.Settings, Valid: len(req.Settings) > 0},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}

	
	details := fmt.Sprintf(`{"name": "%s", "category": "%s"}`, req.Name, req.Category)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "create_community",
		ResourceType: "community",
		ResourceID:   uuid.NullUUID{UUID: community.ID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &community, nil
}

func (s *Service) UpdateCommunity(ctx context.Context, communityID, adminID uuid.UUID, req UpdateCommunityRequest) (*db.Community, error) {
	community, err := s.store.UpdateCommunity(ctx, db.UpdateCommunityParams{
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		CoverImage:  sql.NullString{String: req.CoverImage, Valid: req.CoverImage != ""},
		Category:    req.Category,
		IsPublic:    sql.NullBool{Bool: req.IsPublic, Valid: true},
		Settings:    pqtype.NullRawMessage{RawMessage: req.Settings, Valid: len(req.Settings) > 0},
		ID:          communityID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update community: %w", err)
	}

	
	details := fmt.Sprintf(`{"name": "%s"}`, req.Name)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "update_community",
		ResourceType: "community",
		ResourceID:   uuid.NullUUID{UUID: communityID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &community, nil
}

func (s *Service) DeleteCommunity(ctx context.Context, communityID, adminID uuid.UUID) error {
	
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "delete_community",
		ResourceType: "community",
		ResourceID:   uuid.NullUUID{UUID: communityID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
	})

	err := s.store.DeleteCommunity(ctx, communityID)
	if err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}

	return nil
}

func (s *Service) UpdateCommunityStatus(ctx context.Context, communityID, adminID uuid.UUID, status string) (*db.Community, error) {
	community, err := s.store.UpdateCommunityStatus(ctx, db.UpdateCommunityStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		ID:     communityID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update community status: %w", err)
	}

	
	details := fmt.Sprintf(`{"status": "%s"}`, status)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "update_community_status",
		ResourceType: "community",
		ResourceID:   uuid.NullUUID{UUID: communityID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &community, nil
}

func (s *Service) AssignCommunityModerator(ctx context.Context, communityID, userID, adminID uuid.UUID, permissions []string) error {
	_, err := s.store.AddCommunityModerator(ctx, db.AddCommunityModeratorParams{
		CommunityID: communityID,
		UserID:      userID,
		Permissions: permissions,
	})
	if err != nil {
		return fmt.Errorf("failed to assign moderator: %w", err)
	}

	
	details := fmt.Sprintf(`{"user_id": "%s", "permissions": %v}`, userID.String(), permissions)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "assign_community_moderator",
		ResourceType: "community",
		ResourceID:   uuid.NullUUID{UUID: communityID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return nil
}



func (s *Service) GetAllAnnouncements(ctx context.Context, spaceID uuid.UUID, status, priority string, limit, offset int32) ([]CreateAnnouncementRequest, error) {
	announcements, err := s.store.ListAllAnnouncementsAdmin(ctx, db.ListAllAnnouncementsAdminParams{
		SpaceID:  spaceID,
		Status:   sql.NullString{String: status, Valid: status != ""},
		Priority: sql.NullString{String: priority, Valid: priority != ""},
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get announcements: %w", err)
	}

	result := make([]CreateAnnouncementRequest, len(announcements))

	for i, annoucement := range announcements {
		result[i] = CreateAnnouncementRequest{
			SpaceID:        annoucement.SpaceID,
			Title:          annoucement.Title,
			Content:        annoucement.Content,
			Type:           annoucement.Type,
			TargetAudience: annoucement.TargetAudience,
			Priority:       annoucement.Priority.String,
			ScheduledFor:   &annoucement.ScheduledFor.Time,
			ExpiresAt:      &annoucement.ExpiresAt.Time,
			Attachments:    annoucement.Attachments.RawMessage,
			IsPinned:       annoucement.IsPinned.Bool,
		}
	}
	return result, nil
}

func (s *Service) CreateAnnouncement(ctx context.Context, adminID uuid.UUID, req CreateAnnouncementRequest) (*db.Announcement, error) {
	announcement, err := s.store.CreateAnnouncement(ctx, db.CreateAnnouncementParams{
		SpaceID:        req.SpaceID,
		Title:          req.Title,
		Content:        req.Content,
		Type:           req.Type,
		TargetAudience: req.TargetAudience,
		Priority:       sql.NullString{String: req.Priority, Valid: req.Priority != ""},
		AuthorID:       uuid.NullUUID{UUID: adminID, Valid: true},
		ScheduledFor:   sql.NullTime{Time: *req.ScheduledFor, Valid: req.ScheduledFor != nil},
		ExpiresAt:      sql.NullTime{Time: *req.ExpiresAt, Valid: req.ExpiresAt != nil},
		Attachments:    pqtype.NullRawMessage{RawMessage: req.Attachments, Valid: len(req.Attachments) > 0},
		IsPinned:       sql.NullBool{Bool: req.IsPinned, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create announcement: %w", err)
	}

	
	details := fmt.Sprintf(`{"title": "%s", "type": "%s"}`, req.Title, req.Type)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "create_announcement",
		ResourceType: "announcement",
		ResourceID:   uuid.NullUUID{UUID: announcement.ID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &announcement, nil
}

func (s *Service) UpdateAnnouncement(ctx context.Context, announcementID, adminID uuid.UUID, req UpdateAnnouncementRequest) (*db.Announcement, error) {
	announcement, err := s.store.UpdateAnnouncement(ctx, db.UpdateAnnouncementParams{
		Title:          req.Title,
		Content:        req.Content,
		Type:           req.Type,
		TargetAudience: req.TargetAudience,
		Priority:       sql.NullString{String: req.Priority, Valid: req.Priority != ""},
		ScheduledFor:   sql.NullTime{Time: *req.ScheduledFor, Valid: req.ScheduledFor != nil},
		ExpiresAt:      sql.NullTime{Time: *req.ExpiresAt, Valid: req.ExpiresAt != nil},
		Attachments:    pqtype.NullRawMessage{RawMessage: req.Attachments, Valid: len(req.Attachments) > 0},
		IsPinned:       sql.NullBool{Bool: req.IsPinned, Valid: true},
		ID:             announcementID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update announcement: %w", err)
	}

	
	details := fmt.Sprintf(`{"title": "%s"}`, req.Title)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "update_announcement",
		ResourceType: "announcement",
		ResourceID:   uuid.NullUUID{UUID: announcementID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &announcement, nil
}

func (s *Service) DeleteAnnouncement(ctx context.Context, announcementID, adminID uuid.UUID) error {
	
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "delete_announcement",
		ResourceType: "announcement",
		ResourceID:   uuid.NullUUID{UUID: announcementID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
	})

	err := s.store.DeleteAnnouncement(ctx, announcementID)
	if err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	return nil
}

func (s *Service) UpdateAnnouncementStatus(ctx context.Context, announcementID, adminID uuid.UUID, status string) (*db.Announcement, error) {
	announcement, err := s.store.UpdateAnnouncementStatus(ctx, db.UpdateAnnouncementStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		ID:     announcementID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update announcement status: %w", err)
	}

	
	details := fmt.Sprintf(`{"status": "%s"}`, status)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "publish_announcement",
		ResourceType: "announcement",
		ResourceID:   uuid.NullUUID{UUID: announcementID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &announcement, nil
}



func (s *Service) GetAllEvents(ctx context.Context, spaceID uuid.UUID, status, category string, limit, offset int32) ([]CreateEventRequest, error) {
	events, err := s.store.ListAllEventsAdmin(ctx, db.ListAllEventsAdminParams{
		SpaceID:  spaceID,
		Status:   sql.NullString{String: status, Valid: status != ""},
		Category: category,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	response := make([]CreateEventRequest, len(events))
	for i, event := range events {
		response[i] = CreateEventRequest{
			SpaceID:              event.SpaceID,
			Title:                event.Title,
			Description:          event.Description.String,
			Category:             event.Category,
			Location:             event.Location.String,
			VenueDetails:         event.VenueDetails.String,
			StartDate:            event.StartDate,
			EndDate:              event.EndDate,
			Timezone:             event.Timezone.String,
			Tags:                 event.Tags,
			ImageURL:             event.ImageUrl.String,
			MaxAttendees:         event.MaxAttendees.Int32,
			RegistrationRequired: event.RegistrationRequired.Bool,
			RegistrationDeadline: &event.RegistrationDeadline.Time,
			IsPublic:             event.IsPublic.Bool,
		}
	}
	return response, nil
}

func (s *Service) CreateEvent(ctx context.Context, adminID uuid.UUID, req CreateEventRequest) (*CreateEventRequest, error) {
	event, err := s.store.CreateEvent(ctx, db.CreateEventParams{
		SpaceID:              req.SpaceID,
		Title:                req.Title,
		Description:          sql.NullString{String: req.Description, Valid: req.Description != ""},
		Category:             req.Category,
		Location:             sql.NullString{String: req.Location, Valid: req.Location != ""},
		VenueDetails:         sql.NullString{String: req.VenueDetails, Valid: req.VenueDetails != ""},
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		Timezone:             sql.NullString{String: req.Timezone, Valid: req.Timezone != ""},
		Organizer:            uuid.NullUUID{UUID: adminID, Valid: true},
		Tags:                 req.Tags,
		ImageUrl:             sql.NullString{String: req.ImageURL, Valid: req.ImageURL != ""},
		MaxAttendees:         sql.NullInt32{Int32: req.MaxAttendees, Valid: req.MaxAttendees > 0},
		RegistrationRequired: sql.NullBool{Bool: req.RegistrationRequired, Valid: true},
		RegistrationDeadline: sql.NullTime{Time: *req.RegistrationDeadline, Valid: req.RegistrationDeadline != nil},
		IsPublic:             sql.NullBool{Bool: req.IsPublic, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	
	details := fmt.Sprintf(`{"title": "%s", "category": "%s"}`, req.Title, req.Category)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "create_event",
		ResourceType: "event",
		ResourceID:   uuid.NullUUID{UUID: event.ID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	resp := CreateEventRequest{
		SpaceID:              event.SpaceID,
		Title:                event.Title,
		Description:          event.Description.String,
		Category:             event.Category,
		Location:             event.Location.String,
		VenueDetails:         event.VenueDetails.String,
		StartDate:            event.StartDate,
		EndDate:              event.EndDate,
		Timezone:             event.Timezone.String,
		Tags:                 event.Tags,
		ImageURL:             event.ImageUrl.String,
		MaxAttendees:         event.MaxAttendees.Int32,
		RegistrationRequired: event.RegistrationRequired.Bool,
		RegistrationDeadline: &event.RegistrationDeadline.Time,
		IsPublic:             event.IsPublic.Bool,
	}

	return &resp, nil
}

func (s *Service) UpdateEvent(ctx context.Context, eventID, adminID uuid.UUID, req UpdateEventRequest) (*db.Event, error) {
	event, err := s.store.UpdateEvent(ctx, db.UpdateEventParams{
		Title:                req.Title,
		Description:          sql.NullString{String: req.Description, Valid: req.Description != ""},
		Category:             req.Category,
		Location:             sql.NullString{String: req.Location, Valid: req.Location != ""},
		VenueDetails:         sql.NullString{String: req.VenueDetails, Valid: req.VenueDetails != ""},
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		Timezone:             sql.NullString{String: req.Timezone, Valid: req.Timezone != ""},
		Tags:                 req.Tags,
		ImageUrl:             sql.NullString{String: req.ImageURL, Valid: req.ImageURL != ""},
		MaxAttendees:         sql.NullInt32{Int32: req.MaxAttendees, Valid: req.MaxAttendees > 0},
		RegistrationRequired: sql.NullBool{Bool: req.RegistrationRequired, Valid: true},
		RegistrationDeadline: sql.NullTime{Time: *req.RegistrationDeadline, Valid: req.RegistrationDeadline != nil},
		IsPublic:             sql.NullBool{Bool: req.IsPublic, Valid: true},
		ID:                   eventID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	
	details := fmt.Sprintf(`{"title": "%s"}`, req.Title)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "update_event",
		ResourceType: "event",
		ResourceID:   uuid.NullUUID{UUID: eventID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &event, nil
}

func (s *Service) DeleteEvent(ctx context.Context, eventID, adminID uuid.UUID) error {
	
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "delete_event",
		ResourceType: "event",
		ResourceID:   uuid.NullUUID{UUID: eventID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
	})

	err := s.store.DeleteEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func (s *Service) UpdateEventStatus(ctx context.Context, eventID, adminID uuid.UUID, status string) (*db.Event, error) {
	event, err := s.store.UpdateEventStatus(ctx, db.UpdateEventStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		ID:     eventID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update event status: %w", err)
	}

	
	details := fmt.Sprintf(`{"status": "%s"}`, status)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "cancel_event",
		ResourceType: "event",
		ResourceID:   uuid.NullUUID{UUID: eventID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &event, nil
}

func (s *Service) GetEventRegistrations(ctx context.Context, eventID uuid.UUID) ([]db.GetEventAttendeesRow, error) {
	attendees, err := s.store.GetEventAttendees(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event registrations: %w", err)
	}
	return attendees, nil
}



func (s *Service) CreateUser(ctx context.Context, adminID uuid.UUID, req CreateUserRequest) (*db.User, error) {
	
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		SpaceID:     req.SpaceID,
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
		FullName:    req.FullName,
		Roles:       req.Roles,
		Department:  sql.NullString{String: req.Department, Valid: req.Department != ""},
		Level:       sql.NullString{String: req.Level, Valid: req.Level != ""},
		PhoneNumber: "",
		Interests:   []string{},
		Major:       sql.NullString{},
		Year:        sql.NullInt32{},
		Settings:    pqtype.NullRawMessage{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	
	details := fmt.Sprintf(`{"username": "%s", "email": "%s"}`, req.Username, req.Email)
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "create_user",
		ResourceType: "user",
		ResourceID:   uuid.NullUUID{UUID: user.ID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &user, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID, adminID uuid.UUID, req UpdateUserRequest) (*db.User, error) {
	
	fullName := ""
	if req.FullName != nil {
		fullName = *req.FullName
	}

	user, err := s.store.UpdateUser(ctx, db.UpdateUserParams{
		FullName:   fullName,
		Bio:        sql.NullString{},
		Avatar:     sql.NullString{},
		Department: sql.NullString{String: *req.Department, Valid: req.Department != nil},
		Level:      sql.NullString{String: *req.Level, Valid: req.Level != nil},
		Major:      sql.NullString{},
		Year:       sql.NullInt32{},
		Interests:  []string{},
		Settings:   pqtype.NullRawMessage{},
		ID:         userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	
	if len(req.Roles) > 0 {
		_, err = s.store.UpdateUserRole(ctx, db.UpdateUserRoleParams{
			ID:    userID,
			Roles: req.Roles,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update user roles: %w", err)
		}
	}

	
	if req.Status != nil {
		err = s.store.UpdateUserAccountStatus(ctx, db.UpdateUserAccountStatusParams{
			ID:     userID,
			Status: sql.NullString{String: *req.Status, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update user status: %w", err)
		}
	}

	
	details := fmt.Sprintf(`{"user_id": "%s"}`, userID.String())
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "update_user",
		ResourceType: "user",
		ResourceID:   uuid.NullUUID{UUID: userID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(details), Valid: true},
	})

	return &user, nil
}

func (s *Service) ResetUserPassword(ctx context.Context, userID, adminID uuid.UUID, newPassword string) error {
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.store.ResetUserPassword(ctx, db.ResetUserPasswordParams{
		ID:       userID,
		Password: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	
	_, _ = s.store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminUserID:  adminID,
		Action:       "reset_user_password",
		ResourceType: "user",
		ResourceID:   uuid.NullUUID{UUID: userID, Valid: true},
		Details:      pqtype.NullRawMessage{RawMessage: []byte(`{}`), Valid: true},
	})

	return nil
}


func hashPassword(password string) (string, error) {
	
	
	return "hashed_" + password, nil
}
