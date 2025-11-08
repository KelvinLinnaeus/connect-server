



package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Querier interface {
	AddCommunityModerator(ctx context.Context, arg AddCommunityModeratorParams) (CommunityMember, error)
	AddConversationParticipants(ctx context.Context, arg AddConversationParticipantsParams) error
	AddEventCoOrganizer(ctx context.Context, arg AddEventCoOrganizerParams) (EventAttendee, error)
	AddGroupAdmin(ctx context.Context, arg AddGroupAdminParams) (GroupMember, error)
	AddGroupModerator(ctx context.Context, arg AddGroupModeratorParams) (GroupMember, error)
	AddMentoringSessionMeetingLink(ctx context.Context, arg AddMentoringSessionMeetingLinkParams) error
	AddMessageReaction(ctx context.Context, arg AddMessageReactionParams) error
	AddSessionMeetingLink(ctx context.Context, arg AddSessionMeetingLinkParams) error
	
	AdminCreateUser(ctx context.Context, arg AdminCreateUserParams) (User, error)
	AdminUpdateUser(ctx context.Context, arg AdminUpdateUserParams) (User, error)
	AdvancedSearchPosts(ctx context.Context, arg AdvancedSearchPostsParams) ([]AdvancedSearchPostsRow, error)
	AdvancedSearchUsers(ctx context.Context, arg AdvancedSearchUsersParams) ([]AdvancedSearchUsersRow, error)
	ApplyForProjectRole(ctx context.Context, arg ApplyForProjectRoleParams) (GroupApplication, error)
	CheckAdminPermission(ctx context.Context, id uuid.UUID) (bool, error)
	CheckIfFollowing(ctx context.Context, arg CheckIfFollowingParams) (bool, error)
	CleanupOldLoginAttempts(ctx context.Context, attemptedAt time.Time) error
	CountRecentFailedLoginAttemptsByIP(ctx context.Context, arg CountRecentFailedLoginAttemptsByIPParams) (int64, error)
	CountRecentFailedLoginAttemptsByUsername(ctx context.Context, arg CountRecentFailedLoginAttemptsByUsernameParams) (int64, error)
	CreateAnnouncement(ctx context.Context, arg CreateAnnouncementParams) (Announcement, error)
	
	CreateAuditLog(ctx context.Context, arg CreateAuditLogParams) (AuditLog, error)
	CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error)
	
	CreateCommunity(ctx context.Context, arg CreateCommunityParams) (Community, error)
	
	
	CreateContentReport(ctx context.Context, arg CreateContentReportParams) (Report, error)
	
	CreateConversation(ctx context.Context, arg CreateConversationParams) (Conversation, error)
	
	CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error)
	
	CreateGroup(ctx context.Context, arg CreateGroupParams) (Group, error)
	CreateLoginAttempt(ctx context.Context, arg CreateLoginAttemptParams) (LoginAttempt, error)
	CreateMentorApplication(ctx context.Context, arg CreateMentorApplicationParams) (MentorApplication, error)
	CreateMentorProfile(ctx context.Context, arg CreateMentorProfileParams) (MentorProfile, error)
	CreateMentoringSession(ctx context.Context, arg CreateMentoringSessionParams) (MentoringSession, error)
	CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error)
	
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateProjectRole(ctx context.Context, arg CreateProjectRoleParams) (GroupRole, error)
	CreateReport(ctx context.Context, arg CreateReportParams) (Report, error)
	CreateRepost(ctx context.Context, arg CreateRepostParams) (Post, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (UserSession, error)
	CreateSpace(ctx context.Context, arg CreateSpaceParams) (Space, error)
	
	CreateSpaceActivity(ctx context.Context, arg CreateSpaceActivityParams) (SpaceActivity, error)
	
	CreateTutorApplication(ctx context.Context, arg CreateTutorApplicationParams) (TutorApplication, error)
	CreateTutorProfile(ctx context.Context, arg CreateTutorProfileParams) (TutorProfile, error)
	CreateTutoringSession(ctx context.Context, arg CreateTutoringSessionParams) (TutoringSession, error)
	
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	
	CreateUserSuspension(ctx context.Context, arg CreateUserSuspensionParams) (UserSuspension, error)
	DeactivateUser(ctx context.Context, id uuid.UUID) error
	DecrementFollowersCount(ctx context.Context, id uuid.UUID) error
	DecrementFollowingCount(ctx context.Context, id uuid.UUID) error
	DeleteAnnouncement(ctx context.Context, id uuid.UUID) error
	DeleteCommunity(ctx context.Context, id uuid.UUID) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	DeleteGroup(ctx context.Context, id uuid.UUID) error
	DeleteMessage(ctx context.Context, arg DeleteMessageParams) error
	DeleteNotification(ctx context.Context, id uuid.UUID) error
	DeletePost(ctx context.Context, arg DeletePostParams) error
	DeleteSpace(ctx context.Context, id uuid.UUID) error
	DeleteSystemSetting(ctx context.Context, key string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	
	FollowUser(ctx context.Context, arg FollowUserParams) (Follow, error)
	GetActiveSuspension(ctx context.Context, userID uuid.UUID) (UserSuspension, error)
	GetActivityStats(ctx context.Context, arg GetActivityStatsParams) (GetActivityStatsRow, error)
	
	GetAdminDashboardStats(ctx context.Context, spaceID uuid.UUID) (GetAdminDashboardStatsRow, error)
	
	GetAdminNotifications(ctx context.Context, arg GetAdminNotificationsParams) ([]GetAdminNotificationsRow, error)
	
	
	GetAllAdminUsers(ctx context.Context, arg GetAllAdminUsersParams) ([]GetAllAdminUsersRow, error)
	GetAllMentorApplications(ctx context.Context, arg GetAllMentorApplicationsParams) ([]GetAllMentorApplicationsRow, error)
	
	GetAllSpacesWithStats(ctx context.Context) ([]GetAllSpacesWithStatsRow, error)
	GetAllSystemSettings(ctx context.Context) ([]SystemSetting, error)
	
	GetAllTutorApplications(ctx context.Context, arg GetAllTutorApplicationsParams) ([]GetAllTutorApplicationsRow, error)
	GetAnnouncementByID(ctx context.Context, id uuid.UUID) (GetAnnouncementByIDRow, error)
	GetAuditLogs(ctx context.Context, arg GetAuditLogsParams) ([]GetAuditLogsRow, error)
	GetCommunityAdmins(ctx context.Context, communityID uuid.UUID) ([]GetCommunityAdminsRow, error)
	GetCommunityByID(ctx context.Context, arg GetCommunityByIDParams) (GetCommunityByIDRow, error)
	GetCommunityBySlug(ctx context.Context, arg GetCommunityBySlugParams) (GetCommunityBySlugRow, error)
	GetCommunityCategories(ctx context.Context, spaceID uuid.UUID) ([]string, error)
	GetCommunityMembers(ctx context.Context, communityID uuid.UUID) ([]GetCommunityMembersRow, error)
	GetCommunityModerators(ctx context.Context, communityID uuid.UUID) ([]GetCommunityModeratorsRow, error)
	GetCommunityPosts(ctx context.Context, arg GetCommunityPostsParams) ([]GetCommunityPostsRow, error)
	GetContentGrowthData(ctx context.Context, arg GetContentGrowthDataParams) ([]GetContentGrowthDataRow, error)
	GetContentModerationStats(ctx context.Context, spaceID uuid.UUID) (GetContentModerationStatsRow, error)
	GetContentReportByID(ctx context.Context, id uuid.UUID) (GetContentReportByIDRow, error)
	GetContentReports(ctx context.Context, arg GetContentReportsParams) ([]GetContentReportsRow, error)
	GetConversationByID(ctx context.Context, arg GetConversationByIDParams) (GetConversationByIDRow, error)
	GetConversationByParticipants(ctx context.Context, arg GetConversationByParticipantsParams) (uuid.UUID, error)
	GetConversationMessages(ctx context.Context, arg GetConversationMessagesParams) ([]GetConversationMessagesRow, error)
	GetConversationParticipants(ctx context.Context, conversationID uuid.UUID) ([]GetConversationParticipantsRow, error)
	GetEngagementMetrics(ctx context.Context, spaceID uuid.UUID) ([]GetEngagementMetricsRow, error)
	GetEventAttendees(ctx context.Context, eventID uuid.UUID) ([]GetEventAttendeesRow, error)
	GetEventByID(ctx context.Context, arg GetEventByIDParams) (GetEventByIDRow, error)
	GetEventCategories(ctx context.Context, spaceID uuid.UUID) ([]string, error)
	GetEventCoOrganizers(ctx context.Context, eventID uuid.UUID) ([]GetEventCoOrganizersRow, error)
	GetEventWithRegistrations(ctx context.Context, id uuid.UUID) (GetEventWithRegistrationsRow, error)
	GetGroupByID(ctx context.Context, arg GetGroupByIDParams) (GetGroupByIDRow, error)
	GetGroupJoinRequests(ctx context.Context, groupID uuid.UUID) ([]GetGroupJoinRequestsRow, error)
	GetGroupPosts(ctx context.Context, arg GetGroupPostsParams) ([]GetGroupPostsRow, error)
	
	GetGroupsBySpaceID(ctx context.Context, arg GetGroupsBySpaceIDParams) ([]Group, error)
	GetGroupsByStatus(ctx context.Context, arg GetGroupsByStatusParams) ([]Group, error)
	GetLockedUsers(ctx context.Context) ([]GetLockedUsersRow, error)
	GetLoginAttemptsWithSessions(ctx context.Context, arg GetLoginAttemptsWithSessionsParams) ([]GetLoginAttemptsWithSessionsRow, error)
	GetMentorApplication(ctx context.Context, id uuid.UUID) (GetMentorApplicationRow, error)
	GetMentorProfile(ctx context.Context, userID uuid.UUID) (GetMentorProfileRow, error)
	GetMentorReviews(ctx context.Context, arg GetMentorReviewsParams) ([]GetMentorReviewsRow, error)
	GetMentoringSession(ctx context.Context, id uuid.UUID) (GetMentoringSessionRow, error)
	GetMentoringStats(ctx context.Context, spaceID uuid.UUID) (GetMentoringStatsRow, error)
	GetMessageByID(ctx context.Context, id uuid.UUID) (GetMessageByIDRow, error)
	GetModerationQueue(ctx context.Context, arg GetModerationQueueParams) ([]GetModerationQueueRow, error)
	GetNotification(ctx context.Context, id uuid.UUID) (Notification, error)
	GetOrCreateDirectConversation(ctx context.Context, arg GetOrCreateDirectConversationParams) (uuid.UUID, error)
	GetPendingMentorApplications(ctx context.Context, spaceID uuid.UUID) ([]GetPendingMentorApplicationsRow, error)
	GetPendingReports(ctx context.Context, spaceID uuid.UUID) ([]GetPendingReportsRow, error)
	GetPendingTutorApplications(ctx context.Context, spaceID uuid.UUID) ([]GetPendingTutorApplicationsRow, error)
	GetPopularIndustries(ctx context.Context, spaceID uuid.UUID) ([]GetPopularIndustriesRow, error)
	GetPopularSubjects(ctx context.Context, spaceID uuid.UUID) ([]GetPopularSubjectsRow, error)
	GetPostByID(ctx context.Context, arg GetPostByIDParams) (GetPostByIDRow, error)
	GetPostComments(ctx context.Context, postID uuid.UUID) ([]GetPostCommentsRow, error)
	GetPostLikes(ctx context.Context, postID uuid.NullUUID) ([]GetPostLikesRow, error)
	GetProjectRoles(ctx context.Context, groupID uuid.UUID) ([]GroupRole, error)
	GetRecentFailedLoginAttemptsByIP(ctx context.Context, arg GetRecentFailedLoginAttemptsByIPParams) ([]LoginAttempt, error)
	GetRecentFailedLoginAttemptsByUsername(ctx context.Context, arg GetRecentFailedLoginAttemptsByUsernameParams) ([]LoginAttempt, error)
	GetRecentLoginAttemptsByIP(ctx context.Context, arg GetRecentLoginAttemptsByIPParams) ([]LoginAttempt, error)
	GetRecentLoginAttemptsByUsername(ctx context.Context, arg GetRecentLoginAttemptsByUsernameParams) ([]LoginAttempt, error)
	GetRecommendedMentors(ctx context.Context, arg GetRecommendedMentorsParams) ([]GetRecommendedMentorsRow, error)
	
	GetRecommendedTutors(ctx context.Context, arg GetRecommendedTutorsParams) ([]GetRecommendedTutorsRow, error)
	GetReport(ctx context.Context, id uuid.UUID) (GetReportRow, error)
	GetReportStats(ctx context.Context, spaceID uuid.UUID) (GetReportStatsRow, error)
	GetReportsByContent(ctx context.Context, arg GetReportsByContentParams) ([]Report, error)
	GetRoleApplications(ctx context.Context, groupID uuid.UUID) ([]GetRoleApplicationsRow, error)
	GetSession(ctx context.Context, id uuid.UUID) (UserSession, error)
	GetSpace(ctx context.Context, id uuid.UUID) (Space, error)
	GetSpaceActivities(ctx context.Context, arg GetSpaceActivitiesParams) ([]GetSpaceActivitiesRow, error)
	GetSpaceBySlug(ctx context.Context, slug string) (Space, error)
	GetSpaceStats(ctx context.Context, id uuid.UUID) (GetSpaceStatsRow, error)
	GetSpaceWithStats(ctx context.Context, id uuid.UUID) (GetSpaceWithStatsRow, error)
	GetSuggestedUsers(ctx context.Context, arg GetSuggestedUsersParams) ([]GetSuggestedUsersRow, error)
	
	GetSystemMetrics(ctx context.Context, spaceID uuid.UUID) (GetSystemMetricsRow, error)
	
	GetSystemSetting(ctx context.Context, key string) (SystemSetting, error)
	GetTopCommunities(ctx context.Context, spaceID uuid.UUID) ([]GetTopCommunitiesRow, error)
	GetTopGroups(ctx context.Context, spaceID uuid.UUID) ([]GetTopGroupsRow, error)
	GetTopPosts(ctx context.Context, spaceID uuid.UUID) ([]GetTopPostsRow, error)
	GetTrendingPosts(ctx context.Context, spaceID uuid.UUID) ([]GetTrendingPostsRow, error)
	GetTrendingTopics(ctx context.Context, arg GetTrendingTopicsParams) ([]GetTrendingTopicsRow, error)
	GetTutorApplication(ctx context.Context, id uuid.UUID) (GetTutorApplicationRow, error)
	GetTutorApplicationsByStatus(ctx context.Context, arg GetTutorApplicationsByStatusParams) ([]TutorApplication, error)
	GetTutorProfile(ctx context.Context, userID uuid.UUID) (GetTutorProfileRow, error)
	GetTutorReviews(ctx context.Context, arg GetTutorReviewsParams) ([]GetTutorReviewsRow, error)
	GetTutoringSession(ctx context.Context, id uuid.UUID) (GetTutoringSessionRow, error)
	GetTutoringStats(ctx context.Context, spaceID uuid.UUID) (GetTutoringStatsRow, error)
	GetUnreadCount(ctx context.Context, toUserID uuid.UUID) (int64, error)
	GetUnreadMessageCount(ctx context.Context, arg GetUnreadMessageCountParams) (int64, error)
	GetUnreadNotificationCount(ctx context.Context, toUserID uuid.UUID) (int64, error)
	GetUpcomingEvents(ctx context.Context, spaceID uuid.UUID) ([]GetUpcomingEventsRow, error)
	
	GetUserActivityStats(ctx context.Context, spaceID uuid.UUID) ([]GetUserActivityStatsRow, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (GetUserByIDRow, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserCommunities(ctx context.Context, arg GetUserCommunitiesParams) ([]GetUserCommunitiesRow, error)
	GetUserConversations(ctx context.Context, recipientID uuid.NullUUID) ([]GetUserConversationsRow, error)
	GetUserDetails(ctx context.Context, id uuid.UUID) (GetUserDetailsRow, error)
	GetUserEngagementAnalytics(ctx context.Context, authorID uuid.UUID) (GetUserEngagementAnalyticsRow, error)
	GetUserEngagementRanking(ctx context.Context, spaceID uuid.UUID) ([]GetUserEngagementRankingRow, error)
	GetUserEvents(ctx context.Context, arg GetUserEventsParams) ([]GetUserEventsRow, error)
	GetUserFeed(ctx context.Context, arg GetUserFeedParams) ([]GetUserFeedRow, error)
	GetUserFollowers(ctx context.Context, arg GetUserFollowersParams) ([]GetUserFollowersRow, error)
	GetUserFollowing(ctx context.Context, arg GetUserFollowingParams) ([]GetUserFollowingRow, error)
	GetUserGroups(ctx context.Context, arg GetUserGroupsParams) ([]GetUserGroupsRow, error)
	GetUserGrowth(ctx context.Context, spaceID uuid.UUID) ([]GetUserGrowthRow, error)
	GetUserGrowthData(ctx context.Context, arg GetUserGrowthDataParams) ([]GetUserGrowthDataRow, error)
	GetUserLikedPosts(ctx context.Context, arg GetUserLikedPostsParams) ([]GetUserLikedPostsRow, error)
	GetUserMentorApplicationStatus(ctx context.Context, arg GetUserMentorApplicationStatusParams) (sql.NullString, error)
	GetUserMentorApplicationStatusById(ctx context.Context, id uuid.UUID) (sql.NullString, error)
	GetUserMentoringSessions(ctx context.Context, arg GetUserMentoringSessionsParams) ([]GetUserMentoringSessionsRow, error)
	GetUserNotifications(ctx context.Context, arg GetUserNotificationsParams) ([]GetUserNotificationsRow, error)
	GetUserPosts(ctx context.Context, arg GetUserPostsParams) ([]GetUserPostsRow, error)
	GetUserSessionActivity(ctx context.Context, arg GetUserSessionActivityParams) ([]UserSession, error)
	
	
	GetUserSettings(ctx context.Context, id uuid.UUID) (pqtype.NullRawMessage, error)
	GetUserStats(ctx context.Context, id uuid.UUID) (GetUserStatsRow, error)
	GetUserSuspensions(ctx context.Context, arg GetUserSuspensionsParams) ([]GetUserSuspensionsRow, error)
	GetUserTutorApplicationStatus(ctx context.Context, arg GetUserTutorApplicationStatusParams) (sql.NullString, error)
	GetUserTutorApplicationStatusById(ctx context.Context, id uuid.UUID) (sql.NullString, error)
	GetUserTutoringSessions(ctx context.Context, arg GetUserTutoringSessionsParams) ([]GetUserTutoringSessionsRow, error)
	GetUsersByRole(ctx context.Context, arg GetUsersByRoleParams) ([]GetUsersByRoleRow, error)
	GetUsersWithPendingApplications(ctx context.Context, spaceID uuid.UUID) ([]GetUsersWithPendingApplicationsRow, error)
	IncrementFailedLoginAttempts(ctx context.Context, id uuid.UUID) (User, error)
	IncrementFollowersCount(ctx context.Context, id uuid.UUID) error
	IncrementFollowingCount(ctx context.Context, id uuid.UUID) error
	IncrementPostViews(ctx context.Context, id uuid.UUID) error
	IsCommunityAdmin(ctx context.Context, arg IsCommunityAdminParams) (bool, error)
	IsCommunityModerator(ctx context.Context, arg IsCommunityModeratorParams) (bool, error)
	IsGroupAdmin(ctx context.Context, arg IsGroupAdminParams) (bool, error)
	IsGroupModerator(ctx context.Context, arg IsGroupModeratorParams) (bool, error)
	IsUserSuperAdmin(ctx context.Context, id uuid.UUID) (bool, error)
	JoinCommunity(ctx context.Context, arg JoinCommunityParams) (CommunityMember, error)
	JoinGroup(ctx context.Context, arg JoinGroupParams) (GroupMember, error)
	LeaveCommunity(ctx context.Context, arg LeaveCommunityParams) error
	LeaveConversation(ctx context.Context, arg LeaveConversationParams) error
	LeaveGroup(ctx context.Context, arg LeaveGroupParams) error
	LiftSuspension(ctx context.Context, id uuid.UUID) error
	
	ListAllAnnouncementsAdmin(ctx context.Context, arg ListAllAnnouncementsAdminParams) ([]ListAllAnnouncementsAdminRow, error)
	
	ListAllCommunitiesAdmin(ctx context.Context, arg ListAllCommunitiesAdminParams) ([]ListAllCommunitiesAdminRow, error)
	
	ListAllEventsAdmin(ctx context.Context, arg ListAllEventsAdminParams) ([]ListAllEventsAdminRow, error)
	ListAnnouncements(ctx context.Context, arg ListAnnouncementsParams) ([]ListAnnouncementsRow, error)
	ListCommunities(ctx context.Context, arg ListCommunitiesParams) ([]ListCommunitiesRow, error)
	ListEvents(ctx context.Context, arg ListEventsParams) ([]ListEventsRow, error)
	ListGroups(ctx context.Context, arg ListGroupsParams) ([]ListGroupsRow, error)
	ListSpaces(ctx context.Context, arg ListSpacesParams) ([]Space, error)
	
	ListUsers(ctx context.Context, arg ListUsersParams) ([]ListUsersRow, error)
	MarkAllAsRead(ctx context.Context, toUserID uuid.UUID) error
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkEventAttendance(ctx context.Context, arg MarkEventAttendanceParams) error
	MarkMessagesAsRead(ctx context.Context, arg MarkMessagesAsReadParams) error
	
	
	
	
	
	
	MarkNotificationsAsRead(ctx context.Context, toUserID uuid.UUID) error
	PinPost(ctx context.Context, arg PinPostParams) error
	RateMentoringSession(ctx context.Context, arg RateMentoringSessionParams) (MentoringSession, error)
	RateTutoringSession(ctx context.Context, arg RateTutoringSessionParams) (TutoringSession, error)
	RegisterForEvent(ctx context.Context, arg RegisterForEventParams) (EventAttendee, error)
	RemoveCommunityModerator(ctx context.Context, arg RemoveCommunityModeratorParams) error
	RemoveEventCoOrganizer(ctx context.Context, arg RemoveEventCoOrganizerParams) error
	RemoveGroupAdmin(ctx context.Context, arg RemoveGroupAdminParams) error
	RemoveGroupModerator(ctx context.Context, arg RemoveGroupModeratorParams) error
	RemoveMessageReaction(ctx context.Context, arg RemoveMessageReactionParams) error
	ResetFailedLoginAttempts(ctx context.Context, id uuid.UUID) (User, error)
	ResetUserPassword(ctx context.Context, arg ResetUserPasswordParams) error
	SearchCommunities(ctx context.Context, arg SearchCommunitiesParams) ([]SearchCommunitiesRow, error)
	SearchEvents(ctx context.Context, arg SearchEventsParams) ([]SearchEventsRow, error)
	SearchGroups(ctx context.Context, arg SearchGroupsParams) ([]SearchGroupsRow, error)
	SearchMentors(ctx context.Context, arg SearchMentorsParams) ([]SearchMentorsRow, error)
	SearchPosts(ctx context.Context, arg SearchPostsParams) ([]SearchPostsRow, error)
	SearchTutors(ctx context.Context, arg SearchTutorsParams) ([]SearchTutorsRow, error)
	SearchUsers(ctx context.Context, arg SearchUsersParams) ([]SearchUsersRow, error)
	
	SearchUsersAdmin(ctx context.Context, arg SearchUsersAdminParams) ([]SearchUsersAdminRow, error)
	SendMessage(ctx context.Context, arg SendMessageParams) (Message, error)
	ToggleCommentLike(ctx context.Context, arg ToggleCommentLikeParams) (bool, error)
	TogglePostLike(ctx context.Context, arg TogglePostLikeParams) (sql.NullInt32, error)
	UnfollowUser(ctx context.Context, arg UnfollowUserParams) error
	UnlockExpiredAccounts(ctx context.Context) error
	UnregisterFromEvent(ctx context.Context, arg UnregisterFromEventParams) error
	UpdateAnnouncement(ctx context.Context, arg UpdateAnnouncementParams) (Announcement, error)
	UpdateAnnouncementStatus(ctx context.Context, arg UpdateAnnouncementStatusParams) (Announcement, error)
	UpdateCommunity(ctx context.Context, arg UpdateCommunityParams) (Community, error)
	UpdateCommunityStats(ctx context.Context, communityID uuid.UUID) error
	UpdateCommunityStatus(ctx context.Context, arg UpdateCommunityStatusParams) (Community, error)
	UpdateContentReportPriority(ctx context.Context, arg UpdateContentReportPriorityParams) (Report, error)
	UpdateContentReportStatus(ctx context.Context, arg UpdateContentReportStatusParams) (Report, error)
	UpdateContentReportWithAction(ctx context.Context, arg UpdateContentReportWithActionParams) (Report, error)
	UpdateConversationLastMessage(ctx context.Context, arg UpdateConversationLastMessageParams) error
	UpdateConversationSettings(ctx context.Context, arg UpdateConversationSettingsParams) error
	UpdateEvent(ctx context.Context, arg UpdateEventParams) (Event, error)
	UpdateEventStatus(ctx context.Context, arg UpdateEventStatusParams) (Event, error)
	UpdateGroup(ctx context.Context, arg UpdateGroupParams) (Group, error)
	UpdateGroupMemberRole(ctx context.Context, arg UpdateGroupMemberRoleParams) error
	UpdateGroupStats(ctx context.Context, groupID uuid.UUID) error
	UpdateGroupStatus(ctx context.Context, arg UpdateGroupStatusParams) error
	UpdateMentorApplication(ctx context.Context, arg UpdateMentorApplicationParams) (MentorApplication, error)
	UpdateMentorApplicationStatus(ctx context.Context, arg UpdateMentorApplicationStatusParams) (MentorApplication, error)
	UpdateMentorAvailability(ctx context.Context, arg UpdateMentorAvailabilityParams) (MentorProfile, error)
	UpdateMentorStatus(ctx context.Context, arg UpdateMentorStatusParams) (User, error)
	UpdateMentoringSessionStatus(ctx context.Context, arg UpdateMentoringSessionStatusParams) (MentoringSession, error)
	UpdateParticipantSettings(ctx context.Context, arg UpdateParticipantSettingsParams) error
	UpdateReport(ctx context.Context, arg UpdateReportParams) (Report, error)
	UpdateSessionStatus(ctx context.Context, arg UpdateSessionStatusParams) (TutoringSession, error)
	UpdateSpace(ctx context.Context, arg UpdateSpaceParams) (Space, error)
	UpdateTutorApplication(ctx context.Context, arg UpdateTutorApplicationParams) (TutorApplication, error)
	UpdateTutorApplicationStatus(ctx context.Context, arg UpdateTutorApplicationStatusParams) (TutorApplication, error)
	UpdateTutorAvailability(ctx context.Context, arg UpdateTutorAvailabilityParams) (TutorProfile, error)
	UpdateTutorStatus(ctx context.Context, arg UpdateTutorStatusParams) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserAccountStatus(ctx context.Context, arg UpdateUserAccountStatusParams) error
	UpdateUserLastActive(ctx context.Context, id uuid.UUID) error
	
	UpdateUserLockStatus(ctx context.Context, arg UpdateUserLockStatusParams) (User, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
	UpdateUserRole(ctx context.Context, arg UpdateUserRoleParams) (User, error)
	UpdateUserSettings(ctx context.Context, arg UpdateUserSettingsParams) (User, error)
	UpsertSystemSetting(ctx context.Context, arg UpsertSystemSettingParams) (SystemSetting, error)
}

var _ Querier = (*Queries)(nil)
