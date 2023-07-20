#!/usr/bin/env python

import fileinput


def main():
    max_calories = 0
    cur_calories = 0
    for line in fileinput.input():
        if not line.rstrip():
            max_calories = max(max_calories, cur_calories)
            cur_calories = 0
        else:
            cur_calories += int(line.rstrip())

    print(max_calories)


if __name__ == "__main__":
    main()
