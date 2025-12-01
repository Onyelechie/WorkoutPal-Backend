from api import api_get, api_post
from print_table import print_table


def read_pending_follow_requests():
    print("\n--- Read Pending Follow Requests ---")
    res = api_get("/follow-requests")

    if res.status_code == 200:
        print("✔ Pending follow requests retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def respond_follow_request():
    print("\n--- Respond to Follow Request ---")
    request_id = int(input("Request ID: ").strip())
    action = input("Action (accept/reject): ").strip().lower()

    payload = {
        "requestID": request_id,
        "action": action,
    }

    res = api_post("/follow-requests/respond", json=payload)

    if res.status_code == 200:
        print("✔ Follow request processed.")
    else:
        print("✖ Operation failed:", res.text)


def follow_user():
    print("\n--- Follow User ---")
    target_id = input("Target User ID: ").strip()
    follower_id = input("Follower User ID: ").strip()

    res = api_post(f"/users/{target_id}/follow",
                   params={"follower_id": follower_id})

    if res.status_code == 200:
        print("✔ Successfully followed user.")
    else:
        print("✖ Operation failed:", res.text)


def unfollow_user():
    print("\n--- Unfollow User ---")
    target_id = input("Target User ID: ").strip()
    follower_id = input("Follower User ID: ").strip()

    res = api_post(f"/users/{target_id}/unfollow",
                   params={"follower_id": follower_id})

    if res.status_code == 200:
        print("✔ Successfully unfollowed user.")
    else:
        print("✖ Operation failed:", res.text)


def send_follow_request():
    print("\n--- Send Follow Request ---")
    target_id = input("Target User ID: ").strip()
    requester_id = input("Requester User ID: ").strip()

    res = api_post(f"/users/{target_id}/follow-request",
                   params={"requester_id": requester_id})

    if res.status_code == 200:
        print("✔ Follow request sent.")
    else:
        print("✖ Operation failed:", res.text)


def read_follow_request_status():
    print("\n--- Read Follow Request Status ---")
    target_id = input("Target User ID: ").strip()
    requester_id = input("Requester User ID: ").strip()

    res = api_get(f"/users/{target_id}/follow-request/status",
                  params={"requester_id": requester_id})

    if res.status_code == 200:
        print("✔ Follow request status:")
        print(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_followers():
    print("\n--- Read Followers ---")
    user_id = input("User ID: ").strip()

    res = api_get(f"/users/{user_id}/followers")

    if res.status_code == 200:
        print("✔ Followers retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_following():
    print("\n--- Read Following ---")
    user_id = input("User ID: ").strip()

    res = api_get(f"/users/{user_id}/following")

    if res.status_code == 200:
        print("✔ Following retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)