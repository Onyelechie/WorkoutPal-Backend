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
// @Summary Create a new achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Param request body model.CreateAchievementRequest true "New achievement payload"
// @Success 201 {object} model.Achievement "Achievement created successfully"
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

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ach)
}

// ReadAchievements godoc
// @Summary List achievements
// @Tags Achievements
// @Accept json
// @Produce json
// @Success 200 {array} model.Achievement "Achievements retrieved successfully"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /achievements [get]
func (h *AchievementHandler) ReadAchievements(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ReadAchievements()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

// DeleteAchievement godoc
// @Summary Delete an achievement
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path int true "Achievement ID"
// @Success 200 {object} model.BasicResponse "Achievement deleted successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 404 {object} model.BasicResponse "Achievement not found"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /achievements/{id} [delete]
func (h *AchievementHandler) DeleteAchievement(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	if err := h.svc.DeleteAchievement(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := model.BasicResponse{Message: "Success"}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
