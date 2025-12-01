from commands.relationships import read_pending_follow_requests, respond_follow_request, follow_user, unfollow_user, \
    send_follow_request, read_follow_request_status, read_followers, read_following


def relationships_loop():
    print("Relationships Page")
    print("Type 'pending', 'respond', 'follow', 'unfollow', 'request', 'cancel_request', 'status', 'followers', 'following', 'help', 'quit', or 'exit'")
    while True:
        cmd = input("\nREL> ").strip().lower()

        if cmd == "pending":
            read_pending_follow_requests()

        elif cmd == "respond":
            respond_follow_request()

        elif cmd == "follow":
            follow_user()

        elif cmd == "unfollow":
            unfollow_user()

        elif cmd == "request":
            send_follow_request()

        elif cmd == "status":
            read_follow_request_status()

        elif cmd == "followers":
            read_followers()

        elif cmd == "following":
            read_following()

        elif cmd in ["quit", "exit"]:
            break

        elif cmd == "help":
            print("Type 'pending', 'respond', 'follow', 'unfollow', 'request', 'cancel_request', 'status', 'followers', 'following', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command. Type 'help' for options.")