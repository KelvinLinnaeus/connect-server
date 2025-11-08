package handlers

import (
	"net/http"
	"strconv"

	"github.com/connect-univyn/connect-server/internal/service/mentorship"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MentorshipHandler struct {
	mentorshipService *mentorship.Service
}

func NewMentorshipHandler(mentorshipService *mentorship.Service) *MentorshipHandler {
	return &MentorshipHandler{
		mentorshipService: mentorshipService,
	}
}












func (h *MentorshipHandler) CreateMentorProfile(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req mentorship.CreateMentorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}
	req.UserID = userID

	profile, err := h.mentorshipService.CreateMentorProfile(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(profile))
}







func (h *MentorshipHandler) GetMentorProfile(c *gin.Context) {
	profileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid profile ID"))
		return
	}

	profile, err := h.mentorshipService.GetMentorProfile(c.Request.Context(), profileID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(profile))
}









func (h *MentorshipHandler) UpdateMentorAvailability(c *gin.Context) {
	profileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid profile ID"))
		return
	}

	var req mentorship.UpdateMentorAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	profile, err := h.mentorshipService.UpdateMentorAvailability(c.Request.Context(), profileID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(profile))
}












func (h *MentorshipHandler) SearchMentors(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	params := mentorship.SearchMentorsParams{
		SpaceID: spaceID,
	}

	if industry := c.Query("industry"); industry != "" {
		params.Industry = &industry
	}

	if specialties := c.Query("specialties"); specialties != "" {
		params.Specialties = []string{specialties}
	}

	if minRatingStr := c.Query("min_rating"); minRatingStr != "" {
		minRating, err := strconv.ParseFloat(minRatingStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid min_rating"))
			return
		}
		params.MinRating = &minRating
	}

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid page number"))
			return
		}
		params.Page = int32(page)
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limit < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid limit"))
			return
		}
		params.Limit = int32(limit)
	}

	mentors, err := h.mentorshipService.SearchMentors(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(mentors))
}







func (h *MentorshipHandler) GetMentorReviews(c *gin.Context) {
	mentorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid mentor ID"))
		return
	}

	reviews, err := h.mentorshipService.GetMentorReviews(c.Request.Context(), mentorID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(reviews))
}












func (h *MentorshipHandler) CreateTutorProfile(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req mentorship.CreateTutorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}
	req.UserID = userID

	profile, err := h.mentorshipService.CreateTutorProfile(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(profile))
}







func (h *MentorshipHandler) GetTutorProfile(c *gin.Context) {
	profileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid profile ID"))
		return
	}

	profile, err := h.mentorshipService.GetTutorProfile(c.Request.Context(), profileID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(profile))
}









func (h *MentorshipHandler) UpdateTutorAvailability(c *gin.Context) {
	profileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid profile ID"))
		return
	}

	var req mentorship.UpdateTutorAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	profile, err := h.mentorshipService.UpdateTutorAvailability(c.Request.Context(), profileID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(profile))
}












func (h *MentorshipHandler) SearchTutors(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	params := mentorship.SearchTutorsParams{
		SpaceID: spaceID,
	}

	if subject := c.Query("subject"); subject != "" {
		params.Subjects = []string{subject}
	}

	if maxRate := c.Query("max_rate"); maxRate != "" {
		params.MaxRate = &maxRate
	}

	if minRatingStr := c.Query("min_rating"); minRatingStr != "" {
		minRating, err := strconv.ParseFloat(minRatingStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid min_rating"))
			return
		}
		params.MinRating = &minRating
	}

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid page number"))
			return
		}
		params.Page = int32(page)
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limit < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid limit"))
			return
		}
		params.Limit = int32(limit)
	}

	tutors, err := h.mentorshipService.SearchTutors(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(tutors))
}







func (h *MentorshipHandler) GetTutorReviews(c *gin.Context) {
	tutorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid tutor ID"))
		return
	}

	reviews, err := h.mentorshipService.GetTutorReviews(c.Request.Context(), tutorID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(reviews))
}












func (h *MentorshipHandler) CreateMentoringSession(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req mentorship.CreateMentoringSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}
	req.MenteeID = userID

	session, err := h.mentorshipService.CreateMentoringSession(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(session))
}







func (h *MentorshipHandler) GetMentoringSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid session ID"))
		return
	}

	session, err := h.mentorshipService.GetMentoringSession(c.Request.Context(), sessionID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(session))
}










func (h *MentorshipHandler) GetUserMentoringSessions(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	page := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		p, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || p < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid page number"))
			return
		}
		page = int32(p)
	}

	limit := int32(20)
	if limitStr := c.Query("limit"); limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || l < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid limit"))
			return
		}
		limit = int32(l)
	}

	sessions, err := h.mentorshipService.GetUserMentoringSessions(c.Request.Context(), userID, spaceID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(sessions))
}









func (h *MentorshipHandler) UpdateMentoringSessionStatus(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid session ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=scheduled confirmed in_progress completed cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.mentorshipService.UpdateMentoringSessionStatus(c.Request.Context(), sessionID, req.Status); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]string{"message": "Session status updated successfully"}))
}









