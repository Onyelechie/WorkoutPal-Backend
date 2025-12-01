from commands.schedules import read_schedules, read_schedules_by_day, read_schedule_by_id, create_schedule, \
    update_schedule


def schedules_loop():
    print("Schedules Page")
    print("Type 'read', 'read_day', 'read_one', 'create', 'update', 'delete', 'help', 'quit', or 'exit'")
    while True:
        cmd = input("\nSCHEDULES> ").strip().lower()

        if cmd == "read":
            read_schedules()

        elif cmd == "read_day":
            read_schedules_by_day()

        elif cmd == "read_one":
            read_schedule_by_id()

        elif cmd == "create":
            create_schedule()

        elif cmd == "update":
            update_schedule()

        elif cmd in ["quit", "exit"]:
            break

        elif cmd == "help":
            print("Type 'read', 'read_day', 'read_one', 'create', 'update', 'delete', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command. Type 'help' for options.")