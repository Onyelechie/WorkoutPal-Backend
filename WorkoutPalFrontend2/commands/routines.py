from api import api_post, api_get
from print_table import print_table


def read_user_routines():
    print("\n--- Read User Routines ---")
    user_id = input("User ID: ").strip()

    res = api_get(f"/users/{user_id}/routines")

    if res.status_code == 200:
        print("✔ Routines retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_routine_by_id():
    print("\n--- Read Routine by ID ---")
    routine_id = input("Routine ID: ").strip()

    res = api_get(f"/routines/{routine_id}")

    if res.status_code == 200:
        print("✔ Routine retrieved.")
        print(res.json())
    else:
        print("✖ Read failed:", res.text)


def create_routine():
    print("\n--- Create Routine ---")
    user_id = int(input("User ID: ").strip())
    name = input("Name: ")
    description = input("Description: ")
    exercise_ids_raw = input("Exercise IDs (comma-separated): ")

    exercise_ids = []
    for part in exercise_ids_raw.split(","):
        part = part.strip()
        if not part:
            continue
        try:
            exercise_ids.append(int(part))
        except ValueError:
            print(f"Skipping invalid exercise id: {part}")

    payload = {
        "name": name,
        "description": description,
        "exerciseIds": exercise_ids,
    }

    res = api_post(f"/users/{user_id}/routines", json=payload)

    if res.status_code == 201:
        print("✔ Routine created successfully.")
    else:
        print("✖ Creation failed:", res.text)



def add_exercise_to_routine():
    print("\n--- Add Exercise to Routine ---")
    routine_id = input("Routine ID: ").strip()
    exercise_id = input("Exercise ID: ").strip()

    res = api_post(f"/routines/{routine_id}/exercises",
                   params={"exercise_id": exercise_id})

    if res.status_code == 200:
        print("✔ Exercise added to routine.")
    else:
        print("✖ Operation failed:", res.text)

