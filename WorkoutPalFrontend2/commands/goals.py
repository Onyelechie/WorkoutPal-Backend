from api import api_get, api_post
from print_table import print_table


def read_goals_for_user():
    print("\n--- Read Goals for User ---")
    user_id = input("User ID: ").strip()

    res = api_get(f"/users/{user_id}/goals")

    if res.status_code == 200:
        print("✔ Goals retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def create_goal():
    print("\n--- Create Goal ---")
    user_id = int(input("User ID: ").strip())
    name = input("Name: ")
    description = input("Description: ")
    deadline = input("Deadline (string, e.g. 2025-12-31): ")

    payload = {
        "name": name,
        "description": description,
        "deadline": deadline,
    }

    res = api_post(f"/users/{user_id}/goals", json=payload)

    if res.status_code == 201:
        print("✔ Goal created successfully.")
    else:
        print("✖ Creation failed:", res.text)