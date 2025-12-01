from api import api_get, api_post
from print_table import print_table


def read_achievements_catalog():
    print("\n--- Read Achievements Catalog ---")
    res = api_get("/achievements")

    if res.status_code == 200:
        print("✔ Achievements retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_achievements_feed():
    print("\n--- Read Achievements Feed ---")
    res = api_get("/achievements/feed")

    if res.status_code == 200:
        print("✔ Achievements feed retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_my_unlocked_achievements():
    print("\n--- Read My Unlocked Achievements ---")
    res = api_get("/achievements/unlocked")

    if res.status_code == 200:
        print("✔ Unlocked achievements retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_user_unlocked_achievements():
    print("\n--- Read User's Unlocked Achievements ---")
    user_id = input("User ID: ").strip()

    res = api_get(f"/achievements/unlocked/{user_id}")

    if res.status_code == 200:
        print("✔ Unlocked achievements retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def unlock_achievement():
    print("\n--- Unlock Achievement ---")
    user_id = int(input("User ID: ").strip())
    achievement_id = int(input("Achievement ID: ").strip())

    payload = {
        "userId": user_id,
        "achievementId": achievement_id,
    }

    res = api_post("/achievements", json=payload)

    if res.status_code == 201:
        print("✔ Achievement unlocked successfully.")
    else:
        print("✖ Unlock failed:", res.text)