#!/usr/bin/env python

import fileinput


def main():
    calories = []
    cur_calories = 0
    for line in fileinput.input():
        if not line.rstrip():
            calories.append(cur_calories)
            cur_calories = 0
        else:
            cur_calories += int(line.rstrip())

    calories.append(cur_calories)

    print(sum(sorted(calories)[-3:]))


if __name__ == "__main__":
    main()
