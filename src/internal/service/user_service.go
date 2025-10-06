package service

import (
	"errors"
	"sync"
	"workoutpal/src/internal/model"
)

type UserService struct {
	users       map[int64]*model.User
	goals       map[int64]*model.Goal
	achievements map[int64]*model.UserAchievement
	routines    map[int64]*model.ExerciseRoutine
	nextID      int64
	nextGoalID  int64
	nextAchievementID int64
	nextRoutineID int64
	mutex       sync.RWMutex
}

func NewUserService() *UserService {
	return &UserService{
		users:             make(map[int64]*model.User),
		goals:             make(map[int64]*model.Goal),
		achievements:      make(map[int64]*model.UserAchievement),
		routines:          make(map[int64]*model.ExerciseRoutine),
		nextID:            1,
		nextGoalID:        1,
		nextAchievementID: 1,
		nextRoutineID:     1,
	}
}

func (s *UserService) CreateUser(req *model.CreateUserRequest) (*model.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if username or email already exists
	for _, user := range s.users {
		if user.Username == req.Username {
			return nil, errors.New("username already exists")
		}
		if user.Email == req.Email {
			return nil, errors.New("email already exists")
		}
	}

	user := &model.User{
		ID:           s.nextID,
		Username:     req.Username,
		Email:        req.Email,
		Password:     req.Password, // In real implementation, hash this
		Name:         req.Name,
		Avatar:       req.Avatar,
		Age:          req.Age,
		Height:       req.Height,
		HeightMetric: req.HeightMetric,
		Weight:       req.Weight,
		WeightMetric: req.WeightMetric,
	}

	s.users[s.nextID] = user
	s.nextID++

	return user, nil
}

func (s *UserService) GetAllUsers() []*model.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]*model.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *UserService) GetUserByID(id int64) (*model.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateUser(req *model.UpdateUserRequest) (*model.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[req.ID]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Check if username or email conflicts with other users
	for id, existingUser := range s.users {
		if id != req.ID {
			if existingUser.Username == req.Username {
				return nil, errors.New("username already exists")
			}
			if existingUser.Email == req.Email {
				return nil, errors.New("email already exists")
			}
		}
	}

	// Update fields
	user.Username = req.Username
	user.Email = req.Email
	if req.Password != "" {
		user.Password = req.Password // In real implementation, hash this
	}
	user.Name = req.Name
	user.Avatar = req.Avatar
	user.Age = req.Age
	user.Height = req.Height
	user.HeightMetric = req.HeightMetric
	user.Weight = req.Weight
	user.WeightMetric = req.WeightMetric

	return user, nil
}

func (s *UserService) DeleteUser(id int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(s.users, id)
	return nil
}

// Goal methods
func (s *UserService) CreateGoal(userID int64, req *model.CreateGoalRequest) (*model.Goal, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	goal := &model.Goal{
		ID:          s.nextGoalID,
		UserID:      userID,
		Type:        req.Type,
		TargetValue: req.TargetValue,
		CurrentValue: 0,
		TargetDate:  req.TargetDate,
		CreatedAt:   "2024-01-01", // In real implementation, use time.Now()
		Status:      "active",
	}

	s.goals[s.nextGoalID] = goal
	s.nextGoalID++
	return goal, nil
}

func (s *UserService) GetUserGoals(userID int64) ([]*model.Goal, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, exists := s.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	var goals []*model.Goal
	for _, goal := range s.goals {
		if goal.UserID == userID {
			goals = append(goals, goal)
		}
	}
	return goals, nil
}

// Follower methods
func (s *UserService) FollowUser(followerID, followeeID int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	follower, exists := s.users[followerID]
	if !exists {
		return errors.New("follower not found")
	}
	followee, exists := s.users[followeeID]
	if !exists {
		return errors.New("user to follow not found")
	}

	// Add to following list
	for _, id := range follower.Following {
		if id == followeeID {
			return errors.New("already following this user")
		}
	}
	follower.Following = append(follower.Following, followeeID)

	// Add to followers list
	followee.Followers = append(followee.Followers, followerID)
	return nil
}

func (s *UserService) GetUserFollowers(userID int64) ([]int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user.Followers, nil
}

func (s *UserService) GetUserFollowing(userID int64) ([]int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user.Following, nil
}

// Routine methods
func (s *UserService) CreateRoutine(userID int64, req *model.CreateRoutineRequest) (*model.ExerciseRoutine, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	routine := &model.ExerciseRoutine{
		ID:          s.nextRoutineID,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Exercises:   []model.Exercise{}, // Would populate from exerciseIDs
		CreatedAt:   "2024-01-01",
		IsActive:    true,
	}

	s.routines[s.nextRoutineID] = routine
	s.nextRoutineID++
	return routine, nil
}

func (s *UserService) GetUserRoutines(userID int64) ([]*model.ExerciseRoutine, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, exists := s.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	var routines []*model.ExerciseRoutine
	for _, routine := range s.routines {
		if routine.UserID == userID {
			routines = append(routines, routine)
		}
	}
	return routines, nil
}