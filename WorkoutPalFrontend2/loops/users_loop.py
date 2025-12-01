from commands.users import read_user_by_id, read_users, update_user

def users_loop():
    print("Users Page")
    print("Type 'read', 'read_one', 'create', 'update', 'delete', 'help', 'quit', or 'exit'")
    while True:
        cmd = input("\nUSERS> ").strip().lower()

        if cmd == "read":
            read_users()

        elif cmd == "read_one":
            read_user_by_id()

        elif cmd == "update":
            update_user()

        elif cmd in ["quit", "exit"]:
            break

        elif cmd == "help":
            print("Type 'read', 'read_one', 'create', 'update', 'delete', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command. Type 'help' for options.")