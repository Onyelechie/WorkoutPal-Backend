import getpass

import bcrypt

from api import api_post, api_get, api_patch
from print_table import print_table


def create_user():
    print("\n--- Create User ---")
    name = input("Name: ")
    username = input("Username: ")
    email = input("Email: ")
    password = getpass.getpass("Password: ")

    pw_hash = bcrypt.hashpw(password.encode(), bcrypt.gensalt()).decode()

    payload = {
        "name": name,
        "username": username,
        "email": email,
        "password": pw_hash
    }

    res = api_post("/users", json=payload)

    if res.status_code == 201:
        print("✔ User created successfully.")
    else:
        print("✖ Creation failed:", res.text)


def read_users():
    print("\n--- Read Users ---")
    res = api_get("/users")

    if res.status_code == 200:
        print("✔ Users retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_user_by_id():
    print("\n--- Read User by ID ---")
    user_id = input("User ID: ").strip()

    res = api_get(f"/users/{user_id}")

    if res.status_code == 200:
        print("✔ User retrieved.")
        print(res.json())
    else:
        print("✖ Read failed:", res.text)


def update_user():
    print("\n--- Update User ---")
    user_id = input("User ID: ").strip()

    print("Leave a field blank to skip updating it.")
    name = input("Name: ").strip()
    username = input("Username: ").strip()
    email = input("Email: ").strip()
    avatar = input("Avatar URL: ").strip()
    is_private = input("Is private? (y/n or blank): ").strip().lower()
    show_metrics = input("Show metrics to followers? (y/n or blank): ").strip().lower()

    password = getpass.getpass("New Password (leave blank to keep current): ")

    payload = {}
    if name:
        payload["name"] = name
    if username:
        payload["username"] = username
    if email:
        payload["email"] = email
    if avatar:
        payload["avatar"] = avatar
    if is_private in ("y", "yes", "true", "1"):
        payload["isPrivate"] = True
    elif is_private in ("n", "no", "false", "0"):
        payload["isPrivate"] = False
    if show_metrics in ("y", "yes", "true", "1"):
        payload["showMetricsToFollowers"] = True
    elif show_metrics in ("n", "no", "false", "0"):
        payload["showMetricsToFollowers"] = False
    if password:
        payload["password"] = bcrypt.hashpw(password.encode(), bcrypt.gensalt()).decode()

    res = api_patch(f"/users/{user_id}", json=payload)

    if res.status_code == 200:
        print("✔ User updated.")
        print(res.json())
    else:
        print("✖ Update failed:", res.text)