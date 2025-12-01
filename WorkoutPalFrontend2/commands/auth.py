import getpass
from api import api_post, api_get

def login_user():
    print("\n--- Login ---")
    email = input("Email: ").strip()
    password = getpass.getpass("Password: ")

    payload = {
        "email": email,
        "password": password,
    }

    res = api_post("/login", json=payload)

    if res.status_code == 200:
        print("✔ Login successful.")
        return True
    else:
        print("✖ Login failed:", res.text)
        return False

def logout_user():
    print("\n--- Logout ---")
    res = api_post("/logout")

    if res.status_code == 200:
        print("✔ Logout successful.")
    else:
        print("✖ Logout failed:", res.text)


def read_me():
    print("\n--- Current User (/me) ---")
    res = api_get("/me")

    if res.status_code == 200:
        print("✔ Fetched current user.")
        print(res.json())
    else:
        print("✖ Request failed:", res.text)