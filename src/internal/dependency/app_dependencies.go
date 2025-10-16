package dependency

import (
	"database/sql"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	repository2 "workoutpal/src/internal/repository"
	service2 "workoutpal/src/internal/service"
)

type AppDependencies struct {
	UserRepository         repository.UserRepository
	RelationshipRepository repository.RelationshipRepository
	GoalRepository         repository.GoalRepository
	ExerciseRepository     repository.ExerciseRepository
	RoutineRepository      repository.RoutineRepository
	UserService            service.UserService
	RelationshipService    service.RelationshipService
	GoalService            service.GoalService
	ExerciseService        service.ExerciseService
	RoutineService         service.RoutineService
	AuthService            service.AuthService
}

func NewAppDependencies(db *sql.DB) AppDependencies {
	// --- Init Repositories ---
	userRepository := repository2.NewUserRepository(db)
	relationshipRepository := repository2.NewRelationshipRepository(db)
	goalRepository := repository2.NewGoalRepository(db)
	exerciseRepository := repository2.NewExerciseRepository(db)
	routineRepository := repository2.NewRoutineRepository(db)

	// --- Init Services ---
	userService := service2.NewUserService(userRepository)
	relationshipService := service2.NewRelationshipService(relationshipRepository)
	goalService := service2.NewGoalService(goalRepository)
	exerciseService := service2.NewExerciseService(exerciseRepository)
	routineService := service2.NewRoutineService(routineRepository)
	authService := service2.NewAuthService(userRepository)

	return AppDependencies{
		UserRepository:         userRepository,
		RelationshipRepository: relationshipRepository,
		GoalRepository:         goalRepository,
		ExerciseRepository:     exerciseRepository,
		RoutineRepository:      routineRepository,
		UserService:            userService,
		RelationshipService:    relationshipService,
		GoalService:            goalService,
		ExerciseService:        exerciseService,
		RoutineService:         routineService,
		AuthService:            authService,
	}
}
