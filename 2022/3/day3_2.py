#!/usr/bin/env python

import fileinput


def main():
    total = 0
    common = set()
    rucksacks = 0

    for line in fileinput.input():
        rucksacks += 1
        if not common:
            common = set(line.strip())
        else:
            common &= set(line.strip())

        if rucksacks == 3:
            badge = list(common)[0]
            if badge < "a":
                total += ord(badge) - 38
            else:
                total += ord(badge) - 96
            rucksacks = 0
            common = set()
    print(total)


if __name__ == "__main__":
    main()
