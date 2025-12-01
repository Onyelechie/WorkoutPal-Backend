from commands.exercise_setting import read_exercise_setting, create_exercise_setting, update_exercise_setting


def exercise_settings_loop():
    print("Exercise Settings Page")
    print("Type 'read', 'create', 'update', 'help', 'quit', or 'exit'")
    while True:
        cmd = input("\nEXSET> ").strip().lower()

        if cmd == "read":
            read_exercise_setting()

        elif cmd == "create":
            create_exercise_setting()

        elif cmd == "update":
            update_exercise_setting()

        elif cmd in ["quit", "exit"]:
            break

        elif cmd == "help":
            print("Type 'read', 'create', 'update', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command. Type 'help' for options.")