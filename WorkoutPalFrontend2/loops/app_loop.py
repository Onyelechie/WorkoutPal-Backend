from loops.achievements_loop import achievements_loop
from loops.exercise_loop import exercises_loop
from loops.exercise_setting_loop import exercise_settings_loop
from loops.goals_loop import goals_loop
from loops.posts_loop import posts_loop
from loops.relationships_loop import relationships_loop
from loops.routines_loop import routines_loop
from loops.schedules_loop import schedules_loop
from loops.users_loop import users_loop


def app_loop():
    print("Main Page")
    print("Type 'auth', 'users', 'exercises', 'exercise_settings', 'posts', 'routines', 'schedules', 'goals', 'achievements', 'relationships', or 'quit'")
    while True:
        cmd = input("\nAPP> ").strip().lower()

        if cmd == "users":
            users_loop()
        elif cmd == "exercises":
            exercises_loop()
        elif cmd == "exercise_settings":
            exercise_settings_loop()
        elif cmd == "posts":
            posts_loop()
        elif cmd == "routines":
            routines_loop()
        elif cmd == "schedules":
            schedules_loop()
        elif cmd == "goals":
            goals_loop()
        elif cmd == "achievements":
            achievements_loop()
        elif cmd == "relationships":
            relationships_loop()
        elif cmd in ["quit", "exit"]:
            break
        elif cmd == "help":
            print("Type 'auth', 'users', 'exercises', 'exercise_settings', 'posts', 'routines', 'schedules', 'goals', 'achievements', 'relationships', or 'quit'")
        else:
            print("Unknown command. Type 'help' for options.")