package repository

import (
	"errors"
	"sync"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type userRepository struct {
	users       map[int64]*model.User
	goals       map[int64]*model.Goal
	routines    map[int64]*model.ExerciseRoutine
	nextID      int64
	nextGoalID  int64
	nextRoutineID int64
	mutex       sync.RWMutex
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{
		users:         make(map[int64]*model.User),
		goals:         make(map[int64]*model.Goal),
		routines:      make(map[int64]*model.ExerciseRoutine),
		nextID:        1,
		nextGoalID:    1,
		nextRoutineID: 1,
	}
}

func (u *userRepository) ReadUsers() ([]model.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	users := make([]model.User, 0, len(u.users))
	for _, user := range u.users {
		users = append(users, *user)
	}
	return users, nil
}

func (u *userRepository) GetUserByID(id int64) (model.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, exists := u.users[id]
	if !exists {
		return model.User{}, errors.New("user not found")
	}
	return *user, nil
}

func (u *userRepository) CreateUser(request model.CreateUserRequest) (model.User, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for _, user := range u.users {
		if user.Username == request.Username {
			return model.User{}, errors.New("username already exists")
		}
		if user.Email == request.Email {
			return model.User{}, errors.New("email already exists")
		}
	}

	user := &model.User{
		ID:           u.nextID,
		Username:     request.Username,
		Email:        request.Email,
		Password:     request.Password,
		Name:         request.Name,
		Avatar:       request.Avatar,
		Age:          request.Age,
		Height:       request.Height,
		HeightMetric: request.HeightMetric,
		Weight:       request.Weight,
		WeightMetric: request.WeightMetric,
	}

	u.users[u.nextID] = user
	u.nextID++

	return *user, nil
}

func (u *userRepository) UpdateUser(request model.UpdateUserRequest) (model.User, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	user, exists := u.users[request.ID]
	if !exists {
		return model.User{}, errors.New("user not found")
	}

	user.Username = request.Username
	user.Email = request.Email
	if request.Password != "" {
		user.Password = request.Password
	}
	user.Name = request.Name
	user.Avatar = request.Avatar
	user.Age = request.Age
	user.Height = request.Height
	user.HeightMetric = request.HeightMetric
	user.Weight = request.Weight
	user.WeightMetric = request.WeightMetric

	return *user, nil
}

func (u *userRepository) DeleteUser(request model.DeleteUserRequest) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.users[request.ID]; !exists {
		return errors.New("user not found")
	}

	delete(u.users, request.ID)
	return nil
}

func (u *userRepository) CreateGoal(userID int64, request model.CreateGoalRequest) (model.Goal, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.users[userID]; !exists {
		return model.Goal{}, errors.New("user not found")
	}

	goal := &model.Goal{
		ID:          u.nextGoalID,
		UserID:      userID,
		Type:        request.Type,
		TargetValue: request.TargetValue,
		CurrentValue: 0,
		TargetDate:  request.TargetDate,
		CreatedAt:   "2024-01-01",
		Status:      "active",
	}

	u.goals[u.nextGoalID] = goal
	u.nextGoalID++
	return *goal, nil
}

func (u *userRepository) GetUserGoals(userID int64) ([]model.Goal, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	if _, exists := u.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	var goals []model.Goal
	for _, goal := range u.goals {
		if goal.UserID == userID {
			goals = append(goals, *goal)
		}
	}
	return goals, nil
}

func (u *userRepository) FollowUser(followerID, followeeID int64) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	follower, exists := u.users[followerID]
	if !exists {
		return errors.New("follower not found")
	}
	followee, exists := u.users[followeeID]
	if !exists {
		return errors.New("user to follow not found")
	}

	for _, id := range follower.Following {
		if id == followeeID {
			return errors.New("already following this user")
		}
	}
	follower.Following = append(follower.Following, followeeID)
	followee.Followers = append(followee.Followers, followerID)
	return nil
}

func (u *userRepository) GetUserFollowers(userID int64) ([]int64, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, exists := u.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user.Followers, nil
}

func (u *userRepository) GetUserFollowing(userID int64) ([]int64, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, exists := u.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user.Following, nil
}

func (u *userRepository) CreateRoutine(userID int64, request model.CreateRoutineRequest) (model.ExerciseRoutine, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.users[userID]; !exists {
		return model.ExerciseRoutine{}, errors.New("user not found")
	}

	routine := &model.ExerciseRoutine{
		ID:          u.nextRoutineID,
		UserID:      userID,
		Name:        request.Name,
		Description: request.Description,
		Exercises:   []model.Exercise{},
		CreatedAt:   "2024-01-01",
		IsActive:    true,
	}

	u.routines[u.nextRoutineID] = routine
	u.nextRoutineID++
	return *routine, nil
}

func (u *userRepository) GetUserRoutines(userID int64) ([]model.ExerciseRoutine, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	if _, exists := u.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	var routines []model.ExerciseRoutine
	for _, routine := range u.routines {
		if routine.UserID == userID {
			routines = append(routines, *routine)
		}
	}
	return routines, nil
}
