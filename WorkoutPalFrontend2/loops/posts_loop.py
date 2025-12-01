from commands.posts import read_posts_all, read_posts_followings, read_posts_by_user, create_post, like_post, \
    unlike_post, comment_on_post, reply_to_comment

def posts_loop():
    print("Posts Page")
    print("Type 'read', 'read_followings', 'read_user', 'create', 'like', 'unlike', 'comment', 'reply', 'delete', 'help', 'quit', or 'exit'")
    while True:
        cmd = input("\nPOSTS> ").strip().lower()

        if cmd == "read":
            read_posts_all()

        elif cmd == "read_followings":
            read_posts_followings()

        elif cmd == "read_user":
            read_posts_by_user()

        elif cmd == "create":
            create_post()

        elif cmd == "like":
            like_post()

        elif cmd == "unlike":
            unlike_post()

        elif cmd == "comment":
            comment_on_post()

        elif cmd == "reply":
            reply_to_comment()

        elif cmd in ["quit", "exit"]:
            break

        elif cmd == "help":
            print("Type 'read', 'read_followings', 'read_user', 'create', 'like', 'unlike', 'comment', 'reply', 'delete', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command. Type 'help' for options.")