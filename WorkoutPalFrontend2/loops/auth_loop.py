from commands.auth import logout_user, read_me, login_user
from commands.users import create_user
from loops.app_loop import app_loop


def auth_loop():
    print("Welcome to WorkoutPal CLI")
    print("Type 'login', 'logout', 'me', 'google', 'help', 'quit', or 'exit'")

    while True:
        cmd = input("\nAUTH> ").strip().lower()

        if cmd == "login":
            if login_user():
                app_loop()

        elif cmd == "logout":
            logout_user()

        elif cmd == "me":
            read_me()

        elif cmd == "create user":
            create_user()

        elif cmd in ["quit", "exit"]:
            print("Goodbye.")
            break

        elif cmd == "help":
            print("Type 'login', 'logout', 'me', 'google', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command:", cmd)