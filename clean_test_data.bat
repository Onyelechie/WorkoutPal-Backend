@echo off
echo Cleaning test data from database...
cd src\fitness-db
docker exec fitness-db psql -U user -d workoutpal -c "TRUNCATE TABLE user_achievements, workout_routine, exercises_in_routine, schedule, schedule_routine, posts, post_comments, post_likes, goals, follows, follow_requests RESTART IDENTITY CASCADE;"
echo Done!
