from api import api_put, api_get, api_post
from print_table import print_table

def read_schedules():
    print("\n--- Read All Schedules ---")
    res = api_get("/schedules")

    if res.status_code == 200:
        print("✔ Schedules retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_schedules_by_day():
    print("\n--- Read Schedules by Day ---")
    day = input("Day of week (0-6): ").strip()

    res = api_get(f"/schedules/{day}")

    if res.status_code == 200:
        print("✔ Schedules retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_schedule_by_id():
    print("\n--- Read Schedule by ID ---")
    schedule_id = input("Schedule ID: ").strip()

    res = api_get(f"/schedules/{schedule_id}")

    if res.status_code == 200:
        print("✔ Schedule retrieved.")
        print(res.json())
    else:
        print("✖ Read failed:", res.text)


def create_schedule():
    print("\n--- Create Schedule ---")
    user_id = int(input("User ID: ").strip())
    name = input("Name: ")
    day_of_week = int(input("Day of week (0-6): ").strip())
    time_slot = input("Time slot (e.g. 18:00): ")
    routine_length = int(input("Routine length (minutes): ").strip())
    routine_ids_raw = input("Routine IDs (comma-separated): ")

    routine_ids = []
    for part in routine_ids_raw.split(","):
        part = part.strip()
        if not part:
            continue
        try:
            routine_ids.append(int(part))
        except ValueError:
            print(f"Skipping invalid routine id: {part}")

    payload = {
        "userId": user_id,
        "name": name,
        "dayOfWeek": day_of_week,
        "timeSlot": time_slot,
        "routineLengthMinutes": routine_length,
        "routineIds": routine_ids,
    }

    res = api_post("/schedules", json=payload)

    if res.status_code == 201:
        print("✔ Schedule created successfully.")
    else:
        print("✖ Creation failed:", res.text)


def update_schedule():
    print("\n--- Update Schedule ---")
    schedule_id = input("Schedule ID: ").strip()

    user_id = int(input("User ID: ").strip())
    name = input("Name: ")
    day_of_week = int(input("Day of week (0-6): ").strip())
    time_slot = input("Time slot (e.g. 18:00): ")
    routine_length = int(input("Routine length (minutes): ").strip())
    routine_ids_raw = input("Routine IDs (comma-separated): ")

    routine_ids = []
    for part in routine_ids_raw.split(","):
        part = part.strip()
        if not part:
            continue
        try:
            routine_ids.append(int(part))
        except ValueError:
            print(f"Skipping invalid routine id: {part}")

    payload = {
        "id": int(schedule_id),
        "userId": user_id,
        "name": name,
        "dayOfWeek": day_of_week,
        "timeSlot": time_slot,
        "routineLengthMinutes": routine_length,
        "routineIds": routine_ids,
    }

    res = api_put(f"/schedules/{schedule_id}", json=payload)

    if res.status_code == 200:
        print("✔ Schedule updated successfully.")
    else:
        print("✖ Update failed:", res.text)