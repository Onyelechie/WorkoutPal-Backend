package handler

import (
	"encoding/json"
	"net/http"

	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
	"workoutpal/src/util/constants"
)

type AchievementHandler struct {
	svc service.AchievementService
}

func NewAchievementHandler(svc service.AchievementService) *AchievementHandler {
	return &AchievementHandler{svc: svc}
}

// CreateAchievement godoc
// @Summary Unlock an achievement for a user
// @Description Creates a user_achievement row (i.e., marks the achievement as earned for the user).
// @Tags Achievements
// @Accept json
// @Produce json
// @Param request body model.CreateAchievementRequest true "Unlock achievement payload"
// @Success 201 {object} model.UserAchievement "UserAchievement created successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /achievements [post]
func (h *AchievementHandler) CreateAchievement(w http.ResponseWriter, r *http.Request) {
	var req model.CreateAchievementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ach, err := h.svc.CreateAchievement(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ach)
}

// ReadAllAchievements godoc
// @Summary List all achievements (catalog)
// @Tags Achievements
// @Accept json
// @Produce json
// @Success 200 {array} model.Achievement "Achievements retrieved successfully"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /achievements [get]
func (h *AchievementHandler) ReadAllAchievements(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ReadAllAchievements()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

// ReadUnlockedAchievements godoc
// @Summary List achievements unlocked by the current user
// @Tags Achievements
// @Accept json
// @Produce json
// @Success 200 {array} model.UserAchievement "Unlocked achievements retrieved successfully"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /achievements/unlocked [get]
func (h *AchievementHandler) ReadUnlockedAchievements(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.USER_ID_KEY).(int64)

	list, err := h.svc.ReadUnlockedAchievements(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

// ReadUnlockedAchievementsByUserID godoc
// @Summary List achievements unlocked by a specific user (by path ID)
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} model.UserAchievement "Unlocked achievements retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /achievements/unlocked/{id} [get]
func (h *AchievementHandler) ReadUnlockedAchievementsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.ID_KEY).(int64)

	list, err := h.svc.ReadUnlockedAchievements(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}
