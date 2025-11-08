package analytics

import (
	"context"
	"database/sql"
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






func (s *Service) CreateReport(ctx context.Context, req CreateReportRequest) (*ReportResponse, error) {
	report, err := s.store.CreateReport(ctx, db.CreateReportParams{
		SpaceID:     req.SpaceID,
		ReporterID:  req.ReporterID,
		ContentType: req.ContentType,
		ContentID:   req.ContentID,
		Reason:      req.Reason,
		Description: sqlNullString(req.Description),
		Priority:    sqlNullString(req.Priority),
	})
	if err != nil {
		return nil, err
	}

	return reportToResponse(report), nil
}


func (s *Service) GetReport(ctx context.Context, reportID uuid.UUID) (*ReportDetailResponse, error) {
	report, err := s.store.GetReport(ctx, reportID)
	if err != nil {
		return nil, err
	}

	return &ReportDetailResponse{
		ID:               report.ID,
		SpaceID:          report.SpaceID,
		ReporterID:       report.ReporterID,
		ReporterUsername: report.ReporterUsername,
		ReporterFullName: report.ReporterFullName,
		ContentType:      report.ContentType,
		ContentID:        report.ContentID,
		Reason:           report.Reason,
		Description:      nullStringToPtr(report.Description),
		Status:           nullStringWithDefault(report.Status, "pending"),
		Priority:         nullStringToPtr(report.Priority),
		ReviewedBy:       nullUUIDToPtr(report.ReviewedBy),
		ReviewerUsername: nullStringToPtr(report.ReviewerUsername),
		ReviewedAt:       nullTimeToPtr(report.ReviewedAt),
		ModerationNotes:  nullStringToPtr(report.ModerationNotes),
		ActionsTaken:     nullRawMessageToPtr(report.ActionsTaken),
		CreatedAt:        nullTimeWithDefault(report.CreatedAt),
		UpdatedAt:        nullTimeWithDefault(report.UpdatedAt),
	}, nil
}


func (s *Service) GetReportsByContent(ctx context.Context, contentType string, contentID uuid.UUID) ([]ReportResponse, error) {
	reports, err := s.store.GetReportsByContent(ctx, db.GetReportsByContentParams{
		ContentType: contentType,
		ContentID:   contentID,
	})
	if err != nil {
		return nil, err
	}

	responses := make([]ReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = *reportToResponse(report)
	}

	return responses, nil
}


func (s *Service) GetModerationQueue(ctx context.Context, spaceID uuid.UUID, page, limit int32) ([]ReportDetailResponse, error) {
	if limit == 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	reports, err := s.store.GetModerationQueue(ctx, db.GetModerationQueueParams{
		SpaceID: spaceID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}

	responses := make([]ReportDetailResponse, len(reports))
	for i, report := range reports {
		responses[i] = ReportDetailResponse{
			ID:               report.ID,
			SpaceID:          report.SpaceID,
			ReporterID:       report.ReporterID,
			ReporterUsername: report.ReporterUsername,
			ReporterFullName: report.ReporterFullName,
			ContentType:      report.ContentType,
			ContentID:        report.ContentID,
			Reason:           report.Reason,
			Description:      nullStringToPtr(report.Description),
			Status:           nullStringWithDefault(report.Status, "pending"),
			Priority:         nullStringToPtr(report.Priority),
			ReviewedBy:       nullUUIDToPtr(report.ReviewedBy),
			ReviewerUsername: nullStringToPtr(report.ReviewerUsername),
			ReviewedAt:       nullTimeToPtr(report.ReviewedAt),
			ModerationNotes:  nullStringToPtr(report.ModerationNotes),
			ActionsTaken:     nullRawMessageToPtr(report.ActionsTaken),
			CreatedAt:        nullTimeWithDefault(report.CreatedAt),
			UpdatedAt:        nullTimeWithDefault(report.UpdatedAt),
		}
	}

	return responses, nil
}


func (s *Service) GetPendingReports(ctx context.Context, spaceID uuid.UUID) ([]ReportDetailResponse, error) {
	reports, err := s.store.GetPendingReports(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]ReportDetailResponse, len(reports))
	for i, report := range reports {
		responses[i] = ReportDetailResponse{
			ID:               report.ID,
			SpaceID:          report.SpaceID,
			ReporterID:       report.ReporterID,
			ReporterUsername: report.ReporterUsername,
			ReporterFullName: report.ReporterFullName,
			ContentType:      report.ContentType,
			ContentID:        report.ContentID,
			Reason:           report.Reason,
			Description:      nullStringToPtr(report.Description),
			Status:           nullStringWithDefault(report.Status, "pending"),
			Priority:         nullStringToPtr(report.Priority),
			ReviewedBy:       nullUUIDToPtr(report.ReviewedBy),
			ReviewedAt:       nullTimeToPtr(report.ReviewedAt),
			ModerationNotes:  nullStringToPtr(report.ModerationNotes),
			ActionsTaken:     nullRawMessageToPtr(report.ActionsTaken),
			CreatedAt:        nullTimeWithDefault(report.CreatedAt),
			UpdatedAt:        nullTimeWithDefault(report.UpdatedAt),
		}
	}

	return responses, nil
}


func (s *Service) UpdateReport(ctx context.Context, reportID, reviewerID uuid.UUID, req UpdateReportRequest) (*ReportResponse, error) {
	report, err := s.store.UpdateReport(ctx, db.UpdateReportParams{
		ID:              reportID,
		Status:          sqlNullString(&req.Status),
		ReviewedBy:      uuid.NullUUID{UUID: reviewerID, Valid: true},
		ModerationNotes: sqlNullString(req.ModerationNotes),
		ActionsTaken:    sqlNullRawMessage(req.ActionsTaken),
	})
	if err != nil {
		return nil, err
	}

	return reportToResponse(report), nil
}


func (s *Service) GetContentModerationStats(ctx context.Context, spaceID uuid.UUID) (*ContentModerationStatsResponse, error) {
	stats, err := s.store.GetContentModerationStats(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	return &ContentModerationStatsResponse{
		TotalReports:    stats.TotalReports,
		PendingReports:  stats.PendingReports,
		ApprovedReports: stats.ApprovedReports,
		RejectedReports: stats.RejectedReports,
		UrgentReports:   stats.UrgentReports,
	}, nil
}






func (s *Service) GetSystemMetrics(ctx context.Context, spaceID uuid.UUID) (*SystemMetricsResponse, error) {
	metrics, err := s.store.GetSystemMetrics(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	return &SystemMetricsResponse{
		TotalUsers:                metrics.TotalUsers,
		ActiveUsers:               metrics.ActiveUsers,
		NewUsersToday:             metrics.NewUsersToday,
		DailyPosts:                metrics.DailyPosts,
		TotalGroups:               metrics.TotalGroups,
		TotalCommunities:          metrics.TotalCommunities,
		TotalEvents:               metrics.TotalEvents,
		PendingTutoringSessions:   metrics.PendingTutoringSessions,
		PendingMentoringSessions:  metrics.PendingMentoringSessions,
		PendingReports:            metrics.PendingReports,
		PendingTutorApplications:  metrics.PendingTutorApplications,
		PendingMentorApplications: metrics.PendingMentorApplications,
	}, nil
}


func (s *Service) GetSpaceStats(ctx context.Context, spaceID uuid.UUID) (*SpaceStatsResponse, error) {
	stats, err := s.store.GetSpaceStats(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	return &SpaceStatsResponse{
		Name:           stats.Name,
		Slug:           stats.Slug,
		UserCount:      stats.UserCount,
		PostCount:      stats.PostCount,
		CommunityCount: stats.CommunityCount,
		GroupCount:     stats.GroupCount,
		FirstUserDate:  nullTimeToPtr(stats.FirstUserDate),
	}, nil
}






func (s *Service) GetEngagementMetrics(ctx context.Context, spaceID uuid.UUID) ([]EngagementMetricResponse, error) {
	metrics, err := s.store.GetEngagementMetrics(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]EngagementMetricResponse, len(metrics))
	for i, metric := range metrics {
		responses[i] = EngagementMetricResponse{
			Date:          metric.Date,
			PostCount:     metric.PostCount,
			TotalLikes:    metric.TotalLikes,
			TotalComments: metric.TotalComments,
			TotalViews:    metric.TotalViews,
		}
	}

	return responses, nil
}


func (s *Service) GetUserActivityStats(ctx context.Context, spaceID uuid.UUID) ([]UserActivityStatResponse, error) {
	stats, err := s.store.GetUserActivityStats(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]UserActivityStatResponse, len(stats))
	for i, stat := range stats {
		responses[i] = UserActivityStatResponse{
			Action: stat.Action,
			Count:  stat.Count,
			Date:   stat.Date,
		}
	}

	return responses, nil
}


func (s *Service) GetUserGrowth(ctx context.Context, spaceID uuid.UUID) ([]UserGrowthResponse, error) {
	growth, err := s.store.GetUserGrowth(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]UserGrowthResponse, len(growth))
	for i, g := range growth {
		responses[i] = UserGrowthResponse{
			Date:     g.Date,
			NewUsers: g.NewUsers,
		}
	}

	return responses, nil
}


func (s *Service) GetUserEngagementRanking(ctx context.Context, spaceID uuid.UUID) ([]UserEngagementRankingResponse, error) {
	ranking, err := s.store.GetUserEngagementRanking(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]UserEngagementRankingResponse, len(ranking))
	for i, user := range ranking {
		responses[i] = UserEngagementRankingResponse{
			ID:              user.ID,
			Username:        user.Username,
			FullName:        user.FullName,
			Avatar:          nullStringToPtr(user.Avatar),
			PostCount:       user.PostCount,
			FollowersCount:  user.FollowersCount,
			FollowingCount:  user.FollowingCount,
			EngagementScore: user.EngagementScore,
		}
	}

	return responses, nil
}






