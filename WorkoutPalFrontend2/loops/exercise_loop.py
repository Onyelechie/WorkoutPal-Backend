from commands.exercises import read_exercises, read_exercises_filtered, read_exercise_by_id

def exercises_loop():
    print("Exercises Page")
    print("Type 'read', 'filter', 'read_one', 'create', 'help', 'quit', or 'exit'")
    while True:
        cmd = input("\nEXERCISES> ").strip().lower()

        if cmd == "read":
            read_exercises()

        elif cmd == "filter":
            read_exercises_filtered()

        elif cmd == "read_one":
            read_exercise_by_id()

        elif cmd in ["quit", "exit"]:
            break

        elif cmd == "help":
            print("Type 'read', 'filter', 'read_one', 'create', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command. Type 'help' for options.")