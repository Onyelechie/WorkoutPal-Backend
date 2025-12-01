from api import api_post, api_get
from print_table import print_table

def read_posts_all():
    print("\n--- Read All Posts ---")
    res = api_get("/posts")

    if res.status_code == 200:
        print("✔ Posts retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_posts_followings():
    print("\n--- Read Posts from Followings ---")
    res = api_get("/posts", params={"followings": True})

    if res.status_code == 200:
        print("✔ Posts retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def read_posts_by_user():
    print("\n--- Read Posts by User ---")
    user_id = input("User ID: ").strip()

    res = api_get(f"/posts/user/{user_id}")

    if res.status_code == 200:
        print("✔ Posts retrieved.")
        print_table(res.json())
    else:
        print("✖ Read failed:", res.text)


def create_post():
    print("\n--- Create Post ---")
    title = input("Title: ")
    caption = input("Caption: ")
    body = input("Body: ")
    posted_by = int(input("Posted by (User ID): ").strip())
    status = input("Status (e.g. public/private): ") or "public"

    payload = {
        "title": title,
        "caption": caption,
        "body": body,
        "postedBy": posted_by,
        "status": status,
    }

    res = api_post("/posts", json=payload)

    if res.status_code == 201:
        print("✔ Post created successfully.")
    else:
        print("✖ Creation failed:", res.text)


def like_post():
    print("\n--- Like Post ---")
    post_id = int(input("Post ID: ").strip())
    user_id = int(input("User ID: ").strip())

    payload = {"postId": post_id, "userId": user_id}

    res = api_post("/posts/like", json=payload)

    if res.status_code == 200:
        print("✔ Post liked.")
    else:
        print("✖ Operation failed:", res.text)


def unlike_post():
    print("\n--- Unlike Post ---")
    post_id = int(input("Post ID: ").strip())
    user_id = int(input("User ID: ").strip())

    payload = {"postId": post_id, "userId": user_id}

    res = api_post("/posts/unlike", json=payload)

    if res.status_code == 200:
        print("✔ Post unliked.")
    else:
        print("✖ Operation failed:", res.text)


def comment_on_post():
    print("\n--- Comment on Post ---")
    post_id = int(input("Post ID: ").strip())
    user_id = int(input("User ID: ").strip())
    comment = input("Comment: ")

    payload = {"postId": post_id, "userId": user_id, "comment": comment}

    res = api_post("/posts/comment", json=payload)

    if res.status_code == 200:
        print("✔ Comment added.")
    else:
        print("✖ Operation failed:", res.text)


def reply_to_comment():
    print("\n--- Reply to Comment ---")
    post_id = int(input("Post ID: ").strip())
    comment_id = int(input("Comment ID: ").strip())
    user_id = int(input("User ID: ").strip())
    reply_text = input("Reply: ")

    payload = {
        "postId": post_id,
        "commentId": comment_id,
        "userId": user_id,
        "comment": reply_text,
    }

    res = api_post("/posts/comment/reply", json=payload)

    if res.status_code == 200:
        print("✔ Reply added.")
    else:
        print("✖ Operation failed:", res.text)