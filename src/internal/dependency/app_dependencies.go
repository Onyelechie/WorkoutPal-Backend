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
	ScheduleService        service.ScheduleService
	AuthService            service.AuthService
	PostService            service.PostService
	AchievementService     service.AchievementService
}

func NewAppDependencies(db *sql.DB) AppDependencies {
	// --- Init Repositories ---
	userRepository := repository2.NewUserRepository(db)
	relationshipRepository := repository2.NewRelationshipRepository(db)
	goalRepository := repository2.NewGoalRepository(db)
	exerciseRepository := repository2.NewExerciseRepository(db)
	routineRepository := repository2.NewRoutineRepository(db)
	scheduleRepository := repository2.NewScheduleRepository(db)
	postRepository := repository2.NewPostRepository(db)
	achievementRepository := repository2.NewAchievementRepository(db)

	// --- Init Services ---
	userService := service2.NewUserService(userRepository)
	relationshipService := service2.NewRelationshipService(relationshipRepository, userRepository)
	goalService := service2.NewGoalService(goalRepository)
	exerciseService := service2.NewExerciseService(exerciseRepository)
	routineService := service2.NewRoutineService(routineRepository)
	authService := service2.NewAuthService(userRepository)
	scheduleService := service2.NewScheduleService(scheduleRepository)
	postService := service2.NewPostService(postRepository)
	achievementService := service2.NewAchievementService(achievementRepository)

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
		ScheduleService:        scheduleService,
		AuthService:            authService,
		PostService:            postService,
		AchievementService:     achievementService,
	}
}