func (h *MentorshipHandler) AddMentoringSessionMeetingLink(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid session ID"))
		return
	}

	var req struct {
		MeetingLink string `json:"meeting_link" binding:"required,url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.mentorshipService.AddMentoringSessionMeetingLink(c.Request.Context(), sessionID, req.MeetingLink); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]string{"message": "Meeting link added successfully"}))
}









func (h *MentorshipHandler) RateMentoringSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid session ID"))
		return
	}

	var req mentorship.RateMentoringSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.mentorshipService.RateMentoringSession(c.Request.Context(), sessionID, req); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]string{"message": "Session rated successfully"}))
}












func (h *MentorshipHandler) CreateTutoringSession(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req mentorship.CreateTutoringSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}
	req.StudentID = userID

	session, err := h.mentorshipService.CreateTutoringSession(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(session))
}







func (h *MentorshipHandler) GetTutoringSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid session ID"))
		return
	}

	session, err := h.mentorshipService.GetTutoringSession(c.Request.Context(), sessionID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(session))
}










func (h *MentorshipHandler) GetUserTutoringSessions(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	page := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		p, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil || p < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid page number"))
			return
		}
		page = int32(p)
	}

	limit := int32(20)
	if limitStr := c.Query("limit"); limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || l < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid limit"))
			return
		}
		limit = int32(l)
	}

	sessions, err := h.mentorshipService.GetUserTutoringSessions(c.Request.Context(), userID, spaceID, page, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(sessions))
}









func (h *MentorshipHandler) UpdateTutoringSessionStatus(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid session ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=scheduled confirmed in_progress completed cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.mentorshipService.UpdateSessionStatus(c.Request.Context(), sessionID, req.Status); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]string{"message": "Session status updated successfully"}))
}









func (h *MentorshipHandler) AddTutoringSessionMeetingLink(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid session ID"))
		return
	}

	var req struct {
		MeetingLink string `json:"meeting_link" binding:"required,url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.mentorshipService.AddSessionMeetingLink(c.Request.Context(), sessionID, req.MeetingLink); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]string{"message": "Meeting link added successfully"}))
}









func (h *MentorshipHandler) RateTutoringSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid session ID"))
		return
	}

	var req mentorship.RateTutoringSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.mentorshipService.RateTutoringSession(c.Request.Context(), sessionID, req); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]string{"message": "Session rated successfully"}))
}












func (h *MentorshipHandler) CreateMentorApplication(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req mentorship.CreateMentorApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}
	req.UserID = userID

	application, err := h.mentorshipService.CreateMentorApplication(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(application))
}







func (h *MentorshipHandler) GetMentorApplication(c *gin.Context) {
	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid application ID"))
		return
	}

	application, err := h.mentorshipService.GetMentorApplication(c.Request.Context(), applicationID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(application))
}









func (h *MentorshipHandler) UpdateMentorApplication(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	reviewerID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid reviewer ID"))
		return
	}

	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid application ID"))
		return
	}

	var req mentorship.UpdateMentorApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	application, err := h.mentorshipService.UpdateMentorApplication(c.Request.Context(), applicationID, reviewerID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(application))
}







func (h *MentorshipHandler) GetPendingMentorApplications(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	applications, err := h.mentorshipService.GetPendingMentorApplications(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(applications))
}












func (h *MentorshipHandler) CreateTutorApplication(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	var req mentorship.CreateTutorApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}
	req.UserID = userID

	application, err := h.mentorshipService.CreateTutorApplication(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(application))
}







func (h *MentorshipHandler) GetTutorApplication(c *gin.Context) {
	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid application ID"))
		return
	}

	application, err := h.mentorshipService.GetTutorApplication(c.Request.Context(), applicationID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(application))
}









func (h *MentorshipHandler) UpdateTutorApplication(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	reviewerID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid reviewer ID"))
		return
	}

	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid application ID"))
		return
	}

	var req mentorship.UpdateTutorApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	application, err := h.mentorshipService.UpdateTutorApplication(c.Request.Context(), applicationID, reviewerID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(application))
}







func (h *MentorshipHandler) GetPendingTutorApplications(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	applications, err := h.mentorshipService.GetPendingTutorApplications(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(applications))
}











func (h *MentorshipHandler) GetMyMentorProfile(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	profile, err := h.mentorshipService.GetMentorProfileByUserID(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(profile))
}







func (h *MentorshipHandler) GetMyTutorProfile(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	profile, err := h.mentorshipService.GetTutorProfileByUserID(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(profile))
}







func (h *MentorshipHandler) GetMyMentorApplication(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	spaceId, err := uuid.Parse(authPayload.SpaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	application, err := h.mentorshipService.GetMentorApplicationByUserID(c.Request.Context(), userID, spaceId)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(application))
}







func (h *MentorshipHandler) GetMyTutorApplication(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	spaceID, err := uuid.Parse(authPayload.SpaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	application, err := h.mentorshipService.GetTutorApplicationByUserID(c.Request.Context(), userID, spaceID)

	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(application))
}













func (h *MentorshipHandler) GetRecommendedTutors(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	limit := int32(5)
	if limitStr := c.Query("limit"); limitStr != "" {
		limitVal, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitVal < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid limit"))
			return
		}
		limit = int32(limitVal)
	}

	tutors, err := h.mentorshipService.GetRecommendedTutors(c.Request.Context(), spaceID, userID, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(tutors))
}









func (h *MentorshipHandler) GetRecommendedMentors(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse(http.StatusUnauthorized, "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)

	userID, err := uuid.Parse(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	spaceID, err := uuid.Parse(c.Query("space_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID"))
		return
	}

	limit := int32(5)
	if limitStr := c.Query("limit"); limitStr != "" {
		limitVal, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limitVal < 1 {
			c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid limit"))
			return
		}
		limit = int32(limitVal)
	}

	mentors, err := h.mentorshipService.GetRecommendedMentors(c.Request.Context(), spaceID, userID, limit)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(mentors))
}
