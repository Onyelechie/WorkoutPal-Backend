from commands.achievements import read_achievements_catalog, read_achievements_feed, read_my_unlocked_achievements, \
    read_user_unlocked_achievements, unlock_achievement

def achievements_loop():
    print("Achievements Page")
    print("Type 'catalog', 'feed', 'unlocked_me', 'unlocked_user', 'unlock', 'help', 'quit', or 'exit'")
    while True:
        cmd = input("\nACHIEVEMENTS> ").strip().lower()

        if cmd == "catalog":
            read_achievements_catalog()

        elif cmd == "feed":
            read_achievements_feed()

        elif cmd == "unlocked_me":
            read_my_unlocked_achievements()

        elif cmd == "unlocked_user":
            read_user_unlocked_achievements()

        elif cmd == "unlock":
            unlock_achievement()

        elif cmd in ["quit", "exit"]:
            break

        elif cmd == "help":
            print("Type 'catalog', 'feed', 'unlocked_me', 'unlocked_user', 'unlock', 'help', 'quit', or 'exit'")

        elif cmd == "":
            continue

        else:
            print("Unknown command. Type 'help' for options.")