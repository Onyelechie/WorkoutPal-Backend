package in_memory

import (
	"errors"
	"sync"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type inMemoryUserRepository struct {
	users         map[int64]*model.User
	goals         map[int64]*model.Goal
	routines      map[int64]*model.ExerciseRoutine
	nextID        int64
	nextGoalID    int64
	nextRoutineID int64
	mutex         sync.RWMutex
}

func NewInMemoryUserRepository() repository.UserRepository {
	return &inMemoryUserRepository{
		users:         make(map[int64]*model.User),
		goals:         make(map[int64]*model.Goal),
		routines:      make(map[int64]*model.ExerciseRoutine),
		nextID:        1,
		nextGoalID:    1,
		nextRoutineID: 1,
	}
}

func (u *inMemoryUserRepository) ReadUserByEmail(email string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *inMemoryUserRepository) ReadUsers() ([]*model.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	users := make([]*model.User, 0, len(u.users))
	for _, user := range u.users {
		users = append(users, user)
	}
	return users, nil
}

func (u *inMemoryUserRepository) ReadUserByID(id int64) (*model.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, exists := u.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *inMemoryUserRepository) CreateUser(request model.CreateUserRequest) (*model.User, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for _, user := range u.users {
		if user.Username == request.Username {
			return nil, errors.New("username already exists")
		}
		if user.Email == request.Email {
			return nil, errors.New("email already exists")
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

	return user, nil
}

func (u *inMemoryUserRepository) UpdateUser(request model.UpdateUserRequest) (*model.User, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	user, exists := u.users[request.ID]
	if !exists {
		return nil, errors.New("user not found")
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

	return user, nil
}

func (u *inMemoryUserRepository) DeleteUser(request model.DeleteUserRequest) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.users[request.ID]; !exists {
		return errors.New("user not found")
	}

	delete(u.users, request.ID)
	return nil
}

func (u *inMemoryUserRepository) CreateGoal(userID int64, request model.CreateGoalRequest) (*model.Goal, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	goal := &model.Goal{
		ID:          u.nextGoalID,
		UserID:      userID,
		Name:        request.Name,
		Description: request.Description,
		Deadline:    request.Deadline,
		CreatedAt:   "2024-01-01",
		Status:      "active",
	}

	u.goals[u.nextGoalID] = goal
	u.nextGoalID++
	return goal, nil
}

func (u *inMemoryUserRepository) ReadUserGoals(userID int64) ([]*model.Goal, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	if _, exists := u.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	var goals []*model.Goal
	for _, goal := range u.goals {
		if goal.UserID == userID {
			goals = append(goals, goal)
		}
	}
	return goals, nil
}

func (u *inMemoryUserRepository) UpdateGoal(request model.UpdateGoalRequest) (*model.Goal, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	goal, exists := u.goals[request.ID]
	if !exists {
		return nil, errors.New("goal not found")
	}

	goal.Name = request.Name
	goal.Description = request.Description
	goal.Deadline = request.Deadline
	goal.Status = request.Status

	return goal, nil
}

func (u *inMemoryUserRepository) DeleteGoal(goalID int64) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.goals[goalID]; !exists {
		return errors.New("goal not found")
	}

	delete(u.goals, goalID)
	return nil
}

func (u *inMemoryUserRepository) FollowUser(followerID, followeeID int64) error {
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

func (u *inMemoryUserRepository) UnfollowUser(followerID, followeeID int64) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	follower, exists := u.users[followerID]
	if !exists {
		return errors.New("follower not found")
	}
	followee, exists := u.users[followeeID]
	if !exists {
		return errors.New("user not found")
	}

	// Remove from follower's following list
	for i, id := range follower.Following {
		if id == followeeID {
			follower.Following = append(follower.Following[:i], follower.Following[i+1:]...)
			break
		}
	}

	// Remove from followee's followers list
	for i, id := range followee.Followers {
		if id == followerID {
			followee.Followers = append(followee.Followers[:i], followee.Followers[i+1:]...)
			break
		}
	}

	return nil
}

func (u *inMemoryUserRepository) ReadUserFollowers(userID int64) ([]int64, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, exists := u.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user.Followers, nil
}

func (u *inMemoryUserRepository) ReadUserFollowing(userID int64) ([]int64, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, exists := u.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user.Following, nil
}

func (u *inMemoryUserRepository) CreateRoutine(userID int64, request model.CreateRoutineRequest) (*model.ExerciseRoutine, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.users[userID]; !exists {
		return nil, errors.New("user not found")
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
	return routine, nil
}

func (u *inMemoryUserRepository) ReadUserRoutines(userID int64) ([]*model.ExerciseRoutine, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	if _, exists := u.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	var routines []*model.ExerciseRoutine
	for _, routine := range u.routines {
		if routine.UserID == userID {
			routines = append(routines, routine)
		}
	}
	return routines, nil
}

func (u *inMemoryUserRepository) DeleteRoutine(routineID int64) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.routines[routineID]; !exists {
		return errors.New("routine not found")
	}

	delete(u.routines, routineID)
	return nil
}

func (u *inMemoryUserRepository) ReadRoutineWithExercises(routineID int64) (*model.ExerciseRoutine, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	routine, exists := u.routines[routineID]
	if !exists {
		return nil, errors.New("routine not found")
	}
	return routine, nil
}

func (u *inMemoryUserRepository) AddExerciseToRoutine(routineID, exerciseID int64) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	routine, exists := u.routines[routineID]
	if !exists {
		return errors.New("routine not found")
	}

	// Mock exercise data
	exercise := model.Exercise{ID: exerciseID, Name: "Exercise"}
	routine.Exercises = append(routine.Exercises, exercise)
	return nil
}

func (u *inMemoryUserRepository) RemoveExerciseFromRoutine(routineID, exerciseID int64) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	routine, exists := u.routines[routineID]
	if !exists {
		return errors.New("routine not found")
	}

	for i, exercise := range routine.Exercises {
		if exercise.ID == exerciseID {
			routine.Exercises = append(routine.Exercises[:i], routine.Exercises[i+1:]...)
			break
		}
	}
	return nil
}
