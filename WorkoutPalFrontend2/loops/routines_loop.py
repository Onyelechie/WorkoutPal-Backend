from commands.routines import read_routine_by_id, read_user_routines, create_routine, add_exercise_to_routine


def routines_loop():
    print("Routines Page")
    print("Type 'read_user', 'read_one', 'create', 'delete', 'delete_user', 'add_ex', 'remove_ex', 'help', 'quit', or 'exit'")
    while True:
        cmd = input("\nROUTINES> ").strip().lower()

        if cmd == "read_user":
            read_user_routines()

        elif cmd == "read_one":
            read_routine_by_id()

        elif cmd == "create":
            create_routine()

        elif cmd == "add_ex":
            add_exercise_to_routine()

        elif cmd in ["quit", "exit"]:
            break

        elif cmd == "help":
            print("Type 'read_user', 'read_one', 'create', 'delete', 'delete_user', 'add_ex', 'remove_ex', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command. Type 'help' for options.")