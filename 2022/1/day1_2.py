#!/usr/bin/env python

import fileinput


def main():
    top_calories = []
    cur_calories = 0
    for line in fileinput.input():
        if not line.rstrip():
            if len(top_calories) < 3 or any(
                cur_calories > c for c in top_calories
            ):
                top_calories.append(cur_calories)
                while len(top_calories) > 3:
                    top_calories.remove(min(top_calories))
            cur_calories = 0
        else:
            cur_calories += int(line.rstrip())

    print(sum(top_calories))


if __name__ == "__main__":
    main()
