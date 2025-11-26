package domain

// Services
//go:generate mockgen -destination=../../mock_internal/domain/service/mock_auth_service.go          -package=mock_service workoutpal/src/internal/domain/service AuthService
//go:generate mockgen -destination=../../mock_internal/domain/service/mock_exercise_service.go      -package=mock_service workoutpal/src/internal/domain/service ExerciseService
//go:generate mockgen -destination=../../mock_internal/domain/service/mock_goal_service.go          -package=mock_service workoutpal/src/internal/domain/service GoalService
//go:generate mockgen -destination=../../mock_internal/domain/service/mock_relationship_service.go  -package=mock_service workoutpal/src/internal/domain/service RelationshipService
//go:generate mockgen -destination=../../mock_internal/domain/service/mock_routine_service.go       -package=mock_service workoutpal/src/internal/domain/service RoutineService
//go:generate mockgen -destination=../../mock_internal/domain/service/mock_user_service.go          -package=mock_service workoutpal/src/internal/domain/service UserService
//go:generate mockgen -destination=../../mock_internal/domain/service/mock_schedule_service.go      -package=mock_service workoutpal/src/internal/domain/service ScheduleService
//go:generate mockgen -destination=../../mock_internal/domain/service/mock_post_service.go      	-package=mock_service workoutpal/src/internal/domain/service PostService
//go:generate mockgen -destination=../../mock_internal/domain/service/mock_achievement_service.go   -package=mock_service workoutpal/src/internal/domain/service AchievementService
//go:generate mockgen -destination=../../mock_internal/domain/service/exercise_setting_service.go   -package=mock_service workoutpal/src/internal/domain/service ExerciseSettingService
// Repositories
//go:generate mockgen -destination=../../mock_internal/domain/repository/mock_exercise_repository.go     -package=mock_repository workoutpal/src/internal/domain/repository ExerciseRepository
//go:generate mockgen -destination=../../mock_internal/domain/repository/mock_goal_repository.go         -package=mock_repository workoutpal/src/internal/domain/repository GoalRepository
//go:generate mockgen -destination=../../mock_internal/domain/repository/mock_relationship_repository.go -package=mock_repository workoutpal/src/internal/domain/repository RelationshipRepository
//go:generate mockgen -destination=../../mock_internal/domain/repository/mock_routine_repository.go      -package=mock_repository workoutpal/src/internal/domain/repository RoutineRepository
//go:generate mockgen -destination=../../mock_internal/domain/repository/mock_user_repository.go         -package=mock_repository workoutpal/src/internal/domain/repository UserRepository
//go:generate mockgen -destination=../../mock_internal/domain/repository/mock_schedule_repository.go     -package=mock_repository workoutpal/src/internal/domain/repository ScheduleRepository
//go:generate mockgen -destination=../../mock_internal/domain/repository/mock_post_repository.go		 -package=mock_repository workoutpal/src/internal/domain/repository PostRepository
//go:generate mockgen -destination=../../mock_internal/domain/repository/mock_achievement_repository.go	 -package=mock_repository workoutpal/src/internal/domain/repository AchievementRepository
//go:generate mockgen -destination=../../mock_internal/domain/repository/exercise_setting_repository.go	 -package=mock_repository workoutpal/src/internal/domain/repository ExerciseSettingRepository
