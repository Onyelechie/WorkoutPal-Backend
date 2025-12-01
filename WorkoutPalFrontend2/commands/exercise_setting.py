from api import api_get, api_post, api_put


def read_exercise_setting():
    print("\n--- Read Exercise Setting ---")
    exercise_id = input("Exercise ID: ").strip()
    routine_id = input("Workout Routine ID: ").strip()

    params = {
        "exercise_id": exercise_id,
        "workout_routine_id": routine_id,
    }

    res = api_get("/exercise-settings", params=params)

    if res.status_code == 200:
        print("✔ Setting retrieved.")
        print(res.json())
    else:
        print("✖ Read failed:", res.text)


def create_exercise_setting():
    print("\n--- Create Exercise Setting ---")
    user_id = int(input("User ID: ").strip())
    exercise_id = int(input("Exercise ID: ").strip())
    workout_routine_id = int(input("Workout Routine ID: ").strip())

    def prompt_int(field_name, optional=False):
        val = input(f"{field_name} (integer{' or blank' if optional else ''}): ").strip()
        if not val and optional:
            return None
        try:
            return int(val)
        except ValueError:
            print("Invalid number, using 0.")
            return 0

    def prompt_float(field_name, optional=False):
        val = input(f"{field_name} (number{' or blank' if optional else ''}): ").strip()
        if not val and optional:
            return None
        try:
            return float(val)
        except ValueError:
            print("Invalid number, using 0.")
            return 0.0

    sets_ = prompt_int("Sets", optional=True)
    reps = prompt_int("Reps", optional=True)
    break_interval = prompt_int("Break interval (seconds)", optional=True)
    weight = prompt_float("Weight", optional=True)

    payload = {
        "userId": user_id,
        "exerciseId": exercise_id,
        "workoutRoutineId": workout_routine_id,
        "sets": sets_,
        "reps": reps,
        "breakInterval": break_interval,
        "weight": weight,
    }

    res = api_post("/exercise-settings", json=payload)

    if res.status_code == 201:
        print("✔ Exercise setting created successfully.")
    else:
        print("✖ Creation failed:", res.text)


def update_exercise_setting():
    print("\n--- Update Exercise Setting ---")
    user_id = int(input("User ID: ").strip())
    exercise_id = int(input("Exercise ID: ").strip())
    workout_routine_id = int(input("Workout Routine ID: ").strip())

    def prompt_int(field_name, optional=False):
        val = input(f"{field_name} (integer{' or blank' if optional else ''}): ").strip()
        if not val and optional:
            return None
        try:
            return int(val)
        except ValueError:
            print("Invalid number, using 0.")
            return 0

    def prompt_float(field_name, optional=False):
        val = input(f"{field_name} (number{' or blank' if optional else ''}): ").strip()
        if not val and optional:
            return None
        try:
            return float(val)
        except ValueError:
            print("Invalid number, using 0.")
            return 0.0

    sets_ = prompt_int("Sets", optional=True)
    reps = prompt_int("Reps", optional=True)
    break_interval = prompt_int("Break interval (seconds)", optional=True)
    weight = prompt_float("Weight", optional=True)

    payload = {
        "userId": user_id,
        "exerciseId": exercise_id,
        "workoutRoutineId": workout_routine_id,
        "sets": sets_,
        "reps": reps,
        "breakInterval": break_interval,
        "weight": weight,
    }

    res = api_put("/exercise-settings", json=payload)

    if res.status_code == 200:
        print("✔ Exercise setting updated successfully.")
    else:
        print("✖ Update failed:", res.text)