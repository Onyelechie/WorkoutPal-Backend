from api import api_get
from print_table import print_table


def read_exercises():
    print("\n--- Read Exercises ---")
    res = api_get("/exercises")

    if res.status_code == 200:
        print("✔ Read successful")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_exercises_filtered():
    print("\n--- Read Exercises (Filtered) ---")
    target = input("Target (primary muscle, blank for none): ").strip()
    intensity = input("Intensity (blank for none): ").strip()
    expertise = input("Expertise (blank for none): ").strip()

    params = {}
    if target:
        params["target"] = target
    if intensity:
        params["intensity"] = intensity
    if expertise:
        params["expertise"] = expertise

    res = api_get("/exercises", params=params if params else None)

    if res.status_code == 200:
        print("✔ Read successful")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_exercise_by_id():
    print("\n--- Read Exercise by ID ---")
    ex_id = input("Exercise ID: ").strip()

    res = api_get(f"/exercises/{ex_id}")

    if res.status_code == 200:
        print("✔ Exercise retrieved")
        print(res.json())
    else:
        print("✖ Read failed:", res.text)