func (s *Service) GetTopPosts(ctx context.Context, spaceID uuid.UUID) ([]TopPostResponse, error) {
	posts, err := s.store.GetTopPosts(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]TopPostResponse, len(posts))
	for i, post := range posts {
		responses[i] = TopPostResponse{
			ID:              post.ID,
			AuthorID:        post.AuthorID,
			Username:        post.Username,
			FullName:        post.FullName,
			SpaceID:         post.SpaceID,
			CommunityID:     nullUUIDToPtr(post.CommunityID),
			GroupID:         nullUUIDToPtr(post.GroupID),
			Content:         post.Content,
			Media:           nullRawMessageToPtr(post.Media),
			Tags:            post.Tags,
			LikesCount:      nullInt32WithDefault(post.LikesCount, 0),
			CommentsCount:   nullInt32WithDefault(post.CommentsCount, 0),
			ViewsCount:      nullInt32WithDefault(post.ViewsCount, 0),
			EngagementScore: post.EngagementScore,
			CreatedAt:       nullTimeWithDefault(post.CreatedAt),
		}
	}

	return responses, nil
}


func (s *Service) GetTopCommunities(ctx context.Context, spaceID uuid.UUID) ([]TopCommunityResponse, error) {
	communities, err := s.store.GetTopCommunities(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]TopCommunityResponse, len(communities))
	for i, community := range communities {
		responses[i] = TopCommunityResponse{
			ID:              community.ID,
			SpaceID:         community.SpaceID,
			Name:            community.Name,
			Description:     nullStringToPtr(community.Description),
			Category:        community.Category,
			CoverImage:      nullStringToPtr(community.CoverImage),
			MemberCount:     nullInt32WithDefault(community.MemberCount, 0),
			PostCount:       nullInt32WithDefault(community.PostCount, 0),
			EngagementScore: community.EngagementScore,
			IsPublic:        nullBoolWithDefault(community.IsPublic, true),
			CreatedAt:       nullTimeWithDefault(community.CreatedAt),
		}
	}

	return responses, nil
}


