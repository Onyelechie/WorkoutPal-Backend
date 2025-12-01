from tabulate import tabulate

MAX_WIDTH = 30
MAX_ROWS = 20
MAX_COLS = 6

def truncate(value, width=MAX_WIDTH):
    s = str(value)
    return s if len(s) <= width else s[:width - 3] + "..."

def print_table(json_array):
    if not json_array:
        print("No data.")
        return

    all_headers = list(json_array[0].keys())
    headers = all_headers[:MAX_COLS]

    limited_rows = json_array[:MAX_ROWS]

    rows = [
        [truncate(row.get(h, "")) for h in headers]
        for row in limited_rows
    ]

    print(tabulate(rows, headers, tablefmt="grid"))

    if len(all_headers) > MAX_COLS:
        print(f"\n… {len(all_headers) - MAX_COLS} more columns not shown.")

    if len(json_array) > MAX_ROWS:
        print(f"… {len(json_array) - MAX_ROWS} more rows not shown.")
