package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type goalHandler struct {
	goalService service.GoalService
}

func NewGoalHandler(gs service.GoalService) handler.GoalHandler {
	return &goalHandler{
		goalService: gs,
	}
}

// CreateUserGoal godoc
// @Summary Create a goal for user
// @Tags Goals
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body model.CreateGoalRequest true "Goal payload"
// @Success 201 {object} model.Goal "Goal created successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id}/goals [post]
func (g *goalHandler) CreateUserGoal(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	var req model.CreateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid request body"})
		return
	}

	goal, err := g.goalService.CreateGoal(id, req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, goal)
}

// GetUserGoals godoc
// @Summary Get all goals for a user
// @Tags Goals
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} model.Goal "Goals retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id}/goals [get]
func (g *goalHandler) GetUserGoals(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	goals, err := g.goalService.ReadUserGoals(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, goals)
}