func (s *Service) GetTopGroups(ctx context.Context, spaceID uuid.UUID) ([]TopGroupResponse, error) {
	groups, err := s.store.GetTopGroups(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]TopGroupResponse, len(groups))
	for i, group := range groups {
		responses[i] = TopGroupResponse{
			ID:              group.ID,
			SpaceID:         group.SpaceID,
			CommunityID:     nullUUIDToPtr(group.CommunityID),
			Name:            group.Name,
			Description:     nullStringToPtr(group.Description),
			Category:        group.Category,
			Avatar:          nullStringToPtr(group.Avatar),
			MemberCount:     nullInt32WithDefault(group.MemberCount, 0),
			PostCount:       nullInt32WithDefault(group.PostCount, 0),
			EngagementScore: group.EngagementScore,
			Visibility:      nullStringWithDefault(group.Visibility, "public"),
			CreatedAt:       nullTimeWithDefault(group.CreatedAt),
		}
	}

	return responses, nil
}






func (s *Service) GetMentoringStats(ctx context.Context, spaceID uuid.UUID) (*MentoringStatsResponse, error) {
	stats, err := s.store.GetMentoringStats(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	avgRating := 0.0
	if rating, ok := stats.AverageRating.(float64); ok {
		avgRating = rating
	}

	return &MentoringStatsResponse{
		TotalSessions:     stats.TotalSessions,
		CompletedSessions: stats.CompletedSessions,
		PendingSessions:   stats.PendingSessions,
		AverageRating:     avgRating,
		RatedSessions:     stats.RatedSessions,
	}, nil
}


func (s *Service) GetTutoringStats(ctx context.Context, spaceID uuid.UUID) (*TutoringStatsResponse, error) {
	stats, err := s.store.GetTutoringStats(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	avgRating := 0.0
	if rating, ok := stats.AverageRating.(float64); ok {
		avgRating = rating
	}

	return &TutoringStatsResponse{
		TotalSessions:     stats.TotalSessions,
		CompletedSessions: stats.CompletedSessions,
		PendingSessions:   stats.PendingSessions,
		AverageRating:     avgRating,
		RatedSessions:     stats.RatedSessions,
	}, nil
}


func (s *Service) GetPopularIndustries(ctx context.Context, spaceID uuid.UUID) ([]PopularIndustryResponse, error) {
	industries, err := s.store.GetPopularIndustries(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]PopularIndustryResponse, len(industries))
	for i, industry := range industries {
		responses[i] = PopularIndustryResponse{
			Industry:      industry.Industry,
			SessionCount:  industry.SessionCount,
			AverageRating: industry.AverageRating,
		}
	}

	return responses, nil
}


func (s *Service) GetPopularSubjects(ctx context.Context, spaceID uuid.UUID) ([]PopularSubjectResponse, error) {
	subjects, err := s.store.GetPopularSubjects(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	responses := make([]PopularSubjectResponse, len(subjects))
	for i, subject := range subjects {
		responses[i] = PopularSubjectResponse{
			Subject:       subject.Subject,
			SessionCount:  subject.SessionCount,
			AverageRating: subject.AverageRating,
		}
	}

	return responses, nil
}





func reportToResponse(report db.Report) *ReportResponse {
	return &ReportResponse{
		ID:              report.ID,
		SpaceID:         report.SpaceID,
		ReporterID:      report.ReporterID,
		ContentType:     report.ContentType,
		ContentID:       report.ContentID,
		Reason:          report.Reason,
		Description:     nullStringToPtr(report.Description),
		Status:          nullStringWithDefault(report.Status, "pending"),
		Priority:        nullStringToPtr(report.Priority),
		ReviewedBy:      nullUUIDToPtr(report.ReviewedBy),
		ReviewedAt:      nullTimeToPtr(report.ReviewedAt),
		ModerationNotes: nullStringToPtr(report.ModerationNotes),
		ActionsTaken:    nullRawMessageToPtr(report.ActionsTaken),
		CreatedAt:       nullTimeWithDefault(report.CreatedAt),
		UpdatedAt:       nullTimeWithDefault(report.UpdatedAt),
	}
}

func sqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func sqlNullRawMessage(r *pqtype.NullRawMessage) pqtype.NullRawMessage {
	if r == nil {
		return pqtype.NullRawMessage{Valid: false}
	}
	return *r
}

func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func nullStringWithDefault(ns sql.NullString, def string) string {
	if !ns.Valid {
		return def
	}
	return ns.String
}

func nullInt32ToPtr(ni sql.NullInt32) *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}

func nullInt32WithDefault(ni sql.NullInt32, def int32) int32 {
	if !ni.Valid {
		return def
	}
	return ni.Int32
}

func nullBoolWithDefault(nb sql.NullBool, def bool) bool {
	if !nb.Valid {
		return def
	}
	return nb.Bool
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func nullTimeWithDefault(nt sql.NullTime) time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

func nullUUIDToPtr(nu uuid.NullUUID) *uuid.UUID {
	if !nu.Valid {
		return nil
	}
	return &nu.UUID
}

func nullRawMessageToPtr(nr pqtype.NullRawMessage) *pqtype.NullRawMessage {
	if !nr.Valid {
		return nil
	}
	return &nr
}
