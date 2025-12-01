from commands.goals import read_goals_for_user, create_goal

def goals_loop():
    print("Goals Page")
    print("Type 'read', 'create', 'help', 'quit', or 'exit'")
    while True:
        cmd = input("\nGOALS> ").strip().lower()

        if cmd == "read":
            read_goals_for_user()

        elif cmd == "create":
            create_goal()

        elif cmd in ["quit", "exit"]:
            break

        elif cmd == "help":
            print("Type 'read', 'create', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command. Type 'help' for options.")